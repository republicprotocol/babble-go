package grpc

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/republicprotocol/xoxo-go/foundation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// ErrRateLimitExceeded is returned when a client has attempted to many requests
// over a period of time.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// ErrMalformedTCPAddress is returned when a server cannot determine the TCP
// address of a client.
var ErrMalformedTCPAddress = errors.New("malformed tcp address")

// Dial a `net.Addr` to create an insecure connection to a remote server at that
// `net.Addr`. A `context.Context` can be used to cancel or expire the pending
// connection. A call to `grpc.ClientConn.Close` is required to free all
// allocated resources.
func Dial(ctx context.Context, addr net.Addr) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Backoff calling the `f` function until the `context.Context` is done, or the
// `f` function returns a nil error. The delay increases by 60% but will not
// exceed beyond the `maxBackoffDelayInMs`.
func Backoff(ctx context.Context, f func() error, maxBackoffDelayInMs time.Duration) error {
	timeoutMs := time.Duration(1000)
	for {
		err := f()
		if err == nil {
			return nil
		}
		timer := time.NewTimer(time.Millisecond * timeoutMs)
		select {
		case <-ctx.Done():
			return err
		case <-timer.C:
			timeoutMs = time.Duration(float64(timeoutMs) * 1.6)
			if timeoutMs > maxBackoffDelayInMs {
				timeoutMs = maxBackoffDelayInMs
			}
		}
	}
}

type client struct {
}

// NewClient returns an implementation of the `gossip.Client` interface that
// uses gRPC to invoke RPCs.
func NewClient(addr net.Addr) gossip.Client {
	return &client{}
}

// Send a `message` to the `to` address. A `context.Context` can be used to
// cancel or expire the request. The client will backoff the request with a
// maximum delay of one minute.
func (client *client) Send(ctx context.Context, to net.Addr, message foundation.Message) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := &SendRequest{
		Nonce:     message.Nonce,
		Key:       message.Key,
		Value:     message.Value,
		Signature: message.Signature,
	}

	return Backoff(ctx, func() (err error) {
		_, err = NewXoxoServiceClient(conn).Send(ctx, request)
		return
	}, time.Minute)
}

// Service implements a gRPC Service that accepts RPCs from clients. It
// delegates requests to a `gossip.Server` after enforcing rate limits.
type Service struct {
	server gossip.Server

	rateLimitsMu *sync.Mutex
	rateLimits   map[string]time.Time
	rate         time.Duration
}

// NewService returns a Service that delegates requests to the `server` and uses
// `rate` to enforce rate limits of all RPCs.
func NewService(server gossip.Server, rate time.Duration) Service {
	return Service{
		server: server,

		rateLimitsMu: new(sync.Mutex),
		rateLimits:   make(map[string]time.Time),
		rate:         rate,
	}
}

// Register the service to a `grpc.Server`.
func (service *Service) Register(server *grpc.Server) {
	RegisterXoxoServiceServer(server, service)
}

// Send implements the respective gRPC call.
func (service *Service) Send(ctx context.Context, request *SendRequest) (*SendResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}
	message := foundation.Message{
		Nonce:     request.Nonce,
		Key:       request.Key,
		Value:     request.Value,
		Signature: request.Signature,
	}
	return &SendResponse{}, service.server.Receive(ctx, message)
}

func (service *Service) isRateLimited(ctx context.Context) error {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return ErrMalformedTCPAddress
	}
	if client.Addr == net.Addr(nil) {
		return ErrMalformedTCPAddress
	}

	clientAddr, ok := client.Addr.(*net.TCPAddr)
	if !ok {
		return ErrMalformedTCPAddress
	}
	clientIP := clientAddr.IP.String()

	service.rateLimitsMu.Lock()
	defer service.rateLimitsMu.Unlock()

	if lastPing, ok := service.rateLimits[clientIP]; ok {
		if service.rate > time.Since(lastPing) {
			return ErrRateLimitExceeded
		}
	}

	service.rateLimits[clientIP] = time.Now()
	return nil
}

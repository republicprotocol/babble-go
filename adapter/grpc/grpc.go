package grpc

import (
	"context"
	"errors"
	"fmt"
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

// Dial a net.Addr to create an insecure connection to a remote server at that
// net.Addr. A context can be used to cancel or expire the pending connection. A
// call to grpc.ClientConn.Close is required to free all allocated resources.
func Dial(ctx context.Context, addr net.Addr) (*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, addr.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Backoff calling a function until the context.Context is done, or the function
// returns a nil error.
func Backoff(ctx context.Context, f func() error, maxBackoffDelay time.Duration) error {
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
			if timeoutMs > maxBackoffDelay {
				timeoutMs = maxBackoffDelay
			}
		}
	}
}

type client struct {
}

// NewClient returns an implementation of the gossip.Client interface that uses
// gRPC to invoke RPCs.
func NewClient(addr net.Addr) gossip.Client {
	return &client{}
}

func (client *client) Send(ctx context.Context, to net.Addr, message foundation.Message) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := &BroadcastRequest{
		Nonce:     message.Nonce,
		Key:       message.Key,
		Value:     message.Value,
		Signature: message.Signature,
	}

	return Backoff(ctx, func() (err error) {
		_, err = NewXoxoServiceClient(conn).Broadcast(ctx, request)
		return
	}, time.Minute)
}

// Service is a Service that implements a gRPC Service that accepts RPCs. It
// delegates request to a gossip.Server.
type Service struct {
	server gossip.Server

	rate         time.Duration
	rateLimitsMu *sync.Mutex
	rateLimits   map[string]time.Time
}

func NewService(server gossip.Server, rate time.Duration) Service {
	return Service{
		server: server,

		rate:         rate,
		rateLimitsMu: new(sync.Mutex),
		rateLimits:   make(map[string]time.Time),
	}
}

func (service *Service) Register(server *grpc.Server) {
	RegisterXoxoServiceServer(server, service)
}

func (service *Service) Broadcast(ctx context.Context, request *BroadcastRequest) (*BroadcastResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}
	message := foundation.Message{
		Nonce:     request.Nonce,
		Key:       request.Key,
		Value:     request.Value,
		Signature: request.Signature,
	}
	return &BroadcastResponse{}, service.server.Receive(ctx, message)
}

func (service *Service) isRateLimited(ctx context.Context) error {
	client, ok := peer.FromContext(ctx)
	if !ok {
		return fmt.Errorf("fail to get peer from ctx")
	}
	if client.Addr == net.Addr(nil) {
		return fmt.Errorf("fail to get peer address")
	}

	clientAddr, ok := client.Addr.(*net.TCPAddr)
	if !ok {
		return fmt.Errorf("fail to read peer TCP address")
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

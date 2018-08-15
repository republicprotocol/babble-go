package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/republicprotocol/gossip-network/core"
	"github.com/republicprotocol/gossip-network/foundation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// ErrRateLimitExceeded is returned when the same client sends more than one
// request to the server within a specified rate limit.
var ErrRateLimitExceeded = errors.New("cannot process request, rate limit exceeded")

// Dial creates a client connection to the given net.Addr. A context can be
// used to cancel or expire the pending connection. Once this function returns,
// the cancellation and expiration of the Context will do nothing. Users must
// call grpc.ClientConn.Close to terminate all the pending operations after
// this function returns.
func Dial(ctx context.Context, addr net.Addr) (*grpc.ClientConn, error) {
	clientConn, err := grpc.DialContext(ctx, addr.String(), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return clientConn, nil
}

// Backoff a function call until the context.Context is done, or the function
// returns nil.
func Backoff(ctx context.Context, f func() error) error {
	timeoutMs := time.Duration(1000)
	for {
		err := f()
		if err == nil {
			return nil
		}
		timer := time.NewTimer(time.Millisecond * timeoutMs)
		select {
		case <-ctx.Done():
			return fmt.Errorf("backoff timeout = %v: %v", ctx.Err(), err)
		case <-timer.C:
			timeoutMs = time.Duration(float64(timeoutMs) * 1.6)
		}
	}
}

type gossipClient struct {
	addr net.Addr
}

// NewGossipClient returns an implementation of the Gossip.Client interface that
// uses gRPC and a recycled connection pool.
func NewGossipClient(addr net.Addr) core.Client {
	return &gossipClient{
		addr: addr,
	}
}

func (client *gossipClient) Send(ctx context.Context, to net.Addr, message foundation.Message) error {
	conn, err := Dial(ctx, to)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := &RumorRequest{
		Data:      message.Data(),
		Signature: message.Signature(),
		Nonce:     message.Nonce(),
	}

	return Backoff(ctx, func() error {
		_, err = NewGossipServiceClient(conn).Gossip(ctx, request)
		return err
	})
}

// GossipService is a Service that implements the gRPC GossipService defined in
// protobuf. It delegates responsibility for handling the Ping and Query RPCs
// to a Gossip.Server.
type GossipService struct {
	server core.Server

	rate         time.Duration
	rateLimitsMu *sync.Mutex
	rateLimits   map[string]time.Time
}

// NewGossipService returns a GossipService that uses the Gossip.Server as a
// delegate.
func NewGossipService(server core.Server, rate time.Duration) GossipService {
	return GossipService{
		server: server,

		rate:         rate,
		rateLimitsMu: new(sync.Mutex),
		rateLimits:   make(map[string]time.Time),
	}
}

func (service *GossipService) Send(ctx context.Context, rumor RumorRequest) (*RumorResponse, error) {
	if err := service.isRateLimited(ctx); err != nil {
		return nil, err
	}

	// todo
	return nil, nil
}

func (service *GossipService) isRateLimited(ctx context.Context) error {
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

package rpc

import (
	"context"
	"net"
	"time"

	"github.com/republicprotocol/xoxo-go/core/gossip"
	"google.golang.org/grpc"
)

// Dialer is used to open a connection to a gRPC server.
type Dialer interface {

	// Dial a connection to a remote gRPC server. When dialing fails, it's up
	// to the implementation to backoff and retry or retrun error directly.
	Dial(ctx context.Context, to net.Addr) (*grpc.ClientConn, error)
}

// Caller is used to call gRPC procedures from a client.
type Caller interface {

	// Call the gRPC procedure defined by `f`. When the call fails, it's up to
	// the implementation to backoff and retry or return the error directly.
	Call(ctx context.Context, f func() error) error
}

type client struct {
	Dialer
	Caller
}

// NewClient returns an implementation of the `gossip.Client` interface that
// uses gRPC to invoke RPCs.
func NewClient(dialer Dialer, caller Caller) gossip.Client {
	return &client{dialer, caller}
}

// Send a `message` to the `to` address. A `context.Context` can be used to
// cancel or expire the request. The client will backoff the request with a
// maximum delay of one minute.
func (client *client) Send(ctx context.Context, to net.Addr, message gossip.Message) error {
	conn, err := client.Dial(ctx, to)
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

	return client.Call(ctx, func() error {
		_, err = NewXoxoClient(conn).Send(ctx, request)
		return err
	})
}

// Service implements a gRPC Service that accepts RPCs from clients. It
// delegates requests to a `gossip.Server` after enforcing rate limits.
type Service struct {
	server gossip.Server
}

// NewService returns a Service that delegates requests to the `server` and uses
// `rate` to enforce rate limits of all RPCs.
func NewService(server gossip.Server) Service {
	return Service{
		server: server,
	}
}

// Register the service to a `grpc.Server`.
func (service *Service) Register(server *grpc.Server) {
	RegisterXoxoServer(server, service)
}

// Send implements the respective gRPC call.
func (service *Service) Send(ctx context.Context, request *SendRequest) (*SendResponse, error) {
	message := gossip.Message{
		Nonce:     request.Nonce,
		Key:       request.Key,
		Value:     request.Value,
		Signature: request.Signature,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	return &SendResponse{}, service.server.Receive(ctx, message)
}

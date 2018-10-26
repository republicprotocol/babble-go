package grpc

import (
	"context"
	"net"
	"time"

	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/republicprotocol/xoxo-go/foundation"
	"google.golang.org/grpc"
)

type Dialer interface {
	Dial(ctx context.Context, to net.Addr) (grpc.ClientConn, error)
}

type Caller interface {
	// TODO: Define and use.
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
func (client *client) Send(ctx context.Context, to net.Addr, message foundation.Message) error {
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

	// TODO: Backoff should be part of the GRPC driver and is used by Dialer
	// and Caller.
	return Backoff(ctx, func() (err error) {
		_, err = NewXoxoClient(conn).Send(ctx, request)
		return
	}, time.Minute)
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
	message := foundation.Message{
		Nonce:     request.Nonce,
		Key:       request.Key,
		Value:     request.Value,
		Signature: request.Signature,
	}

	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	return &SendResponse{}, service.server.Receive(ctx, message)
}

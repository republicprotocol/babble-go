package testutils

import (
	"context"
	"net"

	"google.golang.org/grpc"
)

type MockDialer struct {
}

func (dialer MockDialer) Dial(ctx context.Context, to net.Addr) (*grpc.ClientConn, error) {
	return grpc.Dial(to.String(), grpc.WithInsecure())
}

type FaultyDialer struct {
}

func (dialer FaultyDialer) Dial(ctx context.Context, to net.Addr) (*grpc.ClientConn, error) {
	return grpc.Dial(to.String())
}

type MockCaller struct {
}

func (caller MockCaller) Call(ctx context.Context, f func() error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return f()
	}
}

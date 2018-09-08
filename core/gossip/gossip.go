package gossip

import (
	"context"
	"net"
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/xoxo-go/foundation"
)

// A Client sends Messages to a net.Addr.
type Client interface {

	// Send a `message` to the `to` address.
	Send(ctx context.Context, to net.Addr, message foundation.Message) error
}

// A Server receives Messages.
type Server interface {

	// Receive a `message` from a client.
	Receive(ctx context.Context, message foundation.Message) error
}

type Gossiper interface {
	Server
	Broadcast(ctx context.Context, message foundation.Message) error
}

type gossiper struct {
	α        int
	client   Client
	verifier Verifier
	store    Store
}

func NewGossiper(α int, client Client, verifier Verifier, store Store) Gossiper {
	return &gossiper{
		α:        α,
		client:   client,
		verifier: verifier,
		store:    store,
	}
}

func (gossiper *gossiper) Broadcast(ctx context.Context, message foundation.Message) error {

	addrs, err := gossiper.store.Addrs(gossiper.α)
	if err != nil {
		return err
	}

	errs := make([]error, len(addrs))
	co.ForAll(addrs, func(i int) {
		errs[i] = gossiper.client.Send(ctx, addrs[i], message)
	})

	for i := range errs {
		if errs[i] != nil {
			return errs[i]
		}
	}

	return nil
}

func (gossiper *gossiper) Receive(ctx context.Context, message foundation.Message) error {
	if err := gossiper.verifier.Verify(message.Value, message.Signature); err != nil {
		return err
	}

	previousMessage, err := gossiper.store.Message(message.Key)
	if err != nil {
		return err
	}

	if previousMessage.Nonce >= message.Nonce {
		return nil
	}

	if err := gossiper.store.InsertMessage(message); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return gossiper.Broadcast(ctx, message)
}

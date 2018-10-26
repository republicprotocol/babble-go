package gossip

import (
	"context"
	"net"

	"github.com/republicprotocol/co-go"
)

// A Client is used to send Messages to a remote Server.
type Client interface {

	// Send a Message to the a remote `net.Addr`.
	Send(ctx context.Context, to net.Addr, message Message) error
}

// A Server receives Messages.
type Server interface {

	// Receive is called to notify the Server that a Message has been received
	// from a remote Client.
	Receive(ctx context.Context, message Message) error
}

type Gossiper interface {
	Server
	Broadcast(ctx context.Context, message Message) error
}

type gossiper struct {
	α        int
	client   Client
	signer   Signer
	verifier Verifier
	store    Store
}

func NewGossiper(α int, client Client, signer Signer, verifier Verifier, store Store) Gossiper {
	return &gossiper{
		α:        α,
		client:   client,
		verifier: verifier,
		store:    store,
	}
}

func (gossiper *gossiper) Broadcast(ctx context.Context, message Message) error {
	signature, err := gossiper.signer.Sign(message.Value)
	if err != nil {
		return err
	}
	message.Signature = signature

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

func (gossiper *gossiper) Receive(ctx context.Context, message Message) error {
	if err := gossiper.verifier.Verify(message.Value, message.Signature); err != nil {
		return err
	}

	previousMessage, err := gossiper.store.Message(message.Key)
	if err != nil && err != ErrMessageNotFound {
		return err
	}
	if err == nil && previousMessage.Nonce >= message.Nonce {
		return nil
	}
	if err := gossiper.store.InsertMessage(message); err != nil {
		return err
	}

	return gossiper.Broadcast(ctx, message)
}

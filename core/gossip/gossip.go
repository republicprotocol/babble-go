package gossip

import (
	"context"
	"net"
	"time"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/xoxo-go/foundation"
)

type Client interface {
	Send(ctx context.Context, to net.Addr, message foundation.Message) error
}

type Server interface {
	Receive(ctx context.Context, message foundation.Message) error
}

type Gossiper interface {
	Server
	Broadcast(ctx context.Context, message foundation.Message) error
}

type gossiper struct {
	client   Client
	verifier Verifier
	addrs    AddrStore
	messages MessageStore
}

func NewGossiper(client Client, verifier Verifier, addrs AddrStore, messages MessageStore) Gossiper {
	return &gossiper{
		client:   client,
		verifier: verifier,
		addrs:    addrs,
		messages: messages,
	}
}

func (gossiper *gossiper) Broadcast(ctx context.Context, message foundation.Message) error {

	addrs, err := gossiper.addrs.Addrs()
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

	previousMessage, err := gossiper.messages.Message(message.Key)
	if err != nil {
		return err
	}

	if previousMessage.Nonce >= message.Nonce {
		return nil
	}

	if err := gossiper.messages.InsertMessage(message); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return gossiper.Broadcast(ctx, message)
}

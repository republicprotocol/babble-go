package core

import (
	"context"
	"net"
	"math/rand"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/gossip-network/foundation"
)

type Client interface {
	Send(ctx context.Context, to net.Addr, message foundation.Message) error
}

type Gossiper interface {
	Send(ctx context.Context, message foundation.Message) error
}

type gossiper struct {
	client Client
	alpha  int
	store  AddressStorer
}

func newGossiper(client Client, alpha int, store AddressStorer) Gossiper {
	return &gossiper{
		client: client,
		alpha:  alpha,
		store:  store,
	}
}

func (gossiper *gossiper) Send(ctx context.Context, message foundation.Message) error {
	targets, err := randomTargets(gossiper.store, gossiper.alpha)
	if err != nil {
		return err
	}
	errs := make([]error, len(targets))
	co.ForAll(targets, func(i int) {
		errs[i] = gossiper.client.Send(ctx, targets[i], message)
	})

	for i := range errs {
		if errs[i] != nil {
			return errs[i]
		}
	}

	return nil
}

type Server interface {
	Send(ctx context.Context, message foundation.Message) error
}

type server struct {
	gossiper Gossiper
	store    MessageStore
	verifier Verifier
	alpha    int
}

func NewServer(gossiper Gossiper, store MessageStore, verifier Verifier, alpha int) Server {
	return &server{
		gossiper: gossiper,
		store:    store,
		verifier: verifier,
		alpha:    alpha,
	}
}

func (server *server) Send(ctx context.Context, message foundation.Message) error {
	// Verify the message
	if err := server.verifier.Verify(message.Data(), message.Signature()); err != nil {
		return err
	}

	// Check nonce
	// todo : hash of which field?
	old, err := server.store.Message([32]byte{})
	if err != nil {
		return err
	}
	if message.Nonce() > old.Nonce() {
		if err := server.store.InsertMessage(message); err != nil {
			return err
		}

		return server.gossiper.Send(ctx, message)
	}

	return nil
}

// randomTargets randomly selects at most alpha addresses from the storer.
func randomTargets(store AddressStorer, alpha int) ([]net.Addr, error) {
	addrs, err := store.Addresses()
	if err != nil {
		return nil, nil
	}
	if len(addrs) <= alpha {
		return addrs, nil
	}

	// Randomly select α multi-addresses.
	results := make([]net.Addr, 0, alpha)
	randomIndex := rand.Perm(len(addrs))
	for i := 0; i < alpha; i++ {
		results = append(results, addrs[randomIndex[i]])
	}
	return results, nil
}

package testutils

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"

	"github.com/republicprotocol/babble-go/core/addr"
	"github.com/republicprotocol/babble-go/core/gossip"
)

// MockAddr is a mock implementation of the `addr.Store`
type MockAddr struct {
	addresses map[string]net.Addr
}

func NewMockAddrs() addr.Store {
	return MockAddr{
		addresses: map[string]net.Addr{},
	}
}

func (store MockAddr) InsertAddr(addr net.Addr) error {
	store.addresses[addr.String()] = addr
	return nil
}

func (store MockAddr) Addrs() ([]net.Addr, error) {
	addresses := make([]net.Addr, 0, len(store.addresses))
	for _, j := range store.addresses {
		addresses = append(addresses, j)
	}

	return addresses, nil
}

type MockMessages struct {
	messageMu *sync.Mutex
	messages  map[string]gossip.Message
}

func NewMockMessages() MockMessages {
	return MockMessages{
		messageMu: new(sync.Mutex),
		messages:  map[string]gossip.Message{},
	}
}

func (store MockMessages) InsertMessage(message gossip.Message) error {
	store.messageMu.Lock()
	defer store.messageMu.Unlock()
	store.messages[string(message.Key)] = message

	return nil
}

func (store MockMessages) Message(key []byte) (gossip.Message, error) {
	store.messageMu.Lock()
	defer store.messageMu.Unlock()

	return store.messages[string(key)], nil
}

func RandomAddr() net.Addr {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v.%v.%v.%v:%v", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(8000)))
	if err != nil {
		log.Fatal(err)
	}

	return addr
}

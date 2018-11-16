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

// MockAddrStore is a mock implementation of the `addr.Store`
type MockAddrStore struct {
	addresses map[string]net.Addr
}

func NewMockAddrStore() addr.Store {
	return MockAddrStore{
		addresses: map[string]net.Addr{},
	}
}

func (store MockAddrStore) InsertAddr(addr net.Addr) error {
	store.addresses[addr.String()] = addr
	return nil
}

func (store MockAddrStore) Addrs() ([]net.Addr, error) {
	addresses := make([]net.Addr, 0, len(store.addresses))
	for _, j := range store.addresses {
		addresses = append(addresses, j)
	}

	return addresses, nil
}

type MockMessageStore struct {
	messageMu *sync.Mutex
	messages  map[string]gossip.Message
}

func NewMockMessageStore() MockMessageStore {
	return MockMessageStore{
		messageMu: new(sync.Mutex),
		messages:  map[string]gossip.Message{},
	}
}

func (store MockMessageStore) InsertMessage(message gossip.Message) error {
	store.messageMu.Lock()
	defer store.messageMu.Unlock()
	store.messages[string(message.Key)] = message

	return nil
}

func (store MockMessageStore) Message(key []byte) (gossip.Message, error) {
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

package grpc_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/xoxo-go/driver/grpc"

	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/republicprotocol/xoxo-go/foundation"
	"google.golang.org/grpc"
)

var _ = Describe("grpc", func() {

	initService := func(n int) ([]gossip.Client, []gossip.Store, []*grpc.Server, []net.Listener) {
		clients := make([]gossip.Client, n)
		stores := make([]gossip.Store, n)
		servers := make([]*grpc.Server, n)
		listeners := make([]net.Listener, n)

		for i := 0; i < n; i++ {
			clients[i] = NewClient()

			store := NewMockStore()
			for j := 0; j < n; j++ {
				addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", 3000+j))
				Expect(err).ShouldNot(HaveOccurred())
				store.InsertAddr(addr)
			}
			stores[i] = store

			gossiper := gossip.NewGossiper(1+rand.Intn(5), clients[i], mockVerifier{}, store)
			service := NewService(gossiper)
			servers[i] = grpc.NewServer()
			RegisterXoxoServiceServer(servers[i], &service)

			lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", 3000+i))
			Expect(err).ShouldNot(HaveOccurred())
			listeners[i] = lis
		}

		return clients, stores, servers, listeners
	}

	stopService := func(servers []*grpc.Server, listeners []net.Listener) {
		for _, server := range servers {
			server.Stop()
		}
		for _, lis := range listeners {
			lis.Close()
		}

		time.Sleep(100 * time.Millisecond)
	}

	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
	})

	Context("when adding new address", func() {
		It("should store new address", func() {
			numberOfTestObjects := 24
			numberOfMessages := 24
			clients, stores, servers, listens := initService(numberOfTestObjects)
			defer stopService(servers, listens)

			for i := range servers {
				go func(i int) {
					defer GinkgoRecover()

					err := servers[i].Serve(listens[i])
					Expect(err).ShouldNot(HaveOccurred())
				}(i)
			}

			messages := make([]foundation.Message, 0, numberOfMessages)
			for i := 0; i < numberOfMessages; i++ {
				message := randomMessage()
				messages = append(messages, message)
				sender, receiver := rand.Intn(numberOfTestObjects), rand.Intn(numberOfTestObjects)
				for sender == receiver {
					receiver = rand.Intn(numberOfTestObjects)
				}
				to, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", 3000+receiver))
				Expect(err).ShouldNot(HaveOccurred())

				Expect(clients[sender].Send(context.Background(), to, message)).ShouldNot(HaveOccurred())

			}

			time.Sleep(3 * time.Second)

			for _, message := range messages {
				received := 0
				for _, store := range stores {
					msg, err := store.Message(message.Key)
					Expect(err).ShouldNot(HaveOccurred())
					if len(msg.Key) > 0 {
						received++
					}
				}
				Expect(received).Should(BeNumerically(">=", numberOfTestObjects*4/5))
				log.Printf("Total: %v ,received : %v", numberOfTestObjects, received)
			}
		})
	})
})

// A mock verifier will always return true when verifying signature.
type mockVerifier struct {
}

func (verifier mockVerifier) Verify(data []byte, signature []byte) error {
	return nil
}

type mockStore struct {
	addrMu  *sync.Mutex
	address map[string]net.Addr

	messageMu *sync.Mutex
	messages  map[string]foundation.Message
}

func NewMockStore() mockStore {
	return mockStore{
		addrMu:    new(sync.Mutex),
		address:   map[string]net.Addr{},
		messageMu: new(sync.Mutex),
		messages:  map[string]foundation.Message{},
	}
}

func (store mockStore) InsertAddr(addr net.Addr) {
	store.addrMu.Lock()
	defer store.addrMu.Unlock()
	store.address[addr.String()] = addr
}

func (store mockStore) Addrs(α int) ([]net.Addr, error) {
	store.addrMu.Lock()
	defer store.addrMu.Unlock()
	addrs := make([]net.Addr, 0, α)
	for _, j := range store.address {
		if len(addrs) == α {
			break
		}
		addrs = append(addrs, j)
	}

	return addrs, nil
}

func (store mockStore) InsertMessage(message foundation.Message) error {
	store.messageMu.Lock()
	defer store.messageMu.Unlock()
	store.messages[string(message.Key)] = message
	return nil
}

func (store mockStore) Message(key []byte) (foundation.Message, error) {
	store.messageMu.Lock()
	defer store.messageMu.Unlock()

	return store.messages[string(key)], nil
}

// randomMessage returns a random message.
func randomMessage() foundation.Message {
	randomBytes := func() []byte {
		length := rand.Intn(65)
		data := make([]byte, length)
		_, err := rand.Read(data)
		Expect(err).ShouldNot(HaveOccurred())

		return data
	}

	return foundation.Message{
		Nonce:     rand.Uint64(),
		Key:       randomBytes(),
		Value:     randomBytes(),
		Signature: randomBytes(),
	}
}

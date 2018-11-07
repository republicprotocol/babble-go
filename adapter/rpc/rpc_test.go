package rpc_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/xoxo-go/adapter/rpc"

	"github.com/republicprotocol/co-go"
	"github.com/republicprotocol/xoxo-go/core/addr"
	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/republicprotocol/xoxo-go/testutils"
	"google.golang.org/grpc"
)

var _ = Describe("grpc", func() {

	initService := func(α, n int) ([]gossip.Client, []gossip.Store, []*grpc.Server, []net.Listener) {
		clients := make([]gossip.Client, n)
		stores := make([]gossip.Store, n)
		servers := make([]*grpc.Server, n)
		listeners := make([]net.Listener, n)

		for i := 0; i < n; i++ {
			clients[i] = NewClient(testutils.MockDialer{}, testutils.MockCaller{})

			stores[i] = testutils.NewMockStore()
			book, err := addr.NewBook(testutils.NewMockAddrs())
			Expect(err).ShouldNot(HaveOccurred())
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", 8000+j))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(book.InsertAddr(addr)).ShouldNot(HaveOccurred())
			}

			gossiper := gossip.NewGossiper(α, testutils.MockSinger{}, testutils.MockVerifier{}, nil, clients[i], book, stores[i])
			service := NewService(gossiper)
			servers[i] = grpc.NewServer()
			service.Register(servers[i])

			lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", 8000+i))
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

	for _, failureRate := range []int{0, 10, 20} { // percentage
		failureRate := failureRate
		Context("when sending message via grpc", func() {
			It("should receive the message and broadcast the message if it's new", func() {
				numberOfTestNodes := 48
				numberOfMessages := 12
				numberOfFaultyNodes := numberOfTestNodes * failureRate / 100

				shuffle := rand.Perm(numberOfTestNodes)[:numberOfFaultyNodes]
				faultyNodes := map[int]bool{}
				for _, index := range shuffle {
					faultyNodes[index] = true
				}

				clients, stores, servers, listens := initService(7, numberOfTestNodes)
				defer stopService(servers, listens)

				go co.ParForAll(servers, func(i int) {
					defer GinkgoRecover()

					if faultyNodes[i] {
						return
					}
					err := servers[i].Serve(listens[i])
					Expect(err).ShouldNot(HaveOccurred())
				})
				time.Sleep(1 * time.Second)

				// Send message
				messages := make([]gossip.Message, 0, numberOfMessages)
				for i := 0; i < numberOfMessages; i++ {
					message := randomMessage()
					messages = append(messages, message)
					sender, receiver := rand.Intn(numberOfTestNodes), rand.Intn(numberOfTestNodes)
					for {
						if !faultyNodes[sender] {
							break
						}
						sender = rand.Intn(numberOfTestNodes)
					}
					for {
						if !faultyNodes[receiver] && sender != receiver {
							break
						}
						receiver = rand.Intn(numberOfTestNodes)
					}
					to, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%v", 8000+receiver))
					Expect(err).ShouldNot(HaveOccurred())
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					defer cancel()

					Expect(clients[sender].Send(ctx, to, message)).ShouldNot(HaveOccurred())
				}
				time.Sleep(3 * time.Second)

				// Check how many nodes have got the message
				for _, message := range messages {
					received := 0
					for _, store := range stores {
						msg, err := store.Message(message.Key)
						Expect(err).ShouldNot(HaveOccurred())
						if msg.Nonce > 0 {
							received++
						}
					}

					Expect(received).Should(BeNumerically(">=", (numberOfTestNodes-numberOfFaultyNodes)*9/10))
					log.Printf("Total: %v ,received : %v", numberOfTestNodes-numberOfFaultyNodes, received)
				}
			})
		})
	}
})

// randomMessage returns a random message.
func randomMessage() gossip.Message {
	randomBytes := func() []byte {
		length := rand.Intn(65)
		data := make([]byte, length)
		_, err := rand.Read(data)
		Expect(err).ShouldNot(HaveOccurred())

		return data
	}

	return gossip.Message{
		Nonce:     rand.Uint64(),
		Key:       randomBytes(),
		Value:     randomBytes(),
		Signature: randomBytes(),
	}
}

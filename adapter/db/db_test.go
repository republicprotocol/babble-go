package db_test

import (
	"math/rand"
	"os"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/babble-go/adapter/db"

	"github.com/republicprotocol/babble-go/core/gossip"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbDir         = "./tmp"
	addrDbFile    = "./tmp/addrs"
	messageDbFile = "./tmp/messages"
)

var _ = Describe("LevelDB storage", func() {

	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		os.RemoveAll(dbDir)
	})

	Context("when adding new address", func() {
		It("should store new address", func() {
			db, err := leveldb.OpenFile(addrDbFile, nil)
			Expect(err).ShouldNot(HaveOccurred())
			defer db.Close()
			store := NewAddrStore(db)

			addrs := testAddresses()
			lookup := map[string]bool{}
			for i := range addrs {
				lookup[addrs[i].Network()+addrs[i].String()] = true
			}

			for i, addr := range addrs {
				err := store.InsertAddr(addr)
				Expect(err).ShouldNot(HaveOccurred())

				addresses, err := store.Addrs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(addresses)).Should(Equal(i + 1))

				duplicates := make(map[string]struct{})
				for _, address := range addresses {
					Expect(lookup[address.Network()+address.String()]).Should(BeTrue())
					duplicates[address.Network()+address.String()] = struct{}{}
				}

				Expect(len(duplicates)).Should(Equal(i + 1))
			}
		})
	})

	Context("when storing new messages ", func() {
		It("should store new message ", func() {
			db, err := leveldb.OpenFile(messageDbFile, nil)
			Expect(err).ShouldNot(HaveOccurred())
			defer db.Close()
			store := NewMessageStore(db)

			messages := testMessages()
			for _, message := range messages {
				err = store.InsertMessage(message)
				Expect(err).ShouldNot(HaveOccurred())

				msg, err := store.Message(message.Key)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(reflect.DeepEqual(message, msg)).Should(BeTrue())
			}
		})
	})

	Context("when reading messages ", func() {
		It("should return empty message and nil error when reading something not in the store  ", func() {
			db, err := leveldb.OpenFile(messageDbFile, nil)
			Expect(err).ShouldNot(HaveOccurred())
			defer db.Close()
			store := NewMessageStore(db)

			messages := testMessages()
			for _, message := range messages {
				msg, err := store.Message(message.Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(msg.Key)).Should(BeZero())
				Expect(len(msg.Value)).Should(BeZero())
				Expect(len(msg.Signature)).Should(BeZero())
				Expect(msg.Nonce).Should(BeZero())
			}
		})
	})
})

func testAddresses() []Addr {
	nets := []string{"tcp", "ip4", "ip6", "", "invalid", "random"}
	values := []string{"10.0.1.1", "0.0.0.0", "", "abcd", "0000"}
	addrs := make([]Addr, 0, len(nets)*len(values))
	for _, net := range nets {
		for _, value := range values {
			addrs = append(addrs, NewAddr(net, value).(Addr))
		}
	}

	return addrs
}

func testMessages() []gossip.Message {
	nonces := []uint64{0, 1, 5, 100, uint64(time.Now().Unix())}
	keys := [][]byte{{}, randomBytes(), randomBytes()}
	values := [][]byte{{}, randomBytes(), randomBytes()}
	signatures := [][]byte{{}, randomBytes(), randomBytes()}

	messages := make([]gossip.Message, 0)
	for _, nonce := range nonces {
		for _, key := range keys {
			for _, value := range values {
				for _, signature := range signatures {
					messages = append(messages, gossip.Message{
						Nonce:     nonce,
						Key:       key,
						Value:     value,
						Signature: signature,
					})
				}
			}
		}
	}

	return messages
}

func randomBytes() []byte {
	length := rand.Intn(65)
	data := make([]byte, length)
	_, err := rand.Read(data)
	Expect(err).ShouldNot(HaveOccurred())

	return data
}

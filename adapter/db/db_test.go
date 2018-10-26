package leveldb_test

import (
	"math/rand"
	"os"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/xoxo-go/driver/leveldb"

	"github.com/republicprotocol/xoxo-go/foundation"
)

const dbDir = "./tmp"
const dbFile = "./tmp/db"

var _ = Describe("LevelDB storage", func() {

	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
	})

	AfterEach(func() {
		os.RemoveAll(dbDir)
	})

	Context("when adding new address", func() {
		It("should store new address", func() {
			store, err := NewStore(dbFile)
			Expect(err).ShouldNot(HaveOccurred())
			defer store.Close()

			addrs := testAddresses()
			alpha := 1
			for _, addr := range addrs {
				err := store.InsertAddr(addr)
				Expect(err).ShouldNot(HaveOccurred())

				addresses, err := store.Addrs(alpha)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(addresses)).Should(Equal(alpha))

				addresses, err = store.Addrs(alpha + 1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(addresses)).Should(Equal(alpha))

				addresses, err = store.Addrs(alpha - 1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(addresses)).Should(Equal(alpha - 1))

				alpha++
			}
		})
	})

	Context("when storing new messages ", func() {
		It("should store new message ", func() {
			store, err := NewStore(dbFile)
			Expect(err).ShouldNot(HaveOccurred())
			defer store.Close()

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
			store, err := NewStore(dbFile)
			Expect(err).ShouldNot(HaveOccurred())
			defer store.Close()

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

func testMessages() []foundation.Message {
	nonces := []uint64{0, 1, 5, 100, uint64(time.Now().Unix())}
	keys := [][]byte{{}, randomBytes(), randomBytes()}
	values := [][]byte{{}, randomBytes(), randomBytes()}
	signatures := [][]byte{{}, randomBytes(), randomBytes()}

	messages := make([]foundation.Message, 0)
	for _, nonce := range nonces {
		for _, key := range keys {
			for _, value := range values {
				for _, signature := range signatures {
					messages = append(messages, foundation.Message{
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

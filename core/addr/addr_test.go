package addr_test

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/xoxo-go/core/addr"

	"github.com/republicprotocol/co-go"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Store", func() {

	newEmptyBook := func() Book {
		addrs := newMockAddrs()
		book, err := NewBook(addrs)
		Expect(err).ShouldNot(HaveOccurred())

		return book
	}

	resetMap := func(m map[string]int) {
		for i := range m {
			m[i] = 1
		}
	}

	testRetreivingAddrs := func(book Book, numberOfTestAddrs int, lookupMap map[string]int) {
		randAddrs, err := book.Addrs(numberOfTestAddrs)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(len(randAddrs)).Should(Equal(numberOfTestAddrs))
		for i := range randAddrs {
			Expect(lookupMap[randAddrs[i].String()]).Should(Equal(1))
			lookupMap[randAddrs[i].String()]++
		}
		resetMap(lookupMap)

		for i := 0; i < numberOfTestAddrs; i++ {
			randAddrs, err := book.Addrs(i)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(randAddrs)).Should(Equal(i))
			for i := range randAddrs {
				Expect(lookupMap[randAddrs[i].String()]).Should(Equal(1))
				lookupMap[randAddrs[i].String()]++
			}
			resetMap(lookupMap)
		}
	}

	Context("when looking up addresses", func() {

		It("should be able to return α random addresses when initialized with an empty Store", func() {
			book := newEmptyBook()
			lookupMap := map[string]int{}
			numberOfTestAddrs := 100
			for i := 0; i < numberOfTestAddrs; i++ {
				addr := randomAddr()
				lookupMap[addr.String()] = 1
				Expect(book.InsertAddr(addr)).ShouldNot(HaveOccurred())
			}

			testRetreivingAddrs(book, numberOfTestAddrs, lookupMap)
		})

		It("should be able to return α random addresses when initialized with an non-empty Store", func() {
			addrs := newMockAddrs()
			lookupMap := map[string]int{}
			numberOfTestAddrs := 100
			for i := 0; i < numberOfTestAddrs; i++ {
				addr := randomAddr()
				lookupMap[addr.String()] = 1
				Expect(addrs.InsertAddr(addr)).ShouldNot(HaveOccurred())
			}
			book, err := NewBook(addrs)
			Expect(err).ShouldNot(HaveOccurred())

			testRetreivingAddrs(book, numberOfTestAddrs, lookupMap)
		})

		It("should panic when trying to get negative number of addrs", func() {
			book := newEmptyBook()
			Expect(func() { book.Addrs(-1) }).Should(Panic())
			addresses, err := book.Addrs(0)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(addresses)).Should(Equal(0))
		})

		It("should only return what the store have", func() {
			book := newEmptyBook()
			lookupMap := map[string]int{}
			numberOfTestAddrs := 100
			for i := 0; i < numberOfTestAddrs; i++ {
				addr := randomAddr()
				lookupMap[addr.String()] = 1
				Expect(book.InsertAddr(addr)).ShouldNot(HaveOccurred())
			}

			randAddrs, err := book.Addrs(numberOfTestAddrs + 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(randAddrs)).Should(Equal(numberOfTestAddrs))
			for i := range randAddrs {
				Expect(lookupMap[randAddrs[i].String()]).Should(Equal(1))
				lookupMap[randAddrs[i].String()]++
			}
		})

	})

	Context("concurrent use cases", func() {

		It("should be concurrent-safe to inserting and retrieving addrs", func() {
			book := newEmptyBook()
			numberOfTestAddrs := 100

			co.ParBegin(func() { // inserting addrs
				co.ParForAll(numberOfTestAddrs, func(i int) {
					defer GinkgoRecover()

					addr := randomAddr()
					Expect(book.InsertAddr(addr)).ShouldNot(HaveOccurred())
				})

			}, func() { // reading addrs
				co.ParForAll(numberOfTestAddrs, func(i int) {
					defer GinkgoRecover()

					_, err := book.Addrs(i)
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

	})

})

// mockStore is a mock implementation of the `addr.Store`
type mockStore struct {
	addresses map[string]net.Addr
}

func newMockAddrs() Store {
	return mockStore{
		addresses: map[string]net.Addr{},
	}
}

func (store mockStore) InsertAddr(addr net.Addr) error {
	store.addresses[addr.String()] = addr
	return nil
}

func (store mockStore) Addrs() ([]net.Addr, error) {
	addresses := make([]net.Addr, 0, len(store.addresses))
	for _, j := range store.addresses {
		addresses = append(addresses, j)
	}

	return addresses, nil
}

func randomAddr() net.Addr {
	addr, err := net.ResolveTCPAddr("", fmt.Sprintf("%v.%v.%v.%v: %v", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(10000)))
	if err != nil {
		log.Fatal(err)
	}

	return addr
}

package addr

import (
	"math/rand"
	"net"
	"sync"
)

// Addrs is used to store and lookup all known `net.Addrs`. It is not assumed
// that this interface is safe for concurrent use.
type Addrs interface {

	// Insert a new Addr to the store.
	InsertAddr(net.Addr) error

	// Addrs returns all the stored Addrs.
	Addrs() ([]net.Addr, error)
}

// Book is used to provide a fast lookup of random `net.Addrs` that can be used
// by nodes to begin a gossip.
type Book interface {

	// InsertAddr into the Book.
	InsertAddr(net.Addr) error

	// Addrs returns α random `net.Addrs` from the set of all known `net.Addrs`
	// to the Book.
	Addrs(α int) ([]net.Addr, error)
}

type book struct {
	addrsMu    *sync.RWMutex
	addrsCache []net.Addr
	addrs      Addrs
}

func NewBook(addrs Addrs) (Book, error) {
	addresses, err := addrs.Addrs()
	if err != nil {
		return nil, err
	}

	return &book{
		addrsMu:    new(sync.RWMutex),
		addrsCache: addresses,
		addrs:      addrs,
	}, nil
}

func (book *book) InsertAddr(addr net.Addr) error {
	book.addrsMu.Lock()
	defer book.addrsMu.Unlock()

	book.addrsCache = append(book.addrsCache, addr)
	return book.addrs.InsertAddr(addr)
}

func (book *book) Addrs(α int) ([]net.Addr, error) {
	book.addrsMu.RLock()
	defer book.addrsMu.RUnlock()

	addrs := make([]net.Addr, 0, α)
	if α > len(book.addrsCache) {
		α = len(book.addrsCache)
	}
	for i := range rand.Perm(len(book.addrsCache))[:α] {
		addrs = append(addrs, book.addrsCache[i])
	}

	return addrs, nil
}

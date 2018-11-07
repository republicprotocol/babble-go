package addr

import (
	"math/rand"
	"net"
	"sync"
)

// Store is used to store and lookup all known `net.Addr`. It is not assumed
// that this interface is safe for concurrent use.
type Store interface {

	// Insert a new Addr to the store.
	InsertAddr(net.Addr) error

	// Store returns all the stored Store.
	Addrs() ([]net.Addr, error)
}

// Book is used to provide a fast lookup of random `net.Addr` that can be used
// by nodes to begin a gossip.
type Book interface {

	// InsertAddr into the Book.
	InsertAddr(net.Addr) error

	// Store returns α random `net.Addr` from the set of all known `net.Addr`
	// to the Book.
	Addrs(α int) ([]net.Addr, error)
}

type book struct {
	addrsMu    *sync.RWMutex
	addrsCache []net.Addr
	addrs      Store
}

// NewBook returns a new Book with given addr Store.
func NewBook(store Store) (Book, error) {
	addresses, err := store.Addrs()
	if err != nil {
		return nil, err
	}

	return &book{
		addrsMu:    new(sync.RWMutex),
		addrsCache: addresses,
		addrs:      store,
	}, nil
}

// InsertAddr implements Store interface.
func (book *book) InsertAddr(addr net.Addr) error {
	book.addrsMu.Lock()
	defer book.addrsMu.Unlock()

	book.addrsCache = append(book.addrsCache, addr)
	return book.addrs.InsertAddr(addr)
}

// Addrs implements Store interface.
func (book *book) Addrs(α int) ([]net.Addr, error) {
	book.addrsMu.RLock()
	defer book.addrsMu.RUnlock()

	if α > len(book.addrsCache) {
		α = len(book.addrsCache)
	}
	addrs := make([]net.Addr, 0, α)
	for _, index := range rand.Perm(len(book.addrsCache))[:α] {
		addrs = append(addrs, book.addrsCache[index])
	}

	return addrs, nil
}

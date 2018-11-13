package addr

import (
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
	mu    *sync.RWMutex
	store Store
	cache map[string]net.Addr
}

// NewBook returns a new Book with given addr Store.
func NewBook(store Store) (Book, error) {
	addrs, err := store.Addrs()
	if err != nil {
		return nil, err
	}

	cache := make(map[string]net.Addr, len(addrs))
	for _, addr := range addrs {
		cache[addr.String()] = addr
	}

	return &book{
		mu:    new(sync.RWMutex),
		store: store,
		cache: cache,
	}, nil
}

// InsertAddr implements `Store` interface.
func (book *book) InsertAddr(addr net.Addr) error {
	book.mu.Lock()
	defer book.mu.Unlock()

	book.cache[addr.String()] = addr
	return book.store.InsertAddr(addr)
}

// Addrs implements `Store` interface.
func (book *book) Addrs(α int) ([]net.Addr, error) {
	book.mu.RLock()
	defer book.mu.RUnlock()

	addrs := make([]net.Addr, 0, α)
	for _, addr := range book.cache {
		if len(addrs) >= α {
			return addrs, nil
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

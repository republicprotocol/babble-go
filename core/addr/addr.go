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
	addrsMu    *sync.RWMutex
	addrsStore Store
	addrsCache map[string]net.Addr
}

// NewBook returns a new Book with given addr Store.
func NewBook(store Store) (Book, error) {
	addrs, err := store.Addrs()
	if err != nil {
		return nil, err
	}

	addrsCache := make(map[string]net.Addr, len(addrs))
	for _, addr := range addrs {
		addrsCache[addr.String()] = addr
	}

	return &book{
		addrsMu:    new(sync.RWMutex),
		addrsStore: store,
		addrsCache: addrsCache,
	}, nil
}

// InsertAddr implements Store interface.
func (book *book) InsertAddr(addr net.Addr) error {
	book.addrsMu.Lock()
	defer book.addrsMu.Unlock()

	book.addrsCache[addr.String()] = addr
	return book.addrsStore.InsertAddr(addr)
}

// Addrs implements Store interface.
func (book *book) Addrs(α int) ([]net.Addr, error) {
	book.addrsMu.RLock()
	defer book.addrsMu.RUnlock()

	addrs := make([]net.Addr, 0, α)
	for _, addr := range book.addrsCache {
		if len(addrs) >= α {
			return addrs, nil
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

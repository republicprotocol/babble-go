package addr

import (
	"net"
	"sync"
)

// Addrs is used to store and lookup all known `net.Addrs`. It is not assumed
// that this interface is safe for concurrent use.
type Addrs interface {
	InsertAddr(net.Addr) error
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
	addrsMu *sync.RWMutex
	addrs   Addrs
}

func NewBook(addrs Addrs) Book {
	return &book{
		addrsMu: new(sync.RWMutex),
		addrs:   addrs,
	}
}

func (book *book) InsertAddr(addr net.Addr) error {
	panic("unimplemented")
}

func (book *book) Addrs(α int) ([]net.Addr, error) {
	panic("unimplemented")
}

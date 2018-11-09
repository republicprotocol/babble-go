package gossip

import (
	"github.com/republicprotocol/babble-go/core/addr"
	"net"
)

// MessageStore is used to read and write MessageStore that are being
// disseminated throughout the network.
type MessageStore interface {

	// InsertMessage into the MessageStore. If there is an existing Message
	// with the same key, but a lower nonce, then the existing Message will
	// be overwritten.
	InsertMessage(message Message) error

	// Message returns a previously inserted Message associated with the key.
	// It returns an empty message with zero nonce if there is no message with
	// the associated key in the store.
	Message(key []byte) (Message, error)
}

// Store stores the both the Messages and the net.Addrs.
type Store interface {
	MessageStore

	addr.Book
}

type store struct {
	messages MessageStore
	book     addr.Book
}

// NewStore returns a new Store with given MessageStore and addr.Book.
func NewStore(messages MessageStore, book addr.Book) Store {
	return store{
		messages: messages,
		book:     book,
	}
}

// InsertMessage implements the Store interface.
func (store store) InsertMessage(message Message) error {
	return store.messages.InsertMessage(message)
}

// Message implements the Store interface.
func (store store) Message(key []byte) (Message, error) {
	return store.messages.Message(key)
}

// InsertAddr implements the Store interface.
func (store store) InsertAddr(addr net.Addr) error {
	return store.book.InsertAddr(addr)
}

// Addrs implements the Store interface.
func (store store) Addrs(α int) ([]net.Addr, error) {
	return store.book.Addrs(α)
}

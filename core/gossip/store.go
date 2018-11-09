package gossip

import (
	"github.com/republicprotocol/babble-go/core/addr"
	"net"
)

// Store is used to read and write Messages that are being disseminated
// throughout the network.
type Messages interface {

	// InsertMessage into the store. If there is an existing Message with the
	// same key, but a lower nonce, then the existing Message will be
	// overwritten.
	InsertMessage(message Message) error

	// Message returns a previously inserted Message associated with the key.
	// It returns an empty message with zero nonce if there is no message with
	// the associated key in the store.
	Message(key []byte) (Message, error)
}

type Store interface {
	Messages

	addr.Book
}

type store struct {
	messages Messages
	book     addr.Book
}

func NewStore(messages Messages, book addr.Book) Store {
	return store{
		messages: messages,
		book:     book,
	}
}

func (store store) InsertMessage(message Message) error {
	return store.messages.InsertMessage(message)
}

func (store store) Message(key []byte) (Message, error) {
	return store.messages.Message(key)
}

func (store store) InsertAddr(addr net.Addr) error {
	return store.book.InsertAddr(addr)
}

func (store store) Addrs(α int) ([]net.Addr, error) {
	return store.book.Addrs(α)
}

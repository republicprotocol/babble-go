package gossip

import (
	"errors"
	"github.com/republicprotocol/xoxo-go/core/addr"
)

// ErrMessageNotFound is returned when there is no Message associated with a
// key.
var ErrMessageNotFound = errors.New("message not found")

// Messages is used to read and write Messages that are being disseminated
// throughout the network.
type Messages interface {

	// InsertMessage into the store. If there is an existing Message with the
	// same key, but a lower nonce, then the existing Message will be
	// overwritten.
	InsertMessage(message Message) error

	// Message returns a previously inserted Message associated with the key.
	Message(key []byte) (Message, error)
}

// A Store of Addrs and Messages.
type Store interface {
	addr.Book
	Messages
}

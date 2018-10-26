package gossip

import (
	"errors"
	"net"

	"github.com/republicprotocol/xoxo-go/foundation"
)

// ErrMessageNotFound is returned when there is no Message associated with a
// key.
var ErrMessageNotFound = errors.New("message not found")

// Addrs is used to lookup `net.Addrs` of nodes in the network.
type Addrs interface {

	// Addrs returns `net.Addrs` that can be used to start the gossip.
	Addrs(α int) ([]net.Addr, error)
}

// Messages is used to read and write Messages that are being disseminated
// throughout the network.
type Messages interface {

	// InsertMessage into the store. If there is an existing Message with the
	// same key, but a lower nonce, then the existing Message will be
	// overwritten.
	InsertMessage(message foundation.Message) error

	// Message returns a previously inserted Message associated with the key.
	Message(key []byte) (foundation.Message, error)
}

// A Store of Addrs and Messages.
type Store interface {
	Addrs
	Messages
}

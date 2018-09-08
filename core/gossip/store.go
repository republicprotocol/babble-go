package gossip

import (
	"net"

	"github.com/republicprotocol/xoxo-go/foundation"
)

type AddrStore interface {
	Addrs(α int) ([]net.Addr, error)
}

type MessageStore interface {
	InsertMessage(message foundation.Message) error
	Message(key []byte) (foundation.Message, error)
}

type Store interface {
	AddrStore
	MessageStore
}

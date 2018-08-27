package gossip

import (
	"net"

	"github.com/republicprotocol/xoxo-go/foundation"
)

type AddrStore interface {
	InsertAddress(address net.Addr) error
	Addrs() ([]net.Addr, error)
}

type MessageStore interface {
	InsertMessage(message foundation.Message) error
	Message(key []byte) (foundation.Message, error)
}

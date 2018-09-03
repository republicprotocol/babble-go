package gossip

import (
	"net"

	"github.com/republicprotocol/xoxo-go/foundation"
)

type AddrStore interface {
	InsertAddr(addr net.Addr) error
	Addrs(α int) ([]net.Addr, error)
}

type MessageStore interface {
	InsertMessage(message foundation.Message) error
	Message(key []byte) (foundation.Message, error)
}

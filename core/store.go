package core

import (
	"net"

	"github.com/republicprotocol/gossip-network/foundation"
)

// AddressStorer for the ip address
type AddressStorer interface {
	InsertAddress(address net.Addr) error

	Addresses() ([]net.Addr, error)
}

type MessageStore interface {
	InsertMessage(message foundation.Message) error

	Message(hash [32]byte) (foundation.Message, error)
}





package babble

import (
	"github.com/republicprotocol/babble-go/adapter/db"
	"github.com/republicprotocol/babble-go/adapter/rpc"
	"github.com/republicprotocol/babble-go/core/addr"
	"github.com/republicprotocol/babble-go/core/gossip"
)

type (
	Gossiper = gossip.Gossiper
	Observer = gossip.Observer
	Message  = gossip.Message
)

var (
	NewGossiper     = gossip.NewGossiper
	NewRPCClient    = rpc.NewClient
	NewRPCService   = rpc.NewService
	NewBook         = addr.NewBook
	NewMessageStore = db.NewMessageStore
	NewAddrs        = db.NewAddrStore
)

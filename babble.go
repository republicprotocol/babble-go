package babble

import (
	"github.com/republicprotocol/babble-go/adapter/db"
	"github.com/republicprotocol/babble-go/adapter/rpc"
	"github.com/republicprotocol/babble-go/core/addr"
	"github.com/republicprotocol/babble-go/core/gossip"
)

type (
	AddrStore    = addr.Store
	Book         = addr.Book
	MessageStore = gossip.MessageStore
	Gossiper     = gossip.Gossiper
	Message      = gossip.Message
	Client       = gossip.Client
	Notifier     = gossip.Notifier
	Signer       = gossip.Signer
	Verifier     = gossip.Verifier
)

var (
	NewAddrs        = db.NewAddrStore
	NewMessageStore = db.NewMessageStore
	NewBook         = addr.NewBook
	NewGossiper     = gossip.NewGossiper
	NewMessage      = gossip.NewMessage
	NewRPCClient    = rpc.NewClient
	NewRPCService   = rpc.NewService
)

package babble

import (
	"github.com/republicprotocol/babble-go/adapter/db"
	"github.com/republicprotocol/babble-go/adapter/rpc"
	"github.com/republicprotocol/babble-go/core/addr"
	"github.com/republicprotocol/babble-go/core/gossip"
)

type (
	Addrs    = addr.Addrs
	AddrBook = addr.Book
	Messages = gossip.Messages
	Gossiper = gossip.Gossiper
	Message  = gossip.Message
	Client   = gossip.Client
	Observer = gossip.Observer
	Signer   = gossip.Signer
	Verifier = gossip.Verifier
)

var (
	NewDb         = db.New
	NewBook       = addr.NewBook
	NewGossiper   = gossip.NewGossiper
	NewMessage    = gossip.NewMessage
	NewRPCClient  = rpc.NewClient
	NewRPCService = rpc.NewService
)

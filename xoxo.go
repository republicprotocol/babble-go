package xoxo

import (
	"github.com/republicprotocol/xoxo-go/adapter/db"
	"github.com/republicprotocol/xoxo-go/adapter/rpc"
	"github.com/republicprotocol/xoxo-go/core/addr"
	"github.com/republicprotocol/xoxo-go/core/gossip"
)

type (
	Gossiper = gossip.Gossiper
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

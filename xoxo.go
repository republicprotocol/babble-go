package xoxo

import (
	"github.com/republicprotocol/xoxo-go/adapter/grpc"
	"github.com/republicprotocol/xoxo-go/adapter/leveldb"
	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/republicprotocol/xoxo-go/foundation"
)

type (
	Gossiper = gossip.Gossiper
	Message  = foundation.Message
)

var (
	NewGossiper   = gossip.NewGossiper
	NewRPCClient  = grpc.NewClient
	NewRPCService = grpc.NewService
	NewStore      = leveldb.NewStore
)

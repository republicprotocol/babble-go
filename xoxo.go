package xoxo

import (
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/xoxo-go/adapter/grpc"
	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/republicprotocol/xoxo-go/foundation"
)

type Message foundation.Message

type GossipClient gossip.Client

type GossipServer gossip.Server

type Gossiper gossip.Gossiper

var (
	NewRPCClient  = grpc.NewClient
	NewRPCService = grpc.NewService

	NewDB = leveldb.NewStore

	NewGossiper = gossip.NewGossiper
)

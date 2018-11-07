package db

import (
	"encoding/json"
	"net"

	"github.com/republicprotocol/xoxo-go/core/addr"
	"github.com/republicprotocol/xoxo-go/core/gossip"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// Addr stores an address `Value` couple with a `Net`.
type Addr struct {
	Net   string `json:"net"`
	Value string `json:"value"`
}

// NewAddr returns a new `net.Addr`.
func NewAddr(net, value string) net.Addr {
	return Addr{
		net, value,
	}
}

// Network implements the `net.Addr` interface.
func (addr Addr) Network() string {
	return addr.Net
}

// String implements the `net.Addr` interface.
func (addr Addr) String() string {
	return addr.Value
}

// A AddrStore uses LevelDB to store Addrs to persistent storage. It is a basic
// implementation of the `addr.Store` with no explicit in-memory cache, and no
// optimisations for returning random information.
type AddrStore struct {
	db *leveldb.DB
}

// NewAddrStore returns a new AddrStore.
func NewAddrStore(db *leveldb.DB) addr.Store {
	return &AddrStore{db}
}

// InsertAddr implements the `gossip.AddrStore` interface.
func (store *AddrStore) InsertAddr(addr net.Addr) error {
	data, err := json.Marshal(NewAddr(
		addr.Network(),
		addr.String(),
	))
	if err != nil {
		return err
	}
	return store.db.Put(data, []byte{}, nil)
}

// Addrs implements the `gossip.AddrStore` interface.
func (store *AddrStore) Addrs() ([]net.Addr, error) {
	iter := store.db.NewIterator(&util.Range{Start: nil, Limit: nil}, nil)
	defer iter.Release()

	addrs := make([]net.Addr, 0)
	for iter.Next() {
		addr := Addr{}
		if err := json.Unmarshal(iter.Key(), &addr); err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}

	return addrs, iter.Error()
}

// A MessageStore uses LevelDB to store Messages to persistent storage.
type MessageStore struct {
	db *leveldb.DB
}

// NewMessageStore returns a new MessageStore.
func NewMessageStore(db *leveldb.DB) gossip.Store {
	return &AddrStore{db}
}

// InsertMessage implements the `gossip.MessageStore` interface.
func (store *AddrStore) InsertMessage(message gossip.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return store.db.Put(message.Key, data, nil)
}

// Message implements the `gossip.MessageStore` interface.
func (store *AddrStore) Message(key []byte) (gossip.Message, error) {
	message := gossip.Message{}
	data, err := store.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = nil
		}
		return message, err
	}
	err = json.Unmarshal(data, &message)

	return message, err
}

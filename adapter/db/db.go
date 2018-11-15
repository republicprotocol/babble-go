package db

import (
	"encoding/json"
	"net"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/babble-go/core/addr"
	"github.com/republicprotocol/babble-go/core/gossip"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Addr struct {
	Net   string `json:"net"`
	Value string `json:"value"`
}

func NewAddr(net, value string) net.Addr {
	return Addr{
		net, value,
	}
}

func (addr Addr) Network() string {
	return addr.Net
}

func (addr Addr) String() string {
	return addr.Value
}

// A Db uses LevelDB to store Addrs to persistent storage. It is a basic
// implementation of the `addr.Store` with no explicit in-memory cache.
type Db interface {
	addr.Addrs
	gossip.Messages
}

type db struct {
	ldb *leveldb.DB
}

// New Db that uses LevelDB for simple persistent storage.
func New(ldb *leveldb.DB) Db {
	return &db{ldb}
}

// InsertAddr implements the `addr.Addrs` interface.
func (db *db) InsertAddr(addr net.Addr) error {
	data, err := json.Marshal(NewAddr(addr.Network(), addr.String()))
	if err != nil {
		return err
	}
	return db.ldb.Put(keyForAddrs(data), data, nil)
}

// Addrs implements the `addr.Addrs` interface.
func (db *db) Addrs() ([]net.Addr, error) {
	iter := db.ldb.NewIterator(&util.Range{Start: append(keyPrefixForAddrs(), keyIterBegin()...), Limit: append(keyPrefixForAddrs(), keyIterEnd()...)}, nil)
	defer iter.Release()

	addrs := make([]net.Addr, 0)
	for iter.Next() {
		addr := Addr{}
		if err := json.Unmarshal(iter.Value(), &addr); err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}

	return addrs, iter.Error()
}

// InsertMessage implements the `gossip.Messages` interface.
func (db *db) InsertMessage(message gossip.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return db.ldb.Put(keyForMessages(message.Key), data, nil)
}

// Message implements the `gossip.Messages` interface.
func (db *db) Message(key []byte) (gossip.Message, error) {
	message := gossip.Message{}
	data, err := db.ldb.Get(keyForMessages(key), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			err = nil
		}
		return message, err
	}
	err = json.Unmarshal(data, &message)
	return message, err
}

func keyPrefixForMessages() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func keyPrefixForAddrs() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
}

func keyForMessages(key []byte) []byte {
	return append(keyPrefixForMessages(), crypto.Keccak256(key)...)
}

func keyForAddrs(key []byte) []byte {
	return append(keyPrefixForAddrs(), crypto.Keccak256(key)...)
}

func keyIterBegin() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
}

func keyIterEnd() []byte {
	return []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
}

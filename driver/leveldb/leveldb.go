package leveldb

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"

	"github.com/republicprotocol/xoxo-go/foundation"
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

// A Store uses LevelDB to store Addrs and Messages to persistent storage. It is
// a basic implementation of the `gossip.AddrStore` and `gossip.MessageStore`
// with no explicit in-memory cache, and no optimisations for returning random
// information.
type Store struct {
	db *leveldb.DB
}

// NewStore returns a new Store. A call to `Store.Close` is required to release
// all resources.
func NewStore(dir string) (*Store, error) {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Close the Store and release all resources.
func (store *Store) Close() error {
	return store.db.Close()
}

// InsertAddr implements the `gossip.AddrStore` interface.
func (store *Store) InsertAddr(addr net.Addr) error {
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
func (store *Store) Addrs(α int) ([]net.Addr, error) {
	iter := store.db.NewIterator(&util.Range{Start: nil, Limit: nil}, nil)
	defer iter.Release()

	addrs := make([]net.Addr, 0, α)
	for iter.Next() {
		addr := Addr{}
		if err := json.Unmarshal(iter.Value(), &addr); err != nil {
			return nil, err
		}
		addrs = append(addrs, addr)
	}
	shuffle(addrs)

	if len(addrs) <= α {
		return addrs, iter.Error()
	}
	return addrs[:α], iter.Error()
}

// InsertMessage implements the `gossip.MessageStore` interface.
func (store *Store) InsertMessage(message foundation.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return store.db.Put(message.Key, data, nil)
}

// Message implements the `gossip.MessageStore` interface.
func (store *Store) Message(key []byte) (foundation.Message, error) {
	message := foundation.Message{}
	data, err := store.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			message.Nonce = 0
			message.Key = key
			err = nil
		}
		return message, err
	}
	if err := json.Unmarshal(data, &message); err != nil {
		return message, err
	}
	return message, nil
}

func shuffle(values []net.Addr) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(values) > 0 {
		n := len(values)
		randIndex := r.Intn(n)
		values[n-1], values[randIndex] = values[randIndex], values[n-1]
		values = values[:n-1]
	}
}

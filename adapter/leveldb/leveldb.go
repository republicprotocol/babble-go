package leveldb

import (
	"encoding/json"
	"net"

	"github.com/republicprotocol/xoxo-go/foundation"
	"github.com/syndtr/goleveldb/leveldb"
)

type Addr struct {
	net     string `json:"net"`
	address string `json:"address"`
}

func (addr Addr) Network() string {
	return addr.net
}

func (addr Addr) String() string {
	return addr.address
}

type Store struct {
	db    *leveldb.DB
	alpha int
}

func NewStore(dir string, alpha int) (*Store, error) {
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, err
	}
	return &Store{db: db, alpha: alpha}, nil
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) InsertAddr(addr net.Addr) error {
	data, err := json.Marshal(Addr{
		net:     addr.Network(),
		address: addr.String(),
	})
	if err != nil {
		return err
	}
	return store.db.Put(data, []byte{}, nil)
}

func (store *Store) Addrs() ([]net.Addr, error) {
	panic("unimplemented")
}

func (store *Store) InsertMessage(message foundation.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return store.db.Put(message.Key, data, nil)
}

func (store *Store) Message(key []byte) (foundation.Message, error) {
	panic("unimplemented")
}

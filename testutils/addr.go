package testutils

import (
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/republicprotocol/babble-go/core/addr"
)

func RandomAddr() net.Addr {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v.%v.%v.%v:%v", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(8000)))
	if err != nil {
		log.Fatal(err)
	}
	return addr
}

type MockAddrs struct {
	addrs map[string]net.Addr
}

func NewMockAddrs() addr.Addrs {
	return MockAddrs{
		addrs: map[string]net.Addr{},
	}
}

func (addrs MockAddrs) InsertAddr(addr net.Addr) error {
	addrs.addrs[addr.String()] = addr
	return nil
}

func (addrs MockAddrs) Addrs() ([]net.Addr, error) {
	ret := make([]net.Addr, 0, len(addrs.addrs))
	for _, addr := range addrs.addrs {
		ret = append(ret, addr)
	}

	return ret, nil
}

package gossip_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGossip(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gossip Suite")
}

package grpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRPC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RPC Suite")
}

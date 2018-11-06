package addr_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAddr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Addr Suite")
}

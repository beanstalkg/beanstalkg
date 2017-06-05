package testintegration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTestintegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testintegration Suite")
}

package orca_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOrca(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Orca Suite")
}

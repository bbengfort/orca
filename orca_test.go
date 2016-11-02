package orca_test

import (
	. "github.com/bbengfort/orca"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Orca", func() {

	It("should be at version 0.1", func() {
		Ω(Version).Should(Equal("ping"))
	})

})

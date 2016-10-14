package orca_test

import (
	. "github.com/bbengfort/orca"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {

	Describe("Devices", func() {

		It("should implement the Model interface", func() {
			val := &Device{}
			var iface interface{} = val

			_, ok := iface.(Model)
			Ω(ok).Should(BeTrue())

		})

	})

	Describe("Locations", func() {
		It("should implement the Model interface", func() {
			val := &Location{}
			var iface interface{} = val

			_, ok := iface.(Model)
			Ω(ok).Should(BeTrue())

		})
	})

})

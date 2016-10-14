package orca_test

import (
	"time"

	. "github.com/bbengfort/orca"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Echo", func() {

	const (
		secs    = 1464091727
		nsecs   = 1464091727331098063
		nsrest  = 331098063
		moment  = "2016-05-24T08:08:47-04:00"
		nmoment = "2016-05-24T08:08:47.331098063-04:00"
	)

	It("should be able to parse a time message with seconds", func() {
		msg := &Time{secs, 0}
		ts, err := time.Parse(time.RFC3339, moment)
		Ω(err).ShouldNot(HaveOccurred())

		// Use Time.Equal to ensure equality even in different locations
		Ω(ts.Equal(msg.Parse())).Should(BeTrue())
	})

	It("should be able to parse a time message with nanoseconds", func() {
		msg := &Time{0, nsecs}
		ts, err := time.Parse(time.RFC3339Nano, nmoment)
		Ω(err).ShouldNot(HaveOccurred())

		// Use Time.Equal to ensure equality even in different locations
		Ω(ts.Equal(msg.Parse())).Should(BeTrue())
	})

	It("should be able to parse a time message with both seconds and nanoseconds", func() {
		msg := &Time{secs, nsrest}
		ts, err := time.Parse(time.RFC3339Nano, nmoment)
		Ω(err).ShouldNot(HaveOccurred())

		// Use Time.Equal to ensure equality even in different locations
		Ω(ts.Equal(msg.Parse())).Should(BeTrue())
	})

})

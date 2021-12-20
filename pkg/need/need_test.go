package need_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = Describe("Need", func() {
	It("fails", func() {
		gomega.Expect(true).To(gomega.BeTrue())
	})
})

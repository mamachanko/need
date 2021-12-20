package need_test

import (
	"fmt"
	. "github.com/mamachanko/need/pkg/need"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("Need", func() {
	var need *Need

	Describe("Assess", func() {
		Context("When it succeeds", func() {
			BeforeEach(func() {
				need = &Need{AssessCmd: "true"}
			})

			It("returns no error", func() {
				Expect(need.Assess()).To(Not(HaveOccurred()))
			})
		})

		Context("When it fails", func() {
			BeforeEach(func() {
				need = &Need{AssessCmd: "false"}
			})

			It("returns error", func() {
				Expect(need.Assess()).To(HaveOccurred())
			})
		})
	})

	Describe("Fulfill", func() {
		Context("When it succeeds", func() {
			BeforeEach(func() {
				need = &Need{FulfillCmd: "true"}
			})

			It("returns no error", func() {
				Expect(need.Fulfill()).To(Not(HaveOccurred()))
			})
		})

		Context("When it fails", func() {
			BeforeEach(func() {
				need = &Need{FulfillCmd: "false"}
			})

			It("returns error", func() {
				Expect(need.Fulfill()).To(HaveOccurred())
			})
		})
	})

	Describe("Address", func() {
		Context("When already fulfilled", func() {
			BeforeEach(func() {
				need = &Need{AssessCmd: "true"}
			})

			It("returns no error", func() {
				Expect(need.Address()).To(Not(HaveOccurred()))
			})
		})

		Context("When not yet fulfilled", func() {
			var tempDir string

			BeforeEach(func() {
				var err error
				tempDir, err = ioutil.TempDir("", "need_test")
				Expect(err).To(Not(HaveOccurred()))
			})

			AfterEach(func() {
				Expect(os.RemoveAll(tempDir)).To(Not(HaveOccurred()))
			})

			Context("When fulfill-able", func() {
				BeforeEach(func() {
					aFile := filepath.Join(tempDir, "afile")

					need = &Need{
						AssessCmd:  fmt.Sprintf("stat %s", aFile),
						FulfillCmd: fmt.Sprintf("touch %s", aFile),
					}
				})

				It("returns no error", func() {
					Expect(need.Address()).To(Not(HaveOccurred()))
				})
			})

			Context("When not fulfill-able", func() {
				BeforeEach(func() {
					need = &Need{
						AssessCmd:  "false",
						FulfillCmd: "true",
					}
				})

				It("returns an error", func() {
					Expect(need.Address()).To(HaveOccurred())
				})
			})

			Context("When fulfillment fails", func() {
				BeforeEach(func() {
					need = &Need{
						AssessCmd:  "false",
						FulfillCmd: "false",
					}
				})

				It("returns an error", func() {
					Expect(need.Address()).To(HaveOccurred())
				})
			})

		})
	})

})

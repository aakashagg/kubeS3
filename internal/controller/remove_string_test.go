package controller

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("removeString", func() {
	It("removes an existing string from slice", func() {
		slice := []string{"a", "b", "c", "b"}
		result := removeString(slice, "b")
		Expect(result).To(Equal([]string{"a", "c"}))
	})

	It("leaves slice unchanged when string is absent", func() {
		slice := []string{"a", "c"}
		result := removeString(slice, "b")
		Expect(result).To(Equal([]string{"a", "c"}))
	})
})

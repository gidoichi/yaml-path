package main_test

import (
	main "github.com/gidoichi/yaml-path"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	yaml := []byte(`top:
  first:
    - name: myname
      attr1: val1
      attr2: val2
      #       ^
    - value2
    - value3
  second:
    child1: value1
    child1: value2
    child3: value3
`)

	Context("Indicating at some token with Bosh format", func() {
		It("should return the string path to the token in Bosh format.", func() {
			path, err := main.PathAtPoint(5, 14, yaml)
			Expect(err).To(BeNil())
			strpath, err := path.ToString(main.Bosh)
			Expect(err).To(BeNil())
			Expect(strpath).To(Equal("/top/first/name=myname/attr2"))
		})
	})

	Context("Indicating at some token with JsonPath format", func() {
		It("should return the string path to the token in JsonPath format.", func() {
			path, err := main.PathAtPoint(5, 14, yaml)
			Expect(err).To(BeNil())
			strpath, err := path.ToString(main.JsonPath)
			Expect(err).To(BeNil())
			Expect(strpath).To(Equal("$.top.first[0].attr2"))
		})
	})

	Context("Indicating at no token", func() {
		It("should return token not found error.", func() {
			_, err := main.PathAtPoint(6, 14, yaml)
			Expect(err).NotTo(BeNil())
		})
	})
})

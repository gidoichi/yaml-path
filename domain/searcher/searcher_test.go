package searcher_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gidoichi/yaml-path/domain/searcher"
	yaml "gopkg.in/yaml.v3"
)

var _ = Describe("Searcher", func() {
	data := []byte(`top:
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

	Context("When indicating at some token", func() {
		It("should return the path to the token.", func() {
			path, err := searcher.PathAtPoint(5, 14, data)
			Expect(err).To(BeNil())
			Expect(len(path)).To(Equal(9))
			Expect(path[0].Kind).To(Equal(yaml.DocumentNode))
			Expect(path[1].Kind).To(Equal(yaml.MappingNode))
			Expect(path[2].Kind).To(Equal(yaml.ScalarNode))
			Expect(path[2].Value).To(Equal("top"))
			Expect(path[3].Kind).To(Equal(yaml.MappingNode))
			Expect(path[4].Kind).To(Equal(yaml.ScalarNode))
			Expect(path[4].Value).To(Equal("first"))
			Expect(path[5].Kind).To(Equal(yaml.SequenceNode))
			Expect(path[6].Kind).To(Equal(yaml.MappingNode))
			Expect(path[7].Kind).To(Equal(yaml.ScalarNode))
			Expect(path[7].Value).To(Equal("attr2"))
			Expect(path[8].Kind).To(Equal(yaml.ScalarNode))
			Expect(path[8].Value).To(Equal("val2"))
		})
	})

	Context("When indicating at no token", func() {
		It("should return token not found error.", func() {
			_, err := searcher.PathAtPoint(6, 14, data)
			Expect(err).To(BeAssignableToTypeOf(searcher.TokenNotFoundError{}))
		})
	})

	Context("When invalid yaml file is inputted", func() {
		It("should return some error.", func() {
			_, err := searcher.PathAtPoint(1, 1, []byte(`top: -`))
			Expect(err).NotTo(BeNil())
		})
	})
})

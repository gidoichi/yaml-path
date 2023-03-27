package yaml_test

import (
	"bytes"
	"io"

	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	dyaml "github.com/gidoichi/yaml-path/domain/yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	yamlv3 "gopkg.in/yaml.v3"
)

var _ = Describe("YAML", func() {
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
	var reader io.Reader
	BeforeEach(func() {
		reader = bytes.NewReader(data)
	})

	Describe("NewYAML()", func() {
		Context("with single document yaml", func() {
			It("should return parsed yaml", func() {
				_, err := dyaml.NewYAML(reader)

				Expect(err).To(BeNil())
			})
		})

		Context("with multiple document yaml", func() {
			It("should return parsed yaml", func() {
				Expect(nil).NotTo(BeNil())
			})
		})

		Context("with invalid yaml", func() {
			invalid := []byte("top: -")
			BeforeEach(func() {
				reader = bytes.NewReader(invalid)
			})

			It("should return an error", func() {
				_, err := dyaml.NewYAML(reader)

				Expect(err).NotTo(BeNil())
			})
		})
	})

	Describe("PathAtPoint()", func() {
		Context("with single document yaml", func() {
			var yaml *dyaml.YAML
			BeforeEach(func() {
				var err error
				yaml, err = dyaml.NewYAML(reader)
				Expect(err).To(BeNil())
			})

			Context("indicating at mapping value node", func() {
				matcher := dmatcher.NewNodeMatcherByLineAndCol(5, 14)

				It("should return the path to the token", func() {
					path, err := yaml.PathAtPoint(matcher)

					Expect(err).To(BeNil())
					Expect(len(path)).To(Equal(9))
					Expect(path[0].Kind).To(Equal(yamlv3.DocumentNode))
					Expect(path[1].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[2].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[2].Value).To(Equal("top"))
					Expect(path[3].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[4].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[4].Value).To(Equal("first"))
					Expect(path[5].Kind).To(Equal(yamlv3.SequenceNode))
					Expect(path[6].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[6].Value).To(Equal("0"))
					Expect(path[7].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[8].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[8].Value).To(Equal("attr2"))
				})
			})

			Context("indicating at sequence value node", func() {
				matcher := dmatcher.NewNodeMatcherByLineAndCol(7, 7)

				It("should return the path to the token", func() {
					path, err := yaml.PathAtPoint(matcher)

					Expect(err).To(BeNil())
					Expect(len(path)).To(Equal(7))
					Expect(path[0].Kind).To(Equal(yamlv3.DocumentNode))
					Expect(path[1].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[2].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[2].Value).To(Equal("top"))
					Expect(path[3].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[4].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[4].Value).To(Equal("first"))
					Expect(path[5].Kind).To(Equal(yamlv3.SequenceNode))
					Expect(path[6].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[6].Value).To(Equal("1"))
				})
			})

			Context("indicating at no token", func() {
				matcher := dmatcher.NewNodeMatcherByLineAndCol(6, 14)

				It("should return token not found error", func() {
					_, err := yaml.PathAtPoint(matcher)

					Expect(err).To(BeAssignableToTypeOf(dyaml.TokenNotFoundError{}))
				})
			})
		})

		Context("with multiple document yaml", func() {
			Context("indicating at document after the second", func() {
				It("should return the path to the token", func() {
					Expect(nil).NotTo(BeNil())
				})
			})
		})
	})
})

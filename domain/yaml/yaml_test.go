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
	multi := []byte(`first:
  - document
---
second:
  - document
`)

	var reader io.Reader

	Describe("NewYAML()", func() {
		Context("with single document yaml", func() {
			BeforeEach(func() {
				reader = bytes.NewReader(data)
			})

			It("should return parsed yaml", func() {
				_, err := dyaml.NewYAML(reader)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with multiple document yaml", func() {
			BeforeEach(func() {
				reader = bytes.NewReader(multi)
			})

			It("should return parsed yaml", func() {
				yaml, err := dyaml.NewYAML(reader)

				Expect(err).NotTo(HaveOccurred())
				Expect(*yaml).To(HaveLen(2))
			})
		})

		Context("with invalid yaml", func() {
			invalid := []byte("top: -")
			BeforeEach(func() {
				reader = bytes.NewReader(invalid)
			})

			It("should return an error", func() {
				_, err := dyaml.NewYAML(reader)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("PathAtPoint()", func() {
		var yaml *dyaml.YAML

		Context("with single document yaml", func() {
			BeforeEach(func() {
				var err error
				reader = bytes.NewReader(data)
				yaml, err = dyaml.NewYAML(reader)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("indicating at mapping value node", func() {
				matcher := dmatcher.NewNodeMatcherByLineAndCol(5, 14)

				It("should return the path to the token", func() {
					path, err := yaml.PathAtPoint(matcher)

					Expect(err).NotTo(HaveOccurred())
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

					Expect(err).NotTo(HaveOccurred())
					Expect(path).To(HaveLen(7))
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
			BeforeEach(func() {
				var err error
				reader = bytes.NewReader(multi)
				yaml, err = dyaml.NewYAML(reader)
				Expect(err).NotTo(HaveOccurred())
			})

			Context("indicating at document after the second", func() {
				matcher := dmatcher.NewNodeMatcherByLine(5)

				It("should return the path to the token", func() {
					path, err := yaml.PathAtPoint(matcher)

					Expect(err).NotTo(HaveOccurred())
					Expect(path).To(HaveLen(5))
					Expect(path[0].Kind).To(Equal(yamlv3.DocumentNode))
					Expect(path[1].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[2].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[2].Value).To(Equal("second"))
					Expect(path[3].Kind).To(Equal(yamlv3.SequenceNode))
					Expect(path[4].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[4].Value).To(Equal("0"))
				})
			})
		})
	})
})

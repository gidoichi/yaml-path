package yaml_test

import (
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

	Describe("NewYAML()", func() {
		Context("When given valid yaml", func() {
			It("should return parsed yaml.", func() {
				_, err := dyaml.NewYAML(data)

				Expect(err).To(BeNil())
			})
		})

		Context("When given invalid yaml", func() {
			invalid := []byte("top: -")

			It("should return an error.", func() {
				_, err := dyaml.NewYAML(invalid)

				Expect(err).NotTo(BeNil())
			})
		})
	})

	var _ = Describe("PathAtPoint()", func() {
		Context("When given valid yaml", func() {
			var yaml *dyaml.YAML
			BeforeEach(func() {
				var err error
				yaml, err = dyaml.NewYAML(data)
				Expect(err).To(BeNil())
			})

			Context("When indicating at some token", func() {
				matcher := dmatcher.NewNodeMatcherByLineAndCol(5, 14)

				It("should return the path to the token.", func() {
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
					Expect(path[6].Kind).To(Equal(yamlv3.MappingNode))
					Expect(path[7].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[7].Value).To(Equal("attr2"))
					Expect(path[8].Kind).To(Equal(yamlv3.ScalarNode))
					Expect(path[8].Value).To(Equal("val2"))
				})
			})

			Context("When indicating at no token", func() {
				matcher := dmatcher.NewNodeMatcherByLineAndCol(6, 14)

				It("should return token not found error.", func() {
					_, err := yaml.PathAtPoint(matcher)

					Expect(err).To(BeAssignableToTypeOf(dyaml.TokenNotFoundError{}))
				})
			})
		})
	})
})

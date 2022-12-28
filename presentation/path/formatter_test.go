package path_test

import (
	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	dyaml "github.com/gidoichi/yaml-path/domain/yaml"
	ppath "github.com/gidoichi/yaml-path/presentation/path"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	yamlv3 "gopkg.in/yaml.v3"
)

var _ = Describe("Path", func() {
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
	var path *ppath.Path
	BeforeEach(func() {
		matcher := dmatcher.NewNodeMatcherByLineAndCol(5, 14)
		var err error
		path, err = ppath.NewPath(data, matcher)
		Expect(err).To(BeNil())
	})

	Describe("ToString()", func() {
		Context("When invalid path (having invalid kind node) is given", func() {
			invalid := ppath.Path{
				Path: dyaml.Path{
					&yamlv3.Node{
						Kind: 0,
					},
				},
			}
			formatter := &ppath.PathFormatterBosh{}

			It("should fail to convert to string.", func() {
				_, err := invalid.ToString(formatter)

				Expect(err).NotTo(BeNil())
			})
		})

		Context("When invalid path (having no value node next sequence node) is given", func() {
			invalid := ppath.Path{
				Path: dyaml.Path{
					&yamlv3.Node{
						Kind: yamlv3.SequenceNode,
					},
				},
			}
			formatter := &ppath.PathFormatterBosh{}

			It("should fail to convert to string.", func() {
				_, err := invalid.ToString(formatter)

				Expect(err).NotTo(BeNil())
			})
		})

		Context("When path is converted using sequence selector indicating exist node", func() {
			formatter := &ppath.PathFormatterBosh{
				Separator: "/",
				NameAttr:  "name",
			}

			It("should converted to bosh format with selector.", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("/top/first/name=myname/attr2"))
			})
		})

		Context("When path is converted using sequence selector not indicating exist node", func() {
			formatter := &ppath.PathFormatterBosh{
				Separator: "/",
				NameAttr:  "dummy",
			}

			It("should converted to bosh format without selector.", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("/top/first/0/attr2"))
			})
		})

		Context("When path is converted to jsonpath format", func() {
			formatter := &ppath.PathFormatterJSONPath{}

			It("should converted to jsonpath format.", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("$.top.first[0].attr2"))
			})
		})
	})
})

package path_test

import (
	"bytes"

	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	dyaml "github.com/gidoichi/yaml-path/domain/yaml"
	ppath "github.com/gidoichi/yaml-path/presentation/path"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	yamlv3 "gopkg.in/yaml.v3"
)

var _ = Describe("Path", func() {
	var (
		path      *ppath.Path
		formatter ppath.PathFormatter
	)
	BeforeEach(func() {
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
		reader := bytes.NewReader(data)
		matcher := dmatcher.NewNodeMatcherByLineAndCol(5, 14)
		var err error
		path, err = ppath.NewPath(reader, matcher)
		Expect(err).To(BeNil())
	})

	Describe("ToString()", func() {
		Context("with invalid path (having invalid kind node)", func() {
			BeforeEach(func() {
				path = &ppath.Path{
					Path: dyaml.Path{
						&yamlv3.Node{
							Kind: 0,
						},
					},
				}
				formatter = &ppath.PathFormatterBosh{}
			})
			It("should fail to convert to string", func() {
				_, err := path.ToString(formatter)

				Expect(err).NotTo(BeNil())
			})
		})

		Context("with invalid path (having no value node next sequence node)", func() {
			BeforeEach(func() {
				path = &ppath.Path{
					Path: dyaml.Path{
						&yamlv3.Node{
							Kind: yamlv3.SequenceNode,
						},
					},
				}
				formatter = &ppath.PathFormatterBosh{}
			})

			It("should fail to convert to string", func() {
				_, err := path.ToString(formatter)

				Expect(err).NotTo(BeNil())
			})
		})

		Context("with path using sequence selector indicating exist node", func() {
			BeforeEach(func() {
				formatter = &ppath.PathFormatterBosh{
					Separator: "/",
					NameAttr:  "name",
				}
			})
			It("should convert to bosh format with selector", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("/top/first/name=myname/attr2"))
			})
		})

		Context("with path using sequence selector not indicating exist node", func() {
			BeforeEach(func() {
				formatter = &ppath.PathFormatterBosh{
					Separator: "/",
					NameAttr:  "dummy",
				}
			})
			It("should convert to bosh format without selector", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("/top/first/0/attr2"))
			})
		})

		Context("with path using sequence selector indicating exist conflicted node", func() {
			BeforeEach(func() {
				formatter = &ppath.PathFormatterBosh{
					Separator: "/",
					NameAttr:  "name",
				}
				data := []byte(`top:
  first:
    - name: myname
      attr1: val1
      attr2: val2
      #       ^
    - name: myname
      attr1: val1
      attr2: val2
    - value2
    - value3
  second:
    child1: value1
    child1: value2
    child3: value3
`)
				reader := bytes.NewReader(data)
				matcher := dmatcher.NewNodeMatcherByLineAndCol(5, 14)
				var err error
				path, err = ppath.NewPath(reader, matcher)
				Expect(err).To(BeNil())
			})

			It("should convert to bosh format without selector", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("/top/first/0/attr2"))
			})
		})

		Context("converting to jsonpath format", func() {
			BeforeEach(func() {
				formatter = &ppath.PathFormatterJSONPath{}
			})

			It("should convert to jsonpath format", func() {
				strpath, err := path.ToString(formatter)

				Expect(err).To(BeNil())
				Expect(strpath).To(Equal("$.top.first[0].attr2"))
			})
		})
	})
})

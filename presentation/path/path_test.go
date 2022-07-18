package path_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v3"

	dpath "github.com/gidoichi/yaml-path/domain/path"
	dsearcher "github.com/gidoichi/yaml-path/domain/searcher"
	"github.com/gidoichi/yaml-path/presentation/path"
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

	Context("When invalid path (having invalid kind node) is given", func() {
		It("should fail to convert to string.", func() {
			p := path.Path{
				Path: dpath.Path{
					&yaml.Node{
						Kind: 0,
					},
				},
			}

			_, err := p.ToString(path.Bosh)

			Expect(err).NotTo(BeNil())
		})
	})

	Context("When invalid path (having no value node next sequence node) is given", func() {
		It("should fail to convert to string.", func() {
			p := path.Path{
				Path: dpath.Path{
					&yaml.Node{
						Kind: yaml.SequenceNode,
					},
				},
			}

			_, err := p.ToString(path.Bosh)

			Expect(err).NotTo(BeNil())
		})
	})

	Context("When path is converted to bosh format", func() {
		It("should converted to bosh format.", func() {
			matcher := dsearcher.NodeMatcherByLineAndCol{}.New(5, 14)
			p, err := dsearcher.PathAtPoint(matcher, data)
			Expect(err).To(BeNil())

			strpath, err := path.Path{Path: p}.ToString(path.Bosh)

			Expect(err).To(BeNil())
			Expect(strpath).To(Equal("/top/first/name=myname/attr2"))
		})
	})

	Context("When path is converted to jsonpath format", func() {
		It("should converted to jsonpath format.", func() {
			matcher := dsearcher.NodeMatcherByLineAndCol{}.New(5, 14)
			p, err := dsearcher.PathAtPoint(matcher, data)
			Expect(err).To(BeNil())

			strpath, err := path.Path{Path: p}.ToString(path.JsonPath)

			Expect(err).To(BeNil())
			Expect(strpath).To(Equal("$.top.first[0].attr2"))
		})
	})
})

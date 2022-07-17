package path_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"

	"github.com/gidoichi/yaml-path/domain/path"
)

var _ = Describe("Path", func() {
	var p path.Path

	BeforeEach(func() {
		p = path.Path{
			&yaml.Node{},
			&yaml.Node{},
			&yaml.Node{},
		}
	})

	Context("When Reverse is called", func() {
		It("should return reversed path.", func() {
			rev := make(path.Path, len(p))
			copy(rev, p)

			rev = rev.Reverse()

			Expect(rev[0]).To(BeIdenticalTo(p[2]))
			Expect(rev[1]).To(BeIdenticalTo(p[1]))
			Expect(rev[2]).To(BeIdenticalTo(p[0]))
		})
	})
})

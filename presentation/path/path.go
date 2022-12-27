package path

import (
	"fmt"
	"reflect"
	"strings"

	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	dyaml "github.com/gidoichi/yaml-path/domain/yaml"
	yaml "gopkg.in/yaml.v3"
)

type Format string

const (
	Bosh     Format = "bosh"
	JsonPath Format = "jsonpath"
)

type Path struct {
	dyaml.Path
}

func (p *Path) Len() int {
	return len(p.Path)
}

func (p *Path) Get(i int) (node *yaml.Node, err error) {
	if p.Len() <= i {
		return nil, fmt.Errorf("index out of range: %d", i)
	}

	return p.Path[i], nil
}

var (
	Separator = "/"
	NameAttr  = "name"
)

func Configure(sep string, nameAttr string) {
	NameAttr = nameAttr
	Separator = sep
}

func NewPath(in []byte, matcher dmatcher.NodeMatcher) (path *Path, err error) {
	yaml, err := dyaml.NewYAML(in)
	if err != nil {
		return nil, err
	}

	p, err := yaml.PathAtPoint(matcher)
	if err != nil {
		return nil, err
	}

	return &Path{
		Path: p,
	}, nil
}

func (p *Path) get_node_name(node *yaml.Node) string {
	if node.Kind != yaml.MappingNode {
		return ""
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		if keyNode.Value != NameAttr {
			continue
		}
		if valNode.Kind != yaml.ScalarNode {
			return ""
		}
		return valNode.Value
	}

	return ""
}

func (p *Path) String() (strpath string) {
	var arr []string
	for _, node := range p.Path {
		arr = append(arr, fmt.Sprintf("%+v", node))
	}
	str := strings.Join(arr, ",")
	return fmt.Sprintf("%s[%s]", reflect.TypeOf(p), str)
}

func (p *Path) ToString(formatter PathFormatter) (strpath string, err error) {
	return formatter.ToString(p)
}

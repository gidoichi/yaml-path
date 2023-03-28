package path

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	dyaml "github.com/gidoichi/yaml-path/domain/yaml"
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

func (p *Path) Get(i int) (node *dyaml.Node, err error) {
	return p.Path.Get(i)
}

func NewPath(in io.Reader, matcher dmatcher.NodeMatcher) (path *Path, err error) {
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

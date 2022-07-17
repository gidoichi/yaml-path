package path

import (
	"fmt"
	"reflect"
	"strings"

	dpath "github.com/gidoichi/yaml-path/domain/path"
	yaml "gopkg.in/yaml.v3"
)

type Format string

const (
	Bosh     Format = "bosh"
	JsonPath Format = "jsonpath"
)

type Path struct {
	dpath.Path
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

func get_node_name(node *yaml.Node) string {
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

func (p Path) String() (strpath string) {
	var arr []string
	for _, node := range p.Path {
		arr = append(arr, fmt.Sprintf("%+v", node))
	}
	str := strings.Join(arr, ",")
	return fmt.Sprintf("%s[%s]", reflect.TypeOf(p), str)
}

func (p Path) ToString(format Format) (strpath string, err error) {
	switch format {
	case Bosh:
		for i := 0; i < p.Len(); i++ {
			cur, err := p.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			switch cur.Kind {
			case yaml.SequenceNode:
				var j int
				var c *yaml.Node
				next, err := p.Get(i + 1)
				if err != nil {
					return "", fmt.Errorf("get node: %w", err)
				}
				for j, c = range cur.Content {
					if c == next {
						break
					}
				}
				if name := get_node_name(c); name != "" {
					strpath += fmt.Sprintf("%s%s=%s", Separator, NameAttr, name)
				} else {
					strpath += fmt.Sprintf("%s%d", Separator, j)
				}
			case yaml.ScalarNode:
				prev, err := p.Get(i - 1)
				if err != nil {
					return "", fmt.Errorf("get node: %w", err)
				}
				if prev.Kind == yaml.ScalarNode || prev.Kind == yaml.SequenceNode {
					continue
				}
				strpath += Separator + cur.Value
			case yaml.DocumentNode, yaml.MappingNode, yaml.AliasNode:
				continue
			default:
				return "", fmt.Errorf("invalid path: %s", p)
			}
		}

	case JsonPath:
		for i := 0; i < p.Len(); i++ {
			cur, err := p.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			switch cur.Kind {
			case yaml.DocumentNode:
				strpath += "$"
			case yaml.SequenceNode:
				next, err := p.Get(i + 1)
				if err != nil {
					return "", fmt.Errorf("get node: %w", err)
				}
				for j, c := range cur.Content {
					if c == next {
						strpath += fmt.Sprintf("[%d]", j)
						break
					}
				}
			case yaml.ScalarNode:
				prev, err := p.Get(i - 1)
				if err != nil {
					return "", fmt.Errorf("get node: %w", err)
				}
				if prev.Kind == yaml.ScalarNode || prev.Kind == yaml.SequenceNode {
					continue
				}
				strpath += cur.Value
			case yaml.MappingNode, yaml.AliasNode:
				strpath += "."
			default:
				return "", fmt.Errorf("invalid path: %s", p)
			}
		}

	default:
		return "", fmt.Errorf("unsupported path format: %s", format)
	}

	return strpath, nil
}

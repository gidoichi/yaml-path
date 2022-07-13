package path

import (
	"fmt"

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

func (p *Path) Get(i int) *yaml.Node {
	return p.Path[i]
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

func (p Path) ToString(format Format) (strpath string, err error) {
	switch format {
	case Bosh:
		for i := 0; i < p.Len(); i++ {
			switch p.Get(i).Kind {
			case yaml.SequenceNode:
				var j int
				var c *yaml.Node
				for j, c = range p.Get(i).Content {
					if c == p.Get(i+1) {
						break
					}
				}
				if name := get_node_name(c); name != "" {
					strpath += fmt.Sprintf("%s%s=%s", Separator, NameAttr, name)
				} else {
					strpath += fmt.Sprintf("%s%d", Separator, j)
				}
			case yaml.ScalarNode:
				if p.Get(i-1).Kind == yaml.ScalarNode {
					continue
				}
				strpath += Separator + p.Get(i).Value
			case yaml.DocumentNode, yaml.MappingNode, yaml.AliasNode:
				continue
			default:
				panic(fmt.Sprintf("unreachable: Kind=%d", p.Get(i).Kind))
			}
		}

	case JsonPath:
		for i := 0; i < p.Len(); i++ {
			switch p.Get(i).Kind {
			case yaml.DocumentNode:
				strpath += "$"
			case yaml.SequenceNode:
				for j, c := range p.Get(i).Content {
					if c == p.Get(i+1) {
						strpath += fmt.Sprintf("[%d]", j)
						break
					}
				}
			case yaml.ScalarNode:
				if p.Get(i-1).Kind == yaml.ScalarNode {
					continue
				}
				strpath += "." + p.Get(i).Value
			case yaml.MappingNode, yaml.AliasNode:
				continue
			default:
				panic(fmt.Errorf("unreachable: Kind=%d", p.Get(i).Kind))
			}
		}

	default:
		return "", fmt.Errorf("unsupported path format: %s", format)
	}

	return strpath, nil
}

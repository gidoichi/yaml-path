package path

import (
	"fmt"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"
)

type PathFormatter interface {
	ToString(path *Path) (strpath string, err error)
}

type PathFormatterBosh struct {
	Separator string
	NameAttr  string
}

func (f *PathFormatterBosh) ToString(path *Path) (strpath string, err error) {
	var sb strings.Builder
	for i := 0; i < path.Len(); i++ {
		node, err := path.Get(i)
		if err != nil {
			return "", fmt.Errorf("get node: %w", err)
		}
		switch node.Kind {
		case yamlv3.SequenceNode:
			i++
			next, err := path.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			next2, err := path.Get(i + 1)
			if err == nil {
				if name := path.get_node_name((*yamlv3.Node)(next2), f.NameAttr); name != "" {
					sb.WriteString(f.Separator + f.NameAttr + "=" + name)
					continue
				}
			}
			sb.WriteString(f.Separator + next.Value)
		case yamlv3.MappingNode:
			i++
			next, err := path.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			sb.WriteString(f.Separator + next.Value)
		case yamlv3.DocumentNode, yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return sb.String(), nil
}

type PathFormatterJSONPath struct{}

func (f *PathFormatterJSONPath) ToString(path *Path) (strpath string, err error) {
	var sb strings.Builder
	for i := 0; i < path.Len(); i++ {
		node, err := path.Get(i)
		if err != nil {
			return "", fmt.Errorf("get node: %w", err)
		}
		switch node.Kind {
		case yamlv3.DocumentNode:
			sb.WriteString("$")
		case yamlv3.SequenceNode:
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			sb.WriteString("[" + next.Value + "]")
		case yamlv3.MappingNode:
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			sb.WriteString("." + next.Value)
		case yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return sb.String(), nil
}

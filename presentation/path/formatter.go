package path

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type PathFormatter interface {
	ToString(path *Path) (strpath string, err error)
}

type PathFormatterBosh struct{}

func (f *PathFormatterBosh) ToString(path *Path) (strpath string, err error) {
	for i := 0; i < path.Len(); i++ {
		cur, err := path.Get(i)
		if err != nil {
			return "", fmt.Errorf("get node: %w", err)
		}
		switch cur.Kind {
		case yaml.SequenceNode:
			var j int
			var c *yaml.Node
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			for j, c = range cur.Content {
				if c == next {
					break
				}
			}
			if name := path.get_node_name(c); name != "" {
				strpath += fmt.Sprintf("%s%s=%s", Separator, NameAttr, name)
			} else {
				strpath += fmt.Sprintf("%s%d", Separator, j)
			}
		case yaml.ScalarNode:
			prev, err := path.Get(i - 1)
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
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return strpath, nil
}

type PathFormatterJSONPath struct{}

func (f *PathFormatterJSONPath) ToString(path *Path) (strpath string, err error) {
	for i := 0; i < path.Len(); i++ {
		cur, err := path.Get(i)
		if err != nil {
			return "", fmt.Errorf("get node: %w", err)
		}
		switch cur.Kind {
		case yaml.DocumentNode:
			strpath += "$"
		case yaml.SequenceNode:
			next, err := path.Get(i + 1)
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
			prev, err := path.Get(i - 1)
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
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return strpath, nil
}

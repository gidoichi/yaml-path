package path

import (
	"fmt"

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
	for i := 0; i < path.Len(); i++ {
		cur, err := path.Get(i)
		if err != nil {
			return "", fmt.Errorf("get node: %w", err)
		}
		switch cur.Kind {
		case yamlv3.SequenceNode:
			var j int
			var c *yamlv3.Node
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			for j, c = range cur.Content {
				if c == next {
					break
				}
			}
			if name := path.get_node_name(c, f.NameAttr); name != "" {
				strpath += fmt.Sprintf("%s%s=%s", f.Separator, f.NameAttr, name)
			} else {
				strpath += fmt.Sprintf("%s%d", f.Separator, j)
			}
		case yamlv3.ScalarNode:
			prev, err := path.Get(i - 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			if prev.Kind == yamlv3.ScalarNode || prev.Kind == yamlv3.SequenceNode {
				continue
			}
			strpath += f.Separator + cur.Value
		case yamlv3.DocumentNode, yamlv3.MappingNode, yamlv3.AliasNode:
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
		case yamlv3.DocumentNode:
			strpath += "$"
		case yamlv3.SequenceNode:
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
		case yamlv3.ScalarNode:
			prev, err := path.Get(i - 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			if prev.Kind == yamlv3.ScalarNode || prev.Kind == yamlv3.SequenceNode {
				continue
			}
			strpath += cur.Value
		case yamlv3.MappingNode, yamlv3.AliasNode:
			strpath += "."
		default:
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return strpath, nil
}

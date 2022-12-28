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
			i++
			next, err := path.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			next2, err := path.Get(i + 1)
			if err == nil {
				if name := path.get_node_name((*yamlv3.Node)(next2), f.NameAttr); name != "" {
					strpath += f.Separator + f.NameAttr + "=" + name
					continue
				}
			}
			strpath += f.Separator + next.Value
		case yamlv3.MappingNode:
			i++
			next, err := path.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			strpath += f.Separator + next.Value
		case yamlv3.DocumentNode, yamlv3.ScalarNode, yamlv3.AliasNode:
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
			strpath += "[" + next.Value + "]"
		case yamlv3.MappingNode:
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			strpath += "." + next.Value
		case yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return strpath, nil
}

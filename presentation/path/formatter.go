package path

import (
	"fmt"
	"strconv"
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
	var builder strings.Builder
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
			seqidx, err := strconv.ParseInt(next.Value, 10, 0)
			if err != nil {
				return "", fmt.Errorf("invalid number: %w", err)
			}

			if name := node.FindSequenceSelectionByMappingKey(int(seqidx), f.NameAttr); name != "" {
				builder.WriteString(f.Separator + f.NameAttr + "=" + name)
				continue
			}
			builder.WriteString(f.Separator + next.Value)
		case yamlv3.MappingNode:
			i++
			next, err := path.Get(i)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			builder.WriteString(f.Separator + next.Value)
		case yamlv3.DocumentNode, yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return builder.String(), nil
}

type PathFormatterJSONPath struct{}

func (f *PathFormatterJSONPath) ToString(path *Path) (strpath string, err error) {
	var builder strings.Builder
	for i := 0; i < path.Len(); i++ {
		node, err := path.Get(i)
		if err != nil {
			return "", fmt.Errorf("get node: %w", err)
		}
		switch node.Kind {
		case yamlv3.DocumentNode:
			builder.WriteString("$")
		case yamlv3.SequenceNode:
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			builder.WriteString("[" + next.Value + "]")
		case yamlv3.MappingNode:
			next, err := path.Get(i + 1)
			if err != nil {
				return "", fmt.Errorf("get node: %w", err)
			}
			builder.WriteString("." + next.Value)
		case yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %s", path)
		}
	}

	return builder.String(), nil
}

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
	iter := path.Iter()
	for node := iter.Next(); node != nil; node = iter.Next() {
		if node == nil {
			return "", fmt.Errorf("get root node from path: %+v", path)
		}
		switch node.Kind {
		case yamlv3.SequenceNode:
			next := iter.Next()
			if next == nil {
				return "", fmt.Errorf("get next node with iterator: %+v", iter)
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
			next := iter.Next()
			if next == nil {
				return "", fmt.Errorf("get next node with iterator: %+v", iter)
			}
			builder.WriteString(f.Separator + next.Value)
		case yamlv3.DocumentNode, yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %+v", path)
		}
	}

	return builder.String(), nil
}

type PathFormatterJSONPath struct{}

func (f *PathFormatterJSONPath) ToString(path *Path) (strpath string, err error) {
	var builder strings.Builder
	iter := path.Iter()
	for node := iter.Next(); node != nil; node = iter.Next() {
		if node == nil {
			return "", fmt.Errorf("get root node from path: %+v", path)
		}
		switch node.Kind {
		case yamlv3.DocumentNode:
			builder.WriteString("$")
		case yamlv3.SequenceNode:
			next := iter.Next()
			if next == nil {
				return "", fmt.Errorf("get next node with iterator: %+v", iter)
			}
			builder.WriteString("[" + next.Value + "]")
		case yamlv3.MappingNode:
			next := iter.Next()
			if next == nil {
				return "", fmt.Errorf("get next node with iterator: %+v", iter)
			}
			builder.WriteString("." + next.Value)
		case yamlv3.ScalarNode, yamlv3.AliasNode:
			continue
		default:
			return "", fmt.Errorf("invalid path: %+v", path)
		}
	}

	return builder.String(), nil
}

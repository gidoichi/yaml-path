package searcher

import (
	"fmt"

	"github.com/gidoichi/yaml-path/domain/path"
	yaml "gopkg.in/yaml.v3"
)

type TokenNotFoundError struct {
	Line int
	Col  int
	Err  error
}

func (e TokenNotFoundError) Error() string {
	return fmt.Sprintf("token not found at %d:%d", e.Line, e.Col)
}

func node_match(line int, col int, node *yaml.Node) bool {
	if (node.Line == line) && (node.Column <= col) && (node.Column+len(node.Value) > col) {
		return true
	}
	return false
}

// findTokenAtPoint returns token path the arguments indicated.
// Note that returned path is reversed order.
//
// For example, when yaml path is $.top.first[0].attr2, then
// returned value is reversed order of (Document -> Mapping -> Scaler{"top"} ->
// Mapping -> Scaler{"first"} -> Sequence -> Mapping -> Scaler{"attr2"}).
func findTokenAtPoint(line int, col int, node *yaml.Node) (revpath path.Path, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			p, m := findTokenAtPoint(line, col, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			p, m := findTokenAtPoint(line, col, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if node_match(line, col, keyNode) {
				return path.Path{keyNode, node}, true
			}
			valNode := node.Content[i+1]
			p, m := findTokenAtPoint(line, col, valNode)
			if !m {
				continue
			}
			if p == nil {
				return path.Path{keyNode}, true
			} else {
				return append(p, keyNode, node), true
			}
		}

	case yaml.ScalarNode, yaml.AliasNode:
		if node_match(line, col, node) {
			return path.Path{node}, true
		}
	}

	return nil, false
}

func PathAtPoint(line int, col int, in []byte) (path path.Path, err error) {
	node := &yaml.Node{}
	if err := yaml.Unmarshal(in, node); err != nil {
		return nil, fmt.Errorf("cannot unmarshal yaml: %w", err)
	}

	rev, m := findTokenAtPoint(line, col, node)
	if !m {
		e := TokenNotFoundError{
			Line: line,
			Col:  col,
		}
		return nil, e
	}
	return rev.Reverse(), nil
}

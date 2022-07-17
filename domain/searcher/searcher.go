package searcher

import (
	"fmt"

	"github.com/gidoichi/yaml-path/domain/path"
	yaml "gopkg.in/yaml.v3"
)

type TokenNotFoundError struct {
	Matcher NodeMatcher
}

func (e TokenNotFoundError) Error() string {
	return fmt.Sprintf("token not found by %s", e.Matcher)
}

type NodeMatcher interface {
	Match(node *yaml.Node) bool
	String() string
}

type NodeMatcherByLine struct {
	line int
}

func (m NodeMatcherByLine) New(line int) *NodeMatcherByLine {
	return &NodeMatcherByLine{line: line}
}

func (m *NodeMatcherByLine) Match(node *yaml.Node) bool {
	return node.Line == m.line
}

func (m *NodeMatcherByLine) String() string {
	return fmt.Sprintf("{line: %d}", m.line)
}

type NodeMatcherByLineAndCol struct {
	line int
	col  int
}

func (m NodeMatcherByLineAndCol) New(line, col int) *NodeMatcherByLineAndCol {
	return &NodeMatcherByLineAndCol{line: line, col: col}
}

func (m *NodeMatcherByLineAndCol) Match(node *yaml.Node) bool {
	return (node.Line == m.line) &&
		(node.Column <= m.col) && (m.col < node.Column+len(node.Value))
}

func (m *NodeMatcherByLineAndCol) String() string {
	return fmt.Sprintf("{line: %d, col: %d}", m.line, m.col)
}

func node_match(matcher NodeMatcher, node *yaml.Node) bool {
	return matcher.Match(node)
}

// findTokenAtPoint returns token path the arguments indicated.
// Note that returned path is reversed order.
//
// For example, when yaml path is $.top.first[0].attr2, then
// returned value is reversed order of (Document -> Mapping -> Scaler{"top"} ->
// Mapping -> Scaler{"first"} -> Sequence -> Mapping -> Scaler{"attr2"}).
func findMatchedToken(matcher NodeMatcher, node *yaml.Node) (revpath path.Path, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			p, m := findMatchedToken(matcher, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			p, m := findMatchedToken(matcher, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if node_match(matcher, keyNode) {
				return path.Path{keyNode, node}, true
			}
			valNode := node.Content[i+1]
			p, m := findMatchedToken(matcher, valNode)
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
		if node_match(matcher, node) {
			return path.Path{node}, true
		}
	}

	return nil, false
}

func PathAtPoint(matcher NodeMatcher, in []byte) (path path.Path, err error) {
	node := &yaml.Node{}
	if err := yaml.Unmarshal(in, node); err != nil {
		return nil, fmt.Errorf("cannot unmarshal yaml: %w", err)
	}

	rev, m := findMatchedToken(matcher, node)
	if !m {
		e := TokenNotFoundError{
			Matcher: matcher,
		}
		return nil, e
	}
	return rev.Reverse(), nil
}

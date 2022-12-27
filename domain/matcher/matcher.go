package matcher

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type NodeMatcher interface {
	Match(node *yaml.Node) bool
	String() string
}

type NodeMatcherByLine struct {
	line int
}

func (m *NodeMatcherByLine) Match(node *yaml.Node) bool {
	return node.Line == m.line
}

func (m *NodeMatcherByLine) String() string {
	return fmt.Sprintf("{line: %d}", m.line)
}

func NewNodeMatcherByLine(line int) *NodeMatcherByLine {
	return &NodeMatcherByLine{line: line}
}

type NodeMatcherByLineAndCol struct {
	line int
	col  int
}

func (m *NodeMatcherByLineAndCol) Match(node *yaml.Node) bool {
	return (node.Line == m.line) &&
		(node.Column <= m.col) && (m.col < node.Column+len(node.Value))
}

func (m *NodeMatcherByLineAndCol) String() string {
	return fmt.Sprintf("{line: %d, col: %d}", m.line, m.col)
}

func NewNodeMatcherByLineAndCol(line, col int) *NodeMatcherByLineAndCol {
	return &NodeMatcherByLineAndCol{line: line, col: col}
}

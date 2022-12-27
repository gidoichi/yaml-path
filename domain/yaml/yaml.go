package yaml

import (
	"fmt"

	"github.com/gidoichi/yaml-path/domain/matcher"
	"gopkg.in/yaml.v3"
)

type YAML yaml.Node

type TokenNotFoundError struct {
	Matcher matcher.NodeMatcher
}

func (e TokenNotFoundError) Error() string {
	return fmt.Sprintf("token not found by %s", e.Matcher)
}

func NewYAML(in []byte) (dyaml *YAML, err error) {
	node := &yaml.Node{}
	if err := yaml.Unmarshal(in, node); err != nil {
		return nil, fmt.Errorf("cannot unmarshal yaml: %w", err)
	}
	if node.Kind != yaml.DocumentNode {
		return nil, fmt.Errorf("invalid yaml file: top level is not document node")
	}

	return (*YAML)(node), nil
}

func (y *YAML) PathAtPoint(matcher matcher.NodeMatcher) (path Path, err error) {
	rev, m := y.findMatchedToken(matcher, (*yaml.Node)(y))
	if !m {
		e := TokenNotFoundError{
			Matcher: matcher,
		}
		return nil, e
	}
	return y.reverse(rev), nil
}

// findTokenAtPoint returns token path the arguments indicated.
// Note that returned path is reversed order.
//
// For example, when yaml path is $.top.first[0].attr2, then
// returned value is reversed order of (Document -> Mapping -> Scaler{"top"} ->
// Mapping -> Scaler{"first"} -> Sequence -> Mapping -> Scaler{"attr2"}).
func (y *YAML) findMatchedToken(matcher matcher.NodeMatcher, node *yaml.Node) (revpath Path, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			p, m := y.findMatchedToken(matcher, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			p, m := y.findMatchedToken(matcher, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if y.node_match(matcher, keyNode) {
				return Path{keyNode, node}, true
			}
			valNode := node.Content[i+1]
			p, m := y.findMatchedToken(matcher, valNode)
			if !m {
				continue
			}
			if p == nil {
				return Path{keyNode}, true
			} else {
				return append(p, keyNode, node), true
			}
		}

	case yaml.ScalarNode, yaml.AliasNode:
		if y.node_match(matcher, node) {
			return Path{node}, true
		}
	}

	return nil, false
}

func (y *YAML) node_match(matcher matcher.NodeMatcher, node *yaml.Node) bool {
	return matcher.Match(node)
}

func (y *YAML) reverse(p Path) (path Path) {
	len := len(p)
	for i := len/2 - 1; i >= 0; i-- {
		opp := len - i - 1
		p[i], p[opp] = p[opp], p[i]
	}
	return p
}

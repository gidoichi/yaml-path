package yaml

import (
	"fmt"

	"github.com/gidoichi/yaml-path/domain/matcher"
	yamlv3 "gopkg.in/yaml.v3"
)

type YAML yamlv3.Node

type TokenNotFoundError struct {
	Matcher matcher.NodeMatcher
}

func (e TokenNotFoundError) Error() string {
	return fmt.Sprintf("token not found by %s", e.Matcher)
}

func NewYAML(in []byte) (dyaml *YAML, err error) {
	node := &yamlv3.Node{}
	if err := yamlv3.Unmarshal(in, node); err != nil {
		return nil, fmt.Errorf("cannot unmarshal yaml: %w", err)
	}
	if node.Kind != yamlv3.DocumentNode {
		return nil, fmt.Errorf("invalid yaml file: top level is not document node")
	}

	return (*YAML)(node), nil
}

func (y *YAML) PathAtPoint(matcher matcher.NodeMatcher) (path Path, err error) {
	rev, m := y.findMatchedToken(matcher, (*yamlv3.Node)(y))
	if !m {
		e := TokenNotFoundError{
			Matcher: matcher,
		}
		return nil, e
	}
	y.reverse(rev)
	path = rev

	len := path.Len()
	if path[len-3].Kind == yamlv3.MappingNode {
		path = path[:len-1]
	}
	return path, nil
}

// findTokenAtPoint returns token path the arguments indicated.
// Note that returned path is reversed order.
//
// For example, when yaml path is $.top.first[0].attr2, then
// returned value is reversed order of (Document -> Mapping -> Scaler{"top"} ->
// Mapping -> Scaler{"first"} -> Sequence -> Mapping -> Scaler{"attr2"}).
func (y *YAML) findMatchedToken(matcher matcher.NodeMatcher, node *yamlv3.Node) (revpath Path, match bool) {
	switch node.Kind {
	case yamlv3.DocumentNode:
		for _, child := range node.Content {
			p, m := y.findMatchedToken(matcher, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yamlv3.SequenceNode:
		for _, child := range node.Content {
			p, m := y.findMatchedToken(matcher, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yamlv3.MappingNode:
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

	case yamlv3.ScalarNode, yamlv3.AliasNode:
		if y.node_match(matcher, node) {
			return Path{node}, true
		}
	}

	return nil, false
}

func (y *YAML) node_match(matcher matcher.NodeMatcher, node *yamlv3.Node) bool {
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

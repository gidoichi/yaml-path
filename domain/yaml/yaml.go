package yaml

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gidoichi/yaml-path/domain/matcher"
	yamlv3 "gopkg.in/yaml.v3"
)

type YAML []yamlv3.Node

type TokenNotFoundError struct {
	Matcher matcher.NodeMatcher
}

func (e TokenNotFoundError) Error() string {
	return fmt.Sprintf("token not found by %s", e.Matcher)
}

func NewYAML(in io.Reader) (*YAML, error) {
	var dyaml YAML

	decoder := yamlv3.NewDecoder(in)
	for {
		var node yamlv3.Node
		if err := decoder.Decode(&node); err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		dyaml = append(dyaml, node)
	}

	return &dyaml, nil
}

func (y *YAML) PathAtPoint(matcher matcher.NodeMatcher) (Path, error) {
	for _, document := range *y {
		rev, found := y.findMatchedToken(matcher, &document)
		if !found {
			continue
		}

		y.reverse(rev)
		path := rev
		len := path.Len()
		if path[len-3].Kind == yamlv3.MappingNode || path[len-3].Kind == yamlv3.SequenceNode {
			path = path[:len-1]
		}
		return path, nil
	}
	return nil, TokenNotFoundError{
		Matcher: matcher,
	}
}

// findTokenAtPoint returns token path the arguments indicated.
// Note that returned path is reversed order.
//
// For example, when yaml path is $.top.first[0].attr2, then
// returned value is reversed order of (Document -> Mapping -> Scaler{"top"} ->
// Mapping -> Scaler{"first"} -> Sequence -> Scaler{"0"} -> Mapping -> Scaler{"attr2"}).
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
		for i, child := range node.Content {
			p, m := y.findMatchedToken(matcher, child)
			if !m {
				continue
			}
			index := &yamlv3.Node{
				Kind:  yamlv3.ScalarNode,
				Tag:   intTag,
				Value: strconv.Itoa(i),
			}
			return append(p, index, node), true
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

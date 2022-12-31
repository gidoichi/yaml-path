package yaml

import (
	yamlv3 "gopkg.in/yaml.v3"
)

type Node yamlv3.Node

const (
	intTag = "!!int"
)

func (n *Node) GetChildNodeByIndex(idx int) *Node {
	if n.Kind != yamlv3.SequenceNode {
		return nil
	}
	if idx >= len(n.Content) {
		return nil
	}

	return (*Node)(n.Content[idx])
}

func (n *Node) FindChildValueByKey(key string) string {
	if n.Kind != yamlv3.MappingNode {
		return ""
	}

	value := ""
	for i := 0; i < len(n.Content); i += 2 {
		keyNode := n.Content[i]
		if keyNode.Value != key {
			continue
		}
		valNode := n.Content[i+1]
		if valNode.Kind != yamlv3.ScalarNode {
			continue
		}
		if value != "" {
			return ""
		}
		value = valNode.Value
	}

	return value
}

package yaml

import (
	yamlv3 "gopkg.in/yaml.v3"
)

type Node yamlv3.Node

const (
	intTag = "!!int"
)

func (n *Node) FindChildValueByKey(key string) string {
	if n.Kind != yamlv3.MappingNode {
		return ""
	}

	for i := 0; i < len(n.Content); i += 2 {
		keyNode := n.Content[i]
		if keyNode.Value != key {
			continue
		}
		valNode := n.Content[i+1]
		if valNode.Kind != yamlv3.ScalarNode {
			continue
		}
		return valNode.Value
	}

	return ""
}

func (n *Node) FindSequenceSelectionByMappingKey(idx int, key string) string {
	if n.Kind != yamlv3.SequenceNode {
		return ""
	}

	target := (*Node)(n.Content[idx])
	var value string
	if value = target.FindChildValueByKey(key); value == "" {
		return ""
	}

	len := len(n.Content)
	for i := 0; i < len; i++ {
		if i == idx {
			continue
		}
		child := (*Node)(n.Content[i])
		if value := child.FindChildValueByKey(key); value != "" {
			return ""
		}
	}

	return value
}

package yaml

import "gopkg.in/yaml.v3"

type Node yaml.Node

func (n *Node) IsMappingKeyNode() bool {
	return true
}

func (n *Node) IsMappingValueNode() bool {
	return true
}

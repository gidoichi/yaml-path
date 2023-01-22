package yaml

import (
	"fmt"

	yamlv3 "gopkg.in/yaml.v3"
)

type Path []*yamlv3.Node

func (p *Path) Get(i int) (node *Node, err error) {
	if i < 0 || p.Len() <= i {
		return nil, fmt.Errorf("index out of range: %d", i)
	}
	return (*Node)((*p)[i]), nil
}

func (p *Path) Len() int {
	return len(*p)
}

func (p *Path) Iter() *NodeIterator {
	return NewNodeIterator(p)
}

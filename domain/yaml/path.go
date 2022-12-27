package yaml

import (
	yaml "gopkg.in/yaml.v3"
)

type Path []*yaml.Node

func (p Path) Get(i int) *Node {
	if i < 0 && p.Len() < i {
		return nil
	}
	return (*Node)(p[i])
}

func (p Path) Len() int {
	return len(p)
}

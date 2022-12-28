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

func (p *Path) Iter() *nodeIterator {
	return newNodeIterator(p)
}

type nodeIterator struct {
	path *Path
	idx  int
}

func newNodeIterator(path *Path) *nodeIterator {
	return &nodeIterator{
		path: path,
		idx:  -1,
	}
}

func (i *nodeIterator) Next() (node *Node, err error) {
	i.idx++
	node, err = i.path.Get(i.idx)
	if err != nil {
		i.idx--
		return nil, err
	}
	return node, nil
}

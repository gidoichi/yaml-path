package path

import (
	yaml "gopkg.in/yaml.v3"
)

type Path []*yaml.Node

func (p Path) Reverse() (path Path) {
	len := len(p)
	for i := len/2 - 1; i >= 0; i-- {
		opp := len - i - 1
		p[i], p[opp] = p[opp], p[i]
	}
	return p
}

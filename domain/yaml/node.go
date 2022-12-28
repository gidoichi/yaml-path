package yaml

import (
	yamlv3 "gopkg.in/yaml.v3"
)

type Node yamlv3.Node

const (
	intTag = "!!int"
)

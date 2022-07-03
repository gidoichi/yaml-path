package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/urfave/cli/v2"
	yaml "gopkg.in/yaml.v3"
)

type Format string

const (
	Bosh     Format = "bosh"
	JsonPath Format = "jsonpath"
)

type Path []*yaml.Node

var (
	Separator = "/"
	NameAttr  = "name"
)

func node_match(line int, col int, node *yaml.Node) bool {
	if (node.Line == line) && (node.Column <= col) && (node.Column+len(node.Value) > col) {
		return true
	}
	return false
}

func get_node_name(node *yaml.Node) string {
	if node.Kind != yaml.MappingNode {
		return ""
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		if keyNode.Value != NameAttr {
			continue
		}
		if valNode.Kind != yaml.ScalarNode {
			return ""
		}
		return valNode.Value
	}

	return ""
}

// findTokenAtPoint returns token path the arguments indicated.
// Note that returned path is reversed order.
//
// For example, when yaml path is $.top.first[0].attr2, then
// returned value is reversed order of (Document -> Mapping -> Scaler{"top"} ->
// Mapping -> Scaler{"first"} -> Sequence -> Mapping -> Scaler{"attr2"}).
func findTokenAtPoint(line int, col int, node *yaml.Node) (revpath Path, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			p, m := findTokenAtPoint(line, col, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.SequenceNode:
		for _, child := range node.Content {
			p, m := findTokenAtPoint(line, col, child)
			if !m {
				continue
			}
			return append(p, node), true
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if node_match(line, col, keyNode) {
				return Path{keyNode, node}, true
			}
			valNode := node.Content[i+1]
			p, m := findTokenAtPoint(line, col, valNode)
			if !m {
				continue
			}
			if p == nil {
				return Path{keyNode}, true
			} else {
				return append(p, keyNode, node), true
			}
		}

	case yaml.ScalarNode, yaml.AliasNode:
		if node_match(line, col, node) {
			return Path{node}, true
		}
	}

	return nil, false
}

func Configure(sep string, nameAttr string) {
	NameAttr = nameAttr
	Separator = sep
}

func (p Path) Reverse() (path Path) {
	len := len(p)
	for i := len/2 - 1; i >= 0; i-- {
		opp := len - i - 1
		p[i], p[opp] = p[opp], p[i]
	}
	return p
}

func (p Path) ToString(format Format) (strpath string, err error) {
	switch format {
	case Bosh:
		for i := 0; i < len(p); i++ {
			switch p[i].Kind {
			case yaml.SequenceNode:
				var j int
				var c *yaml.Node
				for j, c = range p[i].Content {
					if c == p[i+1] {
						break
					}
				}
				if name := get_node_name(c); name != "" {
					strpath += fmt.Sprintf("%s%s=%s", Separator, NameAttr, name)
				} else {
					strpath += fmt.Sprintf("%s%d", Separator, j)
				}
			case yaml.ScalarNode:
				if p[i-1].Kind == yaml.ScalarNode {
					continue
				}
				strpath += Separator + p[i].Value
			case yaml.DocumentNode, yaml.MappingNode, yaml.AliasNode:
				continue
			default:
				panic(fmt.Sprintf("unreachable: Kind=%d", p[i].Kind))
			}
		}

	case JsonPath:
		for i := 0; i < len(p); i++ {
			switch p[i].Kind {
			case yaml.DocumentNode:
				strpath += "$"
			case yaml.SequenceNode:
				for j, c := range p[i].Content {
					if c == p[i+1] {
						strpath += fmt.Sprintf("[%d]", j)
						break
					}
				}
			case yaml.ScalarNode:
				if p[i-1].Kind == yaml.ScalarNode {
					continue
				}
				strpath += "." + p[i].Value
			case yaml.MappingNode, yaml.AliasNode:
				continue
			default:
				panic(fmt.Sprintf("unreachable: Kind=%d", p[i].Kind))
			}
		}

	default:
		return "", errors.New(fmt.Sprintf("unsupported path format: %s", format))
	}

	return strpath, nil
}

func PathAtPoint(line int, col int, in []byte, format Format) (path string, err error) {
	node := &yaml.Node{}
	if err := yaml.Unmarshal(in, node); err != nil {
		return "", err
	}
	if node != nil {
		revp, m := findTokenAtPoint(line, col, node)
		if !m {
			return "", fmt.Errorf("token not found at %d:%d", line, col)
		}
		return revp.Reverse().ToString(format)
	}

	return "", nil
}

func main() {
	app := cli.NewApp()
	app.ArgsUsage = " "
	app.Usage = "Reads yaml and output a path corresponding to line and column"
	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:  "line",
			Usage: "cursor line",
			Value: 0,
		},
		&cli.UintFlag{
			Name:  "col",
			Usage: "cursor column",
			Value: 0,
		},
		&cli.StringFlag{
			Name:  "path",
			Usage: "set filepath, empty means stdin",
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "output format. \"bosh\" or \"jsonpath\"",
			Value: "bosh",
		},
		&cli.StringFlag{
			Name:  "sep",
			Usage: "set path separator",
			Value: "/",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "set attribut name, empty to disable",
			Value: "name",
		},
	}
	app.HideHelpCommand = true
	if info, ok := debug.ReadBuildInfo(); ok {
		app.Version = info.Main.Version
	}

	app.Action = func(c *cli.Context) error {
		var buf []byte
		var err error
		line := c.Uint("line")
		col := c.Uint("col")
		filePath := c.String("path")
		format := c.String("format")
		sep := c.String("sep")
		attr := c.String("name")

		Configure(sep, attr)
		if filePath != "" {
			if buf, err = ioutil.ReadFile(filePath); err != nil {
				return cli.Exit(err, 1)
			}
		} else {
			if buf, err = ioutil.ReadAll(os.Stdin); err != nil {
				return cli.Exit(err, 1)
			}
		}

		path, err := PathAtPoint(int(line), int(col), buf, Format(format))
		if err != nil {
			return cli.Exit(err, 1)
		}
		fmt.Println(path)

		return nil
	}

	app.Run(os.Args)
}

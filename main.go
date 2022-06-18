package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v3"
)

type Format string

const (
	Bosh     Format = "bosh"
	JsonPath Format = "jsonpath"
)

type Path []*yaml.Node

var (
	line      = kingpin.Flag("line", "Cursor line").Default("0").Int()
	col       = kingpin.Flag("col", "Cursor column").Default("0").Int()
	filePath  = kingpin.Flag("path", "Set filepath, empty means stdin").Default("").String()
	format    = kingpin.Flag("format", "Output format").Default(string(Bosh)).String()
	sep       = kingpin.Flag("sep", "Set path separator").Default("/").String()
	attr      = kingpin.Flag("name", "Set attribut name, empty to disable").Default("name").String()
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

	default:
		panic(fmt.Sprintf("unreachable: Kind=%d", node.Kind))
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
	version.Version = "-"
	version.Revision = "-"
	version.Branch = "-"
	version.BuildUser = "-"
	version.BuildDate = "-"
	if info, ok := debug.ReadBuildInfo(); ok {
		version.Version = info.Main.Version

		m := make(map[string]string, len(info.Settings))
		for _, s := range info.Settings {
			m[s.Key] = s.Value
		}
		if v, ok := m["vcs.revision"]; ok {
			version.Revision = v
		}
	}

	kingpin.Version(version.Print("yaml-path"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	Configure(*sep, *attr)
	var buff []byte
	var err error
	if *filePath != "" {
		buff, err = ioutil.ReadFile(*filePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		buff, _ = ioutil.ReadAll(os.Stdin)
	}
	path, err := PathAtPoint(*line, *col, buff, Format(*format))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(path)
	os.Exit(0)
}

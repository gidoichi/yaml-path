package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v3"
)

var (
	line      = kingpin.Flag("line", "Cursor line").Default("0").Int()
	col       = kingpin.Flag("col", "Cursor column").Default("0").Int()
	sep       = kingpin.Flag("sep", "Set path separator").Default("/").String()
	attr      = kingpin.Flag("name", "Set attribut name, empty to disable").Default("name").String()
	filePath  = kingpin.Flag("path", "Set filepath, empty means stdin").Default("").String()
	Separator = "/"
	NameAttr  = "name"
)

func node_match(line int, col int, node *yaml.Node) bool {
	// fmt.Printf("line: %d=%d,  col : %d <= %d < %d, token=%s\n", node.Line, line, node.Column, col, node.Column+len(node.Value), node.Value)
	if (node.Line == line) && (node.Column <= col) && (node.Column+len(node.Value) > col) {
		// fmt.Printf("match on %s!!\n", node.Value)
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

func findTokenAtPoint(line int, col int, node *yaml.Node) (addr string, match bool) {
	switch node.Kind {
	case yaml.DocumentNode:
		// root node
		for _, child := range node.Content {
			a, m := findTokenAtPoint(line, col, child)
			if !m {
				continue
			}
			// fmt.Printf("return (%s,%t)\n", "(doc)."+a, true)
			return Separator + a, true
		}

	case yaml.MappingNode:
		// map node
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if node_match(line, col, keyNode) {
				// fmt.Printf("return (%s,%t)\n", keyNode.Value, true)
				return keyNode.Value, true
			}
			valNode := node.Content[i+1]
			a, m := findTokenAtPoint(line, col, valNode)
			if !m {
				continue
			}
			if a == "" {
				// fmt.Printf("return (%s,%t)\n", keyNode.Value, true)
				return keyNode.Value, true
			} else {
				// fmt.Printf("return (%s,%t)\n", keyNode.Value+"."+a, true)
				return keyNode.Value + Separator + a, true
			}
		}

	case yaml.SequenceNode:
		// array node
		for idx, child := range node.Content {
			a, m := findTokenAtPoint(line, col, child)
			if !m {
				continue
			}
			name := get_node_name(child)
			if name != "" {
				// fmt.Printf("return (%s,%t)\n", fmt.Sprintf("[name=%s].%s", name, a), true)
				return fmt.Sprintf("%s=%s%s%s", NameAttr, name, Separator, a), true
			}
			// fmt.Printf("return (%s,%t)\n", fmt.Sprintf("[%d].%s", idx, a), true)
			return fmt.Sprintf("%d%s%s", idx, Separator, a), true
		}

	case yaml.ScalarNode:
		// fmt.Printf("%s = %s\n", path, node.Value)
		if node_match(line, col, node) {
			// fmt.Printf("return (%s,%t)\n", "", true)
			return "", true
		}
	}

	// fmt.Printf("return (%s,%t)\n", "", false)
	return "", false
}

func Configure(sep string, nameAttr string) {
	NameAttr = nameAttr
	Separator = sep
}

func PathAtPoint(line int, col int, in []byte) (path string, err error) {
	node := &yaml.Node{}
	yaml.Unmarshal(in, node)
	if node != nil {
		a, m := findTokenAtPoint(line, col, node)
		if false == m {
			return "", fmt.Errorf("token not found at %d:%d", line, col)
		}
		return a, nil
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
	path, err := PathAtPoint(*line, *col, buff)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(path)
	os.Exit(0)
}

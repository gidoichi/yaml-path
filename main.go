package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	//fmt.Printf("line: %d=%d,  col : %d <= %d < %d, token=%s\n", node.line, line, node.column, col, node.column+len(node.value), node.value)
	if (node.Line == line) && (node.Column <= col) && (node.Column+len(node.Value) > col) {
		// fmt.Printf("match on %s!!\n", node.value)
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
	if node.Kind == yaml.DocumentNode {
		// root node
		for _, child := range node.Content {
			a, m := findTokenAtPoint(line, col, child)
			if m == true {
				// fmt.Printf("return (%s,%t)\n", "(doc)."+a, true)
				return Separator + a, true
			}
		}
	} else if node.Kind == yaml.MappingNode {
		// map node
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			if node_match(line, col, keyNode) {
				// fmt.Printf("return (%s,%t)\n", keyNode.value, true)
				return keyNode.Value, true
			}
			valNode := node.Content[i+1]
			a, m := findTokenAtPoint(line, col, valNode)
			if m == true {
				if a == "" {
					// fmt.Printf("return (%s,%t)\n", keyNode.value, true)
					return keyNode.Value, true
				} else {
					// fmt.Printf("return (%s,%t)\n", keyNode.value+"."+a, true)
					return keyNode.Value + Separator + a, true
				}
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		// array node
		for idx, child := range node.Content {
			a, m := findTokenAtPoint(line, col, child)
			if m == true {
				name := get_node_name(child)
				if name != "" {
					// fmt.Printf("return (%s,%t)\n", fmt.Sprintf("[name=%s].%s", name, a), true)
					return fmt.Sprintf("%s=%s%s%s", NameAttr, name, Separator, a), true
				}
				// fmt.Printf("return (%s,%t)\n", fmt.Sprintf("[%d].%s", idx, a), true)
				return fmt.Sprintf("%d%s%s", idx, Separator, a), true
			}
		}
	} else if node.Kind == yaml.ScalarNode {
		// fmt.Printf("%s = %s\n", path, node.value)
		if true == node_match(line, col, node) {
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
	// defer handleErr(&err)
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
	path, err := PathAtPoint(*line-1, *col, buff)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(path)
	os.Exit(0)
}

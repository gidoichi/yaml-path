package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	ppath "github.com/gidoichi/yaml-path/presentation/path"
	"github.com/urfave/cli/v2"
)

func Run() {
	app := cli.NewApp()
	app.ArgsUsage = "--line value"
	app.Usage = "Reads yaml and output a path corresponding to leftmost token at line, or at (line, col)"
	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:     "line",
			Usage:    "cursor line",
			Required: true,
			Hidden:   true,
		},
		&cli.UintFlag{
			Name:  "col",
			Usage: "cursor column, zero to disable",
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
			Name:  "bosh.sep",
			Usage: "set path separator for bosh format",
			Value: "/",
		},
		&cli.StringFlag{
			Name:  "bosh.name",
			Usage: "set attribut name for bosh format, empty to disable",
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

		var formatter ppath.PathFormatter
		switch format {
		case "bosh":
			f := &ppath.PathFormatterBosh{}
			if sep := c.String("bosh.sep"); sep != "" {
				f.Separator = sep
			}
			if attr := c.String("bosh.name"); attr != "" {
				f.NameAttr = attr
			}
			formatter = f
		case "jsonpath":
			formatter = &ppath.PathFormatterJSONPath{}
		default:
			return cli.Exit(fmt.Errorf("unsupported path format: %s", format), 1)
		}

		if filePath != "" {
			if buf, err = ioutil.ReadFile(filePath); err != nil {
				return cli.Exit(fmt.Errorf("read from file: %w", err), 1)
			}
		} else {
			if buf, err = ioutil.ReadAll(os.Stdin); err != nil {
				return cli.Exit(fmt.Errorf("read from stdin: %w", err), 1)
			}
		}

		var matcher dmatcher.NodeMatcher
		if col == 0 {
			matcher = dmatcher.NewNodeMatcherByLine(int(line))
		} else {
			matcher = dmatcher.NewNodeMatcherByLineAndCol(int(line), int(col))
		}
		path, err := ppath.NewPath(buf, matcher)
		if err != nil {
			return cli.Exit(fmt.Errorf("resolve path: %w", err), 1)
		}
		strpath, err := path.ToString(formatter)
		if err != nil {
			return cli.Exit(fmt.Errorf("path formatting error: %s", format), 1)
		}
		fmt.Println(strpath)

		return nil
	}

	app.Run(os.Args)
}

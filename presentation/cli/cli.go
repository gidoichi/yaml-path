package cli

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	dmatcher "github.com/gidoichi/yaml-path/domain/matcher"
	ppath "github.com/gidoichi/yaml-path/presentation/path"
	"github.com/urfave/cli/v3"
)

var version string

func Run() {
	cmd := &cli.Command{
		ArgsUsage: "--line value",
		Usage:     "Reads yaml and output a path corresponding to leftmost token at line, or at (line, col)",
		Flags: []cli.Flag{
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
				Usage: `output format. "bosh" or "jsonpath"`,
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
		},
		HideHelpCommand: true,
		Action: func(ctx context.Context, c *cli.Command) error {
			var file *os.File
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
				if file, err = os.Open(filePath); err != nil {
					return cli.Exit(fmt.Errorf("read from file: %w", err), 1)
				}
			} else {
				file = os.Stdin
			}

			var matcher dmatcher.NodeMatcher
			if col == 0 {
				matcher = dmatcher.NewNodeMatcherByLine(int(line))
			} else {
				matcher = dmatcher.NewNodeMatcherByLineAndCol(int(line), int(col))
			}
			path, err := ppath.NewPath(file, matcher)
			if err != nil {
				return cli.Exit(fmt.Errorf("resolve path: %w", err), 1)
			}
			strpath, err := path.ToString(formatter)
			if err != nil {
				return cli.Exit(fmt.Errorf("path formatting error: %s", format), 1)
			}
			fmt.Println(strpath)

			return nil
		},
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		cmd.Version = info.Main.Version
	}
	if version != "" {
		cmd.Version = version
	}

	cmd.Run(context.Background(), os.Args)
}

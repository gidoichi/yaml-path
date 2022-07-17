package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/gidoichi/yaml-path/domain/searcher"
	"github.com/gidoichi/yaml-path/presentation/path"
	"github.com/urfave/cli/v2"
)

func Run() {
	app := cli.NewApp()
	app.ArgsUsage = "--line value"
	app.Usage = "Reads yaml and output a path corresponding to line (and column)"
	app.Flags = []cli.Flag{
		&cli.UintFlag{
			Name:     "line",
			Usage:    "cursor line",
			Required: true,
			Hidden:   true,
		},
		&cli.UintFlag{
			Name:  "col",
			Usage: "cursor column",
			Value: 1,
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
		sep := c.String("bosh.sep")
		attr := c.String("bosh.name")

		path.Configure(sep, attr)
		if filePath != "" {
			if buf, err = ioutil.ReadFile(filePath); err != nil {
				return cli.Exit(fmt.Errorf("read from file: %w", err), 1)
			}
		} else {
			if buf, err = ioutil.ReadAll(os.Stdin); err != nil {
				return cli.Exit(fmt.Errorf("read from stdin: %w", err), 1)
			}
		}

		var p path.Path
		p.Path, err = searcher.PathAtPoint(int(line), int(col), buf)
		if err != nil {
			if errors.As(err, &searcher.TokenNotFoundError{}) {
				return cli.Exit(err, 1)
			}
			return cli.Exit(fmt.Errorf("specify token path: %w", err), 1)
		}
		strpath, err := p.ToString(path.Format(format))
		if err != nil {
			return cli.Exit(fmt.Errorf("format path: %w", err), 1)
		}
		fmt.Println(strpath)

		return nil
	}

	app.Run(os.Args)
}

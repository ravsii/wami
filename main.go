package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/urfave/cli/v3"
)

var opts options

func main() {
	cmd := cli.Command{
		Name:  "wami",
		Usage: "What are my imports? (wami) is a cli for import analisys for go apps.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive",
				Destination: &opts.parse.recursive,
			},
			&cli.BoolFlag{
				Name:        "ignore-blank",
				Usage:       "ignore blank imports (e.g., '_ fmt')",
				Destination: &opts.parse.ignoreBlank,
			},
			&cli.BoolFlag{
				Name:        "ignore-dot",
				Usage:       "ignore dot imports (e.g., '. fmt')",
				Destination: &opts.parse.ignoreDot,
			},
			&cli.BoolFlag{
				Name:        "ignore-same",
				Usage:       "ignore imports using the same alias as the original package (e.g., 'fmt fmt')",
				Destination: &opts.parse.ignoreSame,
			},

			&cli.UintFlag{
				Name:        "min",
				Usage:       "minimal amount of usages to appear in the output (inclusive)",
				Destination: &opts.output.min,
			},
			&cli.UintFlag{
				Name:        "max",
				Usage:       "maximum amount of usages to appear in the output (inclusive)",
				Destination: &opts.output.max,
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "path",
				Destination: &opts.paths,
				Min:         1,
				Max:         -1,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return run()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := validateOptions(); err != nil {
		return fmt.Errorf("can't validate options: %w", err)
	}

	list, err := parseFiles()
	if err != nil {
		return fmt.Errorf("can't parse: %w", err)
	}

	var printer Printer //nolint
	printer = &TextPrinter{}
	if err := printer.Print(os.Stdout, listToOutput(list)); err != nil {
		return fmt.Errorf("can't print: %w", err)
	}

	return nil
}

// TODO: move
func listToOutput(list importList) ImportsData {
	imports := make(ImportsData, 0, len(list.imports))

	importDataCmp := func(a, b ImportData) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(a.Path, b.Path),
		)
	}
	aliasCmp := func(a, b Alias) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(a.Alias, b.Alias),
		)
	}

	for _, imp := range list.imports {
		// TODO: better filter system

		if opts.output.min > 0 && imp.total < (opts.output.min) ||
			opts.output.max > 0 && imp.total > opts.output.max {
			continue
		}

		aliases := make(Aliases, 0, len(imp.aliases))
		for alias, count := range imp.aliases {
			aliases = append(aliases, Alias{
				Count: count,
				Alias: alias,
			})
		}

		slices.SortFunc(aliases, aliasCmp)

		imports = append(imports, ImportData{
			Path:    imp.path,
			Count:   imp.total,
			Aliases: aliases,
		})
	}

	slices.SortFunc(imports, importDataCmp)

	return imports
}

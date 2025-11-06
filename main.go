package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"unsafe"

	"github.com/urfave/cli/v3"
)

func main() {
	var opts options

	cmd := cli.Command{
		Name:  "wami",
		Usage: "What are my imports? (wami) is a cli for import analisys for go apps.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "aliases-only",
				Aliases:     []string{"a"},
				Usage:       "only output imports that have aliases. Note: all imports will be parsed anyways, for a total amount of usages",
				Destination: &opts.output.aliasesOnly,
			},
			&cli.StringFlag{
				Name:        "format",
				Usage:       "output format (text, json)",
				DefaultText: string(textOutput),
				Aliases:     []string{"f"},
				Destination: (*string)(unsafe.Pointer(&opts.output.format)),
				Config:      cli.StringConfig{TrimSpace: true},
			},
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive",
				Destination: &opts.parse.recursive,
			},
			&cli.StringFlag{
				Name:        "include",
				Usage:       "regexp to include import paths",
				Destination: &opts.parse._includeStr,
				Config:      cli.StringConfig{TrimSpace: true},
			},

			&cli.StringFlag{
				Name:        "include-alias",
				Usage:       "regexp to include import aliases",
				Destination: &opts.parse._includeAliasStr,
				Config:      cli.StringConfig{TrimSpace: true},
			},

			&cli.StringFlag{
				Name:        "ignore",
				Usage:       "regexp to ignore import paths",
				Destination: &opts.parse._ignoreStr,
				Config:      cli.StringConfig{TrimSpace: true},
			},

			&cli.StringFlag{
				Name:        "ignore-alias",
				Usage:       "regexp to ignore import aliases",
				Destination: &opts.parse._ignoreAliasStr,
				Config:      cli.StringConfig{TrimSpace: true},
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
			if err := opts.prepare(); err != nil {
				return fmt.Errorf("can't validate options: %w", err)
			}
			return run(opts)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(opts options) error {
	storage, err := parseFiles(opts)
	if err != nil {
		return fmt.Errorf("can't parse: %w", err)
	}

	var printer Printer
	switch opts.output.format {
	case textOutput:
		printer = &TextPrinter{}
	case jsonOutput:
		printer = &JsonPrinter{}
	}

	if err := printer.Print(os.Stdout, storage.intoOutput()); err != nil {
		return fmt.Errorf("can't print: %w", err)
	}

	return nil
}

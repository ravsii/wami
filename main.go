package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/urfave/cli/v3"
)

func main() {
	var opts options

	cmd := cli.Command{
		Name:  "wami",
		Usage: "What are my imports? (wami) is a cli for import analisys for go apps.",
		Flags: []cli.Flag{
			// just add flags as they're being made.
			// We're sorting them afterwards.
			&cli.BoolFlag{
				Name:        "aliases-only",
				Aliases:     []string{"a"},
				Usage:       "only output imports that have aliases. Note: all imports will be parsed anyways, for a total amount of usages",
				Destination: &opts.output.aliasesOnly,
			},
			&cli.StringFlag{
				Name:        "format",
				Usage:       "output format (text, json)",
				Value:       formatTextColored,
				Aliases:     []string{"f"},
				Destination: &opts.output.format,
				Config:      cli.StringConfig{TrimSpace: true},
				Action: func(_ context.Context, _ *cli.Command, format string) error {
					switch format {
					case formatText, formatTextColored, formatJson:
						return nil
					default:
						return fmt.Errorf("unknown format: %s", format)
					}
				},
			},
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive",
				Destination: &opts.parse.recursive,
			},
			&cli.StringFlag{
				Name:   "include",
				Usage:  "`regexp` to include import paths",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.parse.include),
			},
			&cli.StringFlag{
				Name:   "include-alias",
				Usage:  "`regexp` to include import aliases",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.parse.includeAlias),
			},
			&cli.StringFlag{
				Name:   "ignore",
				Usage:  "`regexp` to ignore import paths",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.parse.ignore),
			},

			&cli.StringFlag{
				Name:   "ignore-alias",
				Usage:  "`regexp` to ignore import aliases",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.parse.ignoreAlias),
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
				UsageText:   "list of directories to parse for imports. For recursion see -r flag",
				Destination: &opts.paths,
				Min:         0,
				Max:         -1,
				Config:      cli.StringConfig{TrimSpace: true},
				Value:       "./...",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if len(opts.paths) == 0 {
				opts.paths = []string{"./..."}
			}

			return run(opts)
		},
	}

	sort.Sort(cli.FlagsByName(cmd.Flags))

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
	case formatText:
		printer = &TextPrinter{}
	case formatTextColored:
		printer = &TextPrinter{colored: true}
	case formatJson:
		printer = &JsonPrinter{}
	}

	if err := printer.Print(os.Stdout, storage.intoOutput()); err != nil {
		return fmt.Errorf("can't print: %w", err)
	}

	return nil
}

func makeParseRegexFunc(dst **regexp.Regexp) func(context.Context, *cli.Command, string) error {
	return func(_ context.Context, _ *cli.Command, regexStr string) error {
		if regexStr == "" {
			return nil
		}

		parsedRegex, err := regexp.Compile(regexStr)
		if err != nil {
			return fmt.Errorf("can't parse regex %q: %w", regexStr, err)
		}

		*dst = parsedRegex
		return nil
	}
}

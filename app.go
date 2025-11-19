package wami

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"

	"github.com/urfave/cli/v3"
)

func Run(args []string) {
	var opts Options

	cmdCount := cli.Command{
		Name:  "count",
		Usage: "Counts import usages and their aliases",
		Flags: []cli.Flag{
			// just add flags as they're being made.
			// We're sorting them afterwards.
			&cli.BoolFlag{
				Name:        "aliases-only",
				Aliases:     []string{"a"},
				Usage:       "only output imports that have aliases. Note: all imports will be parsed anyways, for a total amount of usages",
				Destination: &opts.Output.AliasesOnly,
			},
			&cli.StringFlag{
				Name:        "format",
				Usage:       "output format (text, text-colored, json, csv)",
				Value:       FormatTextColored,
				Aliases:     []string{"f"},
				Destination: &opts.Output.Format,
				Config:      cli.StringConfig{TrimSpace: true},
				Action: func(_ context.Context, _ *cli.Command, format string) error {
					switch format {
					case FormatText, FormatTextColored, FormatJson, FormatCsv:
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
				Destination: &opts.Parse.Recursive,
			},
			&cli.StringFlag{
				Name:   "include",
				Usage:  "`regexp` to include import paths",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.Parse.Include),
			},
			&cli.StringFlag{
				Name:   "include-alias",
				Usage:  "`regexp` to include import aliases",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.Parse.IncludeAlias),
			},
			&cli.StringFlag{
				Name:   "ignore",
				Usage:  "`regexp` to ignore import paths",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.Parse.Ignore),
			},

			&cli.StringFlag{
				Name:   "ignore-alias",
				Usage:  "`regexp` to ignore import aliases",
				Config: cli.StringConfig{TrimSpace: true},
				Action: makeParseRegexFunc(&opts.Parse.IgnoreAlias),
			},
			&cli.BoolFlag{
				Name:        "ignore-blank",
				Usage:       "ignore blank imports (e.g., '_ fmt')",
				Destination: &opts.Parse.IgnoreBlank,
			},
			&cli.BoolFlag{
				Name:        "ignore-dot",
				Usage:       "ignore dot imports (e.g., '. fmt')",
				Destination: &opts.Parse.IgnoreDot,
			},
			&cli.BoolFlag{
				Name:        "ignore-same",
				Usage:       "ignore imports using the same alias as the original package (e.g., 'fmt fmt')",
				Destination: &opts.Parse.IgnoreSame,
			},
			&cli.UintFlag{
				Name:        "min",
				Usage:       "minimal amount of usages to appear in the output (inclusive)",
				Destination: &opts.Output.Min,
			},
			&cli.UintFlag{
				Name:        "max",
				Usage:       "maximum amount of usages to appear in the output (inclusive)",
				Destination: &opts.Output.Max,
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "path",
				UsageText:   "list of directories to parse for imports. For recursion see -r flag",
				Destination: &opts.Paths,
				Min:         0,
				Max:         -1,
				Config:      cli.StringConfig{TrimSpace: true},
				Value:       "./...",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if len(opts.Paths) == 0 {
				opts.Paths = []string{"./..."}
			}

			storage, err := ParseFiles(opts)
			if err != nil {
				return fmt.Errorf("can't parse: %w", err)
			}

			var printer Printer
			switch opts.Output.Format {
			case FormatText:
				printer = &TextPrinter{}
			case FormatTextColored:
				printer = &TextPrinter{Colored: true}
			case FormatJson:
				printer = &JsonPrinter{}
			case FormatCsv:
				printer = &CsvPrinter{}
			}

			if err := printer.Print(os.Stdout, storage.IntoOuput()); err != nil {
				return fmt.Errorf("can't print: %w", err)
			}

			return nil
		},
	}

	cmdGraph := cli.Command{
		Name:  "graph",
		Usage: "Visualizes imports into a graph. Can be used to analyze internal-only imports, external imports and go.mod dependencies",
		Flags: []cli.Flag{},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "path",
				UsageText:   "list of directories",
				Destination: &opts.Paths,
				Min:         0,
				Max:         -1,
				Config:      cli.StringConfig{TrimSpace: true},
				Value:       "./...",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			if len(opts.Paths) == 0 {
				opts.Paths = []string{"./..."}
			}

			links, err := ParseGraphFiles(opts)
			if err != nil {
				return fmt.Errorf("can't parse: %w", err)
			}

			// for pkg, deps := range links {
			// 	_, err := fmt.Fprintf(os.Stdout, "%s imports:\n", pkg)
			// 	if err != nil {
			// 		return fmt.Errorf("package %s: %w", pkg, err)
			// 	}
			// 	for dep := range deps {
			// 		_, err := fmt.Fprintf(os.Stdout, "  -> %s\n", dep)
			// 		if err != nil {
			// 			return fmt.Errorf("dep %s: %w", dep, err)
			// 		}
			// 	}
			// }

			for pkg, deps := range links {
				for dep := range deps {
					_, err := fmt.Fprintf(os.Stdout, "%s -> %s\n", pkg, dep)
					if err != nil {
						return fmt.Errorf("dep %s: %w", dep, err)
					}
				}
			}

			return nil
		},
	}

	cmd := cli.Command{
		Name:  "wami",
		Usage: "What are my imports? (wami) is a cli for import analysis for go apps.",
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "path",
				UsageText:   "list of directories to parse for imports. For recursion see -r flag",
				Destination: &opts.Paths,
				Min:         0,
				Max:         -1,
				Config:      cli.StringConfig{TrimSpace: true},
				Value:       "./...",
			},
		},
		Commands: []*cli.Command{&cmdCount, &cmdGraph},
	}

	sort.Sort(cli.FlagsByName(cmdCount.Flags))
	sort.Sort(cli.FlagsByName(cmdGraph.Flags))
	sort.Sort(cli.FlagsByName(cmd.Flags))

	if err := cmd.Run(context.Background(), args); err != nil {
		log.Fatal(err)
	}
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

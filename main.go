package main

import (
	"cmp"
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

type options struct {
	paths     []string
	recursive bool
}

func main() {
	opts := options{}

	cmd := cli.Command{
		Name:  "wami",
		Usage: "What are my imports? (wami) is a cli for import analisys for go apps.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "recursive",
				Aliases:     []string{"r"},
				Usage:       "enables recursive walking for ALL paths. If disabled, only paths ending with '...' are treated as recursive",
				Destination: &opts.recursive,
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "path",
				Destination: &opts.paths,
				UsageText:   "123123",
				Min:         1,
				Max:         -1,
				Config:      cli.StringConfig{TrimSpace: true},
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			opts, err := validateAndFix(opts)
			if err != nil {
				return fmt.Errorf("can't validate options: %w", err)
			}

			var list importList
			if err := parseFiles(&list, opts); err != nil {
				fmt.Printf("Error: %v\n", err)
			}

			printImports(list)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func validateAndFix(opts options) (options, error) {
	unique := make(map[string]struct{})
	for _, p := range opts.paths {
		unique[p] = struct{}{}
	}

	opts.paths = make([]string, 0, len(unique))
	for p := range unique {
		opts.paths = append(opts.paths, p)
	}

	return opts, nil
}

func parseFiles(list *importList, opts options) error {
	mode := parser.ImportsOnly | parser.SkipObjectResolution

	for _, path := range opts.paths {
		isRecursive := opts.recursive || strings.HasSuffix(path, "...")
		if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, wdErr error) error {
			if !strings.HasSuffix(path, ".go") {
				return nil
			}

			if d.IsDir() {
				if isRecursive {
					return nil
				} else {
					return filepath.SkipDir
				}
			}

			if wdErr != nil {
				fmt.Printf("error when visiting %s: %v", path, wdErr)
				return nil
			}

			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, mode)
			if err != nil {
				fmt.Printf("Error parsing %s: %v\n", path, err)
				return nil
			}

			for _, imp := range file.Imports {
				if imp.Path == nil {
					fmt.Println("found empty import, skipping...")
					continue
				}

				if imp.Name != nil {
					list.addAliased(imp.Path.Value, imp.Name.Name)
				} else {
					list.add(imp.Path.Value)
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

const (
	printItem     = '├'
	printLastItem = '└'
)

type printableAlias struct {
	path    string
	usages  uint
	aliases map[string]uint
}

func printImports(list importList) {
	listArr := make([]printableAlias, 0, len(list.imports))
	for _, imp := range list.imports {
		listArr = append(listArr, printableAlias{
			path:    imp.path,
			usages:  imp.total,
			aliases: imp.aliases,
		})
	}
	slices.SortFunc(listArr, func(a, b printableAlias) int {
		return cmp.Or(
			cmp.Compare(b.usages, a.usages),
			cmp.Compare(a.path, b.path),
		)
	})

	for _, imp := range listArr {
		usage := "usage"
		if imp.usages > 1 {
			usage = "usages"
		}

		fmt.Printf("%s: %d total %s\n", imp.path, imp.usages, usage)
		printAliases(imp.aliases)
	}
}

func printAliases(aliases map[string]uint) {
	if len(aliases) == 0 {
		return
	}

	aliasArr := make([]printableAlias, 0, len(aliases))
	for path, usages := range aliases {
		aliasArr = append(aliasArr, printableAlias{
			path:   path,
			usages: usages,
		})
	}

	slices.SortFunc(aliasArr, func(a, b printableAlias) int {
		return cmp.Or(
			cmp.Compare(b.usages, a.usages),
			cmp.Compare(a.path, b.path),
		)
	})

	for i := range aliasArr {
		c := printItem
		if i == len(aliasArr)-1 {
			c = printLastItem
		}

		fmt.Printf("%4c %d usages as %s\n", c, aliasArr[i].usages, aliasArr[i].path)
	}
}

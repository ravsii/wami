package main

import (
	"cmp"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"."}
	}

	var list importList
	for _, arg := range args {
		if err := parseFiles(arg, &list); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
	}

	printImports(list)
}

func parseFiles(path string, list *importList) error {
	mode := parser.ImportsOnly | parser.SkipObjectResolution

	return filepath.WalkDir(path, func(path string, d fs.DirEntry, wdErr error) error {
		if !strings.HasSuffix(path, ".go") || d.IsDir() {
			return nil
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
	})
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
	listArr := make([]printableAlias, 0, len(list))
	for _, imp := range list {
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

package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func parseFiles() (importList, error) {
	mode := parser.ImportsOnly | parser.SkipObjectResolution

	var list importList

	for _, root := range opts.paths {
		// isRecursive := opts.parse.recursive || strings.HasSuffix(path, "...")
		isRecursive := opts.parse.recursive
		if strings.HasSuffix(root, "...") {
			isRecursive = true
			root = root[:len(root)-3]
		}

		if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, wdErr error) error {
			if d.IsDir() {
				if isRecursive || path == root {
					return nil
				}

				return filepath.SkipDir
			}

			if !strings.HasSuffix(path, ".go") {
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

				addAsAliased := false
				if imp.Name != nil {
					name := imp.Name.Name
					path := imp.Path.Value
					path = strings.Trim(path, `"\`)
					if strings.Contains(path, "/") {
						path = filepath.Base(path)
					}
					addAsAliased = (name != "." || !opts.parse.ignoreDot) && (name != "_" || !opts.parse.ignoreBlank) &&
						(name != path || !opts.parse.ignoreSame)

				}

				if addAsAliased {
					list.addAliased(imp.Path.Value, imp.Name.Name)
				} else {
					list.add(imp.Path.Value)
				}
			}

			return nil
		}); err != nil {
			return importList{}, err
		}
	}

	return list, nil
}

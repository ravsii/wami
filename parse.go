package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func parseFiles(opts options) (*importStorage, error) {
	mode := parser.ImportsOnly | parser.SkipObjectResolution

	fset := token.NewFileSet()
	storage := newStorage(opts)

	for _, root := range opts.paths {
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
				return fmt.Errorf("error when visiting %s: %w", path, wdErr)
			}

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
					storage.addAliased(imp.Path.Value, imp.Name.Name)
				} else {
					storage.add(imp.Path.Value)
				}
			}

			return nil
		}); err != nil {
			return nil, fmt.Errorf("walkDir: %w", err)
		}
	}

	return &storage, nil
}

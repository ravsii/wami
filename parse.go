package wami

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

func ParseFiles(opts Options) (*importStorage, error) {
	mode := parser.ImportsOnly | parser.SkipObjectResolution

	fset := token.NewFileSet()
	storage := NewStorage(opts)

	// struct{} is probably more efficient, but bool is way cleaner
	seen := make(map[string]bool, len(opts.Path))

	root := opts.Path

	isRecursive := opts.Parse.Recursive
	if strings.HasSuffix(root, "/...") {
		isRecursive = true
		root = root[:len(root)-3]
	}

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, wdErr error) error {
		if wdErr != nil {
			return fmt.Errorf("error visiting %s: %w", path, wdErr)
		}

		if d.IsDir() {
			if isRecursive || path == root {
				return nil
			}

			return filepath.SkipDir
		}

		if !strings.HasSuffix(path, ".go") || seen[path] {
			return nil
		}
		seen[path] = true

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
				storage.AddAliased(imp.Path.Value, imp.Name.Name)
			} else {
				storage.Add(imp.Path.Value)
			}
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("walkDir: %w", err)
	}

	return &storage, nil
}

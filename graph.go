package wami

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// importGraph maps "package path" -> set of "imported package paths"
type importGraph map[string]map[string]struct{}

func ParseGraphFiles(opts Options) (importGraph, error) {
	mode := parser.ImportsOnly | parser.SkipObjectResolution
	fset := token.NewFileSet()

	graph := make(importGraph)
	seen := make(map[string]bool)

	for _, root := range opts.Paths {
		isRecursive := opts.Parse.Recursive
		if strings.HasSuffix(root, "/...") {
			isRecursive = true
			root = root[:len(root)-3]
		}

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, wdErr error) error {
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

			// Determine the package key (full import path can be tricky; here
			// we use file path)
			pkg := file.Name.Name
			if _, ok := graph[pkg]; !ok {
				graph[pkg] = make(map[string]struct{})
			}

			// Collect imports
			for _, imp := range file.Imports {
				if imp.Path == nil {
					continue
				}
				impPath := strings.Trim(imp.Path.Value, `"`) // remove quotes
				graph[pkg][impPath] = struct{}{}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return graph, nil
}

func PrintGraph(g importGraph) {
}

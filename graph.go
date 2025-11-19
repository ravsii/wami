package wami

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// importGraph maps "package path" -> set of "imported package paths"
type importGraph map[string]map[string]struct{}

func ParseGraphFiles(opts Options) (importGraph, error) {
	mode := parser.ImportsOnly | parser.SkipObjectResolution
	fset := token.NewFileSet()

	graph := make(importGraph)
	seen := make(map[string]bool)

	packagePrefix := ""
	gomodFound := false

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

			if !gomodFound {
				if gomodPath, err := FindModulePath(path); err == nil {
					packagePrefix = gomodPath
					gomodFound = true
				}
			}

			file, err := parser.ParseFile(fset, path, nil, mode)
			if err != nil {
				fmt.Printf("Error parsing %s: %v\n", path, err)
				return nil
			}

			// Key node is the bare package name of the current file
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

				// If import is part of the same module, strip module prefix
				if packagePrefix != "" && strings.HasPrefix(impPath, packagePrefix) {
					rel := strings.TrimPrefix(impPath, packagePrefix+"/")
					// Use the last path component as the bare package name
					parts := strings.Split(rel, "/")
					impPath = parts[len(parts)-1]
				} else if packagePrefix != "" && impPath == packagePrefix {
					// Special case: importing module root
					impPath = filepath.Base(packagePrefix)
				}

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

func FindModulePath(filePath string) (string, error) {
	dir := filepath.Dir(filePath)

	for {
		modFile := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modFile); err == nil {
			data, err := os.ReadFile(modFile)
			if err != nil {
				return "", fmt.Errorf("failed to read go.mod: %w", err)
			}

			mod, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				return "", fmt.Errorf("failed to parse go.mod: %w", err)
			}

			if mod.Module == nil || mod.Module.Mod.Path == "" {
				return "", errors.New("module path not found in go.mod")
			}

			return mod.Module.Mod.Path, nil
		}

		// go.mod not found, go up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached filesystem root
		}
		dir = parent
	}

	return "", errors.New("go.mod not found in any parent directories")
}

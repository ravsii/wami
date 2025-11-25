package wami

import (
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

	root := opts.Path
	isRecursive := opts.Parse.Recursive
	if strings.HasSuffix(root, "/...") {
		isRecursive = true
		root = root[:len(root)-4]
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("root path %q: %w", root, err)
	}

	gomod, err := findGoMod(root)
	if err != nil {
		return nil, fmt.Errorf("find go.mod: %w", err)
	}
	projectPath := gomod.path

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

		pkg := file.Name.Name
		if _, ok := graph[pkg]; !ok {
			graph[pkg] = make(map[string]struct{})
		}

		for _, imp := range file.Imports {
			if imp.Path == nil {
				continue
			}

			impPath := strings.Trim(imp.Path.Value, `"`)
			if projectPath != "" && strings.HasPrefix(impPath, projectPath) {
				rel := strings.TrimPrefix(impPath, projectPath+"/")
				parts := strings.Split(rel, "/")
				impPath = parts[len(parts)-1]
			} else if projectPath != "" && impPath == projectPath {
				impPath = filepath.Base(projectPath)
			}

			graph[pkg][impPath] = struct{}{}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return graph, nil
}

type gomod struct {
	path string
}

func findGoMod(rootDir string) (gomod, error) {
	dir := rootDir
	for dir != "" {
		modFile := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modFile); err == nil {
			data, err := os.ReadFile(modFile)
			if err != nil {
				return gomod{}, fmt.Errorf("reading %s: %w", modFile, err)
			}

			mod, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				return gomod{}, fmt.Errorf("parsing %s: %w", modFile, err)
			}

			if mod.Module == nil || mod.Module.Mod.Path == "" {
				return gomod{}, fmt.Errorf("no module path in %s", modFile)
			}

			return gomod{path: mod.Module.Mod.Path}, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	return gomod{}, fmt.Errorf("go.mod not found in %s or any parent directories", rootDir)
}

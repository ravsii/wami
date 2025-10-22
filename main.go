package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	test "fmt"
)

// expandEllipsis expands ./... into all subdirectories recursively that contain Go files.
func expandEllipsis(arg string) ([]string, error) {
	var dirs []string

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".go") {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				fmt.Printf("Error parsing %s: %v\n", path, err)
				return nil
			}

			fmt.Printf("Parsed: %s\n", path)

			for _, imp := range file.Imports {
				test.Println(imp.Name)
				test.Println(imp.Path)
			}
		}

		return nil
	})

	return dirs, err
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"."}
	}

	for _, arg := range args {
		dirs, err := expandEllipsis(arg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		for _, dir := range dirs {
			fmt.Println(dir)
			// Here you could parse ASTs or process Go files in `dir`
		}
	}
}

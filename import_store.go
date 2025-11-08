package main

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"
)

type (
	importStorage struct {
		imports map[string]importItem
		opts    options
	}

	importItem struct {
		path    string
		total   uint
		aliases map[string]uint
	}
)

func newStorage(opts options) importStorage {
	return importStorage{
		imports: make(map[string]importItem),
		opts:    opts,
	}
}

func (s *importStorage) add(path string) {
	s.addAliased(path, "")
}

func (s *importStorage) addAliased(path, alias string) {
	path = strings.Trim(path, `"\`)

	if !s.shouldAddPath(path, alias) {
		return
	}

	item, ok := s.imports[path]
	if !ok {
		item = importItem{
			path: path,
		}
	}

	item.total++

	if s.shouldAddAlias(path, alias) {
		if len(item.aliases) == 0 {
			item.aliases = make(map[string]uint, 1)
		}
		item.aliases[alias]++
	}

	s.imports[path] = item
}

func (s *importStorage) shouldAddPath(path, alias string) bool {
	if s.opts.parse.include != nil && !s.opts.parse.include.MatchString(path) ||
		s.opts.parse.ignore != nil && s.opts.parse.ignore.MatchString(path) ||
		s.opts.parse.includeAlias != nil && !s.opts.parse.includeAlias.MatchString(alias) ||
		s.opts.parse.ignoreAlias != nil && s.opts.parse.ignoreAlias.MatchString(alias) {
		return false
	}

	return true
}

func (s *importStorage) shouldAddAlias(path, alias string) bool {
	if alias == "" ||
		alias == "." && s.opts.parse.ignoreDot ||
		alias == "_" && s.opts.parse.ignoreBlank ||
		alias == filepath.Base(path) && s.opts.parse.ignoreSame {
		return false
	}

	return true
}

var (
	importDataCmp = func(a, b OutputImports) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(a.Path, b.Path),
		)
	}
	aliasCmp = func(a, b Alias) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(a.Path, b.Path),
		)
	}
)

func (s *importStorage) intoOutput() []OutputImports {
	imports := make([]OutputImports, 0, len(s.imports))

	for _, imp := range s.imports {
		if s.opts.output.min > 0 && imp.total < (s.opts.output.min) ||
			s.opts.output.max > 0 && imp.total > s.opts.output.max ||
			s.opts.output.aliasesOnly && len(imp.aliases) == 0 {
			continue
		}

		aliases := make([]Alias, 0, len(imp.aliases))
		for alias, count := range imp.aliases {
			aliases = append(aliases, Alias{Count: count, Path: alias})
		}

		slices.SortFunc(aliases, aliasCmp)

		imports = append(imports, OutputImports{
			Path:    imp.path,
			Count:   imp.total,
			Aliases: aliases,
		})
	}

	slices.SortFunc(imports, importDataCmp)

	return imports
}

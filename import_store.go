package wami

import (
	"cmp"
	"path/filepath"
	"slices"
	"strings"
)

type (
	importStorage struct {
		imports map[string]importItem
		opts    Options
	}

	importItem struct {
		path    string
		total   uint
		aliases map[string]uint
	}
)

func NewStorage(opts Options) importStorage {
	return importStorage{
		imports: make(map[string]importItem),
		opts:    opts,
	}
}

func (s *importStorage) Add(path string) {
	s.AddAliased(path, "")
}

func (s *importStorage) AddAliased(path, alias string) {
	path = strings.Trim(path, `"\`)

	if !s.shouldAddPath(path, alias) {
		return
	}

	item, ok := s.imports[path]
	if !ok {
		item = importItem{path: path}
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

var (
	importsCmp = func(a, b OutputImports) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(a.Path, b.Path),
		)
	}
	aliasCmp = func(a, b Alias) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count),
			cmp.Compare(a.Name, b.Name),
		)
	}
)

func (s *importStorage) IntoOuput() []OutputImports {
	results := make([]OutputImports, 0, len(s.imports))

	for _, outputImport := range s.imports {
		// TODO: Better filter func
		if s.opts.Output.Min > 0 && outputImport.total < (s.opts.Output.Min) ||
			s.opts.Output.Max > 0 && outputImport.total > s.opts.Output.Max ||
			s.opts.Output.AliasesOnly && len(outputImport.aliases) == 0 {
			continue
		}

		resultImport := OutputImports{
			Path:  outputImport.path,
			Count: outputImport.total,
		}

		if len(outputImport.aliases) > 0 {
			aliases := make([]Alias, 0, len(outputImport.aliases))

			for alias, count := range outputImport.aliases {
				aliases = append(aliases, Alias{Count: count, Name: alias})
			}

			resultImport.Aliases = aliases
			slices.SortFunc(resultImport.Aliases, aliasCmp)
		}

		results = append(results, resultImport)
	}

	slices.SortFunc(results, importsCmp)

	return results
}

func (s *importStorage) shouldAddPath(path, alias string) bool {
	if s.opts.Parse.Include != nil && !s.opts.Parse.Include.MatchString(path) ||
		s.opts.Parse.Ignore != nil && s.opts.Parse.Ignore.MatchString(path) ||
		s.opts.Parse.IncludeAlias != nil && !s.opts.Parse.IncludeAlias.MatchString(alias) ||
		s.opts.Parse.IgnoreAlias != nil && s.opts.Parse.IgnoreAlias.MatchString(alias) {
		return false
	}

	return true
}

func (s *importStorage) shouldAddAlias(path, alias string) bool {
	if alias == "" ||
		alias == "." && s.opts.Parse.IgnoreDot ||
		alias == "_" && s.opts.Parse.IgnoreBlank ||
		alias == filepath.Base(path) && s.opts.Parse.IgnoreSame {
		return false
	}

	return true
}

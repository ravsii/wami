package wami_test

import (
	"regexp"
	"testing"

	"github.com/ravsii/wami"
	"github.com/stretchr/testify/assert"
)

func TestImportStorage_Add(t *testing.T) {
	s := wami.NewStorage(wami.Options{})
	s.Add("fmt")

	want := []wami.OutputImports{{Path: "fmt", Count: 1}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_AddAliased(t *testing.T) {
	s := wami.NewStorage(wami.Options{})
	s.Add("fmt")
	s.AddAliased("fmt", "f")

	want := []wami.OutputImports{{
		Path:    "fmt",
		Count:   2,
		Aliases: []wami.Alias{{Name: "f", Count: 1}},
	}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IgnoreDotAlias(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.IgnoreDot = true

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", ".")

	want := []wami.OutputImports{{Path: "fmt", Count: 1}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IgnoreBlankAlias(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.IgnoreBlank = true

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", "_")

	want := []wami.OutputImports{{Path: "fmt", Count: 1}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IgnoreSameAlias(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.IgnoreSame = true

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", "fmt")

	want := []wami.OutputImports{{Path: "fmt", Count: 1}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IncludePath(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.Include = regexp.MustCompile(`^fmt$`)

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", "f")
	s.AddAliased("os", "o")

	want := []wami.OutputImports{{
		Path:    "fmt",
		Count:   1,
		Aliases: []wami.Alias{{Name: "f", Count: 1}},
	}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IgnorePath(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.Ignore = regexp.MustCompile(`^os$`)

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", "f")
	s.AddAliased("os", "o")

	want := []wami.OutputImports{{
		Path:    "fmt",
		Count:   1,
		Aliases: []wami.Alias{{Name: "f", Count: 1}},
	}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IncludeAlias(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.IncludeAlias = regexp.MustCompile(`^f$`)

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", "f")
	s.AddAliased("fmt", "x")

	want := []wami.OutputImports{{
		Path:    "fmt",
		Count:   1,
		Aliases: []wami.Alias{{Name: "f", Count: 1}},
	}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_IgnoreAlias(t *testing.T) {
	opts := wami.Options{}
	opts.Parse.IgnoreAlias = regexp.MustCompile(`^ignore$`)

	s := wami.NewStorage(opts)
	s.AddAliased("fmt", "f")
	s.AddAliased("fmt", "ignore")

	want := []wami.OutputImports{{
		Path:    "fmt",
		Count:   1,
		Aliases: []wami.Alias{{Name: "f", Count: 1}},
	}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_MinFilter(t *testing.T) {
	opts := wami.Options{}
	opts.Output.Min = 2

	s := wami.NewStorage(opts)
	s.Add("fmt")

	assert.Empty(t, s.IntoOuput())
}

func TestImportStorage_MaxFilter(t *testing.T) {
	opts := wami.Options{}
	opts.Output.Max = 1

	s := wami.NewStorage(opts)
	s.Add("fmt")
	s.Add("fmt")

	assert.Empty(t, s.IntoOuput())
}

func TestImportStorage_AliasesOnlyFilter(t *testing.T) {
	opts := wami.Options{}
	opts.Output.AliasesOnly = true

	s := wami.NewStorage(opts)
	s.Add("fmt")
	s.AddAliased("os", "o")

	want := []wami.OutputImports{{
		Path:    "os",
		Count:   1,
		Aliases: []wami.Alias{{Name: "o", Count: 1}},
	}}
	got := s.IntoOuput()

	assert.Equal(t, want, got)
}

func TestImportStorage_Sorting(t *testing.T) {
	s := wami.NewStorage(wami.Options{})

	s.AddAliased("b", "b1")
	s.AddAliased("b", "b2")
	s.AddAliased("a", "a1")
	s.AddAliased("a", "a2")
	s.AddAliased("c", "c1")
	s.AddAliased("c", "c1")

	got := s.IntoOuput()

	want := []wami.OutputImports{
		{
			Path:    "a",
			Count:   2,
			Aliases: []wami.Alias{{Name: "a1", Count: 1}, {Name: "a2", Count: 1}},
		},
		{
			Path:    "b",
			Count:   2,
			Aliases: []wami.Alias{{Name: "b1", Count: 1}, {Name: "b2", Count: 1}},
		},
		{
			Path:    "c",
			Count:   2,
			Aliases: []wami.Alias{{Name: "c1", Count: 2}},
		},
	}

	assert.Equal(t, want, got)
}

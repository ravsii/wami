package wami

import (
	"regexp"
)

const (
	FormatText        = "text"
	FormatTextColored = "text-colored"
	FormatJson        = "json"
	FormatCsv         = "csv"
)

type (
	Options struct {
		Paths []string

		// options related to parsing
		Parse ParseOptions

		// options related to Output
		Output OutputOptions
	}

	ParseOptions struct {
		Recursive bool

		Include      *regexp.Regexp
		IncludeAlias *regexp.Regexp
		Ignore       *regexp.Regexp
		IgnoreAlias  *regexp.Regexp

		IgnoreDot   bool
		IgnoreBlank bool
		IgnoreSame  bool
	}

	OutputOptions struct {
		AliasesOnly bool
		Format      string
		Max         uint
		Min         uint
	}
)

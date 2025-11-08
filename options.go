package main

import (
	"regexp"
)

const (
	formatText        = "text"
	formatTextColored = "text-colored"
	formatJson        = "json"
)

type (
	options struct {
		paths []string

		// options related to parsing
		parse parseOptions

		// options related to output
		output outputOptions
	}

	parseOptions struct {
		recursive bool

		include      *regexp.Regexp
		includeAlias *regexp.Regexp
		ignore       *regexp.Regexp
		ignoreAlias  *regexp.Regexp

		ignoreDot   bool
		ignoreBlank bool
		ignoreSame  bool
	}

	outputOptions struct {
		aliasesOnly bool
		format      string
		max         uint
		min         uint
	}
)

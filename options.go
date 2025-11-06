package main

import (
	"fmt"
	"regexp"
)

type outputFormat string

const (
	textOutput outputFormat = "text"
	jsonOutput outputFormat = "json"
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

		_includeStr      string
		_includeAliasStr string
		_ignoreStr       string
		_ignoreAliasStr  string
		include          *regexp.Regexp
		includeAlias     *regexp.Regexp
		ignore           *regexp.Regexp
		ignoreAlias      *regexp.Regexp

		ignoreDot   bool
		ignoreBlank bool
		ignoreSame  bool
	}

	outputOptions struct {
		aliasesOnly bool
		format      outputFormat
		max         uint
		min         uint
	}
)

// prepare checks and fixes any conflicting options (i.e. having both
// "./..." path and --recursive flag).
// If it encounters an option it can't fix, an error is returned.
func (o *options) prepare() error {
	unique := make(map[string]struct{})
	for _, p := range o.paths {
		unique[p] = struct{}{}
	}

	o.paths = make([]string, 0, len(unique))
	for p := range unique {
		o.paths = append(o.paths, p)
	}

	var err error

	// parse options parsing

	if o.parse._includeStr != "" {
		o.parse.include, err = regexp.Compile(o.parse._includeStr)
		if err != nil {
			return fmt.Errorf("parsing include regex %q: %w", o.parse._includeStr, err)
		}
	}

	if o.parse._includeAliasStr != "" {
		o.parse.includeAlias, err = regexp.Compile(o.parse._includeAliasStr)
		if err != nil {
			return fmt.Errorf("parsing include-alias regex %q: %w", o.parse._includeAliasStr, err)
		}
	}

	if o.parse._ignoreStr != "" {
		o.parse.ignore, err = regexp.Compile(o.parse._ignoreStr)
		if err != nil {
			return fmt.Errorf("parsing ignore regex %q: %w", o.parse._ignoreStr, err)
		}
	}

	if o.parse._ignoreAliasStr != "" {
		o.parse.ignoreAlias, err = regexp.Compile(o.parse._ignoreAliasStr)
		if err != nil {
			return fmt.Errorf("parsing ignore-alias regex %q: %w", o.parse._ignoreAliasStr, err)
		}
	}

	// output options parsing

	switch o.output.format {
	case "":
		o.output.format = textOutput
	case textOutput, jsonOutput:
		break
	default:
		return fmt.Errorf("unknown output format: %s", o.output.format)
	}

	return nil
}

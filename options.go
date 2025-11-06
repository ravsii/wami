package main

import (
	"fmt"
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
		recursive   bool
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

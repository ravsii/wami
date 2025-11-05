package main

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
		min uint
		max uint
	}
)

// validateOptions checks and fixes any conflicting options (i.e. having both
// "./..." path and --recursive flag).
// If it encounters an option it can't fix, an error is returned.
func validateOptions() error {
	unique := make(map[string]struct{})
	for _, p := range opts.paths {
		unique[p] = struct{}{}
	}

	opts.paths = make([]string, 0, len(unique))
	for p := range unique {
		opts.paths = append(opts.paths, p)
	}

	return nil
}

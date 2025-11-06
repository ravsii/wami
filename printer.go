package main

import "io"

type Printer interface {
	// Print should output data into w.
	Print(w io.Writer, data []OutputImports) error
}

type (
	OutputImports struct {
		Path    string
		Count   uint
		Aliases []Alias
	}

	Alias struct {
		Count uint
		Alias string
	}
)

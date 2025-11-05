package main

import "io"

type Printer interface {
	// Print should output data into w.
	Print(w io.Writer, data ImportsData) error
}

type (
	ImportsData []ImportData
	ImportData  struct {
		Path    string
		Count   uint
		Aliases Aliases
	}

	Aliases []Alias
	Alias   struct {
		Count uint
		Alias string
	}
)

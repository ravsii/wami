package wami

import "io"

type Printer interface {
	// Print should output data into w.
	Print(w io.Writer, data []OutputImports) error
}

type (
	OutputImports struct {
		Path    string  `json:"path"`
		Count   uint    `json:"count"`
		Aliases []Alias `json:"aliases,omitempty"`
	}

	Alias struct {
		Count uint   `json:"count"`
		Name  string `json:"name"`
	}
)

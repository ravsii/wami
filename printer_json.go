package main

import (
	"encoding/json"
	"io"
)

var _ Printer = (*JsonPrinter)(nil)

type JsonPrinter struct{}

func (JsonPrinter) Print(w io.Writer, imports []OutputImports) error {
	return json.NewEncoder(w).Encode(imports)
}

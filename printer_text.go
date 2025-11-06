package main

import (
	"bytes"
	"fmt"
	"io"
)

const (
	printItem     = '├'
	printLastItem = '└'
)

var _ Printer = (*TextPrinter)(nil)

type TextPrinter struct{}

func (t *TextPrinter) Print(w io.Writer, imports []OutputImports) error {
	var buf bytes.Buffer

	for _, imprt := range imports {
		usageStr := "usage"
		if imprt.Count > 1 {
			usageStr = "usages"
		}

		fmt.Fprintf(&buf, "%q: %d total %s\n", imprt.Path, imprt.Count, usageStr)
		for i, alias := range imprt.Aliases {
			prefix := printItem
			if i == len(imprt.Aliases)-1 {
				prefix = printLastItem
			}

			fmt.Fprintf(&buf, "%4c %d usages as %q\n", prefix, alias.Count, alias.Alias)
		}
	}

	_, err := buf.WriteTo(w)
	return err
}

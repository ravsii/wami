package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fatih/color"
)

const (
	printItem     = '├'
	printLastItem = '└'
)

var _ Printer = (*TextPrinter)(nil)

type TextPrinter struct{ colored bool }

var (
	colorName   = color.New(color.FgHiRed, color.Bold, color.Italic).SprintFunc()
	colorCount  = color.New(color.FgHiYellow).SprintFunc()
	colorPrefix = color.New(color.FgHiBlack).SprintFunc()
	colorAlias  = color.New(color.FgHiBlue, color.Italic).SprintFunc()
)

func (t *TextPrinter) Print(w io.Writer, imports []OutputImports) error {
	var buf bytes.Buffer
	color.NoColor = !t.colored

	for _, imprt := range imports {
		usageStr := "usage"
		if imprt.Count > 1 {
			usageStr = "usages"
		}

		fmt.Fprintf(&buf, "%s: %s total %s\n", colorName(imprt.Path), colorCount(imprt.Count), usageStr)
		for i, alias := range imprt.Aliases {
			prefix := printItem
			if i == len(imprt.Aliases)-1 {
				prefix = printLastItem
			}

			fmt.Fprintf(&buf, "%1c%s %s usages as %s\n", ' ', colorPrefix(string(prefix)), colorCount(alias.Count), colorAlias(alias.Path))
		}
	}

	_, err := buf.WriteTo(w)
	return err
}

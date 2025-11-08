package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var _ Printer = (*CsvPrinter)(nil)

type CsvPrinter struct{}

func (t *CsvPrinter) Print(w io.Writer, imports []OutputImports) error {
	var buf bytes.Buffer

	records := [][]string{
		{"import", "count", "aliases"},
	}

	writer := csv.NewWriter(&buf)

	for _, imprt := range imports {
		record := []string{
			imprt.Path,
			strconv.FormatUint(uint64(imprt.Count), 10),
			"",
		}

		if len(imprt.Aliases) > 0 {
			aliases := make([]string, 0, len(imprt.Aliases))
			for _, alias := range imprt.Aliases {
				count := strconv.FormatUint(uint64(alias.Count), 10)
				aliases = append(aliases, count+","+alias.Path)
			}
			record[2] = strings.Join(aliases, ";")
		}

		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("write csv records: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("flushing: %w", err)
	}

	_, err := buf.WriteTo(w)
	return err
}

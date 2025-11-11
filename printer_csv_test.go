package wami_test

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/ravsii/wami"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCsvPrinter_Print(t *testing.T) {
	p := &wami.CsvPrinter{}
	in := []wami.OutputImports{
		{
			Path:  "fmt",
			Count: 4,
			Aliases: []wami.Alias{
				{Count: 3, Name: "f"},
				{Count: 1, Name: "format"},
			},
		},
		{
			Path:  "os",
			Count: 1,
		},
	}
	want := [][]string{
		{"import", "count", "aliases"},
		{"fmt", "4", "3,f;1,format"},
		{"os", "1", ""},
	}

	var buf bytes.Buffer
	err := p.Print(&buf, in)
	require.NoError(t, err)

	records, err := csv.NewReader(&buf).ReadAll()
	require.NoError(t, err)
	assert.Equal(t, want, records)
}

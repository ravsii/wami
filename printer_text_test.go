package wami_test

import (
	"bytes"
	"testing"

	"github.com/ravsii/wami"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextPrinter_NoColor(t *testing.T) {
	p := &wami.TextPrinter{Colored: false}
	in := []wami.OutputImports{
		{
			Path:  "fmt",
			Count: 2,
			Aliases: []wami.Alias{
				{Count: 3, Name: "f"},
				{Count: 1, Name: "format"},
			},
		},
		{
			Path:  "os",
			Count: 1,
			Aliases: []wami.Alias{
				{Count: 1, Name: "o"},
			},
		},
	}

	var buf bytes.Buffer
	err := p.Print(&buf, in)
	require.NoError(t, err)

	want := `fmt: 2 total usages
 ├ 3 usages as f
 └ 1 usages as format
os: 1 total usage
 └ 1 usages as o
`

	assert.Equal(t, want, buf.String())
}

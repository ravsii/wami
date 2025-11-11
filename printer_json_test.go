package wami_test

import (
	"bytes"
	"testing"

	"github.com/ravsii/wami"
	"github.com/stretchr/testify/assert"
)

func TestJsonPrinter_Print(t *testing.T) {
	p := wami.JsonPrinter{}
	in := []wami.OutputImports{
		{
			Path:  "fmt",
			Count: 2,
			Aliases: []wami.Alias{
				{Count: 3, Name: "f"},
			},
		},
	}
	want := `
	[{
		"path": "fmt",
		"count": 2,
		"aliases": [{
			"name": "f",
			"count": 3
		}]
	}]
	`

	var buf bytes.Buffer
	if err := p.Print(&buf, in); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.JSONEq(t, want, buf.String())
}

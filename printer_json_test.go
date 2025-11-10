package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

// TODO: This is an example test just to generate codecov badge

func TestJsonPrinter_Print(t *testing.T) {
	p := JsonPrinter{}
	data := []OutputImports{
		{
			Path:  "fmt",
			Count: 2,
			Aliases: []Alias{
				{Count: 2, Name: "f"},
			},
		},
	}

	var buf bytes.Buffer
	if err := p.Print(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []OutputImports
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(got) != 1 || got[0].Path != "fmt" || got[0].Count != 2 || got[0].Aliases[0].Name != "f" {
		t.Errorf("unexpected output: %+v", got)
	}
}

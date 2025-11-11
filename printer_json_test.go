package wami_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ravsii/wami"
)

func TestJsonPrinter_Print(t *testing.T) {
	p := wami.JsonPrinter{}
	data := []wami.OutputImports{
		{
			Path:  "fmt",
			Count: 2,
			Aliases: []wami.Alias{
				{Count: 2, Name: "f"},
			},
		},
	}

	var buf bytes.Buffer
	if err := p.Print(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got []wami.OutputImports
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if len(got) != 1 || got[0].Path != "fmt" || got[0].Count != 2 || got[0].Aliases[0].Name != "f" {
		t.Errorf("unexpected output: %+v", got)
	}
}

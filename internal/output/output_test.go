package output_test

import (
	"testing"

	"github.com/shhac/agent-dd/internal/output"
)

func TestResolveFormatDefault(t *testing.T) {
	tests := []struct {
		flag     string
		def      output.Format
		expected output.Format
	}{
		{"", output.FormatJSON, output.FormatJSON},
		{"", output.FormatNDJSON, output.FormatNDJSON},
		{"", output.FormatYAML, output.FormatYAML},
		{"json", output.FormatNDJSON, output.FormatJSON},
		{"jsonl", output.FormatJSON, output.FormatNDJSON},
		{"ndjson", output.FormatJSON, output.FormatNDJSON},
		{"yaml", output.FormatJSON, output.FormatYAML},
		{"bogus", output.FormatJSON, output.FormatJSON},
		{"bogus", output.FormatNDJSON, output.FormatNDJSON},
	}
	for _, tt := range tests {
		got := output.ResolveFormat(tt.flag, tt.def)
		if got != tt.expected {
			t.Errorf("ResolveFormat(%q, %q) = %q, want %q", tt.flag, tt.def, got, tt.expected)
		}
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input   string
		want    output.Format
		wantErr bool
	}{
		{"json", output.FormatJSON, false},
		{"yaml", output.FormatYAML, false},
		{"jsonl", output.FormatNDJSON, false},
		{"ndjson", output.FormatNDJSON, false},
		{"xml", "", true},
	}
	for _, tt := range tests {
		got, err := output.ParseFormat(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseFormat(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if got != tt.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

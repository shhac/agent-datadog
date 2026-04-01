package shared_test

import (
	"testing"
	"time"

	"github.com/shhac/agent-dd/internal/cli/shared"
)

func TestParseTimeRelative(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		checkFn  func(time.Time) bool
	}{
		{"now", false, func(t time.Time) bool { return time.Since(t) < 2*time.Second }},
		{"now-15m", false, func(t time.Time) bool {
			expected := time.Now().Add(-15 * time.Minute)
			return t.Sub(expected).Abs() < 2*time.Second
		}},
		{"now-1h", false, func(t time.Time) bool {
			expected := time.Now().Add(-1 * time.Hour)
			return t.Sub(expected).Abs() < 2*time.Second
		}},
		{"now-7d", false, func(t time.Time) bool {
			expected := time.Now().Add(-7 * 24 * time.Hour)
			return t.Sub(expected).Abs() < 2*time.Second
		}},
		{"now+1h", false, func(t time.Time) bool {
			expected := time.Now().Add(1 * time.Hour)
			return t.Sub(expected).Abs() < 2*time.Second
		}},
	}

	for _, tt := range tests {
		result, err := shared.ParseTime(tt.input)
		if tt.wantErr && err == nil {
			t.Errorf("ParseTime(%q): expected error", tt.input)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("ParseTime(%q): unexpected error: %v", tt.input, err)
		}
		if err == nil && tt.checkFn != nil && !tt.checkFn(result) {
			t.Errorf("ParseTime(%q): time %v didn't pass check", tt.input, result)
		}
	}
}

func TestParseTimeRFC3339(t *testing.T) {
	result, err := shared.ParseTime("2024-01-15T10:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !result.Equal(expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestParseTimeUnixEpoch(t *testing.T) {
	result, err := shared.ParseTime("1705312800")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Unix(1705312800, 0)
	if !result.Equal(expected) {
		t.Errorf("got %v, want %v", result, expected)
	}
}

func TestParseTimeInvalid(t *testing.T) {
	_, err := shared.ParseTime("yesterday")
	if err == nil {
		t.Error("expected error for invalid time string")
	}
}

func TestParseTimeEmpty(t *testing.T) {
	result, err := shared.ParseTime("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsZero() {
		t.Errorf("expected zero time for empty string, got %v", result)
	}
}

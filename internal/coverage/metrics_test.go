package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseMetricsReadsC0AndC1Coverage(t *testing.T) {
	path := writeJSON(t, `{"c0Coverage": 82.1, "c1Coverage": 68.4}`)

	metrics, err := ParseMetrics(path)
	if err != nil {
		t.Fatalf("ParseMetrics() error = %v", err)
	}
	if metrics.C0Coverage == nil {
		t.Fatal("C0Coverage = nil, want value")
	}
	if *metrics.C0Coverage != 82.1 {
		t.Fatalf("C0Coverage = %f, want 82.1", *metrics.C0Coverage)
	}
	if metrics.C1Coverage == nil {
		t.Fatal("C1Coverage = nil, want value")
	}
	if *metrics.C1Coverage != 68.4 {
		t.Fatalf("C1Coverage = %f, want 68.4", *metrics.C1Coverage)
	}
}

func TestParseMetricsRejectsOutOfRangeCoverage(t *testing.T) {
	path := writeJSON(t, `{"c0Coverage": 101}`)

	_, err := ParseMetrics(path)
	if err == nil {
		t.Fatal("ParseMetrics() error = nil, want error")
	}
}

func writeJSON(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "coverage-metrics.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}
	return path
}

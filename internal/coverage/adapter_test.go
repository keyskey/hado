package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAdapterInputReadsGobceJSON(t *testing.T) {
	path := writeAdapterFile(t, "gobce.json", `{
  "language": "go",
  "statementCoverage": 82.1,
  "estimatedBranchCoverage": 68.4,
  "uncoveredBranches": []
}`)

	metrics, err := ParseAdapterInput(AdapterInput{Format: FormatGobceJSON, Path: path})
	if err != nil {
		t.Fatalf("ParseAdapterInput() error = %v", err)
	}
	if metrics.C0Coverage == nil || *metrics.C0Coverage != 82.1 {
		t.Fatalf("C0Coverage = %v, want 82.1", metrics.C0Coverage)
	}
	if metrics.C1Coverage == nil || *metrics.C1Coverage != 68.4 {
		t.Fatalf("C1Coverage = %v, want 68.4", metrics.C1Coverage)
	}
}

func TestParseAdapterInputReadsGoCoverprofile(t *testing.T) {
	path := writeAdapterFile(t, "coverage.out", `mode: set
example.go:1.1,2.2 3 1
example.go:3.1,4.2 2 0
`)

	metrics, err := ParseAdapterInput(AdapterInput{Format: FormatGoCoverprofile, Path: path})
	if err != nil {
		t.Fatalf("ParseAdapterInput() error = %v", err)
	}
	if metrics.C0Coverage == nil || *metrics.C0Coverage != 60 {
		t.Fatalf("C0Coverage = %v, want 60", metrics.C0Coverage)
	}
	if metrics.C1Coverage != nil {
		t.Fatalf("C1Coverage = %v, want nil", metrics.C1Coverage)
	}
}

func TestParseAdapterInputRejectsUnsupportedAdapter(t *testing.T) {
	_, err := ParseAdapterInput(AdapterInput{Format: "lcov", Path: "coverage.info"})
	if err == nil {
		t.Fatal("ParseAdapterInput() error = nil, want unsupported adapter error")
	}
}

func writeAdapterFile(t *testing.T, name, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

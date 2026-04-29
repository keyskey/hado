package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseGoProfileCalculatesC0Coverage(t *testing.T) {
	path := writeProfile(t, `mode: set
example.go:1.1,2.2 3 1
example.go:3.1,4.2 2 0
example.go:5.1,6.2 5 2
`)

	summary, err := ParseGoProfile(path)
	if err != nil {
		t.Fatalf("ParseGoProfile() error = %v", err)
	}
	if summary.CoveredStatements != 8 {
		t.Fatalf("CoveredStatements = %d, want 8", summary.CoveredStatements)
	}
	if summary.TotalStatements != 10 {
		t.Fatalf("TotalStatements = %d, want 10", summary.TotalStatements)
	}
	if summary.C0Coverage != 80 {
		t.Fatalf("C0Coverage = %f, want 80", summary.C0Coverage)
	}
}

func TestParseGoProfileRejectsEmptyProfile(t *testing.T) {
	path := writeProfile(t, "mode: set\n")

	_, err := ParseGoProfile(path)
	if err == nil {
		t.Fatal("ParseGoProfile() error = nil, want error")
	}
}

func writeProfile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "coverage.out")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write coverage profile: %v", err)
	}
	return path
}

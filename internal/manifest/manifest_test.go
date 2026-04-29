package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/keyskey/hado/internal/coverage"
)

func TestLoadReturnsCoverageAdapterInputs(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: coverage-metrics.json
      - adapter: gobce-json
        path: /tmp/gobce.json
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	hadoManifest, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	inputs := hadoManifest.CoverageAdapterInputs()

	want := []coverage.AdapterInput{
		{Format: coverage.FormatHADOJSON, Path: filepath.Join(dir, "coverage-metrics.json")},
		{Format: coverage.FormatGobceJSON, Path: "/tmp/gobce.json"},
	}
	if len(inputs) != len(want) {
		t.Fatalf("len(inputs) = %d, want %d", len(inputs), len(want))
	}
	for i := range want {
		if inputs[i] != want[i] {
			t.Fatalf("inputs[%d] = %+v, want %+v", i, inputs[i], want[i])
		}
	}
}

func TestLoadRejectsCoverageInputWithoutAdapter(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`evidence:
  coverage:
    inputs:
      - path: coverage-metrics.json
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	_, err := Load(manifestPath)
	if err == nil {
		t.Fatal("Load() error = nil, want validation error")
	}
}

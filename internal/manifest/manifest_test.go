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

func TestLoadReturnsOperationsEvidence(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
evidence:
  operations:
    owner: platform-team
    runbook: https://example.com/runbooks/order-api
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	hadoManifest, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if hadoManifest.Evidence.Operations.Owner != "platform-team" {
		t.Fatalf("operations owner = %q, want platform-team", hadoManifest.Evidence.Operations.Owner)
	}
	if hadoManifest.Evidence.Operations.Runbook != "https://example.com/runbooks/order-api" {
		t.Fatalf("operations runbook = %q, want runbook URL", hadoManifest.Evidence.Operations.Runbook)
	}
}

func TestLoadProjectManifest(t *testing.T) {
	hadoManifest, err := Load(filepath.Join("..", "..", "hado.yaml"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	inputs := hadoManifest.CoverageAdapterInputs()
	if len(inputs) != 1 {
		t.Fatalf("len(inputs) = %d, want 1", len(inputs))
	}
	if inputs[0].Format != coverage.FormatGobceJSON {
		t.Fatalf("coverage input format = %q, want %q", inputs[0].Format, coverage.FormatGobceJSON)
	}
	if filepath.Base(inputs[0].Path) != "hado-coverage.json" {
		t.Fatalf("coverage input path = %q, want hado-coverage.json", inputs[0].Path)
	}
	if hadoManifest.Evidence.Operations.Owner != "keyskey" {
		t.Fatalf("operations owner = %q, want keyskey", hadoManifest.Evidence.Operations.Owner)
	}
	if hadoManifest.Evidence.Operations.Runbook != "" {
		t.Fatalf("operations runbook = %q, want empty", hadoManifest.Evidence.Operations.Runbook)
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

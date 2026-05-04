package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsObservabilityEvidence(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
evidence:
  observability:
    slo: slo.yaml
    monitors: monitors.yaml
    dashboard: dash.json
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	hadoManifest, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if hadoManifest.Evidence.Observability == nil {
		t.Fatal("observability evidence is nil")
	}
	if got := hadoManifest.Evidence.Observability.SLO; got != "slo.yaml" {
		t.Fatalf("observability.slo = %q", got)
	}
	if got := hadoManifest.Evidence.Observability.Monitors; got != "monitors.yaml" {
		t.Fatalf("observability.monitors = %q", got)
	}
	if got := hadoManifest.Evidence.Observability.Dashboard; got != "dash.json" {
		t.Fatalf("observability.dashboard = %q", got)
	}
}

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
    slos:
      - name: api availability
        url: https://app.datadoghq.com/slo?slo_id=a
    monitors:
      - name: latency
        url: https://app.datadoghq.com/monitors/1
    dashboards:
      - name: perf
        url: https://app.datadoghq.com/dashboard/bbb
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
	o := hadoManifest.Evidence.Observability
	if len(o.SLOs) != 1 || o.SLOs[0].Name != "api availability" || o.SLOs[0].URL != "https://app.datadoghq.com/slo?slo_id=a" {
		t.Fatalf("slos = %+v", o.SLOs)
	}
	if len(o.Monitors) != 1 || o.Monitors[0].URL != "https://app.datadoghq.com/monitors/1" {
		t.Fatalf("monitors = %+v", o.Monitors)
	}
	if len(o.Dashboards) != 1 || o.Dashboards[0].URL != "https://app.datadoghq.com/dashboard/bbb" {
		t.Fatalf("dashboards = %+v", o.Dashboards)
	}
}

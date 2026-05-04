package evaluate

import (
	"bytes"
	"strings"
	"testing"
)

func TestEvaluateReadyWithObservabilityEvidence(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: observability.slo_exists
    required: true
  - id: observability.monitor_exists
    required: true
  - id: observability.dashboard_exists
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  observability:
    slo: slo.yaml
    monitors: monitors.tf
    dashboard: https://example.com/board/1
`)

	var stdout, stderr bytes.Buffer
	exitCode, err := Run([]string{
		"--standard", standardPath,
		"--manifest", manifestPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "observability.dashboard_exists") {
		t.Fatalf("stdout = %q, want observability.dashboard_exists in output", stdout.String())
	}
}

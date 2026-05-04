package fire

import (
	"bytes"
	"strings"
	"testing"
)

func TestFireReadyWithManifestOperationsEvidence(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: operations.owner_exists
    severity: critical
    required: true
  - id: operations.runbook_exists
    severity: critical
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  operations:
    owner: platform-team
    runbook: https://example.com/runbooks/order-api
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{
		"--standard", standardPath,
		"--manifest", manifestPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run fire: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "operations.runbook_exists") {
		t.Fatalf("stdout = %q, want operations runbook result", stdout.String())
	}
}

func TestFireBlocksWhenOperationsEvidenceIsMissing(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: operations.owner_exists
    severity: critical
    required: true
  - id: operations.runbook_exists
    severity: critical
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  operations:
    owner: platform-team
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{
		"--standard", standardPath,
		"--manifest", manifestPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run fire: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stdout.String(), "Operations runbook is not defined.") {
		t.Fatalf("stdout = %q, want missing runbook message", stdout.String())
	}
}

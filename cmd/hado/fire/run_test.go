package fire

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestFireUsesManifestStandardWhenFlagNotProvided(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatalf("mkdir standards: %v", err)
	}
	standardPath := writeFile(t, stdDir, "standard.yaml", `id: test
gates:
  - id: operations.owner_exists
    severity: critical
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
standard:
  id: standard.yaml
evidence:
  operations:
    owner: platform
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{"--manifest", manifestPath}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run fire: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0 (standard %s)", exitCode, standardPath)
	}
}

func TestFireStandardFlagOverridesManifestStandard(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatalf("mkdir standards: %v", err)
	}
	writeFile(t, stdDir, "default.yaml", `id: default
gates:
  - id: operations.owner_exists
    severity: critical
    required: true
`)
	overridePath := writeFile(t, dir, "override.yaml", `id: override
gates:
  - id: operations.runbook_exists
    severity: critical
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
standard:
  id: default.yaml
evidence:
  operations:
    owner: platform
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{
		"--manifest", manifestPath,
		"--standard", overridePath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run fire: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stdout.String(), "operations.runbook_exists") {
		t.Fatalf("stdout = %q, want override gate output", stdout.String())
	}
}

func TestFireErrorsWithoutResolvableStandard(t *testing.T) {
	dir := t.TempDir()
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence: {}
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{"--manifest", manifestPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run fire error = nil, want error")
	}
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
}

func TestFireRequiresCoverageEvidenceWhenCoverageGateExists(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatalf("mkdir standards: %v", err)
	}
	writeFile(t, stdDir, "standard.yaml", `id: test
gates:
  - id: test.c1_coverage
    severity: critical
    required: true
    threshold:
      min: 70
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
standard:
  id: standard.yaml
evidence: {}
`)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{"--manifest", manifestPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run fire error = nil, want error")
	}
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
}

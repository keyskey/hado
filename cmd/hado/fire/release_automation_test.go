package fire

import (
	"bytes"
	"strings"
	"testing"
)

func TestFireReadyWithReleaseAutomationDeclared(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: release.automation_declared
    severity: critical
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  release:
    automation:
      workflow_refs:
        - .github/workflows/release.yml
      systems:
        - github_actions
`)

	var stdout, stderr bytes.Buffer
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
	if !strings.Contains(stdout.String(), "release.automation_declared") {
		t.Fatalf("stdout = %q, want automation gate in output", stdout.String())
	}
}

func TestFireBlocksWhenReleaseAutomationMissing(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: release.automation_declared
    severity: critical
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  release:
    rollback_plan: docs/rollback.md
`)

	var stdout, stderr bytes.Buffer
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
	if !strings.Contains(stdout.String(), "Release automation workflows are not declared") {
		t.Fatalf("stdout = %q, want automation missing message", stdout.String())
	}
}

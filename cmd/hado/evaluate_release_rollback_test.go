package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestEvaluateReadyWithReleaseRollbackPlan(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: release.rollback_plan_exists
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  release:
    rollback_plan: docs/rollback.md
`)

	var stdout, stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--manifest", manifestPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "release.rollback_plan_exists") {
		t.Fatalf("stdout = %q, want release.rollback_plan_exists in output", stdout.String())
	}
}

func TestEvaluateBlocksWhenRollbackPlanMissing(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: release.rollback_plan_exists
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence: {}
`)

	var stdout, stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--manifest", manifestPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stdout.String(), "Rollback plan is not defined.") {
		t.Fatalf("stdout = %q, want rollback message", stdout.String())
	}
}

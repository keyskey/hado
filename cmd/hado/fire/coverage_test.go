package fire

import (
	"bytes"
	"strings"
	"testing"
)

func TestFireReadyWithNormalizedCoverageMetrics(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
  - id: test.c1_coverage
    severity: critical
    required: true
    threshold:
      min: 65
`)
	writeFile(t, dir, "coverage-metrics.json", `{
  "c0Coverage": 70,
  "c1Coverage": 68.4
}`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: coverage-metrics.json
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
	if !strings.Contains(stdout.String(), "test.c1_coverage") {
		t.Fatalf("stdout = %q, want C1 coverage result", stdout.String())
	}
}

func TestFireBlockedWithGobceAdapter(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    severity: critical
    required: true
    threshold:
      min: 70
  - id: test.c1_coverage
    severity: critical
    required: true
    threshold:
      min: 70
`)
	writeFile(t, dir, "gobce.json", `{
  "language": "go",
  "statementCoverage": 70,
  "estimatedBranchCoverage": 68.4,
  "uncoveredBranches": []
}`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  coverage:
    inputs:
      - adapter: gobce-json
        path: gobce.json
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
	if !strings.Contains(stdout.String(), "C1 coverage is 68.4%, below the required 70.0%") {
		t.Fatalf("stdout = %q, want C1 coverage failure", stdout.String())
	}
}

func TestFireReadyWithGoCoverprofileAdapter(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    severity: critical
    required: true
    threshold:
      min: 70
`)
	writeFile(t, dir, "coverage.out", `mode: set
example.go:1.1,2.1 7 1
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  coverage:
    inputs:
      - adapter: go-coverprofile
        path: coverage.out
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
}

func TestFireBlocked(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    severity: critical
    required: true
    threshold:
      min: 80
`)
	coverprofilePath := writeFile(t, dir, "coverage.out", `mode: set
example.go:1.1,2.1 7 1
example.go:3.1,4.1 3 0
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  coverage:
    inputs:
      - adapter: go-coverprofile
        path: coverage.out
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
	if !strings.Contains(stdout.String(), "HADO: BLOCKED") {
		t.Fatalf("stdout = %q, want blocked status", stdout.String())
	}
	if coverprofilePath == "" {
		t.Fatal("coverprofilePath should not be empty")
	}
}

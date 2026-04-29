package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestEvaluateReady(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
`)
	metricsPath := writeFile(t, dir, "coverage-metrics.json", `{
  "c0Coverage": 70
}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--coverage-input", "hado-json:" + metricsPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "HADO: READY") {
		t.Fatalf("stdout = %q, want ready status", stdout.String())
	}
}

func TestEvaluateReadyWithNormalizedCoverageMetrics(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
  - id: test.c1_coverage
    required: true
    threshold:
      min: 65
`)
	metricsPath := writeFile(t, dir, "coverage-metrics.json", `{
  "c0Coverage": 70,
  "c1Coverage": 68.4
}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--coverage-input", "hado-json:" + metricsPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "test.c1_coverage") {
		t.Fatalf("stdout = %q, want C1 coverage result", stdout.String())
	}
}

func TestEvaluateBlockedWithGobceAdapter(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
  - id: test.c1_coverage
    required: true
    threshold:
      min: 70
`)
	gobcePath := writeFile(t, dir, "gobce.json", `{
  "language": "go",
  "statementCoverage": 70,
  "estimatedBranchCoverage": 68.4,
  "uncoveredBranches": []
}`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--coverage-input", "gobce-json:" + gobcePath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stdout.String(), "C1 coverage is 68.4%, below the required 70.0%") {
		t.Fatalf("stdout = %q, want C1 coverage failure", stdout.String())
	}
}

func TestEvaluateReadyWithGoCoverprofileAdapter(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
`)
	coverprofilePath := writeFile(t, dir, "coverage.out", `mode: set
example.go:1.1,2.1 7 1
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--coverage-input", "go-coverprofile:" + coverprofilePath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
}

func TestEvaluateReadyWithManifestCoverageInputs(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
  - id: test.c1_coverage
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
	if !strings.Contains(stdout.String(), "HADO: READY") {
		t.Fatalf("stdout = %q, want ready status", stdout.String())
	}
}

func TestEvaluateCoverageInputOverridesManifest(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 70
`)
	writeFile(t, dir, "manifest-coverage.json", `{
  "c0Coverage": 10
}`)
	cliMetricsPath := writeFile(t, dir, "cli-coverage.json", `{
  "c0Coverage": 70
}`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: manifest-coverage.json
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--manifest", manifestPath,
		"--coverage-input", "hado-json:" + cliMetricsPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
}

func TestEvaluateRequiresCoverageEvidence(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c1_coverage
    required: true
    threshold:
      min: 70
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{"evaluate", "--standard", standardPath}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run evaluate error = nil, want missing coverage evidence error")
	}
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
}

func TestEvaluateBlocked(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    required: true
    threshold:
      min: 80
`)
	coverprofilePath := writeFile(t, dir, "coverage.out", `mode: set
example.go:1.1,2.1 7 1
example.go:3.1,4.1 3 0
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--coverage-input", "go-coverprofile:" + coverprofilePath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 1 {
		t.Fatalf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stdout.String(), "HADO: BLOCKED") {
		t.Fatalf("stdout = %q, want blocked status", stdout.String())
	}
}

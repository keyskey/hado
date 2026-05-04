package charge

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/keyskey/hado/internal/manifest"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestChargeMergesCoverageInputsWithoutReplacingExisting(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatalf("mkdir standards: %v", err)
	}
	writeFile(t, stdDir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
    severity: critical
    required: true
    threshold:
      min: 70
`)
	writeFile(t, dir, "existing.json", `{"c0Coverage":71}`)
	writeFile(t, dir, "new.json", `{"c0Coverage":72}`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
standard:
  id: standard.yaml
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: existing.json
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{
		"--manifest", manifestPath,
		"--coverage-input", "hado-json:new.json",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run charge: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}

	m, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	if m.Evidence.Coverage == nil || len(m.Evidence.Coverage.Inputs) != 2 {
		t.Fatalf("coverage inputs = %#v, want 2 inputs", m.Evidence.Coverage)
	}
}

func TestChargeRequiresManifest(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run charge error = nil, want error")
	}
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
}

func TestChargeUsesManifestStandardWhenFlagNotProvided(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatalf("mkdir standards: %v", err)
	}
	writeFile(t, stdDir, "standard.yaml", `id: test
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
    owner: team
`)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode, err := Run([]string{"--manifest", manifestPath}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run charge: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "Wrote manifest") {
		t.Fatalf("stdout = %q, want write message", stdout.String())
	}
}

func TestChargeFailsWhenMergedCoverageInputIsInvalid(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatalf("mkdir standards: %v", err)
	}
	writeFile(t, stdDir, "standard.yaml", `id: test
gates:
  - id: test.c0_coverage
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
	exitCode, err := Run([]string{
		"--manifest", manifestPath,
		"--coverage-input", "gobce-json:not-found.json",
	}, &stdout, &stderr)
	if err == nil {
		t.Fatal("run charge error = nil, want parse/read error")
	}
	if exitCode != 2 {
		t.Fatalf("exit code = %d, want 2", exitCode)
	}
}

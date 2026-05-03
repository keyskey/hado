package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunChargeDryRunGapUnsatisfied(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stdPath := filepath.Join(stdDir, "web-service.yaml")
	if err := os.WriteFile(stdPath, []byte(`id: web-service
name: Web
gates:
  - id: test.c0_coverage
    severity: major
    required: true
    threshold:
      min: 80
`), 0o600); err != nil {
		t.Fatal(err)
	}
	mPath := filepath.Join(dir, "hado.yaml")
	if err := os.WriteFile(mPath, []byte(`version: v1
standard:
  id: web-service
`), 0o600); err != nil {
		t.Fatal(err)
	}
	out := &strings.Builder{}
	errOut := &strings.Builder{}
	code, err := runCharge([]string{
		"--manifest", mPath,
		"--standards-dir", stdDir,
		"--output", "text",
	}, out, errOut)
	if err != nil {
		t.Fatalf("runCharge err: %v stderr: %s", err, errOut.String())
	}
	if code != 1 {
		t.Fatalf("want exit 1 when gap not satisfied, got %d out=%s", code, out.String())
	}
}

func TestRunChargeDryRunGapSatisfied(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stdPath := filepath.Join(stdDir, "web-service.yaml")
	if err := os.WriteFile(stdPath, []byte(`id: web-service
name: Web
gates:
  - id: test.c0_coverage
    severity: major
    required: true
    threshold:
      min: 80
`), 0o600); err != nil {
		t.Fatal(err)
	}
	metricsPath := filepath.Join(dir, "m.json")
	if err := os.WriteFile(metricsPath, []byte(`{"c0Coverage": 90}`), 0o600); err != nil {
		t.Fatal(err)
	}
	mPath := filepath.Join(dir, "hado.yaml")
	if err := os.WriteFile(mPath, []byte(`version: v1
standard:
  id: web-service
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: m.json
`), 0o600); err != nil {
		t.Fatal(err)
	}
	out := &strings.Builder{}
	errOut := &strings.Builder{}
	code, err := runCharge([]string{
		"--manifest", mPath,
		"--standards-dir", stdDir,
		"--output", "text",
	}, out, errOut)
	if err != nil {
		t.Fatalf("runCharge err: %v", err)
	}
	if code != 0 {
		t.Fatalf("want exit 0 when satisfied, got %d", code)
	}
}

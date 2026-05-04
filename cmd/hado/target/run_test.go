package target

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/keyskey/hado/internal/manifest"
)

func TestRunTargetNonInteractiveWritesManifest(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "hado.yaml")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	code, err := Run([]string{
		"--manifest", manifestPath,
		"--service-name", "order-api",
		"--service-id", "order-api",
		"--standard-id", "web-service",
		"--rewrite-placeholders=false",
	}, strings.NewReader(""), stdout, stderr)
	if err != nil {
		t.Fatalf("runTarget() err = %v", err)
	}
	if code != 0 {
		t.Fatalf("runTarget() code = %d, stderr = %s", code, stderr.String())
	}

	m, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if m.Service.Name != "order-api" || m.Service.ID != "order-api" {
		t.Fatalf("service = %+v", m.Service)
	}
	if m.Standard.ID != "web-service" {
		t.Fatalf("standard id = %q", m.Standard.ID)
	}
	if m.Version != "v1" {
		t.Fatalf("version = %q", m.Version)
	}
}

func TestRunTargetNonInteractiveMergesExisting(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
service:
  name: keep-name
  id: keep-id
standard:
  id: old-standard
evidence:
  operations:
    owner: team
`), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	code, err := Run([]string{
		"--manifest", manifestPath,
		"--standard-id", "critical-api",
		"--rewrite-placeholders=false",
	}, strings.NewReader(""), stdout, stderr)
	if err != nil {
		t.Fatalf("runTarget() err = %v", err)
	}
	if code != 0 {
		t.Fatalf("runTarget() code = %d, stderr = %s", code, stderr.String())
	}

	m, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if m.Service.Name != "keep-name" || m.Service.ID != "keep-id" {
		t.Fatalf("service = %+v", m.Service)
	}
	if m.Standard.ID != "critical-api" {
		t.Fatalf("standard id = %q", m.Standard.ID)
	}
	if m.Evidence.Operations == nil || m.Evidence.Operations.Owner != "team" {
		t.Fatalf("lost evidence: operations = %+v", m.Evidence.Operations)
	}
}

func TestRunTargetNonInteractiveRequiresField(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "hado.yaml")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	_, err := Run([]string{"--manifest", manifestPath}, strings.NewReader(""), stdout, stderr)
	if err == nil {
		t.Fatal("want error when no service or standard")
	}
}

func TestRunTargetWritesEvidenceScaffold(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	stdYAML := `id: web-service
gates:
  - id: test.c0_coverage
    severity: major
    required: true
    threshold:
      min: 80
  - id: operations.owner_exists
    severity: major
    required: true
`
	if err := os.WriteFile(filepath.Join(stdDir, "web-service.yaml"), []byte(stdYAML), 0o600); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "hado.yaml")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	code, err := Run([]string{
		"--manifest", manifestPath,
		"--service-name", "svc",
		"--standard-id", "web-service",
		"--standards-dir", stdDir,
	}, strings.NewReader(""), stdout, stderr)
	if err != nil {
		t.Fatalf("runTarget: %v", err)
	}
	if code != 0 {
		t.Fatalf("code %d stderr %s", code, stderr.String())
	}
	m, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if m.Evidence.Coverage == nil || len(m.Evidence.Coverage.Inputs) == 0 {
		t.Fatal("expected coverage scaffold input")
	}
	if m.Evidence.Operations == nil || m.Evidence.Operations.Owner != "" {
		t.Fatalf("owner = %+v", m.Evidence.Operations)
	}
}

func TestRunTargetSkipsEvidenceScaffoldWhenDisabled(t *testing.T) {
	dir := t.TempDir()
	stdDir := filepath.Join(dir, "standards")
	if err := os.MkdirAll(stdDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(stdDir, "web-service.yaml"), []byte(`id: web-service
gates:
  - id: operations.owner_exists
    severity: major
    required: true
`), 0o600); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(dir, "hado.yaml")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	_, err := Run([]string{
		"--manifest", manifestPath,
		"--service-name", "svc",
		"--standard-id", "web-service",
		"--standards-dir", stdDir,
		"--rewrite-placeholders=false",
	}, strings.NewReader(""), stdout, stderr)
	if err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if m.Evidence.Coverage != nil && len(m.Evidence.Coverage.Inputs) != 0 {
		t.Fatalf("expected no coverage inputs, got %+v", m.Evidence.Coverage.Inputs)
	}
}

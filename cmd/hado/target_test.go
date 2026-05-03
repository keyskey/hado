package main

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
	code, err := runTarget([]string{
		"--manifest", manifestPath,
		"--service-name", "order-api",
		"--service-id", "order-api",
		"--standard-id", "web-service",
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
	code, err := runTarget([]string{
		"--manifest", manifestPath,
		"--standard-id", "critical-api",
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
	if m.Evidence.Operations.Owner != "team" {
		t.Fatalf("lost evidence: owner = %q", m.Evidence.Operations.Owner)
	}
}

func TestRunTargetNonInteractiveRequiresField(t *testing.T) {
	dir := t.TempDir()
	manifestPath := filepath.Join(dir, "hado.yaml")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	_, err := runTarget([]string{"--manifest", manifestPath}, strings.NewReader(""), stdout, stderr)
	if err == nil {
		t.Fatal("want error when no service or standard")
	}
}

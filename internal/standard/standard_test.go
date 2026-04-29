package standard

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCoverageStandard(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "standard.yaml")
	if err := os.WriteFile(path, []byte(`id: web-service
name: Web Service
gates:
  - id: test.c0_coverage
    severity: major
    required: true
    threshold:
      min: 80
  - id: test.c1_coverage
    severity: major
    required: true
    threshold:
      min: 70
`), 0o600); err != nil {
		t.Fatalf("write standard: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ID != "web-service" {
		t.Fatalf("ID = %q, want web-service", loaded.ID)
	}
	if got := *loaded.Gates[0].Threshold.Min; got != 80 {
		t.Fatalf("threshold min = %v, want 80", got)
	}
	if !loaded.RequiresGate(C0CoverageGateID) {
		t.Fatal("loaded standard should require C0 coverage gate")
	}
	if !loaded.RequiresGate(C1CoverageGateID) {
		t.Fatal("loaded standard should require C1 coverage gate")
	}
}

func TestLoadRejectsCoverageGateWithoutThreshold(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "standard.yaml")
	if err := os.WriteFile(path, []byte(`id: web-service
gates:
  - id: test.c1_coverage
    required: true
`), 0o600); err != nil {
		t.Fatalf("write standard: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("Load() error = nil, want threshold error")
	}
}

func TestLoadOperationGatesWithoutThreshold(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "standard.yaml")
	if err := os.WriteFile(path, []byte(`id: web-service
gates:
  - id: operations.owner_exists
    required: true
  - id: operations.runbook_exists
    required: true
`), 0o600); err != nil {
		t.Fatalf("write standard: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !loaded.RequiresGate(OperationsOwnerExistsGateID) {
		t.Fatal("loaded standard should require operations owner gate")
	}
	if !loaded.RequiresGate(OperationsRunbookExistsGateID) {
		t.Fatal("loaded standard should require operations runbook gate")
	}
}

func TestLoadCLIServiceStandard(t *testing.T) {
	t.Parallel()

	loaded, err := Load("../../standards/cli-service.yaml")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ID != "cli-service" {
		t.Fatalf("ID = %q, want cli-service", loaded.ID)
	}
	if !loaded.RequiresGate(C0CoverageGateID) {
		t.Fatal("cli-service standard should require C0 coverage gate")
	}
	if !loaded.RequiresGate(C1CoverageGateID) {
		t.Fatal("cli-service standard should require C1 coverage gate")
	}
	if !loaded.RequiresGate(OperationsOwnerExistsGateID) {
		t.Fatal("cli-service standard should require operations owner gate")
	}
	if !loaded.RequiresGate(OperationsRunbookExistsGateID) {
		t.Fatal("cli-service standard should require operations runbook gate")
	}
}

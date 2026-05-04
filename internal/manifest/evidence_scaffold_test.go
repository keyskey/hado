package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestApplyEvidenceScaffoldFromStandard(t *testing.T) {
	dir := t.TempDir()
	stdPath := filepath.Join(dir, "web-service.yaml")
	if err := os.WriteFile(stdPath, []byte(`id: web-service
gates:
  - id: test.c0_coverage
    severity: major
    required: true
    threshold:
      min: 80
  - id: operations.owner_exists
    severity: major
    required: true
`), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := standard.Load(stdPath)
	if err != nil {
		t.Fatal(err)
	}
	var m Manifest
	m.Version = "v1"
	ApplyEvidenceScaffold(&m, st, ApplyEvidenceScaffoldOptions{})
	if m.Evidence.Coverage == nil || len(m.Evidence.Coverage.Inputs) != 1 {
		t.Fatalf("coverage inputs: %+v", m.Evidence.Coverage)
	}
	if m.Evidence.Operations == nil || m.Evidence.Operations.Owner != "" {
		t.Fatalf("owner = %+v", m.Evidence.Operations)
	}
}

func TestApplyEvidenceScaffoldMergeKeepsOwner(t *testing.T) {
	dir := t.TempDir()
	stdPath := filepath.Join(dir, "s.yaml")
	if err := os.WriteFile(stdPath, []byte(`id: s
gates:
  - id: operations.owner_exists
    severity: major
    required: true
`), 0o600); err != nil {
		t.Fatal(err)
	}
	st, err := standard.Load(stdPath)
	if err != nil {
		t.Fatal(err)
	}
	m := Manifest{Version: "v1", Evidence: Evidence{Operations: &OperationsEvidence{Owner: "team-a"}}}
	ApplyEvidenceScaffold(&m, st, ApplyEvidenceScaffoldOptions{MergeOnly: true})
	if m.Evidence.Operations == nil || m.Evidence.Operations.Owner != "team-a" {
		t.Fatalf("owner overwritten: %+v", m.Evidence.Operations)
	}
}

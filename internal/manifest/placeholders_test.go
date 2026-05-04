package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestApplyEvidencePlaceholdersFromStandard(t *testing.T) {
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
	ApplyEvidencePlaceholders(&m, st, ApplyEvidencePlaceholdersOptions{})
	if len(m.Evidence.Coverage.Inputs) != 1 {
		t.Fatalf("coverage inputs: %+v", m.Evidence.Coverage.Inputs)
	}
	if m.Evidence.Operations.Owner != EvidencePlaceholder {
		t.Fatalf("owner = %q", m.Evidence.Operations.Owner)
	}
}

func TestApplyEvidencePlaceholdersMergeKeepsOwner(t *testing.T) {
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
	m := Manifest{Version: "v1", Evidence: Evidence{Operations: OperationsEvidence{Owner: "team-a"}}}
	ApplyEvidencePlaceholders(&m, st, ApplyEvidencePlaceholdersOptions{MergeOnly: true})
	if m.Evidence.Operations.Owner != "team-a" {
		t.Fatalf("owner overwritten: %q", m.Evidence.Operations.Owner)
	}
}

func TestResolveStandardPath(t *testing.T) {
	m := Manifest{Standard: StandardRef{ID: "web-service"}}
	p, err := ResolveStandardPath(m, "/r/hado.yaml", "/r/standards", "")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("/r/standards", "web-service.yaml"); p != want {
		t.Fatalf("got %q", p)
	}
}

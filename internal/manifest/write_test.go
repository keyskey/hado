package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadOrEmptyCreatesEmptyManifest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")
	m, err := LoadOrEmpty(path)
	if err != nil {
		t.Fatalf("LoadOrEmpty: %v", err)
	}
	if m.Version != "v1" {
		t.Fatalf("version = %q", m.Version)
	}
}

func TestSaveRoundTripServiceAndStandard(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hado.yaml")
	m := Manifest{
		Version:  "v1",
		Service:  Service{Name: "svc", ID: "svc"},
		Standard: StandardRef{ID: "web-service"},
		Evidence: Evidence{
			Operations: OperationsEvidence{Owner: "o"},
		},
	}
	if err := m.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Service != m.Service || loaded.Standard != m.Standard {
		t.Fatalf("loaded service/standard = %+v / %+v", loaded.Service, loaded.Standard)
	}
	if loaded.Evidence.Operations.Owner != "o" {
		t.Fatalf("operations owner = %q", loaded.Evidence.Operations.Owner)
	}
}

func TestLoadOrEmptyThenSavePreservesEvidence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hado.yaml")
	m, err := LoadOrEmpty(path)
	if err != nil {
		t.Fatal(err)
	}
	m.Service = Service{Name: "n", ID: "n"}
	m.Standard = StandardRef{ID: "std"}
	if err := m.Save(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "service:") {
		t.Fatalf("expected service block in yaml: %s", string(data))
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Service.Name != "n" || loaded.Service.ID != "n" || loaded.Standard.ID != "std" {
		t.Fatalf("loaded = service %+v standard %+v", loaded.Service, loaded.Standard)
	}
}

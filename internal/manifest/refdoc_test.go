package manifest

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestManifestYAMLDocComplete(t *testing.T) {
	paths, err := manifestYAMLPaths()
	if err != nil {
		t.Fatal(err)
	}
	var fromTypes []string
	for _, p := range paths {
		fromTypes = append(fromTypes, p.path)
		if _, ok := manifestYAMLDoc[p.path]; !ok {
			t.Errorf("manifestYAMLDoc missing description for path %q (add to field_docs.go)", p.path)
		}
	}
	var orphan []string
	for k := range manifestYAMLDoc {
		if !slices.Contains(fromTypes, k) {
			orphan = append(orphan, k)
		}
	}
	if len(orphan) > 0 {
		slices.Sort(orphan)
		t.Errorf("manifestYAMLDoc has keys not produced by types walk (remove or fix types): %s", strings.Join(orphan, ", "))
	}
}

func TestWriteManifestReferenceYAML_loads(t *testing.T) {
	var sb strings.Builder
	if err := WriteManifestReferenceYAML(&sb); err != nil {
		t.Fatal(err)
	}
	data := sb.String()
	if !strings.Contains(data, "version:") || !strings.Contains(data, "evidence:") {
		t.Fatalf("unexpected output: %s", truncate(data, 200))
	}
	m, err := parseManifestBytes([]byte(data), t.TempDir())
	if err != nil {
		t.Fatalf("parseManifestBytes: %v\n---\n%s", err, truncate(data, 800))
	}
	if m.Version != "v1" {
		t.Fatalf("Version = %q want v1", m.Version)
	}
	if m.Evidence.Coverage == nil || len(m.Evidence.Coverage.Inputs) != 1 {
		t.Fatalf("expected one coverage input, got %#v", m.Evidence.Coverage)
	}
}

func TestCommittedReferenceYAMLMatchesGenerator(t *testing.T) {
	root := findModuleRoot(t)
	refPath := filepath.Join(root, "docs", "hado.manifest.reference.yaml")
	committed, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatalf("read %s: %v (run from repo root or ensure file exists)", refPath, err)
	}
	var gen strings.Builder
	if err := WriteManifestReferenceYAML(&gen); err != nil {
		t.Fatal(err)
	}
	if normalizeEOL(string(committed)) != normalizeEOL(gen.String()) {
		t.Fatalf("docs/hado.manifest.reference.yaml is out of sync with WriteManifestReferenceYAML.\n"+
			"Run from repo root: make gen-manifest-doc\npath: %s", refPath)
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for i := 0; i < 12; i++ {
		mod := filepath.Join(dir, "go.mod")
		if st, err := os.Stat(mod); err == nil && !st.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	t.Fatalf("go.mod not found from cwd %q (run go test from module root)", wd)
	return ""
}

func normalizeEOL(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

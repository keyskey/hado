package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveUsesTwoSpaceIndent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hado.yaml")
	m := Manifest{
		Version:  "v1",
		Service:  Service{Name: "svc", ID: "svc"},
		Standard: StandardRef{ID: "web-service"},
	}
	if err := m.Save(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	// Nested keys directly under "service:" should be indented with 2 spaces, not yaml.Marshal's default 4.
	if i := strings.Index(s, "service:\n"); i >= 0 {
		after := s[i+len("service:\n"):]
		first := strings.SplitN(after, "\n", 2)[0]
		if strings.HasPrefix(first, "    ") {
			t.Fatalf("first line under service should not use 4-space indent; got %q\n%s", first, s)
		}
		if first != "" && !strings.HasPrefix(first, "  ") {
			t.Fatalf("first line under service should use 2-space indent; got %q\n%s", first, s)
		}
	}
}

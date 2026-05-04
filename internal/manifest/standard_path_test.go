package manifest

import (
	"path/filepath"
	"testing"
)

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

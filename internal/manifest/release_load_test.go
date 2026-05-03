package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsReleaseRollbackPlan(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
evidence:
  release:
    rollback_plan: docs/rollback.md
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	hadoManifest, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := hadoManifest.Evidence.Release.RollbackPlan; got != "docs/rollback.md" {
		t.Fatalf("release.rollback_plan = %q", got)
	}
}

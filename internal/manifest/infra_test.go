package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsInfraEvidence(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
evidence:
  infra:
    deployment_spec: deploy/
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	hadoManifest, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if got := hadoManifest.Evidence.Infra.DeploymentSpec; got != "deploy/" {
		t.Fatalf("infra.deployment_spec = %q", got)
	}
}

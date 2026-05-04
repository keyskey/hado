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

	if hadoManifest.Evidence.Release == nil {
		t.Fatal("release evidence is nil")
	}
	if got := hadoManifest.Evidence.Release.RollbackPlan; got != "docs/rollback.md" {
		t.Fatalf("release.rollback_plan = %q", got)
	}
}

func TestLoadReturnsReleaseAutomationEvidence(t *testing.T) {
	manifestPath := filepath.Join(t.TempDir(), "hado.yaml")
	if err := os.WriteFile(manifestPath, []byte(`version: v1
evidence:
  release:
    automation:
      workflow_refs:
        - ci/release.yaml
      systems:
        - argo_workflow
`), 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	hadoManifest, err := Load(manifestPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if refs := hadoManifest.Evidence.Release.Automation.WorkflowRefs; len(refs) != 1 || refs[0] != "ci/release.yaml" {
		t.Fatalf("release.automation.workflow_refs = %#v", refs)
	}
	if sys := hadoManifest.Evidence.Release.Automation.Systems; len(sys) != 1 || sys[0] != "argo_workflow" {
		t.Fatalf("release.automation.systems = %#v", sys)
	}
	if !hadoManifest.Evidence.Release.AutomationDeclared() {
		t.Fatal("AutomationDeclared() = false, want true")
	}
}

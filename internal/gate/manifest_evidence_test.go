package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/manifest"
)

func TestManifestEvidenceStringPlaceholder(t *testing.T) {
	if ManifestEvidenceString(manifest.EvidencePlaceholder) != "" {
		t.Fatal("placeholder should be empty for gate purposes")
	}
	if ManifestEvidenceString("  real  ") != "real" {
		t.Fatal("trim real value")
	}
}

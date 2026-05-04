package gate

import (
	"strings"

	"github.com/keyskey/hado/internal/manifest"
)

// ManifestEvidenceString returns trimmed s, or empty if unset or placeholder.
func ManifestEvidenceString(s string) string {
	s = strings.TrimSpace(s)
	if manifest.EvidenceUnset(s) {
		return ""
	}
	return s
}

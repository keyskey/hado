package manifest

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ResolveStandardPath returns the filesystem path to the readiness standard YAML.
// If override is set, it is used. Otherwise m.Standard.ID is used:
// if it looks like a path (contains / or \), it is resolved relative to the manifest directory when not absolute;
// otherwise it is treated as a basename under standardsDir (e.g. "web-service" -> standards/web-service.yaml).
func ResolveStandardPath(m Manifest, manifestPath, standardsDir, override string) (string, error) {
	ref := strings.TrimSpace(override)
	if ref == "" {
		ref = strings.TrimSpace(m.Standard.ID)
	}
	if ref == "" {
		return "", fmt.Errorf("readiness standard path: set manifest standard.id or pass override")
	}
	if filepath.IsAbs(ref) || strings.Contains(ref, "/") || strings.Contains(ref, "\\") {
		if filepath.IsAbs(ref) {
			return filepath.Clean(ref), nil
		}
		base := filepath.Dir(manifestPath)
		return filepath.Clean(filepath.Join(base, ref)), nil
	}
	name := ref
	if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
		name = name + ".yaml"
	}
	return filepath.Join(standardsDir, name), nil
}

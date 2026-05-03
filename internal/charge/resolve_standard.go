package charge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/keyskey/hado/internal/manifest"
)

// ResolveStandardPath returns the filesystem path to the readiness standard YAML.
// If opts.StandardPath is set, it is used. Otherwise manifest.Standard.ID is used:
// if it looks like a path (contains / or \ or ends with .yaml/.yml), it is resolved
// relative to the manifest directory when not absolute; otherwise it is treated as
// a basename under standardsDir (e.g. "web-service" -> standards/web-service.yaml).
func ResolveStandardPath(m manifest.Manifest, manifestPath, standardsDir, override string) (string, error) {
	ref := strings.TrimSpace(override)
	if ref == "" {
		ref = strings.TrimSpace(m.Standard.ID)
	}
	if ref == "" {
		return "", fmt.Errorf("readiness standard path: set manifest standard.id or pass --standard")
	}
	if filepath.IsAbs(ref) || strings.Contains(ref, "/") || strings.Contains(ref, "\\") {
		if filepath.IsAbs(ref) {
			return filepath.Clean(ref), nil
		}
		base := filepath.Dir(manifestPath)
		return filepath.Clean(filepath.Join(base, ref)), nil
	}
	// logical id -> standards/<id>.yaml
	name := ref
	if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
		name = name + ".yaml"
	}
	return filepath.Join(standardsDir, name), nil
}

// FileExists reports whether path exists and is not a directory.
func FileExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}

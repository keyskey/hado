package manifest

import (
	"fmt"
	"os"
	"path/filepath"
)

// Load reads a HADO manifest from disk.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	return parseManifestBytes(data, filepath.Dir(path))
}

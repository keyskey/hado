package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadOrEmpty reads a manifest from path. If the file does not exist, it returns a new
// manifest with version v1 and baseDir set from path. Other read or parse errors are returned.
func LoadOrEmpty(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Manifest{
				Version: "v1",
				baseDir: filepath.Dir(path),
			}, nil
		}
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}
	return parseManifestBytes(data, filepath.Dir(path))
}

// parseManifestBytes unmarshals YAML and validates. baseDir is the directory for relative evidence paths.
func parseManifestBytes(data []byte, baseDir string) (Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse manifest: %w", err)
	}
	m.baseDir = baseDir
	if err := m.Validate(); err != nil {
		return Manifest{}, err
	}
	return m, nil
}

// Save writes the manifest to path as YAML. baseDir is not serialized; it is cleared before marshal
// and restored after so in-memory path resolution keeps working if the same value is reused.
func (m Manifest) Save(path string) error {
	saveCopy := m
	saveCopy.baseDir = ""
	data, err := yaml.Marshal(&saveCopy)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create manifest directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

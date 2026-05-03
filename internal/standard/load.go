package standard

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Load reads a YAML readiness standard from disk.
func Load(path string) (Standard, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Standard{}, fmt.Errorf("read standard: %w", err)
	}

	var standard Standard
	if err := yaml.Unmarshal(data, &standard); err != nil {
		return Standard{}, fmt.Errorf("parse standard: %w", err)
	}
	if err := standard.Validate(); err != nil {
		return Standard{}, err
	}

	return standard, nil
}

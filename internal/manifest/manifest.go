package manifest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/keyskey/hado/internal/coverage"
	"gopkg.in/yaml.v3"
)

// Manifest declares the evaluated service and the evidence HADO should read.
type Manifest struct {
	Version  string   `yaml:"version" json:"version,omitempty"`
	Evidence Evidence `yaml:"evidence" json:"evidence,omitempty"`

	baseDir string
}

// Evidence groups evidence declarations by readiness domain.
type Evidence struct {
	Coverage   CoverageEvidence   `yaml:"coverage" json:"coverage,omitempty"`
	Operations OperationsEvidence `yaml:"operations" json:"operations,omitempty"`
}

// CoverageEvidence declares coverage reports and the adapters that parse them.
type CoverageEvidence struct {
	Inputs []CoverageInput `yaml:"inputs" json:"inputs,omitempty"`
}

// CoverageInput identifies one coverage artifact and its adapter.
type CoverageInput struct {
	Adapter string `yaml:"adapter" json:"adapter"`
	Path    string `yaml:"path" json:"path"`
}

// OperationsEvidence declares operational ownership and response evidence.
type OperationsEvidence struct {
	Owner   string `yaml:"owner" json:"owner,omitempty"`
	Runbook string `yaml:"runbook" json:"runbook,omitempty"`
}

// Load reads a HADO manifest from disk.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read manifest: %w", err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return Manifest{}, fmt.Errorf("parse manifest: %w", err)
	}
	manifest.baseDir = filepath.Dir(path)
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}

// Validate checks manifest fields used by the current evaluator.
func (m Manifest) Validate() error {
	for i, input := range m.Evidence.Coverage.Inputs {
		if input.Adapter == "" {
			return fmt.Errorf("evidence.coverage.inputs[%d].adapter is required", i)
		}
		if input.Path == "" {
			return fmt.Errorf("evidence.coverage.inputs[%d].path is required", i)
		}
	}
	return nil
}

// CoverageAdapterInputs returns coverage adapter inputs with manifest-relative paths resolved.
func (m Manifest) CoverageAdapterInputs() []coverage.AdapterInput {
	inputs := make([]coverage.AdapterInput, 0, len(m.Evidence.Coverage.Inputs))
	for _, input := range m.Evidence.Coverage.Inputs {
		path := input.Path
		if !filepath.IsAbs(path) && m.baseDir != "" {
			path = filepath.Join(m.baseDir, path)
		}
		inputs = append(inputs, coverage.AdapterInput{
			Format: input.Adapter,
			Path:   path,
		})
	}
	return inputs
}

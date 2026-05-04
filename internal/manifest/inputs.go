package manifest

import (
	"path/filepath"

	"github.com/keyskey/hado/internal/coverage"
)

// CoverageAdapterInputs returns coverage adapter inputs with manifest-relative paths resolved.
func (m Manifest) CoverageAdapterInputs() []coverage.AdapterInput {
	if m.Evidence.Coverage == nil {
		return nil
	}
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

package manifest

import "fmt"

// Validate checks manifest fields used by the current evaluator.
func (m Manifest) Validate() error {
	if m.Evidence.Coverage == nil {
		return nil
	}
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

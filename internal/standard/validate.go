package standard

import "fmt"

// Validate checks the fields required by the first HADO gate evaluator.
func (s Standard) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("standard id is required")
	}
	if len(s.Gates) == 0 {
		return fmt.Errorf("standard must define at least one gate")
	}
	for i, gate := range s.Gates {
		if gate.ID == "" {
			return fmt.Errorf("gate %d id is required", i)
		}
		if err := gate.Severity.Validate(); err != nil {
			return fmt.Errorf("gate %d: %w", i, err)
		}
		if isCoverageGate(gate.ID) && gate.Threshold.Min == nil {
			return fmt.Errorf("%s gate requires threshold.min", gate.ID)
		}
	}

	return nil
}

// RequiresGate reports whether the standard declares a gate id.
func (s Standard) RequiresGate(id string) bool {
	for _, gate := range s.Gates {
		if gate.ID == id {
			return true
		}
	}
	return false
}

func isCoverageGate(id string) bool {
	return id == C0CoverageGateID || id == C1CoverageGateID
}

// Validate checks whether severity is one of the supported enum values.
// Empty severity is allowed and treated as default (minor) elsewhere.
func (severity Severity) Validate() error {
	switch severity {
	case "":
		return nil
	case SeverityCritical, SeverityMajor, SeverityMinor:
		return nil
	default:
		return fmt.Errorf("severity %q is unsupported", severity)
	}
}

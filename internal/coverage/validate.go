package coverage

import "fmt"

// Validate checks that all present coverage percentages are in range.
func (metrics Metrics) Validate() error {
	if metrics.C0Coverage != nil && (*metrics.C0Coverage < 0 || *metrics.C0Coverage > 100) {
		return fmt.Errorf("c0Coverage must be between 0 and 100")
	}
	if metrics.C1Coverage != nil && (*metrics.C1Coverage < 0 || *metrics.C1Coverage > 100) {
		return fmt.Errorf("c1Coverage must be between 0 and 100")
	}
	return nil
}

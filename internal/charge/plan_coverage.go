package charge

import (
	"github.com/keyskey/hado/internal/coverage"
	"github.com/keyskey/hado/internal/standard"
)

// PlanCoverage builds a coverage plan from the readiness standard.
func PlanCoverage(standardPath string, st standard.Standard) CoveragePlan {
	p := CoveragePlan{
		StandardPath: standardPath,
		RequiresC0:   st.RequiresGate(standard.C0CoverageGateID),
		RequiresC1:   st.RequiresGate(standard.C1CoverageGateID),
	}
	if p.RequiresC1 {
		// C1 needs a producer that supplies C1; gobce-json and hado-json suffice; go-coverprofile does not.
		p.PreferredAdapterFormats = []string{
			coverage.FormatGobceJSON,
			coverage.FormatHADOJSON,
			coverage.FormatGoCoverprofile, // listed last: alone does not satisfy C1 but we still check file for partial gap messaging
		}
		return p
	}
	if p.RequiresC0 {
		p.PreferredAdapterFormats = []string{
			coverage.FormatGobceJSON,
			coverage.FormatHADOJSON,
			coverage.FormatGoCoverprofile,
		}
		return p
	}
	p.PreferredAdapterFormats = []string{}
	return p
}

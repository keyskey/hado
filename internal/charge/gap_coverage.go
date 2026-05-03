package charge

import (
	"fmt"

	"github.com/keyskey/hado/internal/coverage"
	"github.com/keyskey/hado/internal/manifest"
)

// GapCoverage compares the manifest's coverage inputs against the plan.
func GapCoverage(m manifest.Manifest, plan CoveragePlan) CoverageGapReport {
	if !plan.RequiresC0 && !plan.RequiresC1 {
		return CoverageGapReport{
			Satisfied: true,
			Items: []CoverageGapItem{{
				Code:   "no_coverage_gates",
				Detail: "readiness standard does not declare test.c0_coverage or test.c1_coverage",
			}},
		}
	}

	inputs := m.CoverageAdapterInputs()
	for _, in := range inputs {
		metrics, err := coverage.ParseAdapterInput(in)
		if err != nil {
			continue
		}
		if satisfiesPlan(plan, metrics) {
			return CoverageGapReport{
				Satisfied:         true,
				SatisfyingAdapter: in.Format,
				SatisfyingPath:    in.Path,
				Items: []CoverageGapItem{{
					Code:   "satisfied",
					Detail: fmt.Sprintf("adapter %q at %q supplies required coverage metrics", in.Format, in.Path),
				}},
			}
		}
	}

	report := CoverageGapReport{Satisfied: false, Items: nil}
	if len(inputs) == 0 {
		report.Items = append(report.Items, CoverageGapItem{
			Code:   "missing_input",
			Detail: "evidence.coverage.inputs is empty or absent",
		})
	} else {
		for _, in := range inputs {
			_, err := coverage.ParseAdapterInput(in)
			if err != nil {
				report.Items = append(report.Items, CoverageGapItem{
					Code:   "unreadable_input",
					Detail: fmt.Sprintf("adapter %q path %q: %v", in.Format, in.Path, err),
				})
				continue
			}
			metrics, _ := coverage.ParseAdapterInput(in)
			if plan.RequiresC1 && metrics.C1Coverage == nil {
				report.Items = append(report.Items, CoverageGapItem{
					Code:   "insufficient_adapter",
					Detail: fmt.Sprintf("adapter %q at %q does not supply C1 coverage (test.c1_coverage)", in.Format, in.Path),
				})
			} else if plan.RequiresC0 && metrics.C0Coverage == nil {
				report.Items = append(report.Items, CoverageGapItem{
					Code:   "insufficient_adapter",
					Detail: fmt.Sprintf("adapter %q at %q does not supply C0 coverage", in.Format, in.Path),
				})
			}
		}
	}
	if len(report.Items) == 0 {
		report.Items = []CoverageGapItem{{
			Code:   "missing_satisfying_coverage",
			Detail: "no coverage input satisfies the readiness standard",
		}}
	}
	return report
}

func satisfiesPlan(plan CoveragePlan, m coverage.Metrics) bool {
	if plan.RequiresC0 && m.C0Coverage == nil {
		return false
	}
	if plan.RequiresC1 && m.C1Coverage == nil {
		return false
	}
	return plan.RequiresC0 || plan.RequiresC1
}

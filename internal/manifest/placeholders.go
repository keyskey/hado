package manifest

import (
	"github.com/keyskey/hado/internal/standard"
)

// ApplyEvidencePlaceholdersOptions configures how placeholders are merged into the manifest.
type ApplyEvidencePlaceholdersOptions struct {
	// MergeOnly when true only fills fields that are EvidenceUnset; existing non-placeholder values are kept.
	MergeOnly bool
}

// ApplyEvidencePlaceholders sets evidence fields required by the standard's known gates to
// EvidencePlaceholder or minimal scaffold values so humans or charge can fill them later.
// Only gate IDs understood by the evaluator are mapped; unknown gates are skipped.
func ApplyEvidencePlaceholders(m *Manifest, st standard.Standard, opts ApplyEvidencePlaceholdersOptions) {
	merge := opts.MergeOnly
	needsC0 := st.RequiresGate(standard.C0CoverageGateID)
	needsC1 := st.RequiresGate(standard.C1CoverageGateID)
	if needsC0 || needsC1 {
		if !merge {
			m.Evidence.Coverage.Inputs = []CoverageInput{
				{Adapter: "hado-json", Path: "coverage-metrics.json"},
			}
		} else if len(m.Evidence.Coverage.Inputs) == 0 {
			m.Evidence.Coverage.Inputs = []CoverageInput{
				{Adapter: "hado-json", Path: "coverage-metrics.json"},
			}
		}
	}

	setStr := func(dst *string, gateID string) {
		if !st.RequiresGate(gateID) {
			return
		}
		if merge && !EvidenceUnset(*dst) {
			return
		}
		*dst = EvidencePlaceholder
	}

	setStr(&m.Evidence.Operations.Owner, standard.OperationsOwnerExistsGateID)
	setStr(&m.Evidence.Operations.Runbook, standard.OperationsRunbookExistsGateID)
	setStr(&m.Evidence.Observability.SLO, standard.ObservabilitySLOExistsGateID)
	setStr(&m.Evidence.Observability.Monitors, standard.ObservabilityMonitorExistsGateID)
	setStr(&m.Evidence.Observability.Dashboard, standard.ObservabilityDashboardExistsGateID)
	setStr(&m.Evidence.Infra.DeploymentSpec, standard.InfraDeploymentSpecExistsGateID)
	setStr(&m.Evidence.Release.RollbackPlan, standard.ReleaseRollbackPlanExistsGateID)

	if st.RequiresGate(standard.ReleaseAutomationDeclaredGateID) {
		if !merge || len(m.Evidence.Release.Automation.WorkflowRefs) == 0 {
			if !merge {
				m.Evidence.Release.Automation.WorkflowRefs = nil
			}
			m.Evidence.Release.Automation.WorkflowRefs = append(m.Evidence.Release.Automation.WorkflowRefs, EvidencePlaceholder)
		}
	}
}

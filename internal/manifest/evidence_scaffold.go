package manifest

import (
	"strings"

	"github.com/keyskey/hado/internal/standard"
)

// ApplyEvidenceScaffoldOptions configures how scaffold rows are merged into the manifest.
type ApplyEvidenceScaffoldOptions struct {
	// MergeOnly when true only fills fields that are still empty (after trim); non-empty values are kept.
	MergeOnly bool
}

// ApplyEvidenceScaffold ensures evidence blocks and string fields exist for gates declared in the standard.
// String fields are left empty so YAML lists keys (e.g. owner: "") for humans or charge to fill later.
// Only gate IDs understood by the evaluator are mapped; unknown gates are skipped.
func ApplyEvidenceScaffold(m *Manifest, st standard.Standard, opts ApplyEvidenceScaffoldOptions) {
	merge := opts.MergeOnly
	needsC0 := st.RequiresGate(standard.C0CoverageGateID)
	needsC1 := st.RequiresGate(standard.C1CoverageGateID)
	if needsC0 || needsC1 {
		if !merge {
			m.Evidence.Coverage = &CoverageEvidence{
				Inputs: []CoverageInput{{Adapter: "hado-json", Path: "coverage-metrics.json"}},
			}
		} else {
			if m.Evidence.Coverage == nil {
				m.Evidence.Coverage = &CoverageEvidence{}
			}
			if len(m.Evidence.Coverage.Inputs) == 0 {
				m.Evidence.Coverage.Inputs = []CoverageInput{{Adapter: "hado-json", Path: "coverage-metrics.json"}}
			}
		}
	}

	setOps := func() {
		if m.Evidence.Operations == nil {
			m.Evidence.Operations = &OperationsEvidence{}
		}
	}
	if st.RequiresGate(standard.OperationsOwnerExistsGateID) {
		setOps()
		if !(merge && strings.TrimSpace(m.Evidence.Operations.Owner) != "") {
			m.Evidence.Operations.Owner = ""
		}
	}
	if st.RequiresGate(standard.OperationsRunbookExistsGateID) {
		setOps()
		if !(merge && strings.TrimSpace(m.Evidence.Operations.Runbook) != "") {
			m.Evidence.Operations.Runbook = ""
		}
	}

	setObs := func() {
		if m.Evidence.Observability == nil {
			m.Evidence.Observability = &ObservabilityEvidence{}
		}
	}
	if st.RequiresGate(standard.ObservabilitySLOExistsGateID) {
		setObs()
		if !(merge && strings.TrimSpace(m.Evidence.Observability.SLO) != "") {
			m.Evidence.Observability.SLO = ""
		}
	}
	if st.RequiresGate(standard.ObservabilityMonitorExistsGateID) {
		setObs()
		if !(merge && strings.TrimSpace(m.Evidence.Observability.Monitors) != "") {
			m.Evidence.Observability.Monitors = ""
		}
	}
	if st.RequiresGate(standard.ObservabilityDashboardExistsGateID) {
		setObs()
		if !(merge && strings.TrimSpace(m.Evidence.Observability.Dashboard) != "") {
			m.Evidence.Observability.Dashboard = ""
		}
	}

	if st.RequiresGate(standard.InfraDeploymentSpecExistsGateID) {
		if m.Evidence.Infra == nil {
			m.Evidence.Infra = &InfraEvidence{}
		}
		if !(merge && strings.TrimSpace(m.Evidence.Infra.DeploymentSpec) != "") {
			m.Evidence.Infra.DeploymentSpec = ""
		}
	}

	if st.RequiresGate(standard.ReleaseRollbackPlanExistsGateID) {
		if m.Evidence.Release == nil {
			m.Evidence.Release = &ReleaseEvidence{}
		}
		if !(merge && strings.TrimSpace(m.Evidence.Release.RollbackPlan) != "") {
			m.Evidence.Release.RollbackPlan = ""
		}
	}

	if st.RequiresGate(standard.ReleaseAutomationDeclaredGateID) {
		if m.Evidence.Release == nil {
			m.Evidence.Release = &ReleaseEvidence{}
		}
		if !merge || len(m.Evidence.Release.Automation.WorkflowRefs) == 0 {
			m.Evidence.Release.Automation.WorkflowRefs = []string{""}
		}
	}
}

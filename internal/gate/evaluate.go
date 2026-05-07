package gate

import (
	"fmt"

	"github.com/keyskey/hado/internal/standard"
)

// Evaluate compares collected metrics with the gates in a readiness standard.
func Evaluate(s standard.Standard, metrics Metrics) (Evaluation, error) {
	evaluation := Evaluation{Status: DecisionReady}

	for _, gate := range s.Gates {
		resultCount := len(evaluation.Results)
		switch gate.ID {
		case standard.C0CoverageGateID:
			if metrics.C0CoveragePercent == nil {
				return Evaluation{Status: DecisionError}, fmt.Errorf("%s gate requires c0Coverage evidence", standard.C0CoverageGateID)
			}
			result := evaluateCoverageGate(gate, "C0 coverage", *metrics.C0CoveragePercent)
			evaluation.Results = append(evaluation.Results, result)
		case standard.C1CoverageGateID:
			if metrics.C1CoveragePercent == nil {
				return Evaluation{Status: DecisionError}, fmt.Errorf("%s gate requires c1Coverage evidence", standard.C1CoverageGateID)
			}
			result := evaluateCoverageGate(gate, "C1 coverage", *metrics.C1CoveragePercent)
			evaluation.Results = append(evaluation.Results, result)
		case standard.OperationsOwnerExistsGateID:
			result := evaluateExistsGate(gate, metrics.OperationsOwner != "", "Operations owner is defined.", "Operations owner is not defined.")
			evaluation.Results = append(evaluation.Results, result)
		case standard.OperationsRunbookExistsGateID:
			result := evaluateExistsGate(gate, metrics.OperationsRunbook != "", "Operations runbook is defined.", "Operations runbook is not defined.")
			evaluation.Results = append(evaluation.Results, result)
		case standard.ObservabilitySLOExistsGateID:
			result := evaluateExistsGate(gate, metrics.ObservabilitySLOPresent, "SLO / SLI evidence is defined (at least one URL in evidence.observability.slos).", "SLO / SLI evidence is not defined (add slos[].url).")
			evaluation.Results = append(evaluation.Results, result)
		case standard.ObservabilityMonitorExistsGateID:
			result := evaluateExistsGate(gate, metrics.ObservabilityMonitorsPresent, "Monitor evidence is defined (at least one URL in evidence.observability.monitors).", "Monitor evidence is not defined (add monitors[].url).")
			evaluation.Results = append(evaluation.Results, result)
		case standard.ObservabilityDashboardExistsGateID:
			result := evaluateExistsGate(gate, metrics.ObservabilityDashboardPresent, "Dashboard evidence is defined (at least one URL in evidence.observability.dashboards).", "Dashboard evidence is not defined (add dashboards[].url).")
			evaluation.Results = append(evaluation.Results, result)
		case standard.InfraDeploymentSpecExistsGateID:
			result := evaluateExistsGate(gate, metrics.InfraDeploymentSpec != "", "Deployment spec reference is defined.", "Deployment spec reference is not defined.")
			evaluation.Results = append(evaluation.Results, result)
		case standard.ReleaseRollbackPlanExistsGateID:
			result := evaluateExistsGate(gate, metrics.ReleaseRollbackPlan != "", "Rollback plan is defined.", "Rollback plan is not defined.")
			evaluation.Results = append(evaluation.Results, result)
		case standard.ReleaseAutomationDeclaredGateID:
			result := evaluateExistsGate(gate, metrics.ReleaseAutomationDeclared, "Release automation workflows are declared.", "Release automation workflows are not declared (evidence.release.automation.workflow_refs).")
			evaluation.Results = append(evaluation.Results, result)
		default:
			if gate.Required {
				return Evaluation{Status: DecisionError}, fmt.Errorf("unsupported required gate %q", gate.ID)
			}
		}
		if len(evaluation.Results) > resultCount && shouldBlockRelease(gate, evaluation.Results[len(evaluation.Results)-1]) {
			evaluation.Status = DecisionBlocked
		}
	}

	return evaluation, nil
}

func shouldBlockRelease(gate standard.Gate, result Result) bool {
	if !gate.Required || result.Passed {
		return false
	}
	return effectiveSeverity(gate) == standard.SeverityCritical
}

func evaluateCoverageGate(gate standard.Gate, label string, actual float64) Result {
	requiredMin := *gate.Threshold.Min
	result := Result{
		ID:          gate.ID,
		Required:    gate.Required,
		Severity:    effectiveSeverity(gate),
		Actual:      actual,
		RequiredMin: requiredMin,
	}
	if actual >= requiredMin {
		result.Passed = true
		result.Message = fmt.Sprintf("%s is %.1f%%, meeting the required %.1f%%.", label, actual, requiredMin)
		return result
	}

	result.Message = fmt.Sprintf("%s is %.1f%%, below the required %.1f%%.", label, actual, requiredMin)
	return result
}

func evaluateExistsGate(gate standard.Gate, exists bool, passedMessage, failedMessage string) Result {
	result := Result{
		ID:       gate.ID,
		Required: gate.Required,
		Severity: effectiveSeverity(gate),
		Passed:   exists,
	}
	if exists {
		result.Message = passedMessage
		return result
	}

	result.Message = failedMessage
	return result
}

func effectiveSeverity(gate standard.Gate) standard.Severity {
	if gate.Severity == "" {
		return standard.SeverityMinor
	}
	return gate.Severity
}

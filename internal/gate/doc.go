// Package gate evaluates readiness gates and release decisions.
package gate

import (
	"fmt"

	"github.com/keyskey/hado/internal/standard"
)

// DecisionStatus is the final production readiness decision.
type DecisionStatus string

const (
	// DecisionReady means all required gates passed.
	DecisionReady DecisionStatus = "ready"
	// DecisionBlocked means at least one required gate failed.
	DecisionBlocked DecisionStatus = "blocked"
	// DecisionError means HADO could not complete the evaluation.
	DecisionError DecisionStatus = "error"
)

// Metrics are normalized values supplied by modules or evidence parsers.
type Metrics struct {
	C0CoveragePercent *float64
	C1CoveragePercent *float64
	OperationsOwner   string
	OperationsRunbook string
}

// Result captures a single gate evaluation.
type Result struct {
	ID          string  `json:"id"`
	Passed      bool    `json:"passed"`
	Required    bool    `json:"required"`
	Severity    string  `json:"severity"`
	Actual      float64 `json:"actual,omitempty"`
	RequiredMin float64 `json:"requiredMin,omitempty"`
	Message     string  `json:"message"`
}

// Evaluation is the machine-readable readiness decision.
type Evaluation struct {
	Status  DecisionStatus `json:"status"`
	Results []Result       `json:"results"`
}

// Evaluate compares collected metrics with the gates in a readiness standard.
func Evaluate(s standard.Standard, metrics Metrics) (Evaluation, error) {
	evaluation := Evaluation{Status: DecisionReady}

	for _, gate := range s.Gates {
		switch gate.ID {
		case standard.C0CoverageGateID:
			if metrics.C0CoveragePercent == nil {
				return Evaluation{Status: DecisionError}, fmt.Errorf("%s gate requires c0Coverage evidence", standard.C0CoverageGateID)
			}
			result := evaluateCoverageGate(gate, "C0 coverage", *metrics.C0CoveragePercent)
			evaluation.Results = append(evaluation.Results, result)
			if gate.Required && !result.Passed {
				evaluation.Status = DecisionBlocked
			}
		case standard.C1CoverageGateID:
			if metrics.C1CoveragePercent == nil {
				return Evaluation{Status: DecisionError}, fmt.Errorf("%s gate requires c1Coverage evidence", standard.C1CoverageGateID)
			}
			result := evaluateCoverageGate(gate, "C1 coverage", *metrics.C1CoveragePercent)
			evaluation.Results = append(evaluation.Results, result)
			if gate.Required && !result.Passed {
				evaluation.Status = DecisionBlocked
			}
		case standard.OperationsOwnerExistsGateID:
			result := evaluateExistsGate(gate, metrics.OperationsOwner != "", "Operations owner is defined.", "Operations owner is not defined.")
			evaluation.Results = append(evaluation.Results, result)
			if gate.Required && !result.Passed {
				evaluation.Status = DecisionBlocked
			}
		case standard.OperationsRunbookExistsGateID:
			result := evaluateExistsGate(gate, metrics.OperationsRunbook != "", "Operations runbook is defined.", "Operations runbook is not defined.")
			evaluation.Results = append(evaluation.Results, result)
			if gate.Required && !result.Passed {
				evaluation.Status = DecisionBlocked
			}
		default:
			if gate.Required {
				return Evaluation{Status: DecisionError}, fmt.Errorf("unsupported required gate %q", gate.ID)
			}
		}
	}

	return evaluation, nil
}

func evaluateCoverageGate(gate standard.Gate, label string, actual float64) Result {
	requiredMin := *gate.Threshold.Min
	result := Result{
		ID:          gate.ID,
		Required:    gate.Required,
		Severity:    gate.Severity,
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
		Severity: gate.Severity,
		Passed:   exists,
	}
	if exists {
		result.Message = passedMessage
		return result
	}

	result.Message = failedMessage
	return result
}

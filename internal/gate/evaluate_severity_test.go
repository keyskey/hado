package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateRequiredCriticalGateBlocksWhenFailed(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.OperationsOwnerExistsGateID, Severity: standard.SeverityCritical, Required: true},
		},
	}, Metrics{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if evaluation.Status != DecisionBlocked {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionBlocked)
	}
}

func TestEvaluateRequiredMajorGateDoesNotBlockWhenFailed(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.OperationsOwnerExistsGateID, Severity: standard.SeverityMajor, Required: true},
		},
	}, Metrics{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
}

func TestEvaluateRequiredMinorGateDoesNotBlockWhenFailed(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.OperationsOwnerExistsGateID, Severity: standard.SeverityMinor, Required: true},
		},
	}, Metrics{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}
	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
}

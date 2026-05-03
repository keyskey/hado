package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateCoverageAndOperationsGatesTogether(t *testing.T) {
	minimum := 70.0
	actual := 72.5
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.C1CoverageGateID,
				Required: true,
				Threshold: standard.Threshold{
					Min: &minimum,
				},
			},
			{
				ID:       standard.OperationsOwnerExistsGateID,
				Required: true,
			},
		},
	}, Metrics{
		C1CoveragePercent: &actual,
		OperationsOwner:   "platform-team",
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
}

func TestEvaluateRejectsUnsupportedRequiredGate(t *testing.T) {
	_, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       "com.example.custom_gate",
				Required: true,
			},
		},
	}, Metrics{})
	if err == nil {
		t.Fatal("expected error for unsupported required gate")
	}
}

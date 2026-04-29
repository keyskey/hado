package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateReadyWhenOperationsEvidenceExists(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.OperationsOwnerExistsGateID,
				Required: true,
			},
			{
				ID:       standard.OperationsRunbookExistsGateID,
				Required: true,
			},
		},
	}, Metrics{
		OperationsOwner:   "platform-team",
		OperationsRunbook: "https://example.com/runbooks/order-api",
	})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
	if len(evaluation.Results) != 2 {
		t.Fatalf("result count = %d, want 2", len(evaluation.Results))
	}
	for _, result := range evaluation.Results {
		if !result.Passed {
			t.Fatalf("operation gate did not pass: %#v", result)
		}
	}
}

func TestEvaluateBlockedWhenOperationsEvidenceIsMissing(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.OperationsOwnerExistsGateID,
				Required: true,
			},
			{
				ID:       standard.OperationsRunbookExistsGateID,
				Required: true,
			},
		},
	}, Metrics{})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionBlocked {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionBlocked)
	}
	for _, result := range evaluation.Results {
		if result.Passed {
			t.Fatalf("operation gate passed unexpectedly: %#v", result)
		}
	}
}

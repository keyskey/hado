package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateReleaseRollbackPlanReady(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.ReleaseRollbackPlanExistsGateID, Required: true},
		},
	}, Metrics{
		ReleaseRollbackPlan: "docs/rollback.md",
	})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
}

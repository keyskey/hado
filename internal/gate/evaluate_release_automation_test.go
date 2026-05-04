package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateReleaseAutomationDeclared(t *testing.T) {
	t.Parallel()
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.ReleaseAutomationDeclaredGateID, Severity: standard.SeverityCritical, Required: true},
		},
	}, Metrics{ReleaseAutomationDeclared: true})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want ready", evaluation.Status)
	}
}

func TestEvaluateReleaseAutomationDeclaredBlockedWhenFalse(t *testing.T) {
	t.Parallel()
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.ReleaseAutomationDeclaredGateID, Severity: standard.SeverityCritical, Required: true},
		},
	}, Metrics{})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if evaluation.Status != DecisionBlocked {
		t.Fatalf("status = %q, want blocked", evaluation.Status)
	}
}

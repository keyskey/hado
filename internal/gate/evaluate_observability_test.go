package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateObservabilityReady(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.ObservabilitySLOExistsGateID, Required: true},
			{ID: standard.ObservabilityMonitorExistsGateID, Required: true},
			{ID: standard.ObservabilityDashboardExistsGateID, Required: true},
		},
	}, Metrics{
		ObservabilitySLO:       "slo.yaml",
		ObservabilityMonitors:  "monitors.yaml",
		ObservabilityDashboard: "https://example.com/dash",
	})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
}

func TestEvaluateObservabilityBlockedWhenDashboardEmpty(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.ObservabilityDashboardExistsGateID, Severity: standard.SeverityCritical, Required: true},
		},
	}, Metrics{
		ObservabilitySLO:       "ok",
		ObservabilityMonitors:  "ok",
		ObservabilityDashboard: "",
	})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if evaluation.Status != DecisionBlocked {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionBlocked)
	}
}

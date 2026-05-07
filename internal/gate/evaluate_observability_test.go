package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateObservabilityReady(t *testing.T) {
	st := standard.Standard{
		Gates: []standard.Gate{
			{ID: standard.ObservabilitySLOExistsGateID, Required: true},
			{ID: standard.ObservabilityMonitorExistsGateID, Required: true},
			{ID: standard.ObservabilityDashboardExistsGateID, Required: true},
		},
	}
	eval, err := Evaluate(st, Metrics{
		ObservabilitySLOPresent:       true,
		ObservabilityMonitorsPresent:  true,
		ObservabilityDashboardPresent: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range eval.Results {
		if !r.Passed {
			t.Fatalf("%s: want pass, got %s", r.ID, r.Message)
		}
	}
}

func TestEvaluateObservabilityBlockedWhenDashboardEmpty(t *testing.T) {
	st := standard.Standard{
		Gates: []standard.Gate{
			{ID: standard.ObservabilityDashboardExistsGateID, Severity: standard.SeverityCritical, Required: true},
		},
	}
	eval, err := Evaluate(st, Metrics{
		ObservabilitySLOPresent:       true,
		ObservabilityMonitorsPresent:  true,
		ObservabilityDashboardPresent: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(eval.Results) != 1 {
		t.Fatalf("results = %#v", eval.Results)
	}
	if eval.Results[0].Passed {
		t.Fatal("expected dashboard gate to fail")
	}
}

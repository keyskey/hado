package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateReadyWhenC0CoverageMeetsStandard(t *testing.T) {
	minimum := 80.0
	actual := 80.0
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.C0CoverageGateID,
				Severity: standard.SeverityCritical,
				Required: true,
				Threshold: standard.Threshold{
					Min: &minimum,
				},
			},
		},
	}, Metrics{C0CoveragePercent: &actual})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
	if len(evaluation.Results) != 1 {
		t.Fatalf("result count = %d, want 1", len(evaluation.Results))
	}
	if !evaluation.Results[0].Passed {
		t.Fatalf("C0 coverage gate did not pass: %#v", evaluation.Results[0])
	}
}

func TestEvaluateBlockedWhenC0CoverageFallsBelowStandard(t *testing.T) {
	minimum := 80.0
	actual := 79.9
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.C0CoverageGateID,
				Severity: standard.SeverityCritical,
				Required: true,
				Threshold: standard.Threshold{
					Min: &minimum,
				},
			},
		},
	}, Metrics{C0CoveragePercent: &actual})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionBlocked {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionBlocked)
	}
	if evaluation.Results[0].Passed {
		t.Fatalf("C0 coverage gate passed unexpectedly: %#v", evaluation.Results[0])
	}
}

func TestEvaluateReadyWhenC1CoverageMeetsStandard(t *testing.T) {
	minimum := 70.0
	actual := 70.1
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.C1CoverageGateID,
				Severity: standard.SeverityCritical,
				Required: true,
				Threshold: standard.Threshold{
					Min: &minimum,
				},
			},
		},
	}, Metrics{C1CoveragePercent: &actual})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
	if !evaluation.Results[0].Passed {
		t.Fatalf("C1 coverage gate did not pass: %#v", evaluation.Results[0])
	}
}

func TestEvaluateBlockedWhenC1CoverageFallsBelowStandard(t *testing.T) {
	minimum := 70.0
	actual := 68.4
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.C1CoverageGateID,
				Severity: standard.SeverityCritical,
				Required: true,
				Threshold: standard.Threshold{
					Min: &minimum,
				},
			},
		},
	}, Metrics{C1CoveragePercent: &actual})
	if err != nil {
		t.Fatalf("Evaluate returned error: %v", err)
	}

	if evaluation.Status != DecisionBlocked {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionBlocked)
	}
	if evaluation.Results[0].Passed {
		t.Fatalf("C1 coverage gate passed unexpectedly: %#v", evaluation.Results[0])
	}
}

func TestEvaluateErrorsWhenC1CoverageEvidenceIsMissing(t *testing.T) {
	minimum := 70.0
	_, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{
				ID:       standard.C1CoverageGateID,
				Required: true,
				Threshold: standard.Threshold{
					Min: &minimum,
				},
			},
		},
	}, Metrics{})
	if err == nil {
		t.Fatal("expected error for missing C1 coverage evidence")
	}
}

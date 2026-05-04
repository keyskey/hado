package evaluate

import (
	"strings"
	"testing"

	"github.com/keyskey/hado/internal/gate"
	"github.com/keyskey/hado/internal/standard"
)

func TestPrintTextEvaluationPrintsSummaryLastWithSeverityHints(t *testing.T) {
	t.Parallel()

	var b strings.Builder
	printTextEvaluation(&b, gate.Evaluation{
		Status: gate.DecisionBlocked,
		Results: []gate.Result{
			{
				ID:       standard.C0CoverageGateID,
				Passed:   false,
				Required: true,
				Severity: standard.SeverityCritical,
				Message:  "C0 coverage is below threshold.",
			},
			{
				ID:       standard.OperationsRunbookExistsGateID,
				Passed:   false,
				Required: true,
				Severity: standard.SeverityMajor,
				Message:  "Operations runbook is not defined.",
			},
		},
	})

	out := b.String()
	if !strings.Contains(out, "(severity: critical)") {
		t.Fatalf("stdout = %q, want severity column for critical", out)
	}
	if !strings.Contains(out, "release blocked: fix before release") {
		t.Fatalf("stdout = %q, want critical release hint", out)
	}
	if !strings.Contains(out, "release allowed: fix soon after release") {
		t.Fatalf("stdout = %q, want major release hint", out)
	}
	if !strings.HasSuffix(out, "\nHADO: BLOCKED\n") {
		t.Fatalf("stdout = %q, want summary line at end", out)
	}
}

func TestReleaseActionHint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		result gate.Result
		want   string
	}{
		{
			name: "optional gate",
			result: gate.Result{
				Required: false,
				Severity: standard.SeverityCritical,
			},
			want: "optional gate: release allowed",
		},
		{
			name: "critical required",
			result: gate.Result{
				Required: true,
				Severity: standard.SeverityCritical,
			},
			want: "release blocked: fix before release",
		},
		{
			name: "major required",
			result: gate.Result{
				Required: true,
				Severity: standard.SeverityMajor,
			},
			want: "release allowed: fix soon after release",
		},
		{
			name: "minor default",
			result: gate.Result{
				Required: true,
				Severity: standard.SeverityMinor,
			},
			want: "release allowed: fix when appropriate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := releaseActionHint(tt.result); got != tt.want {
				t.Fatalf("releaseActionHint() = %q, want %q", got, tt.want)
			}
		})
	}
}

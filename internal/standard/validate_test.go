package standard

import (
	"strings"
	"testing"
)

func TestValidateRejectsUnsupportedSeverity(t *testing.T) {
	t.Parallel()

	s := Standard{
		ID: "test-standard",
		Gates: []Gate{
			{
				ID:       OperationsOwnerExistsGateID,
				Severity: "urgent",
				Required: true,
			},
		},
	}
	if err := s.Validate(); err == nil {
		t.Fatal("Validate() error = nil, want unsupported severity error")
	} else if !strings.Contains(err.Error(), "severity") {
		t.Fatalf("Validate() error = %v, want severity context", err)
	}
}

func TestValidateAcceptsEnumSeverities(t *testing.T) {
	t.Parallel()

	s := Standard{
		ID: "test-standard",
		Gates: []Gate{
			{ID: OperationsOwnerExistsGateID, Severity: SeverityCritical, Required: true},
			{ID: OperationsRunbookExistsGateID, Severity: SeverityMajor, Required: true},
			{ID: ObservabilityDashboardExistsGateID, Severity: SeverityMinor, Required: false},
		},
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestSeverityValidate(t *testing.T) {
	t.Parallel()

	for _, severity := range []Severity{"", SeverityCritical, SeverityMajor, SeverityMinor} {
		if err := severity.Validate(); err != nil {
			t.Fatalf("Validate(%q) error = %v", severity, err)
		}
	}
	if err := Severity("urgent").Validate(); err == nil {
		t.Fatal("Validate(urgent) error = nil, want unsupported severity")
	}
}

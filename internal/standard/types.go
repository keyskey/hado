package standard

const (
	// C0CoverageGateID is the gate id used for C0 statement coverage.
	C0CoverageGateID = "test.c0_coverage"
	// C1CoverageGateID is the gate id used for C1 condition coverage.
	C1CoverageGateID = "test.c1_coverage"
	// OperationsOwnerExistsGateID is the gate id used for operational owner readiness.
	OperationsOwnerExistsGateID = "operations.owner_exists"
	// OperationsRunbookExistsGateID is the gate id used for operational runbook readiness.
	OperationsRunbookExistsGateID = "operations.runbook_exists"
)

// Standard defines the production readiness gates a service must satisfy.
type Standard struct {
	ID    string `yaml:"id" json:"id"`
	Name  string `yaml:"name" json:"name"`
	Gates []Gate `yaml:"gates" json:"gates"`
}

// Gate is a single readiness check declared by a standard.
type Gate struct {
	ID        string    `yaml:"id" json:"id"`
	Severity  string    `yaml:"severity" json:"severity"`
	Required  bool      `yaml:"required" json:"required"`
	Threshold Threshold `yaml:"threshold" json:"threshold"`
}

// Threshold contains numeric limits for a gate.
type Threshold struct {
	Min *float64 `yaml:"min" json:"min,omitempty"`
}

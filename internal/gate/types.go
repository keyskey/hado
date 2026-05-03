package gate

// DecisionStatus is the final production readiness decision.
type DecisionStatus string

const (
	// DecisionReady means all required gates passed.
	DecisionReady DecisionStatus = "ready"
	// DecisionBlocked means at least one required gate failed.
	DecisionBlocked DecisionStatus = "blocked"
	// DecisionError means HADO could not complete the evaluation.
	DecisionError DecisionStatus = "error"
)

// Metrics are normalized values supplied by modules or evidence parsers.
type Metrics struct {
	C0CoveragePercent *float64
	C1CoveragePercent *float64
	OperationsOwner   string
	OperationsRunbook string
}

// Result captures a single gate evaluation.
type Result struct {
	ID          string  `json:"id"`
	Passed      bool    `json:"passed"`
	Required    bool    `json:"required"`
	Severity    string  `json:"severity"`
	Actual      float64 `json:"actual,omitempty"`
	RequiredMin float64 `json:"requiredMin,omitempty"`
	Message     string  `json:"message"`
}

// Evaluation is the machine-readable readiness decision.
type Evaluation struct {
	Status  DecisionStatus `json:"status"`
	Results []Result       `json:"results"`
}

package standard

// Severity is the importance level assigned to each gate.
type Severity string

const (
	// SeverityCritical blocks release when a required gate fails.
	SeverityCritical Severity = "critical"
	// SeverityMajor does not block release but requires prompt follow-up.
	SeverityMajor Severity = "major"
	// SeverityMinor does not block release and can be addressed opportunistically.
	SeverityMinor Severity = "minor"

	// C0CoverageGateID is the gate id used for C0 statement coverage.
	C0CoverageGateID = "test.c0_coverage"
	// C1CoverageGateID is the gate id used for C1 condition coverage.
	C1CoverageGateID = "test.c1_coverage"
	// OperationsOwnerExistsGateID is the gate id used for operational owner readiness.
	OperationsOwnerExistsGateID = "operations.owner_exists"
	// OperationsRunbookExistsGateID is the gate id used for operational runbook readiness.
	OperationsRunbookExistsGateID = "operations.runbook_exists"

	// ObservabilitySLOExistsGateID gates on declared SLO / SLI evidence (manifest reference non-empty).
	ObservabilitySLOExistsGateID = "observability.slo_exists"
	// ObservabilityMonitorExistsGateID gates on monitor / alert evidence reference.
	ObservabilityMonitorExistsGateID = "observability.monitor_exists"
	// ObservabilityDashboardExistsGateID gates on dashboard evidence reference.
	ObservabilityDashboardExistsGateID = "observability.dashboard_exists"

	// InfraDeploymentSpecExistsGateID gates on a deployment / workload spec reference (path, bundle id, URL, etc.).
	InfraDeploymentSpecExistsGateID = "infra.deployment_spec_exists"

	// ReleaseRollbackPlanExistsGateID gates on a documented rollback plan reference.
	ReleaseRollbackPlanExistsGateID = "release.rollback_plan_exists"
	// ReleaseAutomationDeclaredGateID gates on at least one non-empty workflow_refs entry under evidence.release.automation.
	ReleaseAutomationDeclaredGateID = "release.automation_declared"
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
	Severity  Severity  `yaml:"severity" json:"severity"`
	Required  bool      `yaml:"required" json:"required"`
	Threshold Threshold `yaml:"threshold" json:"threshold"`
}

// Threshold contains numeric limits for a gate.
type Threshold struct {
	Min *float64 `yaml:"min" json:"min,omitempty"`
}

package manifest

// Manifest declares the evaluated service and the evidence HADO should read.
type Manifest struct {
	Version  string   `yaml:"version" json:"version,omitempty"`
	Evidence Evidence `yaml:"evidence" json:"evidence,omitempty"`

	baseDir string
}

// Evidence groups evidence declarations by readiness domain.
type Evidence struct {
	Coverage      CoverageEvidence      `yaml:"coverage" json:"coverage,omitempty"`
	Operations    OperationsEvidence    `yaml:"operations" json:"operations,omitempty"`
	Observability ObservabilityEvidence `yaml:"observability" json:"observability,omitempty"`
	Infra         InfraEvidence         `yaml:"infra" json:"infra,omitempty"`
	Release       ReleaseEvidence       `yaml:"release" json:"release,omitempty"`
}

// CoverageEvidence declares coverage reports and the adapters that parse them.
type CoverageEvidence struct {
	Inputs []CoverageInput `yaml:"inputs" json:"inputs,omitempty"`
}

// CoverageInput identifies one coverage artifact and its adapter.
type CoverageInput struct {
	Adapter string `yaml:"adapter" json:"adapter"`
	Path    string `yaml:"path" json:"path"`
}

// OperationsEvidence declares operational ownership and response evidence.
type OperationsEvidence struct {
	Owner   string `yaml:"owner" json:"owner,omitempty"`
	Runbook string `yaml:"runbook" json:"runbook,omitempty"`
}

// ObservabilityEvidence declares references to SLO, monitors, and dashboard evidence (paths, URLs, or catalog IDs).
type ObservabilityEvidence struct {
	SLO       string `yaml:"slo" json:"slo,omitempty"`
	Monitors  string `yaml:"monitors" json:"monitors,omitempty"`
	Dashboard string `yaml:"dashboard" json:"dashboard,omitempty"`
}

// InfraEvidence declares infrastructure-related evidence references (deployment spec, IaC pointer, etc.).
type InfraEvidence struct {
	DeploymentSpec string `yaml:"deployment_spec" json:"deployment_spec,omitempty"`
}

// ReleaseEvidence declares release and rollback-related references.
type ReleaseEvidence struct {
	RollbackPlan string `yaml:"rollback_plan" json:"rollback_plan,omitempty"`
}

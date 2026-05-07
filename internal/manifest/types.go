package manifest

import "strings"

// Manifest declares the evaluated service and the evidence HADO should read.
type Manifest struct {
	Version  string      `yaml:"version" json:"version,omitempty"`
	Service  Service     `yaml:"service,omitempty" json:"service,omitempty"`
	Standard StandardRef `yaml:"standard,omitempty" json:"standard,omitempty"`
	Evidence Evidence    `yaml:"evidence,omitempty" json:"evidence,omitempty"`

	baseDir string
}

// Service identifies the workload HADO evaluates (minimal fields for targeting).
type Service struct {
	ID   string `yaml:"id,omitempty" json:"id,omitempty"`
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
}

// StandardRef points at the Readiness Standard applied to this service (id or file path).
type StandardRef struct {
	ID string `yaml:"id,omitempty" json:"id,omitempty"`
}

// Evidence groups evidence declarations by readiness domain.
// Sub-blocks are pointers so an empty scaffold (e.g. operations with only empty strings) still serializes under evidence.
type Evidence struct {
	Coverage      *CoverageEvidence      `yaml:"coverage,omitempty" json:"coverage,omitempty"`
	Operations    *OperationsEvidence    `yaml:"operations,omitempty" json:"operations,omitempty"`
	Observability *ObservabilityEvidence `yaml:"observability,omitempty" json:"observability,omitempty"`
	Infra         *InfraEvidence         `yaml:"infra,omitempty" json:"infra,omitempty"`
	Release       *ReleaseEvidence       `yaml:"release,omitempty" json:"release,omitempty"`
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

// ObservabilityLink is a named, browser-openable URL for SLO, monitor, or dashboard evidence (audit / ops).
type ObservabilityLink struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	URL  string `yaml:"url" json:"url"`
}

// ObservabilityEvidence declares observability evidence as lists of links (typically vendor UI URLs).
type ObservabilityEvidence struct {
	SLOs       []ObservabilityLink `yaml:"slos,omitempty" json:"slos,omitempty"`
	Monitors   []ObservabilityLink `yaml:"monitors,omitempty" json:"monitors,omitempty"`
	Dashboards []ObservabilityLink `yaml:"dashboards,omitempty" json:"dashboards,omitempty"`
}

// ObservabilityLinksHaveURL reports whether at least one entry has a non-empty URL after trimming spaces.
func ObservabilityLinksHaveURL(links []ObservabilityLink) bool {
	for _, l := range links {
		if strings.TrimSpace(l.URL) != "" {
			return true
		}
	}
	return false
}

// InfraEvidence declares infrastructure-related evidence references (deployment spec, IaC pointer, etc.).
type InfraEvidence struct {
	DeploymentSpec string `yaml:"deployment_spec" json:"deployment_spec,omitempty"`
}

// ReleaseEvidence declares release and rollback-related references.
type ReleaseEvidence struct {
	RollbackPlan string                    `yaml:"rollback_plan" json:"rollback_plan,omitempty"`
	Automation   ReleaseAutomationEvidence `yaml:"automation" json:"automation,omitempty"`
}

// ReleaseAutomationEvidence declares where automated release / deploy pipelines live (paths, URLs, workflow names).
// Systems (e.g. github_actions, circleci, argo_workflow) are optional metadata for tooling; phase-1 gates use workflow_refs only.
type ReleaseAutomationEvidence struct {
	WorkflowRefs []string `yaml:"workflow_refs" json:"workflow_refs,omitempty"`
	Systems      []string `yaml:"systems,omitempty" json:"systems,omitempty"`
}

package manifest

// Manifest declares the evaluated service and the evidence HADO should read.
type Manifest struct {
	Version  string   `yaml:"version" json:"version,omitempty"`
	Evidence Evidence `yaml:"evidence" json:"evidence,omitempty"`

	baseDir string
}

// Evidence groups evidence declarations by readiness domain.
type Evidence struct {
	Coverage   CoverageEvidence   `yaml:"coverage" json:"coverage,omitempty"`
	Operations OperationsEvidence `yaml:"operations" json:"operations,omitempty"`
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

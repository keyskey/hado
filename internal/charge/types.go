package charge

// Preset identifies a built-in charge strategy (local commands only in MVP).
const (
	PresetGoGobce = "go-gobce"
)

// RunOptions configures a charge run (manifest path, standard resolution, apply mode).
type RunOptions struct {
	ManifestPath string
	// StandardPath overrides manifest standard.id when non-empty (path to standard YAML).
	StandardPath string
	Preset       string
	Apply        bool
}

// CoveragePlan is the output of the plan step for coverage gates.
type CoveragePlan struct {
	StandardPath string
	RequiresC0   bool
	RequiresC1   bool
	// PreferredAdapterFormats lists adapter formats to try in order for satisfying the plan.
	PreferredAdapterFormats []string
}

// CoverageGapReport is the output of the gap step.
type CoverageGapReport struct {
	Satisfied         bool
	SatisfyingAdapter string
	SatisfyingPath    string
	Items             []CoverageGapItem
}

// CoverageGapItem describes one gap or observation.
type CoverageGapItem struct {
	Code   string // e.g. "missing_input", "missing_file", "insufficient_adapter", "satisfied"
	Detail string
}

// ApplyCoverageResult is returned after a successful apply (coverage preset).
type ApplyCoverageResult struct {
	CoverageOutPath string
	GobceJSONPath   string
	WroteManifest   bool
}

// ChargeReport is the full result of a charge run (for JSON output and tests).
type ChargeReport struct {
	Plan  CoveragePlan         `json:"plan"`
	Gap   CoverageGapReport    `json:"gap"`
	Apply *ApplyCoverageResult `json:"apply,omitempty"`
}

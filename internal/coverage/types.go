package coverage

// Metrics contains normalized coverage values supplied to HADO by any producer.
type Metrics struct {
	C0Coverage *float64 `json:"c0Coverage"`
	C1Coverage *float64 `json:"c1Coverage"`
}

// AdapterInput identifies a coverage report and the adapter that can parse it.
type AdapterInput struct {
	Format string
	Path   string
}

const (
	// FormatHADOJSON is HADO's normalized coverage metrics format.
	FormatHADOJSON = "hado-json"
	// FormatGoCoverprofile parses Go's coverprofile format and emits C0 coverage.
	FormatGoCoverprofile = "go-coverprofile"
	// FormatGobceJSON parses keyskey/gobce JSON output and emits C0/C1 coverage.
	FormatGobceJSON = "gobce-json"
)

type gobceResult struct {
	StatementCoverage       float64 `json:"statementCoverage"`
	EstimatedBranchCoverage float64 `json:"estimatedBranchCoverage"`
}

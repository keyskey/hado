// Package coverage parses Go coverage evidence.
package coverage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

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

// ParseMetrics reads normalized coverage metrics from JSON.
func ParseMetrics(path string) (Metrics, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Metrics{}, fmt.Errorf("read coverage metrics: %w", err)
	}

	var metrics Metrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return Metrics{}, fmt.Errorf("parse coverage metrics: %w", err)
	}
	if err := metrics.Validate(); err != nil {
		return Metrics{}, err
	}

	return metrics, nil
}

// ParseAdapterInput parses one coverage report using the requested adapter.
func ParseAdapterInput(input AdapterInput) (Metrics, error) {
	switch input.Format {
	case FormatHADOJSON:
		return ParseMetrics(input.Path)
	case FormatGoCoverprofile:
		summary, err := ParseGoProfile(input.Path)
		if err != nil {
			return Metrics{}, err
		}
		return Metrics{C0Coverage: &summary.C0Coverage}, nil
	case FormatGobceJSON:
		return ParseGobceJSON(input.Path)
	default:
		return Metrics{}, fmt.Errorf("unsupported coverage adapter %q", input.Format)
	}
}

// Merge overlays non-nil metrics from next onto metrics.
func (metrics Metrics) Merge(next Metrics) Metrics {
	if next.C0Coverage != nil {
		metrics.C0Coverage = next.C0Coverage
	}
	if next.C1Coverage != nil {
		metrics.C1Coverage = next.C1Coverage
	}
	return metrics
}

// ParseCoverageSpec parses "<adapter>:<path>" CLI values.
func ParseCoverageSpec(spec string) (AdapterInput, error) {
	format, path, ok := strings.Cut(spec, ":")
	if !ok || format == "" || path == "" {
		return AdapterInput{}, fmt.Errorf("coverage input must be <adapter>:<path>")
	}
	return AdapterInput{Format: format, Path: path}, nil
}

// ParseAdapterInputs parses and merges multiple coverage adapter inputs.
func ParseAdapterInputs(inputs []AdapterInput) (Metrics, error) {
	var merged Metrics
	for _, input := range inputs {
		metrics, err := ParseAdapterInput(input)
		if err != nil {
			return Metrics{}, err
		}
		merged = merged.Merge(metrics)
	}
	if err := merged.Validate(); err != nil {
		return Metrics{}, err
	}
	return merged, nil
}

// ParseCoverageSpecs parses and merges multiple "<adapter>:<path>" values.
func ParseCoverageSpecs(specs []string) (Metrics, error) {
	inputs := make([]AdapterInput, 0, len(specs))
	for _, spec := range specs {
		input, err := ParseCoverageSpec(spec)
		if err != nil {
			return Metrics{}, err
		}
		inputs = append(inputs, input)
	}
	return ParseAdapterInputs(inputs)
}

type gobceResult struct {
	StatementCoverage       float64 `json:"statementCoverage"`
	EstimatedBranchCoverage float64 `json:"estimatedBranchCoverage"`
}

// ParseGobceJSON reads the current keyskey/gobce JSON output format.
func ParseGobceJSON(path string) (Metrics, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Metrics{}, fmt.Errorf("read gobce json: %w", err)
	}

	var result gobceResult
	if err := json.Unmarshal(data, &result); err != nil {
		return Metrics{}, fmt.Errorf("parse gobce json: %w", err)
	}
	metrics := Metrics{
		C0Coverage: &result.StatementCoverage,
		C1Coverage: &result.EstimatedBranchCoverage,
	}
	if err := metrics.Validate(); err != nil {
		return Metrics{}, err
	}
	return metrics, nil
}

// Validate checks that all present coverage percentages are in range.
func (metrics Metrics) Validate() error {
	if metrics.C0Coverage != nil && (*metrics.C0Coverage < 0 || *metrics.C0Coverage > 100) {
		return fmt.Errorf("c0Coverage must be between 0 and 100")
	}
	if metrics.C1Coverage != nil && (*metrics.C1Coverage < 0 || *metrics.C1Coverage > 100) {
		return fmt.Errorf("c1Coverage must be between 0 and 100")
	}
	return nil
}

package coverage

import (
	"encoding/json"
	"fmt"
	"os"
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

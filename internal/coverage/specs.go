package coverage

import (
	"fmt"
	"strings"
)

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

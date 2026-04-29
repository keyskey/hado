// Package standard loads reusable readiness standards.
package standard

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	// C0CoverageGateID is the gate id used for C0 statement coverage.
	C0CoverageGateID = "test.c0_coverage"
	// C1CoverageGateID is the gate id used for C1 condition coverage.
	C1CoverageGateID = "test.c1_coverage"
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
	Severity  string    `yaml:"severity" json:"severity"`
	Required  bool      `yaml:"required" json:"required"`
	Threshold Threshold `yaml:"threshold" json:"threshold"`
}

// Threshold contains numeric limits for a gate.
type Threshold struct {
	Min *float64 `yaml:"min" json:"min,omitempty"`
}

// Load reads a YAML readiness standard from disk.
func Load(path string) (Standard, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Standard{}, fmt.Errorf("read standard: %w", err)
	}

	var standard Standard
	if err := yaml.Unmarshal(data, &standard); err != nil {
		return Standard{}, fmt.Errorf("parse standard: %w", err)
	}
	if err := standard.Validate(); err != nil {
		return Standard{}, err
	}

	return standard, nil
}

// Validate checks the fields required by the first HADO gate evaluator.
func (s Standard) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("standard id is required")
	}
	if len(s.Gates) == 0 {
		return fmt.Errorf("standard must define at least one gate")
	}
	for i, gate := range s.Gates {
		if gate.ID == "" {
			return fmt.Errorf("gate %d id is required", i)
		}
		if isCoverageGate(gate.ID) && gate.Threshold.Min == nil {
			return fmt.Errorf("%s gate requires threshold.min", gate.ID)
		}
	}

	return nil
}

// RequiresGate reports whether the standard declares a gate id.
func (s Standard) RequiresGate(id string) bool {
	for _, gate := range s.Gates {
		if gate.ID == id {
			return true
		}
	}
	return false
}

func isCoverageGate(id string) bool {
	return id == C0CoverageGateID || id == C1CoverageGateID
}

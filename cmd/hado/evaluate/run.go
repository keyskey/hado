package evaluate

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/keyskey/hado/internal/coverage"
	"github.com/keyskey/hado/internal/gate"
	"github.com/keyskey/hado/internal/manifest"
	"github.com/keyskey/hado/internal/standard"
)

type stringList []string

func (values *stringList) String() string {
	return strings.Join(*values, ",")
}

func (values *stringList) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func Run(args []string, stdout, stderr io.Writer) (int, error) {
	fs := flag.NewFlagSet("evaluate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	standardPath := fs.String("standard", "", "path to readiness standard YAML")
	manifestPath := fs.String("manifest", "", "path to HADO manifest YAML")
	var coverageInputs stringList
	fs.Var(&coverageInputs, "coverage-input", "coverage input as <adapter>:<path>; overrides manifest coverage inputs when set; adapters: hado-json, go-coverprofile, gobce-json")
	output := fs.String("output", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	if *standardPath == "" {
		return 2, fmt.Errorf("evaluate requires --standard")
	}

	readinessStandard, err := standard.Load(*standardPath)
	if err != nil {
		return 2, err
	}

	hadoManifest, err := loadManifest(*manifestPath)
	if err != nil {
		return 2, err
	}
	adapterInputs, err := resolveCoverageInputs(coverageInputs, hadoManifest)
	if err != nil {
		return 2, err
	}
	if requiresCoverage(readinessStandard) && len(adapterInputs) == 0 {
		return 2, fmt.Errorf("evaluate requires --coverage-input or --manifest with evidence.coverage.inputs")
	}
	metrics := gate.Metrics{}
	if len(adapterInputs) > 0 {
		coverageMetrics, err := coverage.ParseAdapterInputs(adapterInputs)
		if err != nil {
			return 2, err
		}
		metrics.C0CoveragePercent = coverageMetrics.C0Coverage
		metrics.C1CoveragePercent = coverageMetrics.C1Coverage
	}
	if hadoManifest != nil {
		if op := hadoManifest.Evidence.Operations; op != nil {
			metrics.OperationsOwner = strings.TrimSpace(op.Owner)
			metrics.OperationsRunbook = strings.TrimSpace(op.Runbook)
		}
		if obs := hadoManifest.Evidence.Observability; obs != nil {
			metrics.ObservabilitySLO = strings.TrimSpace(obs.SLO)
			metrics.ObservabilityMonitors = strings.TrimSpace(obs.Monitors)
			metrics.ObservabilityDashboard = strings.TrimSpace(obs.Dashboard)
		}
		if rel := hadoManifest.Evidence.Release; rel != nil {
			metrics.ReleaseRollbackPlan = strings.TrimSpace(rel.RollbackPlan)
			metrics.ReleaseAutomationDeclared = rel.AutomationDeclared()
		}
		if inf := hadoManifest.Evidence.Infra; inf != nil {
			metrics.InfraDeploymentSpec = strings.TrimSpace(inf.DeploymentSpec)
		}
	}

	evaluation, err := gate.Evaluate(readinessStandard, metrics)
	if err != nil {
		return 2, err
	}

	switch *output {
	case "text":
		printTextEvaluation(stdout, evaluation)
	case "json":
		if err := json.NewEncoder(stdout).Encode(evaluation); err != nil {
			return 2, fmt.Errorf("write json evaluation: %w", err)
		}
	default:
		return 2, fmt.Errorf("unsupported output format %q", *output)
	}

	if evaluation.Status == gate.DecisionBlocked {
		return 1, nil
	}
	return 0, nil
}

func loadManifest(manifestPath string) (*manifest.Manifest, error) {
	if manifestPath == "" {
		return nil, nil
	}

	hadoManifest, err := manifest.Load(manifestPath)
	if err != nil {
		return nil, err
	}
	return &hadoManifest, nil
}

func resolveCoverageInputs(coverageInputs []string, hadoManifest *manifest.Manifest) ([]coverage.AdapterInput, error) {
	if len(coverageInputs) > 0 {
		inputs := make([]coverage.AdapterInput, 0, len(coverageInputs))
		for _, spec := range coverageInputs {
			input, err := coverage.ParseCoverageSpec(spec)
			if err != nil {
				return nil, err
			}
			inputs = append(inputs, input)
		}
		return inputs, nil
	}
	if hadoManifest == nil {
		return nil, nil
	}

	return hadoManifest.CoverageAdapterInputs(), nil
}

func requiresCoverage(readinessStandard standard.Standard) bool {
	return readinessStandard.RequiresGate(standard.C0CoverageGateID) || readinessStandard.RequiresGate(standard.C1CoverageGateID)
}

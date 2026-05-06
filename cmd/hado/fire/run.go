package fire

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/keyskey/hado/internal/coverage"
	"github.com/keyskey/hado/internal/gate"
	"github.com/keyskey/hado/internal/manifest"
	"github.com/keyskey/hado/internal/standard"
)

func Run(args []string, stdout, stderr io.Writer) (int, error) {
	fs := flag.NewFlagSet("fire", flag.ContinueOnError)
	fs.SetOutput(stderr)
	standardRef := fs.String("standard", "", "Readiness Standard id or path (optional; defaults to manifest standard.id)")
	manifestPath := fs.String("manifest", "", "path to HADO manifest YAML")
	standardsDir := fs.String("standards-dir", "", "directory containing <id>.yaml standards (default: <manifest-dir>/standards)")
	output := fs.String("output", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	if *manifestPath == "" {
		return 2, fmt.Errorf("fire requires --manifest")
	}

	hadoManifest, err := manifest.Load(*manifestPath)
	if err != nil {
		return 2, err
	}

	stdDir := *standardsDir
	if stdDir == "" {
		stdDir = filepath.Join(filepath.Dir(*manifestPath), "standards")
	}
	standardPath, err := manifest.ResolveStandardPath(hadoManifest, *manifestPath, stdDir, *standardRef)
	if err != nil {
		return 2, err
	}
	readinessStandard, err := standard.Load(standardPath)
	if err != nil {
		return 2, err
	}

	adapterInputs := hadoManifest.CoverageAdapterInputs()
	if requiresCoverage(readinessStandard) && len(adapterInputs) == 0 {
		return 2, fmt.Errorf("fire requires evidence.coverage.inputs in manifest")
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
	applyManifestEvidence(&metrics, hadoManifest)

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

func requiresCoverage(readinessStandard standard.Standard) bool {
	return readinessStandard.RequiresGate(standard.C0CoverageGateID) || readinessStandard.RequiresGate(standard.C1CoverageGateID)
}

func applyManifestEvidence(metrics *gate.Metrics, hadoManifest manifest.Manifest) {
	if op := hadoManifest.Evidence.Operations; op != nil {
		metrics.OperationsOwner = strings.TrimSpace(op.Owner)
		metrics.OperationsRunbook = strings.TrimSpace(op.Runbook)
	}
	if obs := hadoManifest.Evidence.Observability; obs != nil {
		metrics.ObservabilitySLOPresent = manifest.ObservabilityLinksHaveURL(obs.SLOs)
		metrics.ObservabilityMonitorsPresent = manifest.ObservabilityLinksHaveURL(obs.Monitors)
		metrics.ObservabilityDashboardPresent = manifest.ObservabilityLinksHaveURL(obs.Dashboards)
	}
	if rel := hadoManifest.Evidence.Release; rel != nil {
		metrics.ReleaseRollbackPlan = strings.TrimSpace(rel.RollbackPlan)
		metrics.ReleaseAutomationDeclared = rel.AutomationDeclared()
	}
	if inf := hadoManifest.Evidence.Infra; inf != nil {
		metrics.InfraDeploymentSpec = strings.TrimSpace(inf.DeploymentSpec)
	}
}

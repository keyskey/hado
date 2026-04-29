package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/keyskey/hado/internal/coverage"
	"github.com/keyskey/hado/internal/gate"
	"github.com/keyskey/hado/internal/manifest"
	"github.com/keyskey/hado/internal/standard"
)

const version = "dev"

type stringList []string

func (values *stringList) String() string {
	return strings.Join(*values, ",")
}

func (values *stringList) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func main() {
	exitCode, err := run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func run(args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 0 {
		fmt.Fprintln(stdout, "hado: production readiness CLI")
		return 0, nil
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "hado %s\n", version)
		return 0, nil
	case "evaluate":
		return runEvaluate(args[1:], stdout, stderr)
	default:
		return 2, fmt.Errorf("unknown command %q", args[0])
	}
}

func runEvaluate(args []string, stdout, stderr io.Writer) (int, error) {
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
	adapterInputs, err := resolveCoverageInputs(coverageInputs, *manifestPath)
	if err != nil {
		return 2, err
	}
	if len(adapterInputs) == 0 {
		return 2, fmt.Errorf("evaluate requires --coverage-input or --manifest with evidence.coverage.inputs")
	}
	coverageMetrics, err := coverage.ParseAdapterInputs(adapterInputs)
	if err != nil {
		return 2, err
	}
	metrics := gate.Metrics{
		C0CoveragePercent: coverageMetrics.C0Coverage,
		C1CoveragePercent: coverageMetrics.C1Coverage,
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

func resolveCoverageInputs(coverageInputs []string, manifestPath string) ([]coverage.AdapterInput, error) {
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
	if manifestPath == "" {
		return nil, nil
	}

	hadoManifest, err := manifest.Load(manifestPath)
	if err != nil {
		return nil, err
	}
	return hadoManifest.CoverageAdapterInputs(), nil
}

func printTextEvaluation(stdout io.Writer, evaluation gate.Evaluation) {
	fmt.Fprintf(stdout, "HADO: %s\n\n", strings.ToUpper(string(evaluation.Status)))
	for _, result := range evaluation.Results {
		marker := "PASS"
		if !result.Passed {
			marker = "FAIL"
		}
		fmt.Fprintf(stdout, "- [%s] %s: %s\n", marker, result.ID, result.Message)
	}
}

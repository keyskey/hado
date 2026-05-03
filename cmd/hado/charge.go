package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/keyskey/hado/internal/charge"
)

func runCharge(args []string, stdout, stderr io.Writer) (int, error) {
	fs := flag.NewFlagSet("charge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	manifestPath := fs.String("manifest", "", "path to HADO manifest YAML")
	standardPath := fs.String("standard", "", "path to readiness standard YAML (overrides manifest standard.id)")
	apply := fs.Bool("apply", false, "run local preset to fill gaps (writes manifest and artifacts)")
	preset := fs.String("preset", charge.PresetGoGobce, "charge preset (only "+charge.PresetGoGobce+" in MVP)")
	standardsDir := fs.String("standards-dir", "", "directory containing <id>.yaml standards (default: <manifest-dir>/standards)")
	output := fs.String("output", "text", "output format: text or json")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	if *manifestPath == "" {
		return 2, fmt.Errorf("charge requires --manifest")
	}
	stdDir := *standardsDir
	if stdDir == "" {
		stdDir = filepath.Join(filepath.Dir(*manifestPath), "standards")
	}

	rep, err := charge.RunCoverage(charge.RunOptions{
		ManifestPath: *manifestPath,
		StandardPath: *standardPath,
		Preset:       *preset,
		Apply:        *apply,
	}, stdDir)
	if err != nil {
		return 2, err
	}

	switch *output {
	case "json":
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(rep); err != nil {
			return 2, fmt.Errorf("write json: %w", err)
		}
	case "text":
		printChargeText(stdout, rep)
	default:
		return 2, fmt.Errorf("unsupported output format %q", *output)
	}

	if !rep.Gap.Satisfied {
		return 1, nil
	}
	return 0, nil
}

func printChargeText(w io.Writer, rep charge.ChargeReport) {
	fmt.Fprintf(w, "CHARGE (coverage)\n\n")
	fmt.Fprintf(w, "Plan:\n")
	fmt.Fprintf(w, "  standard: %s\n", rep.Plan.StandardPath)
	fmt.Fprintf(w, "  requires test.c0_coverage: %v\n", rep.Plan.RequiresC0)
	fmt.Fprintf(w, "  requires test.c1_coverage: %v\n", rep.Plan.RequiresC1)
	if len(rep.Plan.PreferredAdapterFormats) > 0 {
		fmt.Fprintf(w, "  preferred adapters: %s\n", strings.Join(rep.Plan.PreferredAdapterFormats, ", "))
	}
	fmt.Fprintf(w, "\nGap:\n")
	if rep.Gap.Satisfied {
		fmt.Fprintf(w, "  satisfied: true\n")
		if rep.Gap.SatisfyingAdapter != "" {
			fmt.Fprintf(w, "  via adapter: %s\n", rep.Gap.SatisfyingAdapter)
			fmt.Fprintf(w, "  path: %s\n", rep.Gap.SatisfyingPath)
		}
	} else {
		fmt.Fprintf(w, "  satisfied: false\n")
		for _, it := range rep.Gap.Items {
			fmt.Fprintf(w, "  - [%s] %s\n", it.Code, it.Detail)
		}
	}
	if rep.Apply != nil {
		fmt.Fprintf(w, "\nApply:\n")
		fmt.Fprintf(w, "  coverage.out: %s\n", rep.Apply.CoverageOutPath)
		fmt.Fprintf(w, "  gobce json: %s\n", rep.Apply.GobceJSONPath)
		fmt.Fprintf(w, "  manifest updated: %v\n", rep.Apply.WroteManifest)
	}
}

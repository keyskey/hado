package charge

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/keyskey/hado/internal/coverage"
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
	fs := flag.NewFlagSet("charge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	manifestPath := fs.String("manifest", "", "path to HADO manifest YAML")
	standardRef := fs.String("standard", "", "Readiness Standard id or path (optional; defaults to manifest standard.id)")
	standardsDir := fs.String("standards-dir", "", "directory containing <id>.yaml standards (default: <manifest-dir>/standards)")
	var coverageInputs stringList
	fs.Var(&coverageInputs, "coverage-input", "coverage input as <adapter>:<path>; merges into manifest coverage inputs; adapters: hado-json, go-coverprofile, gobce-json")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	if *manifestPath == "" {
		return 2, fmt.Errorf("charge requires --manifest")
	}

	m, err := manifest.LoadOrEmpty(*manifestPath)
	if err != nil {
		return 2, err
	}

	stdDir := *standardsDir
	if stdDir == "" {
		stdDir = filepath.Join(filepath.Dir(*manifestPath), "standards")
	}
	standardPath, err := manifest.ResolveStandardPath(m, *manifestPath, stdDir, *standardRef)
	if err != nil {
		return 2, err
	}
	st, err := standard.Load(standardPath)
	if err != nil {
		return 2, err
	}

	if err := mergeCoverageInputs(&m, coverageInputs); err != nil {
		return 2, err
	}
	if requiresCoverage(st) && len(m.CoverageAdapterInputs()) == 0 {
		return 2, fmt.Errorf("charge requires --coverage-input or manifest evidence.coverage.inputs for coverage gates")
	}
	if len(m.CoverageAdapterInputs()) > 0 {
		if _, err := coverage.ParseAdapterInputs(m.CoverageAdapterInputs()); err != nil {
			return 2, err
		}
	}
	if err := m.Save(*manifestPath); err != nil {
		return 2, err
	}
	fmt.Fprintf(stdout, "Wrote manifest %s (coverage inputs: %d)\n", *manifestPath, len(m.CoverageAdapterInputs()))
	return 0, nil
}

func mergeCoverageInputs(m *manifest.Manifest, specs []string) error {
	if len(specs) == 0 {
		return nil
	}
	if m.Evidence.Coverage == nil {
		m.Evidence.Coverage = &manifest.CoverageEvidence{}
	}

	existing := make(map[string]struct{}, len(m.Evidence.Coverage.Inputs))
	for _, input := range m.Evidence.Coverage.Inputs {
		key := input.Adapter + "\x00" + input.Path
		existing[key] = struct{}{}
	}
	for _, spec := range specs {
		parsed, err := coverage.ParseCoverageSpec(spec)
		if err != nil {
			return err
		}
		key := parsed.Format + "\x00" + parsed.Path
		if _, ok := existing[key]; ok {
			continue
		}
		m.Evidence.Coverage.Inputs = append(m.Evidence.Coverage.Inputs, manifest.CoverageInput{
			Adapter: parsed.Format,
			Path:    parsed.Path,
		})
		existing[key] = struct{}{}
	}
	return nil
}

func requiresCoverage(st standard.Standard) bool {
	return st.RequiresGate(standard.C0CoverageGateID) || st.RequiresGate(standard.C1CoverageGateID)
}

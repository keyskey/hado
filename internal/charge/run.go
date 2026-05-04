package charge

import (
	"fmt"

	"github.com/keyskey/hado/internal/manifest"
	"github.com/keyskey/hado/internal/standard"
)

// RunCoverage executes plan → gap; if opts.Apply and the gap is not satisfied, runs ApplyCoverageGoGobce then re-loads the manifest and runs gap again.
func RunCoverage(opts RunOptions, standardsDir string) (ChargeReport, error) {
	var rep ChargeReport
	if opts.ManifestPath == "" {
		return rep, fmt.Errorf("manifest path is required")
	}
	m, err := manifest.Load(opts.ManifestPath)
	if err != nil {
		return rep, err
	}

	stdPath, err := manifest.ResolveStandardPath(m, opts.ManifestPath, standardsDir, opts.StandardPath)
	if err != nil {
		return rep, err
	}
	st, err := standard.Load(stdPath)
	if err != nil {
		return rep, err
	}

	rep.Plan = PlanCoverage(stdPath, st)
	rep.Gap = GapCoverage(m, rep.Plan)

	if rep.Gap.Satisfied || !opts.Apply {
		return rep, nil
	}

	switch opts.Preset {
	case "", PresetGoGobce:
	default:
		return rep, fmt.Errorf("unknown preset %q (supported: %s)", opts.Preset, PresetGoGobce)
	}

	mp := m
	res, aerr := ApplyCoverageGoGobce(&mp, opts.ManifestPath)
	if aerr != nil {
		return rep, aerr
	}
	if err := mp.Save(opts.ManifestPath); err != nil {
		return rep, err
	}
	rep.Apply = &res

	m2, err := manifest.Load(opts.ManifestPath)
	if err != nil {
		return rep, err
	}
	rep.Gap = GapCoverage(m2, rep.Plan)
	return rep, nil
}

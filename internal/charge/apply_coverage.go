package charge

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/keyskey/hado/internal/coverage"
	"github.com/keyskey/hado/internal/manifest"
)

const applyCommandTimeout = 10 * time.Minute

// ApplyCoverageGoGobce runs go test -coverprofile and gobce analyze, then sets manifest coverage inputs to gobce-json.
func ApplyCoverageGoGobce(m *manifest.Manifest, manifestPath string) (ApplyCoverageResult, error) {
	var out ApplyCoverageResult
	dir := filepath.Dir(manifestPath)
	coverageOut := filepath.Join(dir, "coverage.out")
	gobceJSON := filepath.Join(dir, "hado-coverage.json")
	out.CoverageOutPath = coverageOut
	out.GobceJSONPath = gobceJSON

	ctx, cancel := context.WithTimeout(context.Background(), applyCommandTimeout)
	defer cancel()

	goTest := exec.CommandContext(ctx, "go", "test", "./...", "-coverprofile", coverageOut)
	goTest.Dir = dir
	goTest.Stdout = os.Stdout
	goTest.Stderr = os.Stderr
	if err := goTest.Run(); err != nil {
		return out, fmt.Errorf("go test -coverprofile: %w", err)
	}

	gobce := exec.CommandContext(ctx, "gobce", "analyze", "--coverprofile", coverageOut, "--format", "json", "--output", gobceJSON)
	gobce.Dir = dir
	gobce.Stdout = os.Stdout
	gobce.Stderr = os.Stderr
	if err := gobce.Run(); err != nil {
		return out, fmt.Errorf("gobce analyze: %w", err)
	}

	relGobce, err := filepath.Rel(dir, gobceJSON)
	if err != nil {
		relGobce = filepath.Base(gobceJSON)
	}
	m.Evidence.Coverage.Inputs = []manifest.CoverageInput{{
		Adapter: coverage.FormatGobceJSON,
		Path:    relGobce,
	}}
	out.WroteManifest = true
	return out, nil
}

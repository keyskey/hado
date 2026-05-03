package charge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/keyskey/hado/internal/manifest"
	"github.com/keyskey/hado/internal/standard"
)

func TestPlanCoverageRequiresC0C1(t *testing.T) {
	st := standard.Standard{ID: "x", Gates: []standard.Gate{
		{ID: standard.C0CoverageGateID, Threshold: standard.Threshold{Min: ptr(80.0)}},
		{ID: standard.C1CoverageGateID, Threshold: standard.Threshold{Min: ptr(70.0)}},
	}}
	p := PlanCoverage("/tmp/std.yaml", st)
	if !p.RequiresC0 || !p.RequiresC1 {
		t.Fatalf("plan = %+v", p)
	}
	if len(p.PreferredAdapterFormats) < 2 {
		t.Fatalf("expected preferred adapters, got %+v", p.PreferredAdapterFormats)
	}
}

func TestGapCoverageSatisfiedWithHadoJSON(t *testing.T) {
	dir := t.TempDir()
	mPath := filepath.Join(dir, "hado.yaml")
	metricsPath := filepath.Join(dir, "m.json")
	if err := os.WriteFile(metricsPath, []byte(`{"c0Coverage": 90, "c1Coverage": 75}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mPath, []byte(`version: v1
evidence:
  coverage:
    inputs:
      - adapter: hado-json
        path: m.json
`), 0o600); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Load(mPath)
	if err != nil {
		t.Fatal(err)
	}
	st := standard.Standard{ID: "x", Gates: []standard.Gate{
		{ID: standard.C0CoverageGateID, Threshold: standard.Threshold{Min: ptr(80.0)}},
		{ID: standard.C1CoverageGateID, Threshold: standard.Threshold{Min: ptr(70.0)}},
	}}
	plan := PlanCoverage("std.yaml", st)
	gap := GapCoverage(m, plan)
	if !gap.Satisfied {
		t.Fatalf("gap = %+v", gap)
	}
}

func TestResolveStandardPathLogicalID(t *testing.T) {
	m := manifest.Manifest{}
	m.Standard.ID = "web-service"
	p, err := ResolveStandardPath(m, "/repo/hado.yaml", "/repo/standards", "")
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join("/repo/standards", "web-service.yaml"); p != want {
		t.Fatalf("got %q want %q", p, want)
	}
}

func TestChargeReportJSON(t *testing.T) {
	rep := ChargeReport{
		Plan: CoveragePlan{StandardPath: "/s.yaml", RequiresC0: true},
		Gap:  CoverageGapReport{Satisfied: true},
	}
	b, err := json.Marshal(rep)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(b) {
		t.Fatalf("invalid json: %s", string(b))
	}
}

func ptr(f float64) *float64 {
	return &f
}

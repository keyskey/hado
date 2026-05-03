package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestEvaluateReadyWithInfraDeploymentSpec(t *testing.T) {
	dir := t.TempDir()
	standardPath := writeFile(t, dir, "standard.yaml", `id: test
gates:
  - id: infra.deployment_spec_exists
    required: true
`)
	manifestPath := writeFile(t, dir, "hado.yaml", `version: v1
evidence:
  infra:
    deployment_spec: k8s/order-api.yaml
`)

	var stdout, stderr bytes.Buffer
	exitCode, err := run([]string{
		"evaluate",
		"--standard", standardPath,
		"--manifest", manifestPath,
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run evaluate: %v", err)
	}
	if exitCode != 0 {
		t.Fatalf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout.String(), "infra.deployment_spec_exists") {
		t.Fatalf("stdout = %q, want infra.deployment_spec_exists in output", stdout.String())
	}
}

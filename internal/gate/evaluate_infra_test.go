package gate

import (
	"testing"

	"github.com/keyskey/hado/internal/standard"
)

func TestEvaluateInfraDeploymentSpecReady(t *testing.T) {
	evaluation, err := Evaluate(standard.Standard{
		ID: "test-standard",
		Gates: []standard.Gate{
			{ID: standard.InfraDeploymentSpecExistsGateID, Required: true},
		},
	}, Metrics{
		InfraDeploymentSpec: "k8s/deployment.yaml",
	})
	if err != nil {
		t.Fatalf("Evaluate: %v", err)
	}
	if evaluation.Status != DecisionReady {
		t.Fatalf("status = %q, want %q", evaluation.Status, DecisionReady)
	}
}

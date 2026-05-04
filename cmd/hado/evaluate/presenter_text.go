package evaluate

import (
	"fmt"
	"io"
	"strings"

	"github.com/keyskey/hado/internal/gate"
	"github.com/keyskey/hado/internal/standard"
)

func printTextEvaluation(stdout io.Writer, evaluation gate.Evaluation) {
	for _, result := range evaluation.Results {
		marker := "PASS"
		if !result.Passed {
			marker = "FAIL"
		}
		severity := string(result.Severity)
		if severity == "" {
			severity = string(standard.SeverityMinor)
		}
		fmt.Fprintf(stdout, "- [%s] %s (severity: %s): %s", marker, result.ID, severity, result.Message)
		if !result.Passed {
			fmt.Fprintf(stdout, " [%s]", releaseActionHint(result))
		}
		fmt.Fprintln(stdout)
	}
	statusLine := fmt.Sprintf("HADO: %s", strings.ToUpper(string(evaluation.Status)))
	fmt.Fprintf(stdout, "\n%s\n", statusLine)
}

func releaseActionHint(result gate.Result) string {
	if !result.Required {
		return "optional gate: release allowed"
	}
	switch result.Severity {
	case standard.SeverityCritical:
		return "release blocked: fix before release"
	case standard.SeverityMajor:
		return "release allowed: fix soon after release"
	default:
		return "release allowed: fix when appropriate"
	}
}

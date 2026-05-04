package evaluate

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/keyskey/hado/internal/gate"
	"github.com/keyskey/hado/internal/standard"
	"golang.org/x/term"
)

const (
	ansiReset       = "\033[0m"
	ansiGreen       = "\033[32m"
	ansiYellow      = "\033[33m"
	ansiBoldRed     = "\033[1;31m"
	ansiBoldGreen   = "\033[1;32m"
	ansiBoldMagenta = "\033[1;35m"
)

func printTextEvaluation(stdout io.Writer, evaluation gate.Evaluation) {
	useColor := shouldColorize(stdout)
	for _, result := range evaluation.Results {
		marker := "PASS"
		if !result.Passed {
			marker = "FAIL"
		}
		markerDisplay := colorizedMarker(result, marker, useColor)
		severity := string(result.Severity)
		if severity == "" {
			severity = string(standard.SeverityMinor)
		}
		severityDisplay := colorizedSeverity(result, severity, useColor)
		fmt.Fprintf(stdout, "- [%s] %s (severity: %s): %s", markerDisplay, result.ID, severityDisplay, result.Message)
		if !result.Passed {
			fmt.Fprintf(stdout, " [%s]", releaseActionHint(result))
		}
		fmt.Fprintln(stdout)
	}
	statusLine := fmt.Sprintf("HADO: %s", strings.ToUpper(string(evaluation.Status)))
	fmt.Fprintf(stdout, "\n%s\n", colorizedSummary(statusLine, evaluation.Status, useColor))
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

func shouldColorize(stdout io.Writer) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	f, ok := stdout.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

func colorizedMarker(result gate.Result, marker string, useColor bool) string {
	if !useColor {
		return marker
	}
	if result.Passed {
		return wrapColor(marker, ansiGreen)
	}
	if result.Required && result.Severity == standard.SeverityCritical {
		return wrapColor(marker, ansiBoldRed)
	}
	return wrapColor(marker, ansiYellow)
}

func colorizedSeverity(result gate.Result, severity string, useColor bool) string {
	if !useColor {
		return severity
	}
	if result.Passed {
		return wrapColor(severity, ansiGreen)
	}
	if result.Required && result.Severity == standard.SeverityCritical {
		return wrapColor(severity, ansiBoldRed)
	}
	return wrapColor(severity, ansiYellow)
}

func colorizedSummary(statusLine string, status gate.DecisionStatus, useColor bool) string {
	if !useColor {
		return statusLine
	}
	switch status {
	case gate.DecisionReady:
		return wrapColor(statusLine, ansiBoldGreen)
	case gate.DecisionBlocked:
		return wrapColor(statusLine, ansiBoldRed)
	default:
		return wrapColor(statusLine, ansiBoldMagenta)
	}
}

func wrapColor(text, color string) string {
	return color + text + ansiReset
}

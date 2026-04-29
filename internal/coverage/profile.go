// Package coverage parses Go coverage evidence.
package coverage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Summary contains C0 coverage metrics derived from a Go coverage profile.
type Summary struct {
	C0Coverage        float64 `json:"c0Coverage"`
	CoveredStatements int     `json:"coveredStatements"`
	TotalStatements   int     `json:"totalStatements"`
}

// ParseGoProfile reads a Go coverprofile and returns C0 coverage percent.
func ParseGoProfile(path string) (Summary, error) {
	file, err := os.Open(path)
	if err != nil {
		return Summary{}, fmt.Errorf("open coverage profile: %w", err)
	}
	defer file.Close()

	var coveredStatements int
	var totalStatements int
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if lineNumber == 1 {
			if !strings.HasPrefix(line, "mode:") {
				return Summary{}, fmt.Errorf("coverage profile line 1 must declare mode")
			}
			continue
		}

		statements, count, err := parseProfileLine(line)
		if err != nil {
			return Summary{}, fmt.Errorf("coverage profile line %d: %w", lineNumber, err)
		}
		totalStatements += statements
		if count > 0 {
			coveredStatements += statements
		}
	}
	if err := scanner.Err(); err != nil {
		return Summary{}, fmt.Errorf("scan coverage profile: %w", err)
	}
	if totalStatements == 0 {
		return Summary{}, fmt.Errorf("coverage profile has no statements")
	}

	return Summary{
		C0Coverage:        float64(coveredStatements) / float64(totalStatements) * 100,
		CoveredStatements: coveredStatements,
		TotalStatements:   totalStatements,
	}, nil
}

func parseProfileLine(line string) (statements int, count int, err error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return 0, 0, fmt.Errorf("expected block, statement count, and execution count")
	}

	statements, err = strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, fmt.Errorf("parse statement count: %w", err)
	}
	count, err = strconv.Atoi(fields[2])
	if err != nil {
		return 0, 0, fmt.Errorf("parse execution count: %w", err)
	}
	if statements < 0 {
		return 0, 0, fmt.Errorf("statement count must be non-negative")
	}
	if count < 0 {
		return 0, 0, fmt.Errorf("execution count must be non-negative")
	}

	return statements, count, nil
}

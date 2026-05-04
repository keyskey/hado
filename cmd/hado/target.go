package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/keyskey/hado/internal/manifest"
	"github.com/keyskey/hado/internal/standard"
	"golang.org/x/term"
)

func runTarget(args []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	fs := flag.NewFlagSet("target", flag.ContinueOnError)
	fs.SetOutput(stderr)
	manifestPath := fs.String("manifest", "", "path to HADO manifest YAML (created if missing)")
	serviceName := fs.String("service-name", "", "service display name")
	serviceID := fs.String("service-id", "", "service id slug (defaults to service-name when empty after merge)")
	standardID := fs.String("standard-id", "", "Readiness Standard id (e.g. web-service) or path to standard YAML")
	standardsDir := fs.String("standards-dir", "", "directory containing <id>.yaml standards (default: <manifest-dir>/standards)")
	rewritePlaceholders := fs.Bool("rewrite-placeholders", true, "add evidence placeholders for gates declared in the resolved standard (merge with existing)")
	if err := fs.Parse(args); err != nil {
		return 2, err
	}
	if *manifestPath == "" {
		return 2, fmt.Errorf("target requires --manifest")
	}

	m, err := manifest.LoadOrEmpty(*manifestPath)
	if err != nil {
		return 2, err
	}

	name := strings.TrimSpace(*serviceName)
	id := strings.TrimSpace(*serviceID)
	stdID := strings.TrimSpace(*standardID)

	useTTY := false
	if f, ok := stdin.(*os.File); ok {
		useTTY = term.IsTerminal(int(f.Fd()))
	}

	allFlagsEmpty := name == "" && id == "" && stdID == ""

	if useTTY && allFlagsEmpty {
		reader := bufio.NewReader(stdin)
		var err error
		name, err = promptLine(reader, stdout, "Service name", strings.TrimSpace(m.Service.Name))
		if err != nil {
			return 2, err
		}
		defID := strings.TrimSpace(m.Service.ID)
		if defID == "" {
			defID = name
		}
		id, err = promptOptionalLine(reader, stdout, "Service id (optional, Enter = same as name)", defID, name)
		if err != nil {
			return 2, err
		}
		stdID, err = promptLine(reader, stdout, "Readiness standard id or path", strings.TrimSpace(m.Standard.ID))
		if err != nil {
			return 2, err
		}
	} else if !useTTY {
		if allFlagsEmpty {
			return 2, fmt.Errorf("target: non-interactive mode requires at least one of --service-name, --service-id, --standard-id")
		}
		if name == "" {
			name = strings.TrimSpace(m.Service.Name)
		}
		if id == "" {
			id = strings.TrimSpace(m.Service.ID)
		}
		if stdID == "" {
			stdID = strings.TrimSpace(m.Standard.ID)
		}
	} else {
		if name == "" {
			name = strings.TrimSpace(m.Service.Name)
		}
		if id == "" {
			id = strings.TrimSpace(m.Service.ID)
		}
		if stdID == "" {
			stdID = strings.TrimSpace(m.Standard.ID)
		}
	}

	name = strings.TrimSpace(name)
	id = strings.TrimSpace(id)
	stdID = strings.TrimSpace(stdID)

	if stdID == "" {
		return 2, fmt.Errorf("target: readiness standard id is required (set in manifest or pass --standard-id)")
	}
	if name == "" && id == "" {
		return 2, fmt.Errorf("target: service name or service id is required (set in manifest or pass --service-name / --service-id)")
	}
	if id == "" {
		id = name
	}
	if name == "" {
		name = id
	}

	if m.Version == "" {
		m.Version = "v1"
	}
	m.Service.Name = name
	m.Service.ID = id
	m.Standard.ID = stdID

	stdDir := *standardsDir
	if stdDir == "" {
		stdDir = filepath.Join(filepath.Dir(*manifestPath), "standards")
	}

	if *rewritePlaceholders {
		stdPath, err := manifest.ResolveStandardPath(m, *manifestPath, stdDir, "")
		if err != nil {
			return 2, err
		}
		st, err := standard.Load(stdPath)
		if err != nil {
			return 2, fmt.Errorf("load standard for placeholders: %w", err)
		}
		manifest.ApplyEvidencePlaceholders(&m, st, manifest.ApplyEvidencePlaceholdersOptions{MergeOnly: true})
	}

	if err := m.Save(*manifestPath); err != nil {
		return 2, err
	}
	fmt.Fprintf(stdout, "Wrote manifest %s (service %q, standard %q", *manifestPath, m.Service.Name, m.Standard.ID)
	if *rewritePlaceholders {
		fmt.Fprintf(stdout, ", evidence placeholders merged from standard")
	}
	fmt.Fprintln(stdout, ")")
	return 0, nil
}

func promptLine(reader *bufio.Reader, stdout io.Writer, label, defaultValue string) (string, error) {
	for {
		if defaultValue != "" {
			fmt.Fprintf(stdout, "%s [%s]: ", label, defaultValue)
		} else {
			fmt.Fprintf(stdout, "%s: ", label)
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("read %s: %w", label, err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if defaultValue != "" {
				return defaultValue, nil
			}
			fmt.Fprintf(stdout, "Value is required.\n")
			continue
		}
		return line, nil
	}
}

func promptOptionalLine(reader *bufio.Reader, stdout io.Writer, label, defaultIfEmpty, fallbackSameAs string) (string, error) {
	fmt.Fprintf(stdout, "%s [%s]: ", label, defaultIfEmpty)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read %s: %w", label, err)
	}
	line = strings.TrimSpace(line)
	if line != "" {
		return line, nil
	}
	if defaultIfEmpty != "" {
		return defaultIfEmpty, nil
	}
	return fallbackSameAs, nil
}

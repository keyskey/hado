package main

import (
	"fmt"
	"io"
	"os"

	chargecmd "github.com/keyskey/hado/cmd/hado/charge"
	firecmd "github.com/keyskey/hado/cmd/hado/fire"
	targetcmd "github.com/keyskey/hado/cmd/hado/target"
)

const version = "dev"

func main() {
	exitCode, err := run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func run(args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 0 {
		fmt.Fprintln(stdout, "hado: production readiness CLI")
		return 0, nil
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Fprintf(stdout, "hado %s\n", version)
		return 0, nil
	case "charge":
		return chargecmd.Run(args[1:], stdout, stderr)
	case "fire":
		return firecmd.Run(args[1:], stdout, stderr)
	case "target":
		return targetcmd.Run(args[1:], os.Stdin, stdout, stderr)
	default:
		return 2, fmt.Errorf("unknown command %q", args[0])
	}
}

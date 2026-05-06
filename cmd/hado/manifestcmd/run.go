package manifestcmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/keyskey/hado/internal/manifest"
)

// Run handles "hado manifest doc" and related subcommands.
func Run(args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 0 {
		return 2, errors.New("usage: hado manifest doc [--out path]  # writes commented reference YAML")
	}
	if args[0] != "doc" {
		return 2, fmt.Errorf("unknown manifest subcommand %q (try: doc)", args[0])
	}
	fs := flag.NewFlagSet("manifest doc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	outPath := fs.String("out", "", "write reference YAML to this file (default: stdout)")
	if err := fs.Parse(args[1:]); err != nil {
		return 2, err
	}

	var w io.Writer = stdout
	var f *os.File
	if *outPath != "" {
		var err error
		f, err = os.Create(*outPath)
		if err != nil {
			return 2, err
		}
		defer f.Close()
		w = f
	}
	if err := manifest.WriteManifestReferenceYAML(w); err != nil {
		return 2, err
	}
	return 0, nil
}

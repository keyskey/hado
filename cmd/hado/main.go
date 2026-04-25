package main

import (
	"fmt"
	"os"
)

const version = "dev"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version", "--version", "-v":
			fmt.Printf("hado %s\n", version)
			return
		}
	}

	fmt.Fprintln(os.Stdout, "hado: production readiness CLI")
}

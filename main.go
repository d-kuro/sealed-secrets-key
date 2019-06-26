package main

import (
	"os"

	"github.com/d-kuro/sealed-secrets-key/cmd"
)

func main() {
	os.Exit(cmd.Execute(os.Stdout, os.Stderr))
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.0.1"

var Revision = "development"

func NewVersionCmd(o *GenerateOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			runVersionCmd(o)
		},
	}
}

func runVersionCmd(o *GenerateOptions) {
	fmt.Fprintf(o.outStream, "version: %s (rev: %s)\n", Version, Revision)
}

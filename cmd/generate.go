package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/d-kuro/sealed-secrets-key/pkg/generate"
	"github.com/spf13/cobra"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 1
)

type GenerateOptions struct {
	Name      string
	Namespace string
}

func Execute(outstream, errstream io.Writer) int {
	o := NewGenerateOptions()
	cmd := NewRootCommand(o)
	addCommands(cmd)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(errstream, "error: %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}

func addCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		NewVersionCmd(),
	)
}

func NewGenerateOptions() *GenerateOptions {
	return &GenerateOptions{}
}

func NewRootCommand(option *GenerateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sealed-secrets-key",
		Short: "Generate sealed-secrets-key",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return func() error {
				if len(args) != 1 {
					return errors.New("invalid arguments")
				}
				return nil
			}()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], option)
		},
	}

	fset := cmd.Flags()

	fset.StringVar(&option.Name, "name", "sealed-secrets-key", "Secret name")
	fset.StringVar(&option.Namespace, "namespace", "kube-system", "Namespace")

	return cmd
}

func run(output string, o *GenerateOptions) error {
	return generate.Generate(o.Name, o.Namespace, output)
}

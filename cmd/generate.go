package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/d-kuro/sealed-secrets-key/pkg/generate"
	"github.com/spf13/cobra"
)

const (
	exitCodeOK  = 0
	exitCodeErr = 1
)

const example = `
$ sealed-secrets-key -o secret.yaml`

type GenerateOptions struct {
	Name      string
	Namespace string
	Output    string
	KeySize   int
	KeyTTL    time.Duration
	outStream io.Writer
	errStream io.Writer
}

func Execute(outStream, errStream io.Writer) int {
	o := NewGenerateOptions(outStream, errStream)
	cmd := NewRootCommand(o)
	addCommands(cmd, o)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(errStream, "error: %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}

func addCommands(rootCmd *cobra.Command, o *GenerateOptions) {
	rootCmd.AddCommand(
		NewVersionCmd(o),
	)
}

func NewGenerateOptions(outStream, errStream io.Writer) *GenerateOptions {
	return &GenerateOptions{
		outStream: outStream,
		errStream: errStream,
	}
}

func NewRootCommand(o *GenerateOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "sealed-secrets-key",
		Short:         "Generate sealed-secrets-key",
		Example:       example,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(o)
		},
	}

	fset := cmd.Flags()

	fset.StringVar(&o.Name, "secret-name", "sealed-secrets-key", "Name of Secret containing public/private key.")
	fset.StringVar(&o.Namespace, "namespace", "kube-system", "Namespace of Secret.")
	fset.StringVarP(&o.Output, "output", "o", "", "Output file.")
	fset.IntVar(&o.KeySize, "key-size", 4096, "Size of encryption key.")
	fset.DurationVar(&o.KeyTTL, "key-ttl", 10*365*24*time.Hour, "Duration that certificate is valid for.")

	return cmd
}

func run(o *GenerateOptions) error {
	out, err := generate.Generate(o.Name, o.Namespace, o.KeySize, o.KeyTTL)
	if err != nil {
		return err
	}

	if o.Output != "" {
		f, err := os.OpenFile(o.Output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
		if err != nil {
			if os.IsExist(err) {
				return errors.New("file already exists")
			}
			return err
		}
		defer f.Close()
		f.Write(out)
		return nil
	}

	fmt.Fprint(o.outStream, string(out))
	return nil
}

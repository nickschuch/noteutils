package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"

	"github.com/nickschuch/noteutils/cmd/noteutils/validate"
)

const cmdExample = `
  # Validate notes configuration using a spec.
  noteutils validate --spec ./spec.yaml /path/to/file`

var cmd = &cobra.Command{
	Use:     "noteutils",
	Short:   "Utilities for working with ELF notes (readelf -n)",
	Example: cmdExample,
}

func main() {
	cmd.AddCommand(validate.NewCommand())

	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

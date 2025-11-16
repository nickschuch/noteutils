package validate

import (
	"fmt"
	"os"
	"runtime"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/nickschuch/noteutils/internal/notes"
)

const cmdExample = `
  # Validate notes configuration using a spec.
  noteutils validate --spec ./spec.yaml /path/to/file`

// Options holds the command line options for the validate command
type Options struct {
	Spec string
}

// NewCommand creates a new cobra.Command for 'validate' sub command
func NewCommand() *cobra.Command {
	var options Options

	cmd := &cobra.Command{
		Use:                   "validate <path to file>",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "Validate notes configuration using a spec",
		Long:                  cmdExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				binary = args[0]
				arch   = runtime.GOARCH
			)

			fmt.Println("Validating", binary, "for arch:", arch, "using spec:", options.Spec)

			have, err := notes.LoadFromFile(binary)
			if err != nil {
				return fmt.Errorf("loading notes from file: %w", err)
			}

			spec, err := loadSpecification(options.Spec)
			if err != nil {
				return fmt.Errorf("failed to load specification: %w", err)
			}

			if _, ok := spec.Probes[arch]; !ok {
				return fmt.Errorf("architecture not found in spec file: %s", arch)
			}

			if diff := cmp.Diff(have, spec.Probes[arch], cmpopts.SortSlices(func(a, b notes.Probe) bool {
				return a.Name < b.Name
			})); diff != "" {
				fmt.Println(diff)
				return fmt.Errorf("mistmatch found")
			}

			fmt.Println("Specification matches notes")

			return nil
		},
	}

	cmd.Flags().StringVar(&options.Spec, "spec", "", "Path to the spec file for validation")

	return cmd
}

func loadSpecification(path string) (*Specification, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var spec Specification

	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &spec, nil
}

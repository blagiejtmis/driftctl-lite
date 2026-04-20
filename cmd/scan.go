package cmd

import (
	"fmt"

	"github.com/driftctl-lite/drift"
	"github.com/spf13/cobra"
)

var (
	statePath string
	outputFmt string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for drift between live AWS resources and a Terraform state file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateOutputFmt(outputFmt); err != nil {
			return err
		}
		result, err := drift.Scan(statePath)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		result.Print(outputFmt)
		return nil
	},
}

func init() {
	scanCmd.Flags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	scanCmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "Output format: text or json")
}

// validateOutputFmt returns an error if the provided output format is not supported.
func validateOutputFmt(fmt string) error {
	switch fmt {
	case "text", "json":
		return nil
	default:
		return fmt.Errorf("unsupported output format %q: must be \"text\" or \"json\"", fmt)
	}
}

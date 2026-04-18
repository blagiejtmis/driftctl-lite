package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var statePath string
	var outputPath string

	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Save current managed resources as a drift baseline",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(statePath)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			report := drift.Compare(resources, resources)

			if err := drift.SaveBaseline(outputPath, report); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Baseline saved to %s (%d resources)\n", outputPath, len(report.Managed))
			return nil
		},
	}

	baselineCmd.Flags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	baselineCmd.Flags().StringVarP(&outputPath, "output", "o", "baseline.json", "Path to write baseline file")

	rootCmd.AddCommand(baselineCmd)
}

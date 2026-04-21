package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var tfDir string

	rollupCmd := &cobra.Command{
		Use:   "rollup",
		Short: "Show per-type drift rollup summary",
		Long:  "Aggregates drift results by resource type and prints a summary table.",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			iacResources, err := drift.Scan(tfDir)
			if err != nil {
				// tolerate missing IaC dir — treat as empty
				iacResources = nil
			}

			report := drift.Compare(resources, iacResources)
			entries := drift.Rollup(report)
			drift.FprintRollup(os.Stdout, entries)

			if drift.RollupHasDrift(entries) {
				os.Exit(1)
			}
			return nil
		},
	}

	rollupCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate",
		"Path to the Terraform state file")
	rollupCmd.Flags().StringVarP(&tfDir, "iac", "i", ".",
		"Path to IaC definitions (directory or state file)")

	rootCmd.AddCommand(rollupCmd)
}

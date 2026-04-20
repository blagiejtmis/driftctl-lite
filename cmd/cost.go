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

	costCmd := &cobra.Command{
		Use:   "cost",
		Short: "Estimate monthly costs for managed and untracked resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			report := drift.Compare(resources)

			costReport := drift.EstimateCosts(report)
			drift.FprintCost(os.Stdout, costReport)

			if costReport.TotalUntracked > 0 {
				fmt.Fprintf(os.Stderr, "\nWARN: untracked resources contribute $%.2f/mo\n", costReport.TotalUntracked)
			}
			return nil
		},
	}

	costCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	costCmd.Flags().StringVarP(&tfDir, "dir", "d", ".", "Directory containing .tf files")

	rootCmd.AddCommand(costCmd)
}

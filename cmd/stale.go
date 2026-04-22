package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var thresholdDays int

	staleCmd := &cobra.Command{
		Use:   "stale",
		Short: "Detect resources that have not been modified recently",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			report := drift.EvaluateStale(resources, thresholdDays, time.Now())
			drift.FprintStale(os.Stdout, report)

			if drift.StaleHasEntries(report) {
				os.Exit(1)
			}
			return nil
		},
	}

	staleCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate",
		"Path to the Terraform state file")
	staleCmd.Flags().IntVarP(&thresholdDays, "days", "d", 90,
		"Age threshold in days to consider a resource stale")

	rootCmd.AddCommand(staleCmd)
}

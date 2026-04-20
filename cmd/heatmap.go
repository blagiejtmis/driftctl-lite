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

	heatmapCmd := &cobra.Command{
		Use:   "heatmap",
		Short: "Show drift intensity heatmap grouped by resource type",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			report := drift.Compare(resources)
			heatmap := drift.BuildHeatmap(report)
			drift.FprintHeatmap(os.Stdout, heatmap)

			if anyDrifted(heatmap) {
				os.Exit(1)
			}
			return nil
		},
	}

	heatmapCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate",
		"Path to Terraform state file")
	heatmapCmd.Flags().StringVarP(&tfDir, "dir", "d", ".",
		"Directory containing Terraform IaC definitions")

	rootCmd.AddCommand(heatmapCmd)
}

func anyDrifted(h drift.HeatmapResult) bool {
	for _, e := range h.Entries {
		if e.Drifted > 0 {
			return true
		}
	}
	return false
}

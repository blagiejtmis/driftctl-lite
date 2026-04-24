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

	classifyCmd := &cobra.Command{
		Use:   "classify",
		Short: "Classify drifted resources by security severity",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			report := drift.Compare(resources)
			cr := drift.Classify(report)
			drift.FprintClassify(os.Stdout, cr)
			if drift.ClassifyHasCritical(cr) {
				os.Exit(1)
			}
			return nil
		},
	}

	classifyCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	classifyCmd.Flags().StringVarP(&tfDir, "dir", "d", ".", "Directory containing Terraform configs")

	rootCmd.AddCommand(classifyCmd)
}

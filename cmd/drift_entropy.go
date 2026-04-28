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
	var exitCode bool

	entropyCmd := &cobra.Command{
		Use:   "entropy",
		Short: "Compute Shannon entropy of drift distribution per resource type",
		Long: `Analyses how evenly drift is distributed across resource types.

A high entropy value indicates drift is spread across many types;
a low value means drift is concentrated in a few types.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			var baseline []drift.Resource
			report := drift.Compare(resources, baseline)
			entropyReport := drift.BuildDriftEntropy(report)

			drift.FprintEntropy(os.Stdout, entropyReport)

			if exitCode && drift.EntropyHasDrift(entropyReport) {
				os.Exit(1)
			}
			return nil
		},
	}

	entropyCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate",
		"Path to Terraform state file")
	entropyCmd.Flags().StringVarP(&tfDir, "tf-dir", "t", ".",
		"Directory containing Terraform definitions")
	entropyCmd.Flags().BoolVar(&exitCode, "exit-code", false,
		"Exit with code 1 when drift entropy is non-zero")

	rootCmd.AddCommand(entropyCmd)
}

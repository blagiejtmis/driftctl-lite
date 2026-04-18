package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var policyFile string
	var tfDir string

	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Evaluate drift report against a policy file",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			report := drift.Compare(resources.Live, resources.IaC)

			pf, err := drift.LoadPolicy(policyFile)
			if err != nil {
				return fmt.Errorf("load policy: %w", err)
			}

			results := drift.EvaluatePolicy(pf, report)
			drift.FprintPolicy(os.Stdout, results)

			if drift.PolicyHasErrors(results) {
				os.Exit(1)
			}
			return nil
		},
	}

	policyCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	policyCmd.Flags().StringVarP(&tfDir, "dir", "d", ".", "Directory containing .tf files")
	policyCmd.Flags().StringVarP(&policyFile, "policy", "p", "policy.json", "Path to policy JSON file")

	rootCmd.AddCommand(policyCmd)
}

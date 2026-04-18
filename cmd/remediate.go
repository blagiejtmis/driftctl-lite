package cmd

import (
	"fmt"
	"os"

	"github.com/example/driftctl-lite/drift"
	"github.com/spf13/cobra"
)

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Show remediation hints for detected drift",
	RunE: func(cmd *cobra.Command, args []string) error {
		statePath, _ := cmd.Flags().GetString("state")
		includeTypes, _ := cmd.Flags().GetStringSlice("include-types")
		excludeTypes, _ := cmd.Flags().GetStringSlice("exclude-types")

		resources, err := drift.Scan(statePath)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		filtered := drift.Filter(resources, drift.FilterOptions{
			IncludeTypes: includeTypes,
			ExcludeTypes: excludeTypes,
		})

		report := drift.Compare(filtered, nil)
		hints := drift.Remediate(report)
		drift.FprintRemediation(os.Stdout, hints)
		return nil
	},
}

func init() {
	remediateCmd.Flags().String("state", "terraform.tfstate", "Path to Terraform state file")
	remediateCmd.Flags().StringSlice("include-types", nil, "Only include these resource types")
	remediateCmd.Flags().StringSlice("exclude-types", nil, "Exclude these resource types")
	RootCmd.AddCommand(remediateCmd)
}

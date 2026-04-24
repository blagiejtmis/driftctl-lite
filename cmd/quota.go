package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var quotaStateFile string
var quotaUseDefaults bool

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Evaluate resource counts against quota limits",
	RunE: func(cmd *cobra.Command, args []string) error {
		resources, err := drift.Scan(quotaStateFile)
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}

		limits := drift.DefaultQuotaLimits()

		report := drift.EvaluateQuota(resources, limits)
		drift.FprintQuota(os.Stdout, report)

		if drift.QuotaHasViolations(report) {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	quotaCmd.Flags().StringVarP(&quotaStateFile, "state", "s", "terraform.tfstate",
		"Path to Terraform state file")
	quotaCmd.Flags().BoolVar(&quotaUseDefaults, "defaults", true,
		"Use built-in default quota limits")
	rootCmd.AddCommand(quotaCmd)
}

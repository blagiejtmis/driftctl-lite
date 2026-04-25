package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var (
	budgetMaxMissing   int
	budgetMaxUntracked int
	budgetState        string
	budgetTFState      string
)

var budgetCmd = &cobra.Command{
	Use:   "budget",
	Short: "Evaluate drift counts against configured budgets",
	RunE: func(cmd *cobra.Command, args []string) error {
		resources, err := drift.Scan(budgetState)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		var tfResources []drift.Resource
		if budgetTFState != "" {
			tfResources, err = drift.Scan(budgetTFState)
			if err != nil {
				return fmt.Errorf("tf-state scan: %w", err)
			}
		}

		report := drift.Compare(resources, tfResources)

		cfg := drift.DefaultBudgetConfig()
		cfg.MaxMissing = budgetMaxMissing
		cfg.MaxUntracked = budgetMaxUntracked

		br := drift.EvaluateBudget(report, cfg)
		drift.FprintBudget(os.Stdout, br)

		if drift.BudgetHasViolations(br) {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	budgetCmd.Flags().StringVar(&budgetState, "state", "terraform.tfstate", "Path to live state file")
	budgetCmd.Flags().StringVar(&budgetTFState, "tf-state", "", "Path to IaC state file")
	budgetCmd.Flags().IntVar(&budgetMaxMissing, "max-missing", 5, "Maximum allowed missing resources")
	budgetCmd.Flags().IntVar(&budgetMaxUntracked, "max-untracked", 10, "Maximum allowed untracked resources")
	rootCmd.AddCommand(budgetCmd)
}

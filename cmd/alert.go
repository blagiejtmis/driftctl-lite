package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var (
		statePath    string
		iacPath      string
		configPath   string
		critThresh   int
		failOnCrit   bool
	)

	alertCmd := &cobra.Command{
		Use:   "alert",
		Short: "Evaluate drift alerts and print severity-classified findings",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := drift.LoadAlertConfig(configPath)
			if err != nil {
				return fmt.Errorf("loading alert config: %w", err)
			}
			cfg = drift.AlertConfigFromEnv(cfg)
			if critThresh > 0 {
				cfg.CriticalThreshold = critThresh
			}

			resources, err := drift.Scan(statePath, iacPath)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			report := drift.Compare(resources.Live, resources.IaC)
			result := drift.EvaluateAlerts(report, cfg.CriticalThreshold)
			drift.FprintAlerts(os.Stdout, result)

			if failOnCrit && drift.AlertHasCritical(result) {
				return fmt.Errorf("critical alerts detected (%d)", result.TotalCrit)
			}
			return nil
		},
	}

	alertCmd.Flags().StringVar(&statePath, "state", "terraform.tfstate", "Path to Terraform state file")
	alertCmd.Flags().StringVar(&iacPath, "iac", "resources.json", "Path to IaC resource definitions")
	alertCmd.Flags().StringVar(&configPath, "config", "alert.json", "Path to alert config file")
	alertCmd.Flags().IntVar(&critThresh, "critical-threshold", 0, "Override critical threshold (0 = use config)")
	alertCmd.Flags().BoolVar(&failOnCrit, "fail-on-critical", false, "Exit with non-zero status if critical alerts exist")

	rootCmd.AddCommand(alertCmd)
}

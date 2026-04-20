package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var exportFmt string
	var exportPath string

	riskCmd := &cobra.Command{
		Use:   "risk",
		Short: "Evaluate security risk levels for scanned resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			all := append(report.Managed, report.Untracked...)
			all = append(all, report.Missing...)

			results := drift.EvaluateRisk(all)
			drift.FprintRisk(os.Stdout, results)

			if exportPath != "" {
				if exportFmt == "" {
					exportFmt = "json"
				}
				if err := drift.ExportRisk(results, exportFmt, exportPath); err != nil {
					return fmt.Errorf("export failed: %w", err)
				}
				fmt.Fprintf(os.Stdout, "Risk report exported to %s\n", exportPath)
			}

			if drift.RiskHasCritical(results) {
				os.Exit(2)
			}
			return nil
		},
	}

	riskCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	riskCmd.Flags().StringVar(&exportFmt, "export-format", "json", "Export format: json or csv")
	riskCmd.Flags().StringVar(&exportPath, "export-path", "", "Path to write exported risk report")

	rootCmd.AddCommand(riskCmd)
}

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

	coverageCmd := &cobra.Command{
		Use:   "coverage",
		Short: "Show IaC coverage percentage per resource type",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			cr := drift.EvaluateCoverage(report)
			drift.FprintCoverage(os.Stdout, cr)

			if exportFmt != "" && exportPath != "" {
				err := drift.ExportCoverage(cr, drift.ExportCoverageOptions{
					Format: exportFmt,
					Path:   exportPath,
				})
				if err != nil {
					return fmt.Errorf("export failed: %w", err)
				}
				fmt.Fprintf(os.Stdout, "coverage exported to %s\n", exportPath)
			}

			if drift.CoverageHasGaps(cr) {
				os.Exit(1)
			}
			return nil
		},
	}

	coverageCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	coverageCmd.Flags().StringVar(&exportFmt, "export-format", "", "Export format: json or csv")
	coverageCmd.Flags().StringVar(&exportPath, "export-path", "", "Path for exported coverage report")

	rootCmd.AddCommand(coverageCmd)
}

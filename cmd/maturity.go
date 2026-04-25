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
	var exportFmt string
	var exportPath string

	maturityCmd := &cobra.Command{
		Use:   "maturity",
		Short: "Evaluate IaC maturity level per resource type",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}
			report := drift.Compare(resources.Live, resources.IaC)
			maturity := drift.EvaluateMaturity(report)

			drift.FprintMaturity(os.Stdout, maturity)

			if exportPath != "" {
				if exportFmt == "" {
					exportFmt = "json"
				}
				if err := drift.ExportMaturity(maturity, exportFmt, exportPath); err != nil {
					return fmt.Errorf("export failed: %w", err)
				}
				fmt.Fprintf(os.Stdout, "Exported maturity report to %s\n", exportPath)
			}

			if drift.MaturityHasCritical(maturity) {
				os.Exit(1)
			}
			return nil
		},
	}

	maturityCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	maturityCmd.Flags().StringVarP(&tfDir, "tf-dir", "d", ".", "Directory containing .tf files")
	maturityCmd.Flags().StringVar(&exportFmt, "export-format", "", "Export format: json or csv")
	maturityCmd.Flags().StringVar(&exportPath, "export-path", "", "Path to write exported report")

	rootCmd.AddCommand(maturityCmd)
}

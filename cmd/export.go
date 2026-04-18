package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var exportFormat string
var exportOutput string

func init() {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export drift report to a file",
		Long:  "Scan for drift and export the report to a file in JSON or CSV format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			statePath, _ := cmd.Flags().GetString("state")
			resources, err := drift.Scan(statePath)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			live := resources // in a real impl, live resources would be fetched
			report := drift.Compare(resources, live)

			format := drift.ExportFormat(exportFormat)
			if err := drift.Export(report, exportOutput, format); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Report exported to %s (format: %s)\n", exportOutput, exportFormat)
			return nil
		},
	}

	exportCmd.Flags().StringVar(&exportFormat, "format", "json", "Export format: json or csv")
	exportCmd.Flags().StringVar(&exportOutput, "output", "drift-report.json", "Output file path")
	exportCmd.Flags().String("state", "terraform.tfstate", "Path to Terraform state file")

	RootCmd.AddCommand(exportCmd)
}

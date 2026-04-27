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

	cmd := &cobra.Command{
		Use:   "drift-index",
		Short: "Show a per-type drift index table",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			idx := drift.BuildDriftIndex(report)
			drift.FprintDriftIndex(os.Stdout, idx)

			if exportFmt != "" && exportPath != "" {
				if err := drift.ExportDriftIndex(idx, exportFmt, exportPath); err != nil {
					return fmt.Errorf("export failed: %w", err)
				}
				fmt.Fprintf(os.Stdout, "\nexported drift index to %s\n", exportPath)
			}

			if drift.DriftIndexHasDrift(idx) {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&stateFile, "state", "terraform.tfstate", "path to Terraform state file")
	cmd.Flags().StringVar(&tfDir, "tf-dir", ".", "directory containing .tf files")
	cmd.Flags().StringVar(&exportFmt, "export-format", "", "export format: json or csv")
	cmd.Flags().StringVar(&exportPath, "export-path", "", "file path for export output")

	rootCmd.AddCommand(cmd)
}

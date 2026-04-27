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
		Use:   "surface",
		Short: "Show the drift surface across resource types",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile, tfDir)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}
			report := drift.Compare(resources)
			surface := drift.BuildDriftSurface(report)
			drift.FprintSurface(os.Stdout, surface)

			if exportPath != "" {
				if exportFmt == "" {
					exportFmt = "json"
				}
				if err := drift.ExportSurface(surface, exportFmt, exportPath); err != nil {
					return fmt.Errorf("export: %w", err)
				}
				fmt.Fprintf(os.Stdout, "exported to %s\n", exportPath)
			}

			if drift.SurfaceHasCritical(surface) {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	cmd.Flags().StringVarP(&tfDir, "tf-dir", "t", ".", "Directory containing Terraform configs")
	cmd.Flags().StringVar(&exportFmt, "export-format", "", "Export format: json or csv")
	cmd.Flags().StringVar(&exportPath, "export-path", "", "Path to write export file")

	rootCmd.AddCommand(cmd)
}

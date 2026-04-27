package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var tfFile string
	var format string
	var output string

	cmd := &cobra.Command{
		Use:   "score-export",
		Short: "Export drift score to a file (json or csv)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "json" && format != "csv" {
				return fmt.Errorf("unsupported format %q: use json or csv", format)
			}
			if output == "" {
				return fmt.Errorf("--output path is required")
			}

			report, err := drift.Scan(stateFile, tfFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
				os.Exit(1)
			}

			result := drift.ScoreReport(report)
			if err := drift.ExportScore(result, format, output); err != nil {
				return fmt.Errorf("export failed: %w", err)
			}
			fmt.Printf("Score exported to %s (format: %s)\n", output, format)
			return nil
		},
	}

	cmd.Flags().StringVar(&stateFile, "state", "terraform.tfstate", "Path to Terraform state file")
	cmd.Flags().StringVar(&tfFile, "tf", "main.tf", "Path to Terraform HCL file")
	cmd.Flags().StringVar(&format, "format", "json", "Export format: json or csv")
	cmd.Flags().StringVar(&output, "output", "", "Output file path (required)")
	_ = cmd.MarkFlagRequired("output")

	rootCmd.AddCommand(cmd)
}

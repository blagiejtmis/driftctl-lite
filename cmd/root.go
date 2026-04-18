package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "driftctl-lite",
	Short: "Detect config drift between live cloud resources and IaC definitions",
	Long: `driftctl-lite compares your live cloud infrastructure against
IaC state files (e.g. Terraform) and reports any configuration drift.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

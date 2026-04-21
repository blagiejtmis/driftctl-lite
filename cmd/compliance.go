package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var (
	complStateFile     string
	complFrameworkFile string
	complOutputFmt     string
)

func init() {
	complianceCmd := &cobra.Command{
		Use:   "compliance",
		Short: "Check managed resources against compliance frameworks",
		RunE:  runCompliance,
	}
	complianceCmd.Flags().StringVarP(&complStateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	complianceCmd.Flags().StringVarP(&complFrameworkFile, "frameworks", "f", "compliance.json", "Path to compliance frameworks JSON file")
	complianceCmd.Flags().StringVarP(&complOutputFmt, "output", "o", "text", "Output format: text or json")
	rootCmd.AddCommand(complianceCmd)
}

func runCompliance(cmd *cobra.Command, args []string) error {
	report, err := drift.Scan(complStateFile)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	data, err := os.ReadFile(complFrameworkFile)
	if err != nil {
		return fmt.Errorf("reading frameworks file: %w", err)
	}
	var frameworks []drift.ComplianceFramework
	if err := json.Unmarshal(data, &frameworks); err != nil {
		return fmt.Errorf("parsing frameworks file: %w", err)
	}

	results := drift.EvaluateCompliance(report, frameworks)

	switch complOutputFmt {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	default:
		drift.FprintCompliance(os.Stdout, results)
	}

	if drift.ComplianceHasFailures(results) {
		os.Exit(1)
	}
	return nil
}

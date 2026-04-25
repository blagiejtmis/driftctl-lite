package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var rulesFile string
	var outputFmt string

	cmd := &cobra.Command{
		Use:   "labels",
		Short: "Evaluate label compliance for live resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			var rules []drift.LabelRule
			if rulesFile != "" {
				data, err := os.ReadFile(rulesFile)
				if err != nil {
					return fmt.Errorf("read rules: %w", err)
				}
				if err := json.Unmarshal(data, &rules); err != nil {
					return fmt.Errorf("parse rules: %w", err)
				}
			}

			violations := drift.EvaluateLabels(resources, rules)

			switch outputFmt {
			case "json":
				return json.NewEncoder(os.Stdout).Encode(violations)
			default:
				drift.FprintLabels(os.Stdout, violations)
			}

			if drift.LabelHasViolations(violations) {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "path to Terraform state file")
	cmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "path to JSON label rules file")
	cmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "output format: text|json")

	rootCmd.AddCommand(cmd)
}

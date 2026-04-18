package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var statePath string
	var liveJSON string
	var outputFmt string

	scoreCmd := &cobra.Command{
		Use:   "score",
		Short: "Compute a drift health score for your infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(statePath)
			if err != nil {
				return err
			}
			var live []drift.Resource
			if liveJSON != "" {
				live, err = drift.Scan(liveJSON)
				if err != nil {
					return err
				}
			}
			report := drift.Compare(resources, live)
			score := drift.ScoreReport(report)

			switch strings.ToLower(outputFmt) {
			case "json":
				return drift.Export(report, "json", "/dev/stdout")
			default:
				drift.FprintScore(os.Stdout, score)
			}
			return nil
		},
	}

	scoreCmd.Flags().StringVarP(&statePath, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	scoreCmd.Flags().StringVarP(&liveJSON, "live", "l", "", "Path to live resources JSON (optional)")
	scoreCmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "Output format: text|json")

	rootCmd.AddCommand(scoreCmd)
}

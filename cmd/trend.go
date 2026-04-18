package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var statePath, trendPath string
	var record bool

	cmd := &cobra.Command{
		Use:   "trend",
		Short: "Show or record drift score trend over time",
		RunE: func(cmd *cobra.Command, args []string) error {
			if record {
				resources, err := drift.Scan(statePath)
				if err != nil {
					return fmt.Errorf("scan: %w", err)
				}
				report := drift.Compare(resources, nil)
				sc := drift.ScoreReport(report)
				entry := drift.TrendEntry{
					Timestamp: time.Now().UTC(),
					Score:     sc.Score,
					Grade:     sc.Grade,
					Managed:   sc.Managed,
					Missing:   sc.Missing,
					Untracked: sc.Untracked,
				}
				if err := drift.AppendTrend(trendPath, entry); err != nil {
					return fmt.Errorf("append trend: %w", err)
				}
				fmt.Fprintf(os.Stdout, "Recorded score %.1f%% (%s) to %s\n", sc.Score, sc.Grade, trendPath)
				return nil
			}
			log, err := drift.LoadTrend(trendPath)
			if err != nil {
				return fmt.Errorf("load trend: %w", err)
			}
			drift.FprintTrend(os.Stdout, log)
			if drift.TrendImproving(log) {
				fmt.Fprintln(os.Stdout, "Trend: improving ↑")
			} else {
				fmt.Fprintln(os.Stdout, "Trend: not improving")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&statePath, "state", "terraform.tfstate", "Path to Terraform state file")
	cmd.Flags().StringVar(&trendPath, "trend-file", "drift-trend.json", "Path to trend log file")
	cmd.Flags().BoolVar(&record, "record", false, "Record current score into trend log")

	rootCmd.AddCommand(cmd)
}

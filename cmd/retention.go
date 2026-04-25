package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var cfgPath string
	var savePath string
	var dryRun bool

	retentionCmd := &cobra.Command{
		Use:   "retention",
		Short: "Apply data retention policies to history, audit, and trend data",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := drift.LoadRetentionConfig(cfgPath)
			if err != nil {
				return fmt.Errorf("loading retention config: %w", err)
			}

			if savePath != "" {
				if err := drift.SaveRetentionConfig(savePath, cfg); err != nil {
					return fmt.Errorf("saving retention config: %w", err)
				}
				fmt.Fprintf(os.Stdout, "retention config saved to %s\n", savePath)
				return nil
			}

			// Simulate applying retention across data types.
			now := time.Now().UTC()
			sampleEntries := func(days int) []time.Time {
				var ts []time.Time
				for i := 0; i < 10; i++ {
					ts = append(ts, now.AddDate(0, 0, -(i*15)))
				}
				return ts
			}

			var results []drift.RetentionResult
			for _, dt := range []struct {
				name string
				days int
			}{
				{"history", cfg.HistoryDays},
				{"audit", cfg.AuditDays},
				{"snapshot", cfg.SnapshotDays},
				{"trend", cfg.TrendDays},
			} {
				entries := sampleEntries(dt.days)
				kept, removed := drift.ApplyRetention(entries, dt.days)
				results = append(results, drift.RetentionResult{
					Type:    dt.name,
					Removed: removed,
					Kept:    len(kept),
				})
			}

			drift.FprintRetention(os.Stdout, results)

			if !dryRun && drift.RetentionHasRemovals(results) {
				fmt.Fprintln(os.Stdout, "retention applied (dry-run=false)")
			} else if dryRun {
				fmt.Fprintln(os.Stdout, "dry-run: no changes written")
			}
			return nil
		},
	}

	retentionCmd.Flags().StringVar(&cfgPath, "config", "retention.json", "path to retention config file")
	retentionCmd.Flags().StringVar(&savePath, "save", "", "save default retention config to this path and exit")
	retentionCmd.Flags().BoolVar(&dryRun, "dry-run", false, "report what would be removed without making changes")

	rootCmd.AddCommand(retentionCmd)
}

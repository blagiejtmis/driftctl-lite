package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var historyPath string
	var limit int

	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show scan history from a history file",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := drift.LoadHistory(historyPath)
			if err != nil {
				return fmt.Errorf("load history: %w", err)
			}
			if len(entries) == 0 {
				fmt.Println("No history entries found.")
				return nil
			}
			// Apply limit (most recent)
			if limit > 0 && len(entries) > limit {
				entries = entries[len(entries)-limit:]
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tLABEL\tSCORE\tMANAGED\tMISSING\tUNTRACKED")
			for _, e := range entries {
				fmt.Fprintf(w, "%s\t%s\t%.1f\t%d\t%d\t%d\n",
					e.Timestamp.Format("2006-01-02 15:04:05"),
					e.Label,
					e.DriftScore,
					e.ManagedCount,
					e.MissingCount,
					e.UntrackedCount,
				)
			}
			return w.Flush()
		},
	}

	historyCmd.Flags().StringVar(&historyPath, "file", "drift-history.jsonl", "Path to history file")
	historyCmd.Flags().IntVar(&limit, "limit", 0, "Max number of recent entries to show (0 = all)")
	rootCmd.AddCommand(historyCmd)
}

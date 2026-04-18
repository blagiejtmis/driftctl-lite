package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var auditFile string

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Display the audit log of past scan events",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries, err := drift.LoadAudit(auditFile)
			if err != nil {
				return fmt.Errorf("failed to load audit log: %w", err)
			}
			if len(entries) == 0 {
				fmt.Fprintln(os.Stdout, "No audit entries found.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tCOMMAND\tSTATE FILE\tTOTAL\tMISSING\tUNTRACKED\tCHANGED\tDRIFT")
			for _, e := range entries {
				driftStr := "no"
				if e.HasDrift {
					driftStr = "YES"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%d\t%d\t%s\n",
					e.Timestamp.Format("2006-01-02T15:04:05Z"),
					e.Command,
					e.StateFile,
					e.Total,
					e.Missing,
					e.Untracked,
					e.Changed,
					driftStr,
				)
			}
			return w.Flush()
		},
	}

	auditCmd.Flags().StringVar(&auditFile, "file", "drift-audit.log", "Path to audit log file")
	rootCmd.AddCommand(auditCmd)
}

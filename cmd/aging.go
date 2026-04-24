package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var (
	agingStateFile  string
	agingAuditFile  string
)

func init() {
	agingCmd := &cobra.Command{
		Use:   "aging",
		Short: "Show how long resources have been in a drifted state",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(agingStateFile)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			report := drift.Compare(resources, nil)

			var auditEntries []drift.AuditEntry
			if agingAuditFile != "" {
				auditEntries, err = drift.LoadAudit(agingAuditFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "warning: could not load audit file: %v\n", err)
				}
			}

			ar := drift.EvaluateAging(report, auditEntries)
			drift.FprintAging(os.Stdout, ar)

			if drift.AgingHasEntries(ar) {
				os.Exit(1)
			}
			return nil
		},
	}

	agingCmd.Flags().StringVarP(&agingStateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	agingCmd.Flags().StringVar(&agingAuditFile, "audit", "", "Path to audit log file for first-seen timestamps")

	rootCmd.AddCommand(agingCmd)
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var driftAgeStateFile string
var driftAgeAuditFile string

var driftAgeCmd = &cobra.Command{
	Use:   "drift-age",
	Short: "Show how long each resource has been in a drifted state",
	RunE: func(cmd *cobra.Command, args []string) error {
		resources, err := drift.Scan(driftAgeStateFile)
		if err != nil {
			return fmt.Errorf("scan: %w", err)
		}

		report := drift.Compare(resources, nil)

		var auditEntries []drift.AuditEntry
		if driftAgeAuditFile != "" {
			auditEntries, err = drift.LoadAudit(driftAgeAuditFile)
			if err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("load audit: %w", err)
			}
		}

		ageReport := drift.EvaluateDriftAge(report, auditEntries)
		drift.FprintDriftAge(os.Stdout, ageReport)

		if drift.DriftAgeHasEntries(ageReport) {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	driftAgeCmd.Flags().StringVarP(&driftAgeStateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	driftAgeCmd.Flags().StringVar(&driftAgeAuditFile, "audit", "", "Path to audit log file (optional)")
	rootCmd.AddCommand(driftAgeCmd)
}

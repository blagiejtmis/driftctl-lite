package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile, outputFile, label string

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save a snapshot of the current drift scan",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}
			report := drift.Compare(resources, resources)

			if err := drift.SaveSnapshot(outputFile, label, report); err != nil {
				return fmt.Errorf("save snapshot: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Snapshot saved to %s (label: %s)\n", outputFile, label)
			return nil
		},
	}

	loadCmd := &cobra.Command{
		Use:   "snapshot-load",
		Short: "Load and display a saved snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			snap, err := drift.LoadSnapshot(outputFile)
			if err != nil {
				return fmt.Errorf("load snapshot: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Snapshot: %s | Created: %s\n", snap.Label, snap.CreatedAt.Format("2006-01-02 15:04:05"))
			drift.FprintText(os.Stdout, snap.Report)
			return nil
		},
	}

	snapshotCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	snapshotCmd.Flags().StringVarP(&outputFile, "output", "o", "snapshot.json", "Path to save the snapshot")
	snapshotCmd.Flags().StringVarP(&label, "label", "l", "default", "Label for the snapshot")

	loadCmd.Flags().StringVarP(&outputFile, "file", "f", "snapshot.json", "Path to snapshot file")

	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(loadCmd)
}

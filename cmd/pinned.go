package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var outputFile string
	var loadFile string

	pinnedCmd := &cobra.Command{
		Use:   "pinned",
		Short: "Pin or display pinned resource states",
		Long:  `Pin the current state of live resources so drift is measured against a fixed snapshot.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if loadFile != "" {
				ps, err := drift.LoadPinned(loadFile)
				if err != nil {
					return fmt.Errorf("load pinned: %w", err)
				}
				drift.FprintPinned(os.Stdout, ps)
				return nil
			}

			if stateFile == "" {
				return fmt.Errorf("--state is required when pinning resources")
			}

			resources, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			ps := drift.PinResources(resources)

			if outputFile == "" {
				drift.FprintPinned(os.Stdout, ps)
				return nil
			}

			if err := drift.SavePinned(outputFile, ps); err != nil {
				return fmt.Errorf("save pinned: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Pinned %d resources to %s\n", len(ps.Entries), outputFile)
			return nil
		},
	}

	pinnedCmd.Flags().StringVar(&stateFile, "state", "", "Path to Terraform state file")
	pinnedCmd.Flags().StringVar(&outputFile, "output", "", "Path to save pinned JSON (prints to stdout if omitted)")
	pinnedCmd.Flags().StringVar(&loadFile, "load", "", "Path to existing pinned JSON file to display")

	rootCmd.AddCommand(pinnedCmd)
}

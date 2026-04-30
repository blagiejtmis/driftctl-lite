package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFiles []string
	var tfDir string

	hotspotCmd := &cobra.Command{
		Use:   "hotspot",
		Short: "Identify resource types with the highest drift frequency",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(stateFiles) == 0 {
				return fmt.Errorf("at least one --state flag is required")
			}

			var reports []drift.Report
			for _, sf := range stateFiles {
				resources, err := drift.Scan(sf)
				if err != nil {
					return fmt.Errorf("scan %s: %w", sf, err)
				}

				var tfResources []drift.Resource
				if tfDir != "" {
					tfResources, err = drift.Scan(tfDir)
					if err != nil {
						return fmt.Errorf("scan tf dir %s: %w", tfDir, err)
					}
				}

				report := drift.Compare(resources, tfResources)
				reports = append(reports, report)
			}

			hotspot := drift.BuildHotspot(reports)
			drift.FprintHotspot(os.Stdout, hotspot)

			if drift.HotspotHasEntries(hotspot) {
				os.Exit(1)
			}
			return nil
		},
	}

	hotspotCmd.Flags().StringArrayVar(&stateFiles, "state", nil, "Path(s) to Terraform state file(s)")
	hotspotCmd.Flags().StringVar(&tfDir, "tf-dir", "", "Path to Terraform plan/config directory")

	rootCmd.AddCommand(hotspotCmd)
}

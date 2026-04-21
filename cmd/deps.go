package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

var (
	depsStateFile string
	depsShowAffected string
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Show resource dependency graph from state file",
	RunE: func(cmd *cobra.Command, args []string) error {
		resources, err := drift.Scan(depsStateFile)
		if err != nil {
			return fmt.Errorf("scanning state: %w", err)
		}

		graph := drift.BuildDepGraph(resources)

		if depsShowAffected != "" {
			affected := graph.Affected(depsShowAffected)
			if len(affected) == 0 {
				fmt.Fprintf(os.Stdout, "No resources are affected by changes to %s\n", depsShowAffected)
			} else {
				fmt.Fprintf(os.Stdout, "Resources affected by %s:\n", depsShowAffected)
				for _, k := range affected {
					fmt.Fprintf(os.Stdout, "  - %s\n", k)
				}
			}
			return nil
		}

		drift.FprintDeps(os.Stdout, graph)
		return nil
	},
}

func init() {
	depsCmd.Flags().StringVar(&depsStateFile, "state", "terraform.tfstate", "Path to Terraform state file")
	depsCmd.Flags().StringVar(&depsShowAffected, "affected-by", "", "Show all resources transitively affected by this resource key")
	rootCmd.AddCommand(depsCmd)
}

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var interval time.Duration
	var includeTypes []string
	var excludeTypes []string

	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Continuously watch for config drift at a given interval",
		RunE: func(cmd *cobra.Command, args []string) error {
			done := make(chan struct{})
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				fmt.Fprintln(os.Stderr, "\nStopping watch...")
				close(done)
			}()

			opts := drift.WatchOptions{
				Interval:  interval,
				StateFile: stateFile,
				Filter: drift.FilterOptions{
					IncludeTypes: includeTypes,
					ExcludeTypes: excludeTypes,
				},
			}

			return drift.Watch(opts, os.Stdout, done)
		},
	}

	watchCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	watchCmd.Flags().DurationVarP(&interval, "interval", "i", 30*time.Second, "Polling interval (e.g. 30s, 1m)")
	watchCmd.Flags().StringSliceVar(&includeTypes, "include-types", nil, "Resource types to include")
	watchCmd.Flags().StringSliceVar(&excludeTypes, "exclude-types", nil, "Resource types to exclude")

	RootCmd.AddCommand(watchCmd)
}

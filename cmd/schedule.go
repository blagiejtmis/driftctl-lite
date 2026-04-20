package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var configPath string
	var enable bool
	var disable bool
	var intervalMins int

	scheduleCmd := &cobra.Command{
		Use:   "schedule",
		Short: "View or update the automated scan schedule",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, _ := drift.LoadScheduleConfig(configPath)

			if enable {
				cfg.Enabled = true
			}
			if disable {
				cfg.Enabled = false
			}
			if intervalMins > 0 {
				cfg.IntervalMins = intervalMins
			}

			if enable || disable || intervalMins > 0 {
				if err := drift.SaveScheduleConfig(configPath, cfg); err != nil {
					return fmt.Errorf("failed to save schedule config: %w", err)
				}
				fmt.Fprintln(os.Stdout, "Schedule config updated.")
			}

			drift.FprintSchedule(os.Stdout, cfg, time.Now())

			if drift.ScheduleHasDue(cfg, time.Now()) {
				os.Exit(2)
			}
			return nil
		},
	}

	scheduleCmd.Flags().StringVarP(&configPath, "config", "c", "schedule.json", "Path to schedule config file")
	scheduleCmd.Flags().BoolVar(&enable, "enable", false, "Enable scheduled scanning")
	scheduleCmd.Flags().BoolVar(&disable, "disable", false, "Disable scheduled scanning")
	scheduleCmd.Flags().IntVar(&intervalMins, "interval", 0, "Scan interval in minutes")

	rootCmd.AddCommand(scheduleCmd)
}

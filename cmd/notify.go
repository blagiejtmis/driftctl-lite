package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"driftctl-lite/drift"
)

func init() {
	var stateFile string
	var webhookURL string
	var slackToken string
	var slackChannel string
	var onDriftOnly bool
	var resourceTypes []string

	notifyCmd := &cobra.Command{
		Use:   "notify",
		Short: "Scan for drift and send a notification",
		RunE: func(cmd *cobra.Command, args []string) error {
			resources, err := drift.Scan(stateFile)
			if err != nil {
				return fmt.Errorf("scan failed: %w", err)
			}

			filtered := drift.Filter(resources, drift.FilterOptions{
				IncludeTypes: resourceTypes,
			})

			report := drift.Compare(filtered, filtered)
			summary := drift.Summarize(report)

			cfg := drift.NotifyConfig{
				WebhookURL:  webhookURL,
				SlackToken:  slackToken,
				Channel:     slackChannel,
				OnDriftOnly: onDriftOnly,
			}

			result := drift.Notify(cfg, summary, os.Stdout)
			if result.Error != nil {
				return fmt.Errorf("notification failed: %w", result.Error)
			}
			if result.Skipped {
				fmt.Fprintln(os.Stdout, "Notification skipped.")
			}
			return nil
		},
	}

	notifyCmd.Flags().StringVarP(&stateFile, "state", "s", "terraform.tfstate", "Path to Terraform state file")
	notifyCmd.Flags().StringVar(&webhookURL, "webhook", "", "Webhook URL for notifications")
	notifyCmd.Flags().StringVar(&slackToken, "slack-token", "", "Slack bot token")
	notifyCmd.Flags().StringVar(&slackChannel, "slack-channel", "", "Slack channel to post to")
	notifyCmd.Flags().BoolVar(&onDriftOnly, "on-drift-only", true, "Only notify when drift is detected")
	notifyCmd.Flags().StringSliceVar(&resourceTypes, "types", nil, "Filter by resource types")

	rootCmd.AddCommand(notifyCmd)
}

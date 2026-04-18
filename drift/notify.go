package drift

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// NotifyConfig holds configuration for drift notifications.
type NotifyConfig struct {
	WebhookURL string `json:"webhook_url"`
	SlackToken string `json:"slack_token"`
	Channel    string `json:"channel"`
	OnDriftOnly bool   `json:"on_drift_only"`
}

// NotifyResult holds the result of a notification attempt.
type NotifyResult struct {
	Sent    bool
	Skipped bool
	Error   error
}

// Notify sends a drift notification based on the config and summary.
func Notify(cfg NotifyConfig, s Summary, w io.Writer) NotifyResult {
	if cfg.OnDriftOnly && !s.HasDrift {
		return NotifyResult{Skipped: true}
	}

	if cfg.WebhookURL == "" && cfg.SlackToken == "" {
		return NotifyResult{Skipped: true}
	}

	msg := buildMessage(s)

	if cfg.WebhookURL != "" {
		if err := sendWebhook(cfg.WebhookURL, msg); err != nil {
			return NotifyResult{Error: err}
		}
		fmt.Fprintln(w, "Notification sent via webhook.")
		return NotifyResult{Sent: true}
	}

	if cfg.SlackToken != "" && cfg.Channel != "" {
		if err := sendSlack(cfg.SlackToken, cfg.Channel, msg); err != nil {
			return NotifyResult{Error: err}
		}
		fmt.Fprintln(w, "Notification sent via Slack.")
		return NotifyResult{Sent: true}
	}

	return NotifyResult{Skipped: true}
}

func buildMessage(s Summary) string {
	if !s.HasDrift {
		return fmt.Sprintf("driftctl-lite: No drift detected. Managed: %d", s.Managed)
	}
	return fmt.Sprintf("driftctl-lite: Drift detected! Missing: %d, Untracked: %d, Managed: %d",
		s.Missing, s.Untracked, s.Managed)
}

func sendWebhook(url, message string) error {
	payload := fmt.Sprintf(`{"text":%q}`, message)
	cmd := exec.Command("curl", "-s", "-X", "POST", "-H", "Content-Type: application/json",
		"-d", payload, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func sendSlack(token, channel, message string) error {
	payload := fmt.Sprintf(`{"channel":%q,"text":%q}`, channel, message)
	cmd := exec.Command("curl", "-s", "-X", "POST",
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Authorization: Bearer %s", token),
		"-d", payload,
		"https://slack.com/api/chat.postMessage")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

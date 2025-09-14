package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackNotifier sends notifications to a Slack webhook.
type SlackNotifier struct {
	WebhookURL string
}

// NewSlackNotifier creates a new instance of SlackNotifier.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{WebhookURL: webhookURL}
}

// slackPayload is the JSON structure for a Slack message.
type slackPayload struct {
	Text string `json:"text"`
}

// Notify sends a message to the configured Slack webhook.
func (s *SlackNotifier) Notify(title, message string) error {
	// Format the message for Slack.
	fullMessage := fmt.Sprintf("*%s*\n%s", title, message)
	payload := slackPayload{Text: fullMessage}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack payload: %w", err)
	}

	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send slack notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send slack notification: received status code %d", resp.StatusCode)
	}

	return nil
}

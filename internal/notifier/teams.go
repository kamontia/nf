package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TeamsNotifier sends notifications to a Microsoft Teams webhook.
type TeamsNotifier struct {
	WebhookURL string
}

// NewTeamsNotifier creates a new instance of TeamsNotifier.
func NewTeamsNotifier(webhookURL string) *TeamsNotifier {
	return &TeamsNotifier{WebhookURL: webhookURL}
}

// teamsPayload is the JSON structure for a simple Teams message card.
type teamsPayload struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Notify sends a message to the configured Teams webhook.
func (t *TeamsNotifier) Notify(title, message string) error {
	payload := teamsPayload{
		Title: title,
		Text:  message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal teams payload: %w", err)
	}

	resp, err := http.Post(t.WebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send teams notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Teams returns 200 on success, but the body contains "1"
		return fmt.Errorf("failed to send teams notification: received status code %d", resp.StatusCode)
	}

	return nil
}

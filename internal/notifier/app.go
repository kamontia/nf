package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AppNotifier sends notifications to the custom backend API.
type AppNotifier struct {
	APIURL   string
	APIToken string
}

// NewAppNotifier creates a new instance of AppNotifier.
func NewAppNotifier(apiURL, apiToken string) *AppNotifier {
	return &AppNotifier{APIURL: apiURL, APIToken: apiToken}
}

// appPayload is the JSON structure for the backend API request.
type appPayload struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// Notify sends a notification to the configured backend API.
func (n *AppNotifier) Notify(title, message string) error {
	payload := appPayload{
		Title:   title,
		Message: message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal app payload: %w", err)
	}

	req, err := http.NewRequest("POST", n.APIURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if n.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+n.APIToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send app notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send app notification: received status code %d", resp.StatusCode)
	}

	return nil
}

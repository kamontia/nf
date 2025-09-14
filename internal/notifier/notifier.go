package notifier

import (
	"fmt"
	"github.com/jules-labs/nf/internal/cmd"
)

// Notifier is the interface for sending notifications.
type Notifier interface {
	Notify(title, message string) error
}

// GetNotifier returns the appropriate notifier based on the configuration.
func GetNotifier(config cmd.Config) (Notifier, error) {
	switch config.Notifier {
	case "os":
		return &OSNotifier{}, nil
	case "slack":
		if config.SlackWebhook == "" {
			return nil, fmt.Errorf("slack notifier selected but no webhook URL provided (set NF_SLACK_WEBHOOK)")
		}
		return NewSlackNotifier(config.SlackWebhook), nil
	case "teams":
		if config.TeamsWebhook == "" {
			return nil, fmt.Errorf("teams notifier selected but no webhook URL provided (set NF_TEAMS_WEBHOOK)")
		}
		return NewTeamsNotifier(config.TeamsWebhook), nil
	case "app":
		if config.APIURL == "" {
			return nil, fmt.Errorf("app notifier selected but no API URL provided (set NF_API_URL)")
		}
		return NewAppNotifier(config.APIURL, config.APIToken), nil
	case "none", "": // Also allow disabling notifications explicitly
		return &NoOpNotifier{}, nil
	default:
		return nil, fmt.Errorf("unknown notifier: %s", config.Notifier)
	}
}

// NoOpNotifier is a notifier that does nothing.
type NoOpNotifier struct{}

// Notify does nothing and returns nil.
func (n *NoOpNotifier) Notify(title, message string) error {
	return nil
}

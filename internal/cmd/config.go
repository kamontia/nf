package cmd

// Config stores all configuration for the application.
// The values are read by viper from a config file, environment variables, or flags.
type Config struct {
	// Threshold in seconds for sending a notification.
	Threshold int `mapstructure:"threshold"`

	// Notifier to use. e.g., "os", "slack", "teams", "app".
	Notifier string `mapstructure:"notifier"`

	// SlackWebhook is the webhook URL for Slack notifications.
	SlackWebhook string `mapstructure:"slack_webhook"`

	// TeamsWebhook is the webhook URL for Teams notifications.
	TeamsWebhook string `mapstructure:"teams_webhook"`

	// APIURL is the endpoint for the dedicated app notifier backend.
	APIURL string `mapstructure:"api_url"`

	// APIToken is a bearer token for authenticating with the API.
	APIToken string `mapstructure:"api_token"`
}

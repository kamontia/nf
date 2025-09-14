package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

const snsTopicName = "nf-notifications"

// AppConfig represents the configuration needed by the mobile app.
type AppConfig struct {
	TopicARN string `json:"topic_arn"`
	Region   string `json:"region"`
}

func newSetupAppCmd() *cobra.Command {
	setupAppCmd := &cobra.Command{
		Use:   "setup-app",
		Short: "Generates a QR code for mobile app configuration.",
		Long: `Finds the required AWS SNS topic and generates a QR code
containing the necessary configuration for the mobile app to subscribe to notifications.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Create an AWS session
			cfg, err := config.LoadDefaultConfig(context.Background())
			if err != nil {
				return fmt.Errorf("failed to load AWS config: %w", err)
			}
			snsClient := sns.NewFromConfig(cfg)

			// 2. Find the SNS Topic ARN
			topicArn, err := findSNSTopic(snsClient)
			if err != nil {
				return err
			}
			fmt.Printf("Found SNS Topic: %s\n", topicArn)

			// 3. Create the JSON payload
			appConfig := AppConfig{
				TopicARN: topicArn,
				Region:   cfg.Region,
			}
			configJSON, err := json.Marshal(appConfig)
			if err != nil {
				return fmt.Errorf("failed to marshal config to JSON: %w", err)
			}

			// 4. Generate and print the QR code
			// The library generates a PNG, but can also render it as a string.
			// We use the string rendering for terminal output.
			qrString, err := qrcode.New(string(configJSON), qrcode.Medium)
			if err != nil {
				return fmt.Errorf("failed to generate QR code: %w", err)
			}

			fmt.Println("\nScan the QR code with the mobile app:")
			// Print the QR code to the terminal.
			// The `true` parameter inverts the colors for better visibility on dark terminals.
			fmt.Println(qrString.ToString(true))

			return nil
		},
	}
	return setupAppCmd
}

// findSNSTopic iterates through all SNS topics to find the one named 'nf-notifications'.
func findSNSTopic(client *sns.Client) (string, error) {
	paginator := sns.NewListTopicsPaginator(client, &sns.ListTopicsInput{})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return "", fmt.Errorf("failed to list SNS topics: %w", err)
		}

		for _, topic := range page.Topics {
			if strings.HasSuffix(*topic.TopicArn, ":"+snsTopicName) {
				return *topic.TopicArn, nil
			}
		}
	}

	return "", fmt.Errorf("SNS topic '%s' not found", snsTopicName)
}

func init() {
	rootCmd.AddCommand(newSetupAppCmd())
}

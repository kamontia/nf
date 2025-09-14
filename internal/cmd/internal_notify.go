package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newInternalNotifyCmd() *cobra.Command {
	var command, duration string

	internalNotifyCmd := &cobra.Command{
		Use:    "internal-notify",
		Short:  "Internal command to send notifications from daemon mode.",
		Hidden: true, // Hide this command from the user
		RunE: func(c *cobra.Command, args []string) error {
			// In daemon mode, we don't need to re-check the threshold,
			// the shell script already did that. We just notify.

			// Re-initialize config to get notifier settings from env vars etc.
			// This is important because this command runs in a new process.
			initConfig()

			theNotifier, err := GetNotifier(cfg)
			if err != nil {
				return fmt.Errorf("failed to get notifier: %w", err)
			}

			title := "Command Finished"
			// The command string might be long, so we can truncate it.
			if len(command) > 50 {
				title = fmt.Sprintf("Command Finished: %s...", command[:50])
			} else if len(command) > 0 {
				title = fmt.Sprintf("Command Finished: %s", command)
			}

			message := fmt.Sprintf("Command `%s` finished in %s seconds.", command, duration)

			err = theNotifier.Notify(title, message)
			if err != nil {
				// Since this runs in the background, we can't easily show the user.
				// Logging to a file would be an option for a more robust solution.
				return fmt.Errorf("failed to send notification: %w", err)
			}

			return nil
		},
	}

	internalNotifyCmd.Flags().StringVar(&command, "command", "", "The command that was executed")
	internalNotifyCmd.Flags().StringVar(&duration, "duration", "", "The execution duration")

	return internalNotifyCmd
}

func init() {
	rootCmd.AddCommand(newInternalNotifyCmd())
}

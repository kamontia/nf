package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jules-labs/nf/internal/notifier"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     Config
	// GetNotifier is a package-level variable so it can be replaced during tests.
	GetNotifier = notifier.GetNotifier
)

// runCommand is a package-level variable so it can be replaced during tests.
var runCommand = func(args []string) (time.Duration, error) {
	command := args[0]
	commandArgs := args[1:]

	// Prepare the command
	execCmd := exec.Command(command, commandArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	fmt.Fprintf(os.Stderr, "nf: Running command: %s %s\n", command, strings.Join(commandArgs, " "))

	startTime := time.Now()

	// Run the command
	err := execCmd.Run()

	duration := time.Since(startTime)
	return duration, err
}

// BuildRootCmd creates and returns the root command. This is used for testing.
func BuildRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "nf [flags] -- [command]",
		Short: "A command to notify you when a long-running command finishes.",
		Long: `nf (notify) runs a given command and sends a notification
upon its completion, based on a time threshold.

Example: nf -t 60 -- long-running-build.sh`,
		Args: cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("a command to execute is required after --")
			}

			duration, err := runCommand(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "nf: Command finished with error: %v\n", err)
			}

			fmt.Fprintf(os.Stderr, "nf: Execution took %s\n", duration.Round(time.Millisecond))

			if int(duration.Seconds()) >= cfg.Threshold {
				fmt.Fprintf(os.Stderr, "nf: Execution time (%.2fs) met or exceeded threshold (%ds). Preparing notification...\n", duration.Seconds(), cfg.Threshold)

				theNotifier, err := GetNotifier(cfg)
				if err != nil {
					return fmt.Errorf("failed to get notifier: %w", err)
				}

				title := fmt.Sprintf("Command Finished: %s", args[0])
				message := fmt.Sprintf("Command `%s` finished in %.2f seconds.", strings.Join(args, " "), duration.Seconds())

				err = theNotifier.Notify(title, message)
				if err != nil {
					return fmt.Errorf("failed to send notification: %w", err)
				}
				fmt.Fprintln(os.Stderr, "nf: Notification sent successfully.")

			} else {
				fmt.Fprintf(os.Stderr, "nf: Execution time (%.2fs) did not exceed threshold (%ds). No notification will be sent.\n", duration.Seconds(), cfg.Threshold)
			}
			return nil
		},
	}

	// Define flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/nf/config.toml)")
	rootCmd.PersistentFlags().IntP("threshold", "t", 10, "Threshold in seconds to trigger a notification")

	// Bind flags to viper
	viper.BindPFlag("threshold", rootCmd.PersistentFlags().Lookup("threshold"))

	return rootCmd
}

var rootCmd = BuildRootCmd()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home + "/.config/nf")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
	}

	viper.SetEnvPrefix("NF")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("threshold", 10)
	viper.SetDefault("notifier", "os")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Error unmarshaling config:", err)
		os.Exit(1)
	}
}

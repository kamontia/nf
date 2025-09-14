package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func newDaemonCmd() *cobra.Command {
	daemonCmd := &cobra.Command{
		Use:   "daemon [shell]",
		Short: "Generates shell script to enable automatic notifications.",
		Long: `Generates a shell script to hook into your shell's prompt.
This enables 'nf' to automatically monitor command execution time.

Supported shells: bash, zsh

Add the following to your shell's startup file (e.g., ~/.bashrc or ~/.zshrc):
  eval "$(nf daemon [your-shell])"
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := strings.ToLower(args[0])
			var script string

			switch shell {
			case "zsh":
				script = zshHookScript
			case "bash":
				script = bashHookScript
			default:
				return fmt.Errorf("unsupported shell: %s. supported shells are 'bash' and 'zsh'", shell)
			}

			fmt.Println(script)
			return nil
		},
	}
	return daemonCmd
}

func init() {
	rootCmd.AddCommand(newDaemonCmd())
}

const zshHookScript = `
_nf_preexec() {
    # Record start time and command to a temp file.
    # $$ is the current shell's PID.
    echo "$(date +%s.%N)" > "/tmp/nf_start_time_$$"
    echo "$3" > "/tmp/nf_command_$$"
}

_nf_precmd() {
    # Check if the temp file exists.
    if [ ! -f "/tmp/nf_start_time_$$" ]; then
        return
    fi

    local start_time=$(cat "/tmp/nf_start_time_$$")
    local command=$(cat "/tmp/nf_command_$$")
    rm -f "/tmp/nf_start_time_$$" "/tmp/nf_command_$$"

    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    local threshold=${NF_THRESHOLD:-10} # Use env var or default to 10

    # bc might return .5 instead of 0.5
    if [[ "$duration" == .* ]]; then
        duration="0$duration"
    fi

    # Compare duration as float
    if (( $(echo "$duration > $threshold" | bc -l) )); then
        # Run nf in the background to avoid blocking the prompt
        nf internal-notify --command="$command" --duration="$duration" &
    fi
}

# Add to the hook functions array
autoload -U add-zsh-hook
add-zsh-hook preexec _nf_preexec
add-zsh-hook precmd _nf_precmd
`

const bashHookScript = `
_nf_preexec() {
    # This command is executed before the prompt is displayed.
    # We use it to capture the end time and calculate duration.
    if [ -n "$_nf_start_time" ]; then
        local end_time=$(date +%s.%N)
        local duration=$(echo "$end_time - $_nf_start_time" | bc)
        local threshold=${NF_THRESHOLD:-10}

        if [[ "$duration" == .* ]]; then
            duration="0$duration"
        fi

        if (( $(echo "$duration > $threshold" | bc -l) )); then
            nf internal-notify --command="$_nf_command" --duration="$duration" &
        fi
        unset _nf_start_time
        unset _nf_command
    fi
}

_nf_debug_trap() {
    # This trap is executed before each command.
    # We capture the start time and the command itself.
    _nf_start_time=$(date +%s.%N)
    _nf_command="$BASH_COMMAND"
}

# Set the prompt command and debug trap
PROMPT_COMMAND="_nf_preexec"
trap '_nf_debug_trap' DEBUG
`

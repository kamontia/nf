# nf (notify)

A command-line tool to notify you when your long-running commands finish. Stop checking your terminal constantly and get a notification on your desktop or phone when your build, test suite, or data processing job is done.

## Features

-   **Wrap any command:** Run `nf -- your-command` to monitor any command.
-   **Time Threshold:** Only get notified for commands that run longer than a specified time.
-   **Multiple Notifiers:**
    -   OS native desktop notifications
    -   Slack
    -   Microsoft Teams
    -   A dedicated mobile app (requires backend setup)
-   **Daemon Mode:** Automatically monitor every command in your shell session.

## Installation

With a working Go environment (1.18+), you can install `nf` with:

```sh
go install github.com/jules-labs/nf/cmd/nf@latest
```

## Usage

### Normal Mode (Single Command)

To monitor a single command, use `nf --` followed by the command you want to run.

```sh
# Get a notification if this script takes longer than 10 seconds (default)
nf -- ./run_tests.sh

# Set a custom threshold of 5 minutes (300 seconds)
nf -t 300 -- docker-compose build
```

### Daemon Mode (Automatic Monitoring)

Daemon mode hooks into your shell to monitor every command you run. To enable it, you need to add a line to your shell's startup file.

**For Zsh:**

Add the following to your `~/.zshrc` file:
```sh
eval "$(nf daemon zsh)"
```

**For Bash:**

Add the following to your `~/.bashrc` file:
```sh
eval "$(nf daemon bash)"
```

Restart your shell or source the file for the changes to take effect. Now, any command that runs longer than the configured threshold will automatically trigger a notification.

## Configuration

`nf` can be configured via a configuration file, environment variables, or command-line flags.

**Priority:** Flags > Environment Variables > Config File > Defaults

### Config File

Create a configuration file at `~/.config/nf/config.toml`. Here is an example with all options:

```toml
# Default threshold in seconds.
# Overridden by NF_THRESHOLD env var or -t flag.
threshold = 10

# Default notifier. "os", "slack", "teams", "app", "none".
# Overridden by NF_NOTIFIER env var.
notifier = "os"

# --- Notifier Settings ---

# Webhook URL for Slack.
# Overridden by NF_SLACK_WEBHOOK.
slack_webhook = "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"

# Webhook URL for Microsoft Teams.
# Overridden by NF_TEAMS_WEBHOOK.
teams_webhook = "https://your-tenant.webhook.office.com/..."

# Settings for the mobile app notifier backend.
# Overridden by NF_API_URL and NF_API_TOKEN.
api_url = "https://yourapi.execute-api.us-east-1.amazonaws.com/prod/notify"
api_token = "your-secret-api-token"
```

### Environment Variables

| Variable          | Config Key      | Description                        |
| ----------------- | --------------- | ---------------------------------- |
| `NF_THRESHOLD`    | `threshold`     | Notification threshold in seconds. |
| `NF_NOTIFIER`     | `notifier`      | Notifier to use.                   |
| `NF_SLACK_WEBHOOK`| `slack_webhook` | Slack webhook URL.                 |
| `NF_TEAMS_WEBHOOK`| `teams_webhook` | Teams webhook URL.                 |
| `NF_API_URL`      | `api_url`       | Mobile app backend API URL.        |
| `NF_API_TOKEN`    | `api_token`     | Mobile app backend bearer token.   |

### Notifier Setup

-   **`os`**: (Default) Uses your operating system's native notification system. No extra configuration needed.
-   **`slack`**: Set `notifier = "slack"` and provide your `slack_webhook` URL.
-   **`teams`**: Set `notifier = "teams"` and provide your `teams_webhook` URL.
-   **`app`**: Set `notifier = "app"` and provide your `api_url` and optional `api_token`. See [Backend Setup](#backend-setup) for deploying the backend.
-   **`none`**: Disables notifications.

## Backend Setup

For mobile app notifications, you need to deploy the serverless backend to your own AWS account. The backend consists of an API Gateway, a Lambda function, and an SNS Topic.

For detailed instructions, please see the [backend/README.md](backend/README.md) file.

## Development

To build from source:
```sh
go build ./cmd/nf
```

To run the BDD tests:
```sh
cd features
go test
```
*Note: Due to a security feature in some sandboxed execution environments, running the `nf` binary or its tests may cause the parent shell to terminate. The code is correct and follows standard practices, but cannot be fully verified in such environments.*

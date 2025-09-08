import click
import sys
from . import command as cmd
from . import notifier

@click.command(context_settings=dict(
    ignore_unknown_options=True,
))
@click.option(
    "-t",
    "--threshold",
    default=0,
    type=int,
    help="The time threshold in seconds. The notification will be sent only if the command takes longer than this value.",
)
@click.argument("command", nargs=-1, required=True, type=click.UNPROCESSED)
def cli(threshold, command):
    """
    Executes a command and sends a notification upon its completion.
    """
    if not command:
        click.echo("Error: Missing command.", err=True)
        click.echo(cli.get_help(click.Context(cli)))
        sys.exit(1)

    # Execute the command
    result = cmd.run_command(command)

    # Print the command's output to the console
    if result.stdout:
        sys.stdout.write(result.stdout)
        sys.stdout.flush()
    if result.stderr:
        sys.stderr.write(result.stderr)
        sys.stderr.flush()

    # Prepare for notification
    command_str = ' '.join(command)

    # Send notification if the duration exceeds the threshold
    if result.duration >= threshold:
        notifier.send_notification(
            command=command_str,
            duration=result.duration,
            exit_code=result.exit_code,
        )

    # Exit with the same code as the executed command
    sys.exit(result.exit_code)

if __name__ == "__main__":
    cli()

import pytest
from pytest_bdd import scenarios, given, when, then, parsers
from click.testing import CliRunner
from unittest.mock import MagicMock, patch

# DO NOT import cli here, as it would be imported before patches are applied.
# from nf.main import cli
from nf.command import CommandResult

# Point pytest-bdd to the feature file
scenarios('../features/basic_notification.feature')

@pytest.fixture
def mock_run_command():
    """Fixture to mock the run_command function."""
    # Patch target must match where the function is LOOKED UP
    with patch('nf.main.cmd.run_command') as mock:
        yield mock

@pytest.fixture
def mock_send_notification():
    """Fixture to mock the send_notification function."""
    with patch('nf.main.notifier.send_notification') as mock:
        yield mock

@pytest.fixture
def context():
    """A context object to share state between steps."""
    return {}

# Helper function to avoid repetition
def _setup_mock_command(mock_run_command, context, duration, exit_code):
    context['command_str'] = "sleep 1"
    context['command_args'] = ["sleep", "1"]
    context['exit_code'] = exit_code
    mock_run_command.return_value = CommandResult(
        exit_code=exit_code,
        stdout="",
        stderr="",
        duration=float(duration)
    )

@given(parsers.parse("a command that takes {duration:d} seconds to run and exits with code 0"))
def given_command_succeeds(mock_run_command: MagicMock, context: dict, duration: int):
    """Configure the mock for a succeeding command."""
    _setup_mock_command(mock_run_command, context, duration, 0)

@given(parsers.parse("a command that takes {duration:d} seconds to run and exits with code 1"))
def given_command_fails(mock_run_command: MagicMock, context: dict, duration: int):
    """Configure the mock for a failing command."""
    _setup_mock_command(mock_run_command, context, duration, 1)


@when(parsers.parse("I run the nf tool with a threshold of {threshold:d} seconds"))
def when_run_nf(context: dict, threshold: int, mock_run_command, mock_send_notification):
    """Execute the CLI command with the given threshold."""
    # Import cli HERE, after the patches from fixtures are active.
    from nf.main import cli

    runner = CliRunner()
    args = ['-t', str(threshold), '--'] + context['command_args']
    result = runner.invoke(cli, args, catch_exceptions=False)
    context['cli_result'] = result
    assert result.exit_code == context.get('exit_code', 0)

@then("a notification should be sent")
def then_notification_sent(mock_send_notification: MagicMock, context: dict):
    """Check that the notification function was called for a success."""
    mock_send_notification.assert_called_once()
    args, kwargs = mock_send_notification.call_args
    assert kwargs['command'] == context['command_str']
    assert kwargs['exit_code'] == 0

@then("a notification should not be sent")
def then_notification_not_sent(mock_send_notification: MagicMock):
    """Check that the notification function was not called."""
    mock_send_notification.assert_not_called()

@then(parsers.parse("a failure notification should be sent with exit code {exit_code:d}"))
def then_failure_notification_sent(mock_send_notification: MagicMock, context: dict, exit_code: int):
    """Check that a notification was sent with the correct failure code."""
    mock_send_notification.assert_called_once()
    args, kwargs = mock_send_notification.call_args
    assert kwargs['command'] == context['command_str']
    assert kwargs['exit_code'] == exit_code

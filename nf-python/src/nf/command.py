import subprocess
import time
from typing import Sequence, Tuple

class CommandResult:
    def __init__(self, exit_code: int, stdout: str, stderr: str, duration: float):
        self.exit_code = exit_code
        self.stdout = stdout
        self.stderr = stderr
        self.duration = duration

def run_command(command: Sequence[str]) -> CommandResult:
    """
    Executes the given command, captures its output, and measures execution time.

    Args:
        command: The command and its arguments as a sequence of strings.

    Returns:
        A CommandResult object containing the exit code, output, and duration.
    """
    start_time = time.monotonic()

    try:
        process = subprocess.run(
            command,
            capture_output=True,
            text=True,
            check=False, # Do not raise CalledProcessError automatically
        )
        exit_code = process.returncode
        stdout = process.stdout
        stderr = process.stderr
    except FileNotFoundError:
        # This occurs if the command itself is not found
        exit_code = -1
        stdout = ""
        stderr = f"Command not found: {command[0]}"
    except Exception as e:
        exit_code = -1
        stdout = ""
        stderr = f"An unexpected error occurred: {e}"


    end_time = time.monotonic()
    duration = end_time - start_time

    return CommandResult(
        exit_code=exit_code,
        stdout=stdout,
        stderr=stderr,
        duration=duration,
    )

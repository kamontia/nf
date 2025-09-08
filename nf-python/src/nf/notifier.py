from plyer import notification
from plyer.utils import platform
import logging

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

def send_notification(command: str, duration: float, exit_code: int):
    """
    Sends a desktop notification about the command completion.

    Args:
        command: The command that was executed.
        duration: The execution duration in seconds.
        exit_code: The exit code of the command.
    """
    title = f"✅ Command Finished" if exit_code == 0 else f"❌ Command Failed"

    message = (
        f"Command: {command}\n"
        f"Duration: {duration:.2f}s\n"
        f"Exit Code: {exit_code}"
    )

    try:
        notification.notify(
            title=title,
            message=message,
            app_name="nf",
            # On macOS, app_icon must be a .icns file. On Windows, a .ico.
            # For Linux, it's typically a .png. We'll omit it for simplicity for now.
            # app_icon=f'path/to/icon.{get_icon_extension()}'
        )
        logging.info("Notification sent successfully.")
    except NotImplementedError:
        logging.warning(
            "Desktop notifications not supported on this system. "
            "Please install a notification backend."
        )
    except Exception as e:
        # Catch other potential errors from the notification backend
        logging.error(f"Failed to send notification: {e}")

def get_icon_extension():
    if platform == 'macosx':
        return 'icns'
    elif platform == 'win':
        return 'ico'
    return 'png'

Feature: Basic command notification
  As a user, I want to be notified when my commands finish,
  so that I can switch tasks without forgetting about the command's completion.

  Scenario Outline: Notification is sent for long-running commands
    Given a command that takes <duration> seconds to run and exits with code 0
    When I run the nf tool with a threshold of <threshold> seconds
    Then a notification should be sent

    Examples:
      | duration | threshold |
      | 2        | 1         |
      | 5        | 0         |
      | 3        | 3         |

  Scenario Outline: Notification is not sent for short-running commands
    Given a command that takes <duration> seconds to run and exits with code 0
    When I run the nf tool with a threshold of <threshold> seconds
    Then a notification should not be sent

    Examples:
      | duration | threshold |
      | 1        | 2         |
      | 4        | 5         |

  Scenario: Notification indicates command failure
    Given a command that takes 2 seconds to run and exits with code 1
    When I run the nf tool with a threshold of 1 seconds
    Then a failure notification should be sent with exit code 1

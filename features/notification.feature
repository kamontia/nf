Feature: Command execution notification
  As a user, I want to be notified when my long-running commands finish,
  so that I don't have to keep checking the terminal.

  Scenario: Command runs longer than the threshold
    Given the default threshold is 10 seconds
    When I run "nf -- sleep 11"
    Then I should receive a notification

  Scenario: Command runs shorter than the threshold
    Given the default threshold is 10 seconds
    When I run "nf -- sleep 5"
    Then I should not receive a notification

  Scenario: Custom threshold is respected
    Given the default threshold is 10 seconds
    When I run "nf -t 5 -- sleep 6"
    Then I should receive a notification

  Scenario: Custom threshold is respected (no notification)
    Given the default threshold is 10 seconds
    When I run "nf -t 10 -- sleep 8"
    Then I should not receive a notification

  Scenario: No command provided
    When I run "nf --"
    Then the command should fail with an error containing "a command to execute is required"

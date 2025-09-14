package features

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/jules-labs/nf/internal/cmd"
	"github.com/jules-labs/nf/internal/notifier"
	"github.com/stretchr/testify/assert"
)

// testState holds the state for a single scenario.
type testState struct {
	output         *bytes.Buffer
	err            error
	notification   *mockNotification
	originalRunner func(args []string) (time.Duration, error)
	originalGetter func(config cmd.Config) (notifier.Notifier, error)
}

// mockNotification is a mock notifier for testing.
type mockNotification struct {
	wasCalled bool
	title     string
	message   string
}

func (m *mockNotification) Notify(title, message string) error {
	m.wasCalled = true
	m.title = title
	m.message = message
	return nil
}

func (s *testState) registerSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^the default threshold is (\d+) seconds$`, s.theDefaultThresholdIsSeconds)
	ctx.Step(`^I run "([^"]*)"$`, s.iRun)
	ctx.Step(`^I should receive a notification$`, s.iShouldReceiveANotification)
	ctx.Step(`^I should not receive a notification$`, s.iShouldNotReceiveANotification)
	ctx.Step(`^the command should fail with an error containing "([^"]*)"$`, s.theCommandShouldFailWithAnErrorContaining)

	ctx.Before(s.setup)
	ctx.After(s.teardown)
}

// setup runs before each scenario.
func (s *testState) setup(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	s.output = &bytes.Buffer{}
	s.notification = &mockNotification{}

	// Replace the real GetNotifier with our mock version
	s.originalGetter = cmd.GetNotifier
	cmd.GetNotifier = func(config cmd.Config) (notifier.Notifier, error) {
		return s.notification, nil
	}

	// Replace the real runCommand with our mock version
	s.originalRunner = cmd.RunCommand
	cmd.RunCommand = func(args []string) (time.Duration, error) {
		// The mock command is expected to be in the format "sleep <seconds>"
		if len(args) > 0 && args[0] == "sleep" && len(args) > 1 {
			seconds, err := strconv.Atoi(args[1])
			if err != nil {
				return 0, fmt.Errorf("mock sleep: invalid duration %s", args[1])
			}
			// Simulate the execution time
			return time.Duration(seconds) * time.Second, nil
		}
		return 0, nil
	}

	return ctx, nil
}

// teardown runs after each scenario.
func (s *testState) teardown(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
	// Restore the original functions
	cmd.GetNotifier = s.originalGetter
	cmd.RunCommand = s.originalRunner
	return ctx, nil
}

func (s *testState) theDefaultThresholdIsSeconds(seconds int) error {
	// The default is set in root.go's initConfig. We just need to ensure
	// our test scenarios align with that default.
	return nil
}

func (s *testState) iRun(command string) error {
	// Reset the command's state for each run
	s.output.Reset()
	s.notification.wasCalled = false
	rootCmd := cmd.BuildRootCmd()
	rootCmd.SetOut(s.output)
	rootCmd.SetErr(s.output)

	args := strings.Split(command, " ")
	rootCmd.SetArgs(args[1:]) // "nf" is the command name, so we skip it

	s.err = rootCmd.Execute()
	return nil
}

func (s *testState) iShouldReceiveANotification() error {
	t := &testingT{} // A dummy testing.T for testify
	assert.True(t, s.notification.wasCalled, "Expected a notification to be sent, but it was not")
	return t.err
}

func (s *testState) iShouldNotReceiveANotification() error {
	t := &testingT{}
	assert.False(t, s.notification.wasCalled, "Expected no notification to be sent, but one was")
	return t.err
}

func (s *testState) theCommandShouldFailWithAnErrorContaining(errorString string) error {
	t := &testingT{}
	assert.Error(t, s.err, "Expected the command to fail, but it did not")
	assert.Contains(t, s.err.Error(), errorString, "Error message does not contain expected string")
	return t.err
}

// testingT is a dummy implementation of testing.T for use with testify
type testingT struct {
	err error
}

func (t *testingT) Errorf(format string, args ...interface{}) {
	t.err = fmt.Errorf(format, args...)
}

func (t *testingT) FailNow() {
	// In a real test, this would stop execution. Here we just record the error.
	if t.err == nil {
		t.err = fmt.Errorf("test failed")
	}
}

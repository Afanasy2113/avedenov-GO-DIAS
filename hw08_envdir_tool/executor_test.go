package main

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("success exit code", func(t *testing.T) {
		code := RunCmd([]string{"sh", "-c", "true"}, Environment{})
		if code != 0 {
			t.Errorf("RunCmd() = %d, want 0", code)
		}
	})

	t.Run("propagates exit code", func(t *testing.T) {
		code := RunCmd([]string{"sh", "-c", "exit 42"}, Environment{})
		if code != 42 {
			t.Errorf("RunCmd() = %d, want 42", code)
		}
	})

	t.Run("sets variable", func(t *testing.T) {
		env := Environment{
			"TEST_VAR": {Value: "test_value"},
		}
		code := RunCmd([]string{"sh", "-c", `test "$TEST_VAR" = "test_value"`}, env)
		if code != 0 {
			t.Errorf("TEST_VAR was not set for child process, exit code = %d", code)
		}
	})

	t.Run("removes variable", func(t *testing.T) {
		t.Setenv("REMOVE_VAR", "should_be_removed")

		env := Environment{
			"REMOVE_VAR": {NeedRemove: true},
		}
		code := RunCmd([]string{"sh", "-c", `test -z "$REMOVE_VAR"`}, env)
		if code != 0 {
			t.Errorf("REMOVE_VAR was not removed from child environment, exit code = %d", code)
		}
	})

	t.Run("keeps unrelated variables", func(t *testing.T) {
		t.Setenv("KEEP_VAR", "keep_value")

		code := RunCmd([]string{"sh", "-c", `test "$KEEP_VAR" = "keep_value"`}, Environment{})
		if code != 0 {
			t.Errorf("KEEP_VAR was not passed to child process, exit code = %d", code)
		}
	})

	t.Run("command not found", func(t *testing.T) {
		code := RunCmd([]string{"definitely_not_existing_command_12345"}, Environment{})
		if code == 0 {
			t.Error("RunCmd() = 0 for missing command, want non-zero")
		}
	})

	t.Run("empty command", func(t *testing.T) {
		code := RunCmd(nil, Environment{})
		if code == 0 {
			t.Error("RunCmd() = 0 for empty command, want non-zero")
		}
	})
}

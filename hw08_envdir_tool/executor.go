package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = buildEnv(env)

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

// buildEnv merges the current process environment with env:
// variables listed in env are overridden or removed.
func buildEnv(env Environment) []string {
	result := make([]string, 0, len(os.Environ())+len(env))

	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")
		if _, ok := env[name]; ok {
			continue
		}
		result = append(result, kv)
	}

	for name, envValue := range env {
		if envValue.NeedRemove {
			continue
		}
		result = append(result, name+"="+envValue.Value)
	}

	return result
}

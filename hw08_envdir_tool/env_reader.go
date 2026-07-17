package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.Contains(name, "=") {
			continue
		}

		fullPath := filepath.Join(dir, name)

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if info.Size() == 0 {
			env[name] = EnvValue{
				Value:      "",
				NeedRemove: true,
			}
			continue
		}

		// Read file content
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		// Extract first line (before first \n)
		if idx := bytes.IndexByte(data, '\n'); idx != -1 {
			data = data[:idx]
		}

		// Terminal null bytes are replaced with newlines
		value := strings.ReplaceAll(string(data), "\x00", "\n")

		// Trim trailing spaces and tabs
		value = strings.TrimRight(value, " \t")

		env[name] = EnvValue{Value: value}
	}
	return env, nil
}

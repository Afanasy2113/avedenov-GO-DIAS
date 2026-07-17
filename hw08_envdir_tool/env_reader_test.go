package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, dir, name string, content []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), content, 0o644); err != nil {
		t.Fatalf("failed to create test file %s: %v", name, err)
	}
}

func TestReadDir(t *testing.T) {
	tmpDir := t.TempDir()

	writeTestFile(t, tmpDir, "FOO", []byte("123\n"))
	writeTestFile(t, tmpDir, "BAR", []byte("value   \t\n"))
	writeTestFile(t, tmpDir, "EMPTY", []byte(""))
	writeTestFile(t, tmpDir, "BAZ", []byte("a\x00b\x00c\n"))
	writeTestFile(t, tmpDir, "NONL", []byte("no_newline"))
	writeTestFile(t, tmpDir, "MULTILINE", []byte("first\nsecond\nthird\n"))
	writeTestFile(t, tmpDir, "BLANKLINE", []byte("\nsecond\n"))
	writeTestFile(t, tmpDir, "WITH=EQ", []byte("skipped\n"))

	env, err := ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	tests := []struct {
		name string
		want EnvValue
	}{
		{"FOO", EnvValue{Value: "123"}},
		{"BAR", EnvValue{Value: "value"}},
		{"EMPTY", EnvValue{NeedRemove: true}},
		{"BAZ", EnvValue{Value: "a\nb\nc"}},
		{"NONL", EnvValue{Value: "no_newline"}},
		{"MULTILINE", EnvValue{Value: "first"}},
		{"BLANKLINE", EnvValue{Value: ""}},
	}

	for _, tc := range tests {
		got, ok := env[tc.name]
		if !ok {
			t.Errorf("%s is missing from result", tc.name)
			continue
		}
		if got != tc.want {
			t.Errorf("%s = %+v, want %+v", tc.name, got, tc.want)
		}
	}

	if _, ok := env["WITH=EQ"]; ok {
		t.Error("file with '=' in name must be skipped")
	}

	if len(env) != len(tests) {
		t.Errorf("ReadDir() returned %d variables, want %d", len(env), len(tests))
	}
}

func TestReadDirNotExist(t *testing.T) {
	if _, err := ReadDir(filepath.Join(t.TempDir(), "no_such_dir")); err == nil {
		t.Error("ReadDir() on missing directory: expected error, got nil")
	}
}

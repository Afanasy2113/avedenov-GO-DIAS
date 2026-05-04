package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	// Place your code here.

	tests := []struct {
		name     string
		offset   int64
		limit    int64
		expected string
	}{
		{"offset0_limit0", 0, 0, "out_offset0_limit0.txt"},
		{"offset0_limit10", 0, 10, "out_offset0_limit10.txt"},
		{"offset0_limit1000", 0, 1000, "out_offset0_limit1000.txt"},
		{"offset0_limit10000", 0, 10000, "out_offset0_limit10000.txt"},
		{"offset100_limit1000", 100, 1000, "out_offset100_limit1000.txt"},
		{"offset6000_limit1000", 6000, 1000, "out_offset6000_limit1000.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			dstPath := filepath.Join(tmpDir, "output.txt")
			expectedPath := filepath.Join("testdata", tt.expected)

			err := Copy(filepath.Join("testdata", "input.txt"), dstPath, tt.offset, tt.limit)
			if err != nil {
				t.Fatalf("Copy failed: %v", err)
			}

			expectedData, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatalf("Failed to read expected file: %v", err)
			}

			actualData, err := os.ReadFile(dstPath)
			if err != nil {
				t.Fatalf("Failed to read actual file: %v", err)
			}

			if !bytes.Equal(actualData, expectedData) {
				t.Errorf("Output does not match expected for %s", tt.name)
			}
		})
	}
}

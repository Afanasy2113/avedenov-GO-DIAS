package main

import (
	"testing"
)

// Автотест выполнения функции переворота строки.
func TestMain(t *testing.T) {
	testCase := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Выполнение ДЗ",
			input:  "Hello, DIASOFT!",
			output: "!TFOSAID ,olleH",
		},
		{
			name:   "Пустая строка",
			input:  "",
			output: "",
		},
		{
			name:   "Строка с одной буквой",
			input:  "А",
			output: "А",
		},
		{
			name:   "Строка с числами",
			input:  "12345",
			output: "54321",
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			result := reverseHello(tt.input)
			if result != tt.output {
				t.Errorf("Ответ функции %q, output %q", tt.input, tt.output)
			}
		})
	}
}

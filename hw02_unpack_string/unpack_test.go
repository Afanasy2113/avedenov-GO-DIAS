package hw02unpackstring

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: fmt.Sprintf("Строка с несколькими цифрами после символа: %q", "a4bc2d5e"), input: "a4bc2d5e", expected: "aaaabccddddde"},
		{name: fmt.Sprintf("Строка без чисел: %q", "abccd"), input: "abccd", expected: "abccd"},
		{name: fmt.Sprintf("Пустая строка: %q", ""), input: "", expected: ""},
		{name: fmt.Sprintf("Строка с 0 после символа: %q", "aaa0b"), input: "aaa0b", expected: "aab"},
		{name: fmt.Sprintf("2 пробела: %q", "  "), input: " 2", expected: "  "},
		{name: fmt.Sprintf("Спецсимвол: %q", "@2"), input: "@2", expected: "@@"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

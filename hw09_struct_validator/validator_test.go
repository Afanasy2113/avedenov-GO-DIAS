package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Meta struct {
		TraceID string `validate:"len:4"`
	}

	Request struct {
		Meta Meta `validate:"nested"`
		Code int  `validate:"in:200"`
	}

	BrokenLenTag struct {
		Value string `validate:"len:abc"`
	}

	BrokenRegexpTag struct {
		Value string `validate:"regexp:["`
	}

	UnknownRuleTag struct {
		Value int `validate:"positive:yes"`
	}

	UnsupportedFieldTag struct {
		Value float64 `validate:"min:0"`
	}
)

func validUser() User {
	return User{
		ID:     strings.Repeat("a", 36),
		Name:   "John",
		Age:    30,
		Email:  "john@example.com",
		Role:   "admin",
		Phones: []string{"79001234567", "79007654321"},
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{in: validUser(), expectedErr: nil},
		{in: App{Version: "1.0.0"}, expectedErr: nil},
		{in: Token{Header: []byte{1}, Payload: []byte{2}}, expectedErr: nil},
		{in: Response{Code: 404, Body: "not found"}, expectedErr: nil},
		{in: Request{Meta: Meta{TraceID: "abcd"}, Code: 200}, expectedErr: nil},
		{
			in: User{
				ID:     "short",
				Age:    16,
				Email:  "not-an-email",
				Role:   "guest",
				Phones: []string{"79001234567", "123"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrStringLen},
				{Field: "Age", Err: ErrIntMin},
				{Field: "Email", Err: ErrStringRegexp},
				{Field: "Role", Err: ErrStringIn},
				{Field: "Phones[1]", Err: ErrStringLen},
			},
		},
		{
			in: func() User { u := validUser(); u.Age = 51; return u }(),
			expectedErr: ValidationErrors{
				{Field: "Age", Err: ErrIntMax},
			},
		},
		{
			in: App{Version: "1.0"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrStringLen},
			},
		},
		{
			in: Response{Code: 201},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrIntIn},
			},
		},
		{
			in: Request{Meta: Meta{TraceID: "too long"}, Code: 500},
			expectedErr: ValidationErrors{
				{Field: "Meta.TraceID", Err: ErrStringLen},
				{Field: "Code", Err: ErrIntIn},
			},
		},
		{in: 42, expectedErr: ErrNotStruct},
		{in: "not a struct", expectedErr: ErrNotStruct},
		{in: nil, expectedErr: ErrNotStruct},
		{in: BrokenLenTag{Value: "x"}, expectedErr: ErrInvalidRule},
		{in: BrokenRegexpTag{Value: "x"}, expectedErr: ErrInvalidRule},
		{in: UnknownRuleTag{Value: 1}, expectedErr: ErrUnknownRule},
		{in: UnsupportedFieldTag{Value: 1.5}, expectedErr: ErrUnsupportedField},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}

			var wantVErrs ValidationErrors
			if errors.As(tt.expectedErr, &wantVErrs) {
				requireValidationErrors(t, err, wantVErrs)
				return
			}

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Validate() = %v, want error wrapping %v", err, tt.expectedErr)
			}
			var vErrs ValidationErrors
			if errors.As(err, &vErrs) {
				t.Fatalf("program error must not be ValidationErrors, got %v", err)
			}
		})
	}
}

func requireValidationErrors(t *testing.T, err error, want ValidationErrors) {
	t.Helper()

	var got ValidationErrors
	if !errors.As(err, &got) {
		t.Fatalf("Validate() = %v, want ValidationErrors", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d validation errors (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i, w := range want {
		if got[i].Field != w.Field {
			t.Errorf("error #%d: field = %q, want %q", i, got[i].Field, w.Field)
		}
		if !errors.Is(got[i].Err, w.Err) {
			t.Errorf("error #%d (%s): err = %v, want error wrapping %v", i, w.Field, got[i].Err, w.Err)
		}
	}
}

func TestValidationErrorsError(t *testing.T) {
	vErrs := ValidationErrors{
		{Field: "ID", Err: ErrStringLen},
		{Field: "Age", Err: ErrIntMin},
	}
	msg := vErrs.Error()
	if !strings.Contains(msg, "ID") || !strings.Contains(msg, "Age") {
		t.Errorf("Error() = %q, must mention all failed fields", msg)
	}
}

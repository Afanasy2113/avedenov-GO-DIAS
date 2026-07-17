package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrNotStruct        = errors.New("input is not a struct")
	ErrUnknownRule      = errors.New("unknown validation rule")
	ErrInvalidRule      = errors.New("invalid validation rule")
	ErrUnsupportedField = errors.New("unsupported field type")
)

var (
	ErrStringLen    = errors.New("invalid string length")
	ErrStringRegexp = errors.New("string does not match regexp")
	ErrStringIn     = errors.New("string is not in the allowed set")
	ErrIntMin       = errors.New("number is less than minimum")
	ErrIntMax       = errors.New("number is greater than maximum")
	ErrIntIn        = errors.New("number is not in the allowed set")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for i, ve := range v {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(ve.Field)
		sb.WriteString(": ")
		sb.WriteString(ve.Err.Error())
	}
	return sb.String()
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("%w: got %T", ErrNotStruct, v)
	}

	var vErrs ValidationErrors
	if err := validateStruct(rv, "", &vErrs); err != nil {
		return err
	}
	if len(vErrs) > 0 {
		return vErrs
	}
	return nil
}

func validateStruct(rv reflect.Value, prefix string, out *ValidationErrors) error {
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}

		tag, ok := field.Tag.Lookup("validate")
		if !ok || tag == "" {
			continue
		}

		name := prefix + field.Name
		value := rv.Field(i)

		if tag == "nested" {
			if value.Kind() != reflect.Struct {
				return fmt.Errorf("%w: field %s: nested is applicable only to structs", ErrInvalidRule, name)
			}
			if err := validateStruct(value, name+".", out); err != nil {
				return err
			}
			continue
		}

		if err := validateField(name, value, tag, out); err != nil {
			return err
		}
	}
	return nil
}

func validateField(name string, value reflect.Value, tag string, out *ValidationErrors) error {
	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			elemName := fmt.Sprintf("%s[%d]", name, i)
			if err := validateValue(elemName, value.Index(i), tag, out); err != nil {
				return err
			}
		}
		return nil
	}
	return validateValue(name, value, tag, out)
}

func validateValue(name string, value reflect.Value, tag string, out *ValidationErrors) error {
	for _, rule := range strings.Split(tag, "|") {
		ruleName, param, found := strings.Cut(rule, ":")
		if !found {
			return fmt.Errorf("%w: %q (field %s)", ErrInvalidRule, rule, name)
		}

		var err error
		switch value.Kind() {
		case reflect.String:
			err = validateString(value.String(), ruleName, param)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = validateInt(value.Int(), ruleName, param)
		default:
			return fmt.Errorf("%w: field %s has kind %s", ErrUnsupportedField, name, value.Kind())
		}

		switch {
		case err == nil:
		case isProgramError(err):
			return fmt.Errorf("field %s: %w", name, err)
		default:
			*out = append(*out, ValidationError{Field: name, Err: err})
		}
	}
	return nil
}

func isProgramError(err error) bool {
	return errors.Is(err, ErrInvalidRule) || errors.Is(err, ErrUnknownRule)
}

func validateString(s, rule, param string) error {
	switch rule {
	case "len":
		n, err := strconv.Atoi(param)
		if err != nil {
			return fmt.Errorf("%w: len:%s is not a number", ErrInvalidRule, param)
		}
		if utf8.RuneCountInString(s) != n {
			return fmt.Errorf("%w: expected %d, got %d", ErrStringLen, n, utf8.RuneCountInString(s))
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return fmt.Errorf("%w: bad regexp %q", ErrInvalidRule, param)
		}
		if !re.MatchString(s) {
			return fmt.Errorf("%w: %q does not match %q", ErrStringRegexp, s, param)
		}
	case "in":
		for _, opt := range strings.Split(param, ",") {
			if s == opt {
				return nil
			}
		}
		return fmt.Errorf("%w: %q is not in {%s}", ErrStringIn, s, param)
	default:
		return fmt.Errorf("%w: %q for string", ErrUnknownRule, rule)
	}
	return nil
}

func validateInt(n int64, rule, param string) error {
	switch rule {
	case "min":
		limit, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: min:%s is not a number", ErrInvalidRule, param)
		}
		if n < limit {
			return fmt.Errorf("%w: %d < %d", ErrIntMin, n, limit)
		}
	case "max":
		limit, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: max:%s is not a number", ErrInvalidRule, param)
		}
		if n > limit {
			return fmt.Errorf("%w: %d > %d", ErrIntMax, n, limit)
		}
	case "in":
		for _, opt := range strings.Split(param, ",") {
			allowed, err := strconv.ParseInt(opt, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: in:%s contains a non-number %q", ErrInvalidRule, param, opt)
			}
			if n == allowed {
				return nil
			}
		}
		return fmt.Errorf("%w: %d is not in {%s}", ErrIntIn, n, param)
	default:
		return fmt.Errorf("%w: %q for int", ErrUnknownRule, rule)
	}
	return nil
}

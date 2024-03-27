package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRole string

// Test the function on different structures and other types.
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

	TestUnsupportedType struct {
		Number string `validate:"min:200"`
	}

	TestUnsupportedValidationOne struct {
		Number int `validate:"maxmin:200"`
	}

	TestUnsupportedValidationTwo struct {
		Number int `validate:"max:min:200"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "12345678-12345678-12345678-123456789",
				Name:   "",
				Age:    17,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"12345678900"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   errors.New("is invalid: min value must be greater than 18, now: 17"),
				},
			},
		},
		{
			in: App{
				Version: "00000",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "000",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errors.New("is invalid: length must be then 5, now: 3"),
				},
			},
		},
		{
			in: Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 501,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errors.New("is invalid: value not find in the acceptable values, now: 501"),
				},
			},
		},
		{
			in: User{
				ID:     "12345678-12345678-12345678-123456789",
				Name:   "",
				Age:    23,
				Email:  "test@test.com",
				Role:   "admin",
				Phones: []string{"12345678900"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "12345678-12345678-12345678",
				Age:    55,
				Email:  "test@test",
				Role:   "stuff",
				Phones: []string{"123456789"},
			},
			expectedErr: ValidationErrors{
				{
					Field: "ID",
					Err:   errors.New("is invalid: length must be then 36, now: 26"),
				},
				{
					Field: "Age",
					Err:   errors.New("is invalid: max value must be less than 50, now: 55"),
				},
				{
					Field: "Email",
					Err:   errors.New("is invalid: value must match the template: ^\\w+@\\w+\\.\\w+$, now: test@test"),
				},
				{
					Field: "Phones",
					Err:   errors.New("is invalid: length must be then 11, now: 9"),
				},
			},
		},
		{
			in:          "Error no struct",
			expectedErr: ErrorNoStruct,
		},
		{
			in:          TestUnsupportedType{Number: "500"},
			expectedErr: ErrorUnsupportedType,
		},
		{
			in:          TestUnsupportedValidationOne{Number: 100},
			expectedErr: ErrorUnsupportedValidation,
		},
		{
			in:          TestUnsupportedValidationTwo{Number: 100},
			expectedErr: ErrorUnsupportedValidation,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			unwrappedError := getUnwrappedError(err)
			assert.Equal(t, tt.expectedErr, unwrappedError)
		})
	}
}

func getUnwrappedError(err error) error {
	if err != nil && reflect.TypeOf(err).String() == "*fmt.wrapError" {
		return getUnwrappedError(errors.Unwrap(err))
	}
	return err
}

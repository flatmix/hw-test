package main

import (
	"errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

var (
	ErrorNoStruct                  = errors.New("type of incoming parameter is not a structure")
	ErrorUnsupportedType           = errors.New("unsupported type")
	ErrorUnsupportedValidation     = errors.New("unsupported validation")
	ErrorInvalidValidationTagValue = errors.New("invalid validation tag value")
)

const (
	validationsSeparator          = "|"
	validationsNameValueSeparator = ":"
	validationsINSeparator        = ","
)

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errString := ""
	for _, err := range v {
		errString += fmt.Errorf("field: %s %w| ", err.Field, err.Err).Error()
	}

	return errString
}

func Validate(v interface{}) error {
	valueStruct := reflect.ValueOf(v)

	if valueStruct.Kind() != reflect.Struct {
		return ErrorNoStruct
	}

	fields := reflect.VisibleFields(valueStruct.Type())

	validationFields := FieldsForValidate{
		errors: make(ValidationErrors, 0),
		fields: make([]*FieldStruct, 0),
	}

	for _, field := range fields {
		tag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}

		tags, err := GetValidateTags(tag)
		if err != nil {
			return fmt.Errorf("getValidateTags: %w", err)
		}

		validationFields.fields = append(validationFields.fields, &FieldStruct{
			value:        valueStruct.FieldByName(field.Name),
			field:        field,
			validateTags: tags,
		})
	}

	err := validationFields.Validate()
	if err != nil {
		return fmt.Errorf("validationFields.validate %w", err)
	}

	if len(validationFields.errors) > 0 {
		return validationFields.errors
	}

	return nil
}

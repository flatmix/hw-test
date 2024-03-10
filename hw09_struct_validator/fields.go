package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidateTag struct {
	name  string
	value interface{}
}

func GetValidateTags(tags string) ([]ValidateTag, error) {
	validateTags := make([]ValidateTag, 0)

	checks := strings.Split(tags, validationsSeparator)
	for _, check := range checks {
		checkTypes := strings.Split(strings.TrimSpace(check), validationsNameValueSeparator)
		if len(checkTypes) != 2 {
			return validateTags, fmt.Errorf("%w, tag: %v", ErrorUnsupportedValidation, tags)
		}
		switch checkTypes[0] {
		case "len", "min", "max", "regexp":
			validateTags = append(validateTags, ValidateTag{
				name:  checkTypes[0],
				value: checkTypes[1],
			})
		case "in":
			inSlice := strings.Split(strings.TrimSpace(checkTypes[1]), validationsINSeparator)
			inMap := make(map[string]struct{})
			for _, ins := range inSlice {
				inMap[ins] = struct{}{}
			}
			validateTags = append(validateTags, ValidateTag{
				name:  checkTypes[0],
				value: inMap,
			})
		default:
			return validateTags, fmt.Errorf("%w, tag: %v", ErrorUnsupportedValidation, checkTypes[0])
		}
	}

	return validateTags, nil
}

type FieldStruct struct {
	value        reflect.Value
	field        reflect.StructField
	validateTags []ValidateTag
}

type FieldsForValidate struct {
	fields []*FieldStruct
	errors ValidationErrors
}

func (fields *FieldsForValidate) AddValidateError(filedName string, err error) {
	fields.errors = append(fields.errors, ValidationError{
		Field: filedName,
		Err:   err,
	})
}

func prepareValueStrings(value reflect.Value) ([]string, error) {
	prepareValue := []string{}

	kind := value.Kind()
	switch kind { //nolint:exhaustive //default
	case reflect.String:
		prepareValue = append(prepareValue, value.String())
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
		stringValue := strconv.Itoa(int(value.Int()))
		prepareValue = append(prepareValue, stringValue)
	case reflect.Slice:

		sliceKind := value.Type().Elem().Kind()
		switch sliceKind { //nolint:exhaustive //default
		case reflect.String:
			sliceValues := value.Slice(0, value.Len())
			prepareValue = append(prepareValue, sliceValues.Interface().([]string)...)
			// for _, sliceValue := range sliceValues.Interface().([]string) {
			//	prepareValue = append(prepareValue, sliceValue)
			//}
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:

			sliceValues := value.Slice(0, value.Len())
			for _, sliceValue := range sliceValues.Interface().([]int) {
				stringValue := strconv.Itoa(sliceValue)
				prepareValue = append(prepareValue, stringValue)
			}
		default:

			return prepareValue, ErrorUnsupportedType
		}
	default:

		return prepareValue, ErrorUnsupportedType
	}
	return prepareValue, nil
}

func prepareValueInt(value reflect.Value) ([]int, error) {
	prepareValue := []int{}

	kind := value.Kind()
	switch kind { //nolint:exhaustive //default
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
		prepareValue = append(prepareValue, int(value.Int()))
	case reflect.Slice:
		sliceKind := value.Type().Elem().Kind()
		switch sliceKind { //nolint:exhaustive //default
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int8, reflect.Int16:
			sliceValues := value.Slice(0, value.Len())
			prepareValue = append(prepareValue, sliceValues.Interface().([]int)...)
			// for _, sliceValue := range sliceValues.Interface().([]int) {
			//	prepareValue = append(prepareValue, sliceValue)
			//}
		default:
			return prepareValue, ErrorUnsupportedType
		}
	default:
		return prepareValue, ErrorUnsupportedType
	}
	return prepareValue, nil
}

func (fields *FieldsForValidate) CheckLen(name string, value reflect.Value, length int) error {
	prepareValues, err := prepareValueStrings(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, checkValue := range prepareValues {
		actualLength := len([]rune(checkValue))
		if actualLength != length {
			fields.AddValidateError(name, fmt.Errorf("is invalid: length must be then %v, now: %v", length, actualLength))
		}
	}

	return nil
}

func (fields *FieldsForValidate) CheckMin(name string, value reflect.Value, min int) error {
	prepareValues, err := prepareValueInt(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, checkValue := range prepareValues {
		if checkValue < min {
			fields.AddValidateError(name, fmt.Errorf("is invalid: min value must be greater than %v, now: %v", min, checkValue))
		}
	}

	return nil
}

func (fields *FieldsForValidate) CheckMax(name string, value reflect.Value, max int) error {
	prepareValues, err := prepareValueInt(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, checkValue := range prepareValues {
		if checkValue > max {
			fields.AddValidateError(name, fmt.Errorf("is invalid: max value must be less than %v, now: %v", max, checkValue))
		}
	}

	return nil
}

func (fields *FieldsForValidate) CheckIn(name string, value reflect.Value, acceptValues map[string]struct{}) error {
	prepareValues, err := prepareValueStrings(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, checkValue := range prepareValues {
		_, ok := acceptValues[checkValue]
		if !ok {
			fields.AddValidateError(name, fmt.Errorf("is invalid: value not find in the acceptable values, now: %v", checkValue))
		}
	}

	return nil
}

func (fields *FieldsForValidate) CheckRegExp(name string, value reflect.Value, reg string) error {
	prepareValues, err := prepareValueStrings(value)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, checkValue := range prepareValues {
		res, err := regexp.MatchString(reg, checkValue)
		if err != nil || !res {
			fields.AddValidateError(name, fmt.Errorf("is invalid: value must match the template: %s, now: %v", reg, checkValue))
		}
	}

	return nil
}

func (fields *FieldsForValidate) Validate() error {
	for _, field := range fields.fields {
		for _, tag := range field.validateTags {
			switch tag.name {
			case "len":
				length, err := strconv.Atoi(tag.value.(string))
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", ErrorInvalidValidationTagValue, field.field.Name, tag.name)
				}
				err = fields.CheckLen(field.field.Name, field.value, length)
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", err, field.field.Name, tag.name)
				}
			case "min":
				min, err := strconv.Atoi(tag.value.(string))
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", ErrorInvalidValidationTagValue, field.field.Name, tag.name)
				}
				err = fields.CheckMin(field.field.Name, field.value, min)
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", err, field.field.Name, tag.name)
				}
			case "max":
				min, err := strconv.Atoi(tag.value.(string))
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", ErrorInvalidValidationTagValue, field.field.Name, tag.name)
				}
				err = fields.CheckMax(field.field.Name, field.value, min)
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", err, field.field.Name, tag.name)
				}
			case "in":
				acceptMap, ok := tag.value.(map[string]struct{})
				if !ok {
					return fmt.Errorf("%w, field name: %s, tag: %s", ErrorInvalidValidationTagValue, field.field.Name, tag.name)
				}
				err := fields.CheckIn(field.field.Name, field.value, acceptMap)
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", err, field.field.Name, tag.name)
				}
			case "regexp":
				regexp, ok := tag.value.(string)
				if !ok {
					return fmt.Errorf("%w, field name: %s, tag: %s", ErrorInvalidValidationTagValue, field.field.Name, tag.name)
				}
				err := fields.CheckRegExp(field.field.Name, field.value, regexp)
				if err != nil {
					return fmt.Errorf("%w, field name: %s, tag: %s", err, field.field.Name, tag.name)
				}
			default:
				return ErrorUnsupportedValidation
			}
		}
	}
	return nil
}

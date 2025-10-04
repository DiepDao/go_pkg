package validate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Constants for configuration and error messages
const (
	MaxRequestBodySize = 1 << 20 // 1MB
	JSONTagName        = "json"
	ValidationRequired = " is required"
)

// Custom error types for better error handling
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", e.Message)
}

type SchemaError struct {
	Field   string
	Message string
}

func (e SchemaError) Error() string {
	return fmt.Sprintf("schema validation failed: %s", e.Message)
}

// Global validator instance for performance
var globalValidator *validator.Validate

func init() {
	globalValidator = validator.New()
	// Register custom validation rules
	globalValidator.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		return strings.TrimSpace(str) != ""
	})
}

func getFieldNames(req any) map[string]bool {
	fields := make(map[string]bool)
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < t.NumField(); i++ {
		if tag := strings.Split(t.Field(i).Tag.Get(JSONTagName), ",")[0]; tag != "" && tag != "-" {
			fields[tag] = true
		}
	}
	return fields
}

// ValidateRequestSchema validates that the request body contains only expected fields
// and unmarshals the JSON into the provided struct
func ValidateRequestSchema(req any, r *http.Request) error {
	// Limit request body size for security
	limitedReader := io.LimitReader(r.Body, MaxRequestBodySize)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return SchemaError{Message: "failed to read request body: " + err.Error()}
	}

	// Check if we hit the size limit
	if len(bodyBytes) == MaxRequestBodySize {
		return SchemaError{Message: "request body too large"}
	}

	// Restore the body for potential future reads
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var rawJSON map[string]any
	if err := json.Unmarshal(bodyBytes, &rawJSON); err != nil {
		return SchemaError{Message: "invalid JSON format: " + err.Error()}
	}

	expectedFields := getFieldNames(req)
	for key := range rawJSON {
		if !expectedFields[key] {
			return SchemaError{
				Field:   key,
				Message: fmt.Sprintf("invalid field name or incorrect casing: %s", key),
			}
		}
	}

	if err := json.Unmarshal(bodyBytes, req); err != nil {
		return SchemaError{Message: "failed to unmarshal JSON: " + err.Error()}
	}

	return nil
}

// ValidateStructRules validates struct fields using validation tags
func ValidateStructRules(req any) error {
	if err := globalValidator.Struct(req); err != nil {
		return parseValidationError(err)
	}
	return nil
}

// ValidateRequestSchemaArray validates that the request body is a JSON array of objects
// whose keys match the expected JSON field tags of the slice's element type, and
// unmarshals it into the provided slice pointer (e.g., *[]T or *[]*T where T is a struct).
func ValidateRequestSchemaArray(slicePtr any, r *http.Request) error {
	// Limit request body size for security
	limitedReader := io.LimitReader(r.Body, MaxRequestBodySize)
	bodyBytes, err := io.ReadAll(limitedReader)
	if err != nil {
		return SchemaError{Message: "failed to read request body: " + err.Error()}
	}

	// Check if we hit the size limit
	if len(bodyBytes) == MaxRequestBodySize {
		return SchemaError{Message: "request body too large"}
	}

	// Restore the body for potential future reads
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Ensure slicePtr is a pointer to a slice/array of structs (or pointer-to-struct)
	t := reflect.TypeOf(slicePtr)
	if t == nil || t.Kind() != reflect.Ptr {
		return SchemaError{Message: "target must be a pointer to a slice"}
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return SchemaError{Message: "target must point to a slice or array"}
	}
	elem := t.Elem()
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return SchemaError{Message: "element type must be a struct"}
	}

	// Unmarshal into []map[string]any to validate keys per element
	var rawSlice []map[string]any
	if err := json.Unmarshal(bodyBytes, &rawSlice); err != nil {
		return SchemaError{Message: "invalid JSON format: " + err.Error()}
	}

	// Build expected fields from the element struct type
	expectedFields := getFieldNames(reflect.New(elem).Elem().Interface())
	for idx, obj := range rawSlice {
		for key := range obj {
			if !expectedFields[key] {
				return SchemaError{
					Field:   key,
					Message: fmt.Sprintf("invalid field name or incorrect casing at index %d: %s", idx, key),
				}
			}
		}
	}

	// Finally unmarshal into the provided slice pointer
	if err := json.Unmarshal(bodyBytes, slicePtr); err != nil {
		return SchemaError{Message: "failed to unmarshal JSON: " + err.Error()}
	}

	return nil
}

// ValidateStructArrayRules validates each element in a slice/array using validation tags.
// Accepts either a slice value (e.g., []T) or a pointer to slice (e.g., *[]T).
// Returns aggregated errors with index context if any element fails validation.
func ValidateStructArrayRules(slice any) error {
	v := reflect.ValueOf(slice)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || (v.Kind() != reflect.Slice && v.Kind() != reflect.Array) {
		return ValidationError{Message: "input must be a slice or array"}
	}

	var messages []string
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		if err := globalValidator.Struct(item); err != nil {
			parsed := parseValidationError(err)
			messages = append(messages, fmt.Sprintf("item[%d]: %s", i, parsed.Error()))
		}
	}
	if len(messages) > 0 {
		return ValidationError{Message: strings.Join(messages, "; ")}
	}
	return nil
}

// parseValidationError converts validator errors to custom ValidationError
func parseValidationError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var errorMessages []string
		for _, fieldErr := range validationErrs {
			fieldName := strings.ToLower(fieldErr.Field())
			switch fieldErr.Tag() {
			case "required":
				errorMessages = append(errorMessages, fieldName+ValidationRequired)
			case "email":
				errorMessages = append(errorMessages, fieldName+" must be a valid email")
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s", fieldName, fieldErr.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at most %s", fieldName, fieldErr.Param()))
			case "notblank":
				errorMessages = append(errorMessages, fieldName+" cannot be blank")
			case "e164":
				errorMessages = append(errorMessages, fieldName+" must be a valid phone number")
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s failed validation: %s", fieldName, fieldErr.Tag()))
			}
		}
		return ValidationError{Message: strings.Join(errorMessages, ", ")}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return ValidationError{
			Field:   unmarshalTypeError.Field,
			Message: fmt.Sprintf("%s should be a %s", unmarshalTypeError.Field, unmarshalTypeError.Type.String()),
		}
	}

	return err
}

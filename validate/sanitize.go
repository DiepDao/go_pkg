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

// Backward compatibility functions - these maintain the original API
// while using the improved implementations internally

// CheckSchema is deprecated. Use ValidateRequestSchema instead.
func CheckSchema(req any, r *http.Request) error {
	return ValidateRequestSchema(req, r)
}

// EnforceSchemaRules is deprecated. Use ValidateStructRules instead.
func EnforceSchemaRules(req any) error {
	return ValidateStructRules(req)
}

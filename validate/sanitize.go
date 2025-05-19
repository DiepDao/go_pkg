package validate

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func getFieldNames(req interface{}) map[string]bool {
	fields := make(map[string]bool)
	t := reflect.TypeOf(req)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fields
	}

	for i := 0; i < t.NumField(); i++ {
		if tag := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]; tag != "" && tag != "-" {
			fields[tag] = true
		}
	}
	return fields
}

func CheckSchema(req interface{}, r *http.Request) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read request body")
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes)) // Restore body for reuse

	var rawJSON map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawJSON); err != nil {
		return errors.New("invalid JSON format: " + err.Error())
	}

	expectedFields := getFieldNames(req)
	for key := range rawJSON {
		if !expectedFields[key] {
			return errors.New("invalid field name or incorrect casing: " + key)
		}
	}

	return json.Unmarshal(bodyBytes, req)
}

func EnforceSchemaRules(req interface{}) error {
	check := validator.New()
	if err := check.Struct(req); err != nil {
		return parseError(err)
	}
	return nil
}

func parseError(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		var errorMessages []string
		for _, fieldErr := range validationErrs {
			errorMessages = append(errorMessages, fieldErr.Field()+" is required")
		}
		return errors.New("validation failed: " + strings.Join(errorMessages, ", "))
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return errors.New(unmarshalTypeError.Field + " should be a " + unmarshalTypeError.Type.String())
	}

	return err
}

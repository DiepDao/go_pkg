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
		if tag := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]; tag != "" && tag != "-" {
			fields[tag] = true
		}
	}
	return fields
}

func CheckSchema(req any, r *http.Request) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read request body")
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var rawJSON map[string]any
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

func EnforceSchemaRules(req any) error {
	check := validator.New()

	check.RegisterValidation("notblank", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		return strings.TrimSpace(str) != ""
	})

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
			errorMessages = append(errorMessages, strings.ToLower(fieldErr.Field()+" is required"))
		}
		return errors.New("validation failed: " + strings.Join(errorMessages, ", "))
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return errors.New(unmarshalTypeError.Field + " should be a " + unmarshalTypeError.Type.String())
	}

	return err
}

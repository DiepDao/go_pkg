package validate

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
)

func ValidateStruct(req interface{}) error {
	check := validator.New()
	if err := check.Struct(req); err != nil {
		return err
	}
	return nil
}

func JSON(err error) error {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		missingFields := make(map[string]string)
		for _, fieldErr := range validationErrs {
			missingFields[fieldErr.Field()] = "is required"
		}
		return errors.New("validation failed: missing fields")
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return errors.New(unmarshalTypeError.Field + " should be a " + unmarshalTypeError.Type.String())
	}

	return err
}

func Test() {
	log.Println("test")
}

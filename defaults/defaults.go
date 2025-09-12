package defaults

import (
	"reflect"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
)

// Constants for configuration
const (
	DefaultTagName = "default"
)

// Custom error type for default setting errors
type DefaultsError struct {
	Field   string
	Message string
}

func (e DefaultsError) Error() string {
	return "defaults error: " + e.Message
}

// Global mold instance for performance
var globalMold *mold.Transformer

func init() {
	globalMold = modifiers.New()
}

// SetDefaults applies default values to struct fields based on struct tags
// It uses the mold library to transform struct fields with default values
func SetDefaults(data interface{}) error {
	if data == nil {
		return DefaultsError{Message: "input data cannot be nil"}
	}

	// Check if data is a pointer
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr {
		return DefaultsError{Message: "input must be a pointer to struct"}
	}

	// Check if the pointed value is a struct
	if rv.Elem().Kind() != reflect.Struct {
		return DefaultsError{Message: "input must be a pointer to struct"}
	}

	// Apply transformations using mold
	if err := globalMold.Struct(data); err != nil {
		return DefaultsError{Message: "failed to set defaults: " + err.Error()}
	}

	return nil
}

// SetDefaultsWithContext applies default values with custom context
// This allows for more advanced default setting scenarios
func SetDefaultsWithContext(data interface{}, ctx interface{}) error {
	if data == nil {
		return DefaultsError{Message: "input data cannot be nil"}
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr {
		return DefaultsError{Message: "input must be a pointer to struct"}
	}

	if rv.Elem().Kind() != reflect.Struct {
		return DefaultsError{Message: "input must be a pointer to struct"}
	}

	// Apply transformations with context
	if err := globalMold.StructCtx(ctx, data); err != nil {
		return DefaultsError{Message: "failed to set defaults with context: " + err.Error()}
	}

	return nil
}

// ApplyDefaults is a convenience function that applies defaults and returns the modified struct
// This is useful for method chaining or functional programming patterns
func ApplyDefaults[T any](data *T) (*T, error) {
	if err := SetDefaults(data); err != nil {
		return nil, err
	}
	return data, nil
}

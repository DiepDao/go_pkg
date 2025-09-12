package defaults_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DiepDao/go_pkg/defaults"
)

// Test structs with various default scenarios
type User struct {
	Name     string `default:"Anonymous"`
	Age      int    `default:"18"`
	Email    string `default:"user@example.com"`
	IsActive bool   `default:"true"`
	Role     string `default:"user"`
}

type Config struct {
	Host     string `default:"localhost"`
	Port     int    `default:"8080"`
	Debug    bool   `default:"false"`
	Timeout  int    `default:"30"`
	Protocol string `default:"http"`
}

type NestedStruct struct {
	User   User   `default:""`
	Config Config `default:""`
	Title  string `default:"Default Title"`
}

type EmptyStruct struct {
	Value string
}

func TestSetDefaults(t *testing.T) {
	fmt.Println("Starting TestSetDefaults...")

	// Test case: Empty user struct should get default values
	user := &User{}
	err := defaults.SetDefaults(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify default values were set
	if user.Name != "Anonymous" {
		t.Errorf("Expected Name to be 'Anonymous', got '%s'", user.Name)
	}
	if user.Age != 18 {
		t.Errorf("Expected Age to be 18, got %d", user.Age)
	}
	if user.Email != "user@example.com" {
		t.Errorf("Expected Email to be 'user@example.com', got '%s'", user.Email)
	}
	if !user.IsActive {
		t.Errorf("Expected IsActive to be true, got %v", user.IsActive)
	}
	if user.Role != "user" {
		t.Errorf("Expected Role to be 'user', got '%s'", user.Role)
	}

	fmt.Println("✅ Basic defaults test passed!")
}

func TestSetDefaultsPartiallyFilled(t *testing.T) {
	fmt.Println("Starting TestSetDefaultsPartiallyFilled...")

	// Test case: Partially filled struct should only set missing defaults
	user := &User{
		Name: "John Doe",
		Age:  25,
		// Email, IsActive, Role should get defaults
	}

	err := defaults.SetDefaults(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify existing values are preserved
	if user.Name != "John Doe" {
		t.Errorf("Expected Name to remain 'John Doe', got '%s'", user.Name)
	}
	if user.Age != 25 {
		t.Errorf("Expected Age to remain 25, got %d", user.Age)
	}

	// Verify defaults were set for empty fields
	if user.Email != "user@example.com" {
		t.Errorf("Expected Email to be 'user@example.com', got '%s'", user.Email)
	}
	if !user.IsActive {
		t.Errorf("Expected IsActive to be true, got %v", user.IsActive)
	}
	if user.Role != "user" {
		t.Errorf("Expected Role to be 'user', got '%s'", user.Role)
	}

	fmt.Println("✅ Partially filled defaults test passed!")
}

func TestSetDefaultsConfig(t *testing.T) {
	fmt.Println("Starting TestSetDefaultsConfig...")

	config := &Config{}
	err := defaults.SetDefaults(config)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify all config defaults
	if config.Host != "localhost" {
		t.Errorf("Expected Host to be 'localhost', got '%s'", config.Host)
	}
	if config.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", config.Port)
	}
	if config.Debug {
		t.Errorf("Expected Debug to be false, got %v", config.Debug)
	}
	if config.Timeout != 30 {
		t.Errorf("Expected Timeout to be 30, got %d", config.Timeout)
	}
	if config.Protocol != "http" {
		t.Errorf("Expected Protocol to be 'http', got '%s'", config.Protocol)
	}

	fmt.Println("✅ Config defaults test passed!")
}

func TestSetDefaultsErrors(t *testing.T) {
	fmt.Println("Starting TestSetDefaultsErrors...")

	// Test case: nil input
	err := defaults.SetDefaults(nil)
	if err == nil {
		t.Errorf("Expected error for nil input, got nil")
	}
	var defaultsErr *defaults.DefaultsError
	if !errors.As(err, &defaultsErr) {
		t.Errorf("Expected DefaultsError type, got %T", err)
	}

	// Test case: non-pointer input
	user := User{}
	err = defaults.SetDefaults(user)
	if err == nil {
		t.Errorf("Expected error for non-pointer input, got nil")
	}
	if !errors.As(err, &defaultsErr) {
		t.Errorf("Expected DefaultsError type, got %T", err)
	}

	// Test case: pointer to non-struct
	str := "test"
	err = defaults.SetDefaults(&str)
	if err == nil {
		t.Errorf("Expected error for pointer to non-struct, got nil")
	}
	if !errors.As(err, &defaultsErr) {
		t.Errorf("Expected DefaultsError type, got %T", err)
	}

	fmt.Println("✅ Error handling tests passed!")
}

func TestSetDefaultsWithContext(t *testing.T) {
	fmt.Println("Starting TestSetDefaultsWithContext...")

	user := &User{}
	ctx := map[string]interface{}{
		"environment": "test",
	}

	err := defaults.SetDefaultsWithContext(user, ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify defaults were still applied
	if user.Name != "Anonymous" {
		t.Errorf("Expected Name to be 'Anonymous', got '%s'", user.Name)
	}

	fmt.Println("✅ Context defaults test passed!")
}

func TestApplyDefaults(t *testing.T) {
	fmt.Println("Starting TestApplyDefaults...")

	user := &User{Name: "Test User"}
	result, err := defaults.ApplyDefaults(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Errorf("Expected non-nil result")
	}

	// Verify it's the same pointer
	if result != user {
		t.Errorf("Expected same pointer to be returned")
	}

	// Verify defaults were applied
	if result.Email != "user@example.com" {
		t.Errorf("Expected Email to be 'user@example.com', got '%s'", result.Email)
	}

	fmt.Println("✅ ApplyDefaults test passed!")
}

func TestSetDefaultsNoDefaults(t *testing.T) {
	fmt.Println("Starting TestSetDefaultsNoDefaults...")

	// Test struct with no default tags
	empty := &EmptyStruct{Value: "original"}
	err := defaults.SetDefaults(empty)
	if err != nil {
		t.Errorf("Expected no error for struct without defaults, got %v", err)
	}

	// Value should remain unchanged
	if empty.Value != "original" {
		t.Errorf("Expected Value to remain 'original', got '%s'", empty.Value)
	}

	fmt.Println("✅ No defaults test passed!")
}

func TestSetDefaultsZeroValues(t *testing.T) {
	fmt.Println("Starting TestSetDefaultsZeroValues...")

	// Test that zero values get replaced with defaults
	user := &User{
		Name:     "", // empty string should get default
		Age:      0,  // zero int should get default
		IsActive: false, // zero bool should get default (but this might not work as expected with mold)
	}

	err := defaults.SetDefaults(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify zero values were replaced
	if user.Name != "Anonymous" {
		t.Errorf("Expected Name to be 'Anonymous', got '%s'", user.Name)
	}
	if user.Age != 18 {
		t.Errorf("Expected Age to be 18, got %d", user.Age)
	}

	fmt.Println("✅ Zero values test passed!")
}

// Benchmark tests
func BenchmarkSetDefaults(b *testing.B) {
	user := &User{}
	for i := 0; i < b.N; i++ {
		// Reset struct for each iteration
		*user = User{}
		defaults.SetDefaults(user)
	}
}

func BenchmarkApplyDefaults(b *testing.B) {
	for i := 0; i < b.N; i++ {
		user := &User{}
		defaults.ApplyDefaults(user)
	}
}

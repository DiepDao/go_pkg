package validate_test

import (
	"fmt"
	"testing"

	"github.com/DiepDao/pkg/validate"
)

// Define a nested struct with validation tags
type Address struct {
	City string `validate:"required"`
	Zip  int    `validate:"required,min=10000,max=99999"`
}

type User struct {
	Name    string  `validate:"required"`
	Email   string  `validate:"required,email"`
	Age     int     `validate:"required,min=10,max=20"`
	Phone   string  `validate:"omitempty,e164"` // Optional field
	Address Address `validate:"required"`
}

func TestValidateStruct(t *testing.T) {
	fmt.Println("Starting TestValidateStruct...") // Log start

	// Test case: Valid input (without phone)
	validUser := User{
		Name:  "Alice",
		Email: "alice@example.com",
		Age:   15,
		Address: Address{
			City: "New York",
			Zip:  12345,
		},
		Phone: "", // ✅ Empty phone field (optional)
	}
	err := validate.ValidateStruct(validUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test case: Valid phone number
	validUserWithPhone := User{
		Name:  "Bob",
		Email: "bob@example.com",
		Age:   18,
		Address: Address{
			City: "Boston",
			Zip:  54321,
		},
		Phone: "+1234567890", // ✅ Valid E.164 phone format
	}
	err = validate.ValidateStruct(validUserWithPhone)
	if err != nil {
		t.Errorf("Expected no error, got nil")
	}

	// Test case: Invalid phone format
	invalidUserWithPhone := User{
		Name:  "Charlie",
		Email: "charlie@example.com",
		Age:   20,
		Address: Address{
			City: "Chicago",
			Zip:  12345,
		},
		Phone: "invalid-phone", // ❌ Not a valid E.164 format
	}
	err = validate.ValidateStruct(invalidUserWithPhone)
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

	fmt.Println("Finished TestValidateStruct.") // Log end
}

func TestValidateStructFailures(t *testing.T) {
	fmt.Println("Starting TestValidateStructFailures...") // Log start

	// Test case: Completely empty user (should fail)
	emptyUser := User{}
	err := validate.ValidateStruct(emptyUser)
	fmt.Println("Validation error (Empty user):", err) // Print error
	if err == nil {
		t.Errorf("Expected validation error for empty user, got nil")
	} else {
		fmt.Println("✅ Empty user test passed:", err.Error())
	}

	// Test case: Missing required fields
	invalidUser := User{
		Name:  "", // ❌ Missing Name
		Email: "", // ❌ Missing Email
		Age:   15, // ✅ Valid Age
		Address: Address{
			City: "",  // ❌ Missing City
			Zip:  123, // ❌ Invalid Zip Code (below min)
		},
	}
	err = validate.ValidateStruct(invalidUser)
	fmt.Println("Validation error (Missing fields):", err) // Print error
	if err == nil {
		t.Errorf("Expected validation error due to missing fields, got nil")
	}

	// Test case: Invalid email format
	invalidEmailUser := User{
		Name:  "John",
		Email: "invalid-email", // ❌ Not a valid email format
		Age:   16,
		Address: Address{
			City: "Los Angeles",
			Zip:  90210,
		},
		Phone: "+1234567890",
	}
	err = validate.ValidateStruct(invalidEmailUser)
	fmt.Println("Validation error (Invalid email):", err) // Print error
	if err == nil {
		t.Errorf("Expected validation error due to invalid email, got nil")
	}

	// Test case: Age out of range
	outOfRangeAgeUser := User{
		Name:  "David",
		Email: "david@example.com",
		Age:   25, // ❌ Exceeds max age limit
		Address: Address{
			City: "Seattle",
			Zip:  98101,
		},
	}
	err = validate.ValidateStruct(outOfRangeAgeUser)
	fmt.Println("Validation error (Age out of range):", err) // Print error
	if err == nil {
		t.Errorf("Expected validation error due to age being out of range, got nil")
	}

	fmt.Println("Finished TestValidateStructFailures.") // Log end
}

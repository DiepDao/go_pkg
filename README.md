# Go Package Collection

A collection of utility packages for Go applications.

## Packages

### validate

A comprehensive validation package for HTTP request validation and struct validation with custom rules.

### defaults

A utility package for setting default values in Go structs using struct tags and the mold library.

## Installation

```bash
go get github.com/DiepDao/go_pkg
```

## validate Package

The `validate` package provides robust validation for HTTP requests and Go structs with support for custom validation rules and detailed error messages.

### Features

- **Request Schema Validation**: Validates JSON request bodies against struct definitions
- **Struct Field Validation**: Validates struct fields using validation tags
- **Custom Error Types**: Structured error handling with detailed messages
- **Security**: Built-in request body size limits (1MB)
- **Performance**: Singleton validator instance for optimal performance
- **Custom Validation Rules**: Support for custom validation like `notblank`
- **Array Payload Validation**: Validate JSON arrays of objects using `ValidateRequestSchemaArray` and per-item rules with `ValidateStructArrayRules` (aggregates errors with element indices)

### Quick Start

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
    
    "github.com/DiepDao/go_pkg/validate"
)

// Define your struct with validation tags
type User struct {
    Name    string `json:"name" validate:"required,notblank"`
    Email   string `json:"email" validate:"required,email"`
    Age     int    `json:"age" validate:"required,min=18,max=100"`
    Phone   string `json:"phone,omitempty" validate:"omitempty,e164"`
}

func handleUser(w http.ResponseWriter, r *http.Request) {
    var user User
    
    // Validate request schema and unmarshal JSON
    if err := validate.ValidateRequestSchema(&user, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Validate struct fields
    if err := validate.ValidateStructRules(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process valid user...
    fmt.Fprintf(w, "User %s created successfully!", user.Name)
}
```

### Basic Usage Examples

#### 1. HTTP Request Validation

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    
    // This validates the JSON schema AND unmarshals the data
    if err := validate.ValidateRequestSchema(&user, r); err != nil {
        // Handle schema errors (unknown fields, invalid JSON, etc.)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // This validates the struct fields using validation tags
    if err := validate.ValidateStructRules(user); err != nil {
        // Handle validation errors (required fields, format validation, etc.)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // User is valid, proceed with business logic
    saveUser(user)
    w.WriteHeader(http.StatusCreated)
}
```

#### 2. Struct-Only Validation

```go
func validateUserData() {
    user := User{
        Name:  "John Doe",
        Email: "john@example.com",
        Age:   25,
        Phone: "+1234567890",
    }
    
    if err := validate.ValidateStructRules(user); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
        return
    }
    
    fmt.Println("User data is valid!")
}
```

#### 3. Array of Objects (HTTP + per-item rules)

```go
// Reuse the User struct from above

func bulkCreateUsers(w http.ResponseWriter, r *http.Request) {
    var users []User

    // 1) Validate schema for an array payload and unmarshal into the slice
    if err := validate.ValidateRequestSchemaArray(&users, r); err != nil {
        // SchemaError for unknown fields, bad JSON, wrong casing, etc.
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 2) Validate each element using struct tags, with indexed aggregation
    if err := validate.ValidateStructArrayRules(users); err != nil {
        // ValidationError with messages like: "item[2]: validation failed: ..."
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    fmt.Fprintf(w, "Created %d users", len(users))
}
```

### Advanced Usage Examples

#### 1. Nested Struct Validation

```go
type Address struct {
    Street  string `json:"street" validate:"required"`
    City    string `json:"city" validate:"required"`
    ZipCode string `json:"zip_code" validate:"required,len=5"`
}

type UserWithAddress struct {
    Name    string  `json:"name" validate:"required,notblank"`
    Email   string  `json:"email" validate:"required,email"`
    Address Address `json:"address" validate:"required"`
}

func handleUserWithAddress(w http.ResponseWriter, r *http.Request) {
    var user UserWithAddress
    
    // Validates nested structures automatically
    if err := validate.ValidateRequestSchema(&user, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    if err := validate.ValidateStructRules(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Both user and address are validated
    fmt.Fprintf(w, "User %s at %s is valid!", user.Name, user.Address.City)
}
```

#### 2. Optional Fields with Validation

```go
type Product struct {
    Name        string  `json:"name" validate:"required,notblank"`
    Price       float64 `json:"price" validate:"required,min=0"`
    Description string  `json:"description,omitempty" validate:"omitempty,max=500"`
    Category    string  `json:"category,omitempty" validate:"omitempty,oneof=electronics books clothing"`
}
```

### Error Handling

The package provides structured error types for better error handling:

```go
func handleValidationErrors(err error) {
    switch e := err.(type) {
    case validate.ValidationError:
        // Field validation errors
        fmt.Printf("Validation Error: %s\n", e.Message)
        if e.Field != "" {
            fmt.Printf("Field: %s\n", e.Field)
        }
        
    case validate.SchemaError:
        // Schema validation errors (unknown fields, JSON parsing, etc.)
        fmt.Printf("Schema Error: %s\n", e.Message)
        if e.Field != "" {
            fmt.Printf("Field: %s\n", e.Field)
        }
        
    default:
        // Other errors
        fmt.Printf("Error: %v\n", err)
    }
}
```

Note:
- For array schema checks, `ValidateRequestSchemaArray` reports the first offending element's index in the error message.
- For per-item validation, `ValidateStructArrayRules` aggregates failures and prefixes each with `item[i]`.

### Common Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must be present and not zero value | `validate:"required"` |
| `notblank` | String field cannot be empty or whitespace only | `validate:"notblank"` |
| `email` | Must be valid email format | `validate:"email"` |
| `min` | Minimum value for numbers or length for strings | `validate:"min=18"` |
| `max` | Maximum value for numbers or length for strings | `validate:"max=100"` |
| `len` | Exact length for strings/slices | `validate:"len=5"` |
| `oneof` | Value must be one of specified options | `validate:"oneof=red blue green"` |
| `e164` | Valid E.164 phone number format | `validate:"e164"` |
| `omitempty` | Skip validation if field is empty | `validate:"omitempty,email"` |

### Example JSON Requests

#### Valid Request
```json
{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "age": 25,
    "phone": "+1234567890"
}
```

#### Invalid Requests

**Unknown field:**
```json
{
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "age": 25,
    "unknown_field": "value"
}
```
*Error: `schema validation failed: invalid field name or incorrect casing: unknown_field`*

**Invalid email:**
```json
{
    "name": "Alice Johnson",
    "email": "invalid-email",
    "age": 25
}
```
*Error: `validation failed: email must be a valid email`*

**Missing required field:**
```json
{
    "email": "alice@example.com",
    "age": 25
}
```
*Error: `validation failed: name is required`*

#### Valid Array Request
```json
[
  {"name": "Alice", "email": "alice@example.com", "age": 25, "address": {"city": "NYC", "zip": 10001}},
  {"name": "Bob",   "email": "bob@example.com",   "age": 30, "address": {"city": "Boston", "zip": 21008}}
]
```

#### Invalid Array Request (unknown field in one element)
```json
[
  {"name": "Alice", "email": "alice@example.com", "age": 25, "address": {"city": "NYC", "zip": 10001}},
  {"name": "Bob",   "email": "bob@example.com",   "age": 30, "address": {"city": "Boston", "zip": 21008}, "extra": "nope"}
]
```
*Error: `schema validation failed: invalid field name or incorrect casing at index 1: extra`*

### Migration from Old API

If you're using the deprecated functions, here's how to migrate:

```go
// Old way (still works but deprecated)
err := validate.CheckSchema(&user, request)
err = validate.EnforceSchemaRules(user)

// New recommended way
err := validate.ValidateRequestSchema(&user, request)
err = validate.ValidateStructRules(user)
```

```go
// Arrays of objects
var users []User
err := validate.ValidateRequestSchemaArray(&users, request)
err = validate.ValidateStructArrayRules(users)
```

### Security Features

- **Request Size Limiting**: Automatically limits request body to 1MB to prevent DoS attacks
- **Schema Validation**: Prevents injection of unexpected fields
- **Input Sanitization**: Built-in validation rules help prevent malformed data

## Development

### Running Tests

```bash
go test -v ./validate
```

### Publishing New Versions

```bash
git tag v1.0.0
git push origin v1.0.0
```

## defaults Package

The `defaults` package provides an easy way to set default values for struct fields using struct tags and the powerful mold library.

### Features

- **Struct Tag Support**: Use `mod` tags with `default=...` to specify default values
- **Type Safety**: Supports all basic Go types (string, int, bool, float, etc.)
- **Performance**: Singleton transformer instance for optimal performance
- **Error Handling**: Custom error types for better debugging
- **Generic Support**: Generic helper functions for modern Go patterns
- **Context Support**: Advanced default setting with custom context

### Quick Start

```go
package main

import (
    "fmt"
    "github.com/DiepDao/go_pkg/defaults"
)

type Config struct {
    Host     string `mod:"default=localhost"`
    Port     int    `mod:"default=8080"`
    Debug    bool   `mod:"default=false"`
    Timeout  int    `mod:"default=30"`
    Protocol string `mod:"default=http"`
}

func main() {
    config := &Config{}
    
    // Set default values
    if err := defaults.SetDefaults(config); err != nil {
        fmt.Printf("Error setting defaults: %v\n", err)
        return
    }
    
    fmt.Printf("Config: %+v\n", config)
    // Output: Config: {Host:localhost Port:8080 Debug:false Timeout:30 Protocol:http}
}
```

### Usage Examples

#### 1. Basic Default Setting

```go
type User struct {
    Name     string `mod:"default=Anonymous"`
    Age      int    `mod:"default=18"`
    Email    string `mod:"default=user@example.com"`
    IsActive bool   `mod:"default=true"`
    Role     string `mod:"default=user"`
}

func createUser() {
    user := &User{}
    
    if err := defaults.SetDefaults(user); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("User: %+v\n", user)
    // Output: User: {Name:Anonymous Age:18 Email:user@example.com IsActive:true Role:user}
}
```

#### 2. Partial Default Setting

```go
func createPartialUser() {
    // Some fields already have values
    user := &User{
        Name: "John Doe",
        Age:  25,
        // Email, IsActive, Role will get defaults
    }
    
    if err := defaults.SetDefaults(user); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("User: %+v\n", user)
    // Output: User: {Name:John Doe Age:25 Email:user@example.com IsActive:true Role:user}
}
```

#### 3. Using Generic Helper

```go
func createUserWithGeneric() {
    user := &User{Name: "Jane"}
    
    // ApplyDefaults returns the modified struct
    result, err := defaults.ApplyDefaults(user)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Result: %+v\n", result)
    // result is the same pointer as user, but with defaults applied
}
```

#### 4. Advanced Usage with Context

```go
func createUserWithContext() {
    user := &User{}
    // Use a real context; replace the key with your own typed key if needed
    ctx := context.WithValue(context.Background(), "environment", "production")
    
    if err := defaults.SetDefaultsWithContext(user, ctx); err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("User: %+v\n", user)
}
```

#### 5. Integration with HTTP Handlers

```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    
    // First, validate and unmarshal the request
    if err := validate.ValidateRequestSchema(&user, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Set defaults for any missing fields
    if err := defaults.SetDefaults(&user); err != nil {
        http.Error(w, "Failed to set defaults: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Validate the final struct (with defaults)
    if err := validate.ValidateStructRules(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Process user with defaults applied
    fmt.Fprintf(w, "User created: %+v", user)
}
```

### Supported Default Types

| Go Type | Example Tag | Description |
|---------|-------------|-------------|
| `string` | `mod:"default=hello"` | String literals |
| `int`, `int32`, `int64` | `mod:"default=42"` | Integer numbers |
| `float32`, `float64` | `mod:"default=3.14"` | Floating point numbers |
| `bool` | `mod:"default=true"` | Boolean values (true/false) |
| `uint`, `uint32`, `uint64` | `mod:"default=100"` | Unsigned integers |

### Error Handling

The defaults package provides structured error handling:

```go
func handleDefaultsErrors(err error) {
    var defaultsErr defaults.DefaultsError
    if errors.As(err, &defaultsErr) {
        fmt.Printf("Defaults Error: %s\n", defaultsErr.Message)
        if defaultsErr.Field != "" {
            fmt.Printf("Field: %s\n", defaultsErr.Field)
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
}
```

### Common Use Cases

#### Configuration Structs
```go
type DatabaseConfig struct {
    Host         string `mod:"default=localhost"`
    Port         int    `mod:"default=5432"`
    Database     string `mod:"default=myapp"`
    MaxConns     int    `mod:"default=10"`
    Timeout      int    `mod:"default=30"`
    SSLMode      string `mod:"default=disable"`
    AutoMigrate  bool   `mod:"default=true"`
}
```

#### API Response Defaults
```go
type APIResponse struct {
    Status  string `mod:"default=success"`
    Code    int    `mod:"default=200"`
    Message string `mod:"default=Operation completed successfully"`
    Data    interface{} // No default needed
}
```

#### User Preferences
```go
type UserPreferences struct {
    Theme         string `mod:"default=light"`
    Language      string `mod:"default=en"`
    Notifications bool   `mod:"default=true"`
    PageSize      int    `mod:"default=20"`
    Timezone      string `mod:"default=UTC"`
}
```

### Running Tests

```bash
go test -v ./defaults
```

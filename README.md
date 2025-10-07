# Go Package Collection

A collection of utility packages for Go applications.

## Packages

### validate

A comprehensive validation package for HTTP request validation and struct validation with custom rules.

### defaults

A utility package for setting default values in Go structs using struct tags and the mold library.

### array

A generic utility package for array/slice operations with type-safe filtering functions, similar to JavaScript's Array methods.

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

#### Run All Package Tests (Detailed)

Use the comprehensive test runner for detailed failure analysis:

```bash
./test-all.sh
```

**Features:**
- 🔍 Auto-discovers all packages with tests
- 📊 Shows detailed test results per package
- ❌ Highlights failed tests with:
  - Function name that failed
  - File and line number
  - Failure reason/error message
- 📈 Provides summary statistics
- 🎨 Color-coded output for easy reading

**Example output for failures:**
```
❌ Test: TestValidateRequestSchema
   Location: sanitize_test.go:287
   Reason:
     Expected SchemaError type, got validate.SchemaError
```

#### Run All Package Tests (Quick)

For a quick test run without detailed analysis:

```bash
./test-quick.sh
```

#### Run Individual Package Tests

```bash
# Test specific package
go test -v ./validate
go test -v ./defaults
go test -v ./array

# Run with benchmarks
go test -v -bench=. ./array
```

#### Test Coverage Report

Generate detailed coverage reports for all packages:

```bash
# Generate coverage report
./test-coverage.sh

# Generate coverage report with HTML output
./test-coverage.sh --html
```

**Features:**
- 📊 Coverage percentage for each package
- 🎨 Color-coded results (green ≥80%, yellow 60-79%, red <60%)
- 📁 Saves coverage files in `coverage/` directory
- 🌐 Optional HTML report for detailed visualization

#### Watch Mode (Auto-run tests on file changes)

Automatically run tests when Go files change:

```bash
./test-watch.sh
```

**Note:** Requires `inotify-tools` package:
```bash
# Install on Ubuntu/Debian
sudo apt-get install inotify-tools

# Install on macOS
brew install fswatch
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

## array Package

The `array` package provides generic, type-safe utility functions for slice operations in Go, inspired by JavaScript's Array methods. All functions work with any type using Go generics.

### Features

- **Generic Type Support**: Works with any type (structs, strings, ints, floats, custom types, etc.)
- **Type Safety**: Full compile-time type checking with Go generics
- **Immutable Operations**: Original slices are never modified
- **JavaScript-like API**: Familiar filter functions similar to Array.filter()
- **Performance**: Optimized implementations with pre-allocated capacity
- **Zero Dependencies**: Pure Go implementation

### Quick Start

```go
package main

import (
    "fmt"
    "github.com/DiepDao/go_pkg/array"
)

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // Filter even numbers
    evens := array.Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // Output: [2 4 6 8 10]
    
    // Map: double each number
    doubled := array.Map(numbers, func(n int) int {
        return n * 2
    })
    fmt.Println(doubled) // Output: [2 4 6 8 10 12 14 16 18 20]
    
    // Chain: filter evens, then double them
    result := array.Map(evens, func(n int) int {
        return n * 2
    })
    fmt.Println(result) // Output: [4 8 12 16 20]
}
```

### Functions

#### Filter

`Filter[T any](in []T, predicate func(T) bool) []T`

Returns a new slice containing only elements for which the predicate returns true. Similar to JavaScript's `Array.filter()`.

**Parameters:**
- `in`: Input slice of any type
- `predicate`: Function that tests each element (returns true to keep, false to filter out)

**Returns:** New filtered slice

**Example - Filter Integers:**
```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Get even numbers
evens := array.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// Result: [2, 4, 6, 8, 10]

// Get numbers greater than 5
large := array.Filter(numbers, func(n int) bool {
    return n > 5
})
// Result: [6, 7, 8, 9, 10]
```

**Example - Filter Strings:**
```go
words := []string{"apple", "banana", "apricot", "cherry", "avocado"}

// Get words starting with 'a'
aWords := array.Filter(words, func(s string) bool {
    return len(s) > 0 && s[0] == 'a'
})
// Result: ["apple", "apricot", "avocado"]

// Get words longer than 6 characters
longWords := array.Filter(words, func(s string) bool {
    return len(s) > 6
})
// Result: ["apricot", "avocado"]
```

**Example - Filter Structs:**
```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 17},
    {Name: "Charlie", Age: 30},
    {Name: "Diana", Age: 16},
}

// Filter adults (age >= 18)
adults := array.Filter(people, func(p Person) bool {
    return p.Age >= 18
})
// Result: [{Alice 25}, {Charlie 30}]
```

**Example - Filter with Type Any (Unknown Struct Types):**
```go
// When you have a slice of any and don't know the struct type
mixed := []any{1, "hello", 2.5, "world", 42, true, "test"}

// Filter only strings using type assertion
strings := array.Filter(mixed, func(v any) bool {
    _, ok := v.(string)
    return ok
})
// Result: ["hello", "world", "test"]

// Filter only integers
integers := array.Filter(mixed, func(v any) bool {
    _, ok := v.(int)
    return ok
})
// Result: [1, 42]
```

#### FilterIndexed

`FilterIndexed[T any](in []T, predicate func(v T, index int, arr []T) bool) []T`

Like `Filter`, but the predicate receives additional parameters: the element's index and the original array. This matches JavaScript's `Array.filter((item, index, array) => boolean)` signature.

**Parameters:**
- `in`: Input slice of any type
- `predicate`: Function that receives (value, index, original array) and returns boolean

**Returns:** New filtered slice

**Example - Filter by Index:**
```go
numbers := []int{10, 20, 30, 40, 50}

// Get elements at even indices (0, 2, 4)
evenIndexElements := array.FilterIndexed(numbers, func(v int, index int, arr []int) bool {
    return index%2 == 0
})
// Result: [10, 30, 50]
```

**Example - Filter Using Original Array:**
```go
numbers := []int{5, 10, 15, 20, 25}

// Get elements greater than the first element
result := array.FilterIndexed(numbers, func(v int, index int, arr []int) bool {
    return v > arr[0]
})
// Result: [10, 15, 20, 25]

// Get elements less than the average
sum := 0
for _, n := range numbers {
    sum += n
}
avg := sum / len(numbers)

belowAvg := array.FilterIndexed(numbers, func(v int, index int, arr []int) bool {
    return v < avg
})
// Result: [5, 10]
```

**Example - Complex Filtering Logic:**
```go
type Product struct {
    ID    int
    Name  string
    Price float64
}

products := []Product{
    {ID: 1, Name: "Laptop", Price: 999.99},
    {ID: 2, Name: "Mouse", Price: 25.50},
    {ID: 3, Name: "Keyboard", Price: 75.00},
    {ID: 4, Name: "Monitor", Price: 299.99},
}

// Filter: price > $50 AND at odd index positions
result := array.FilterIndexed(products, func(p Product, index int, arr []Product) bool {
    return p.Price > 50 && index%2 == 1
})
// Result: [{Monitor 299.99}] (only index 3 matches both conditions)
```

#### Map

`Map[T any, R any](in []T, mapper func(T) R) []R`

Transforms each element in the input slice using the provided mapper function and returns a new slice with the transformed elements. Similar to JavaScript's `Array.map()`. The input and output types can be different.

**Parameters:**
- `in`: Input slice of any type
- `mapper`: Function that transforms each element (required, panics if nil)

**Returns:** New slice with transformed elements

**Example - Transform Same Type:**
```go
numbers := []int{1, 2, 3, 4, 5}

// Double each number
doubled := array.Map(numbers, func(n int) int {
    return n * 2
})
// Result: [2, 4, 6, 8, 10]

// Square each number
squared := array.Map(numbers, func(n int) int {
    return n * n
})
// Result: [1, 4, 9, 16, 25]
```

**Example - Type Conversion (Int to String):**
```go
numbers := []int{1, 2, 3, 4, 5}

// Convert to strings
strings := array.Map(numbers, func(n int) string {
    return fmt.Sprintf("Number: %d", n)
})
// Result: ["Number: 1", "Number: 2", "Number: 3", "Number: 4", "Number: 5"]
```

**Example - Type Conversion (String to Int):**
```go
words := []string{"hello", "world", "go", "test"}

// Get length of each string
lengths := array.Map(words, func(s string) int {
    return len(s)
})
// Result: [5, 5, 2, 4]
```

**Example - Extract Struct Field:**
```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {Name: "Alice", Age: 25},
    {Name: "Bob", Age: 30},
    {Name: "Charlie", Age: 35},
}

// Extract just the names
names := array.Map(people, func(p Person) string {
    return p.Name
})
// Result: ["Alice", "Bob", "Charlie"]

// Extract just the ages
ages := array.Map(people, func(p Person) int {
    return p.Age
})
// Result: [25, 30, 35]
```

**Example - Transform Struct to Different Struct:**
```go
type User struct {
    ID   int
    Name string
}

type UserDTO struct {
    UserID   int
    FullName string
}

users := []User{
    {ID: 1, Name: "Alice"},
    {ID: 2, Name: "Bob"},
}

// Transform to DTOs
dtos := array.Map(users, func(u User) UserDTO {
    return UserDTO{
        UserID:   u.ID,
        FullName: u.Name,
    }
})
// Result: [{UserID: 1, FullName: "Alice"}, {UserID: 2, FullName: "Bob"}]
```

#### MapIndexed

`MapIndexed[T any, R any](in []T, mapper func(v T, index int, arr []T) R) []R`

Like `Map`, but the mapper receives additional parameters: the element's index and the original array. This matches JavaScript's `Array.map((item, index, array) => newValue)` signature.

**Parameters:**
- `in`: Input slice of any type
- `mapper`: Function that receives (value, index, original array) and returns transformed value

**Returns:** New slice with transformed elements

**Example - Use Index in Transformation:**
```go
numbers := []int{10, 20, 30, 40, 50}

// Multiply each number by its index
result := array.MapIndexed(numbers, func(v int, index int, arr []int) int {
    return v * index
})
// Result: [0, 20, 60, 120, 200]
```

**Example - Format with Index:**
```go
words := []string{"apple", "banana", "cherry"}

// Create numbered list
numbered := array.MapIndexed(words, func(v string, index int, arr []string) string {
    return fmt.Sprintf("%d: %s", index+1, v)
})
// Result: ["1: apple", "2: banana", "3: cherry"]
```

**Example - Use Original Array in Transformation:**
```go
numbers := []int{5, 10, 15, 20, 25}

// Calculate difference from first element
differences := array.MapIndexed(numbers, func(v int, index int, arr []int) int {
    return v - arr[0]
})
// Result: [0, 5, 10, 15, 20]

// Calculate percentage of total
sum := 0
for _, n := range numbers {
    sum += n
}
percentages := array.MapIndexed(numbers, func(v int, index int, arr []int) float64 {
    return float64(v) / float64(sum) * 100
})
// Result: [6.67, 13.33, 20.0, 26.67, 33.33] (approximately)
```

### Special Cases

#### Nil Predicate
If you pass `nil` as the predicate, both `Filter` and `FilterIndexed` return a shallow copy of the input slice:

```go
numbers := []int{1, 2, 3, 4, 5}
copy := array.Filter(numbers, nil)
// Result: [1, 2, 3, 4, 5] (new slice, same values)
```

#### Empty Slice
Filtering an empty slice returns an empty slice:

```go
empty := []int{}
result := array.Filter(empty, func(n int) bool {
    return n > 0
})
// Result: []
```

#### No Matches
If no elements match the predicate, an empty slice is returned:

```go
numbers := []int{1, 3, 5, 7, 9}
evens := array.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// Result: []
```

### Common Use Cases

#### API Response Filtering
```go
type User struct {
    ID       int
    Name     string
    IsActive bool
    Role     string
}

users := []User{
    {ID: 1, Name: "Alice", IsActive: true, Role: "admin"},
    {ID: 2, Name: "Bob", IsActive: false, Role: "user"},
    {ID: 3, Name: "Charlie", IsActive: true, Role: "user"},
}

// Get only active users
activeUsers := array.Filter(users, func(u User) bool {
    return u.IsActive
})

// Get only admins
admins := array.Filter(users, func(u User) bool {
    return u.Role == "admin"
})
```

#### Data Validation and Cleanup
```go
type Record struct {
    ID    int
    Value string
}

records := []Record{
    {ID: 1, Value: "valid"},
    {ID: 2, Value: ""},
    {ID: 3, Value: "also valid"},
    {ID: 4, Value: ""},
}

// Remove records with empty values
validRecords := array.Filter(records, func(r Record) bool {
    return r.Value != ""
})
```

#### Combining Filter and Map
```go
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Chain: filter evens, then double them
evens := array.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
doubled := array.Map(evens, func(n int) int {
    return n * 2
})
// Result: [4, 8, 12, 16, 20]

// Chain: filter > 5, then convert to strings
large := array.Filter(numbers, func(n int) bool {
    return n > 5
})
strings := array.Map(large, func(n int) string {
    return fmt.Sprintf("Value: %d", n)
})
// Result: ["Value: 6", "Value: 7", "Value: 8", "Value: 9", "Value: 10"]
```

#### Data Transformation Pipelines
```go
type Product struct {
    ID    int
    Name  string
    Price float64
}

products := []Product{
    {ID: 1, Name: "Laptop", Price: 999.99},
    {ID: 2, Name: "Mouse", Price: 25.50},
    {ID: 3, Name: "Keyboard", Price: 75.00},
}

// Filter expensive products, then extract names
expensive := array.Filter(products, func(p Product) bool {
    return p.Price > 50
})
names := array.Map(expensive, func(p Product) string {
    return p.Name
})
// Result: ["Laptop", "Keyboard"]
```

### Performance Notes

- **Pre-allocated Slices**: Map functions pre-allocate the exact size needed; Filter functions pre-allocate with capacity for better performance
- **Immutable**: Original slices are never modified; new slices are always returned
- **Type Flexibility**: Map supports transforming between different types (e.g., `[]int` to `[]string`)
- **Efficient Memory Usage**: Minimal allocations and efficient memory reuse

### Running Tests

```bash
# Run all tests
go test -v ./array

# Run with benchmarks
go test -v -bench=. ./array
```

# Go Package Collection

A collection of generic utility packages for Go applications, providing validation, defaults, and array/slice helpers.

## Table of Contents

- [Installation](#installation)
- [Packages Overview](#packages-overview)
- [validate](#validate-package) - HTTP request & struct validation
- [defaults](#defaults-package) - Struct default values via tags
- [array](#array-package) - Generic slice utilities (Filter, Map, Includes)
- [Development](#development)

---

## Installation

### Requirements

- **Go 1.23+** (uses generics)

### Install the module

```bash
go get github.com/DiepDao/go_pkg@latest
```

This adds the dependency to your `go.mod` and downloads it. To pin a specific version:

```bash
go get github.com/DiepDao/go_pkg@v1.0.19
```

Then import the sub-packages you need:

```go
import (
    "github.com/DiepDao/go_pkg/array"
    "github.com/DiepDao/go_pkg/defaults"
    "github.com/DiepDao/go_pkg/validate"
)
```

After adding imports, run:

```bash
go mod tidy
```

---

## Packages Overview

| Package      | Import Path                          | Description                                                   |
| ------------ | ------------------------------------ | ------------------------------------------------------------- |
| **validate** | `github.com/DiepDao/go_pkg/validate` | HTTP request body & struct field validation with custom rules |
| **defaults** | `github.com/DiepDao/go_pkg/defaults` | Set default values on struct fields via `mod` tags            |
| **array**    | `github.com/DiepDao/go_pkg/array`    | Generic `Filter`, `Map`, `Includes` for slices (JS-style API) |

---

## validate Package

Validates HTTP JSON request bodies against struct definitions and validates struct fields using tags. Built on [go-playground/validator](https://github.com/go-playground/validator).

### Exported API

| Function                     | Signature                               | Description                                 |
| ---------------------------- | --------------------------------------- | ------------------------------------------- |
| `ValidateRequestSchema`      | `(req any, r *http.Request) error`      | Validate & unmarshal a JSON object body     |
| `ValidateRequestSchemaArray` | `(slicePtr any, r *http.Request) error` | Validate & unmarshal a JSON array body      |
| `ValidateStructRules`        | `(req any) error`                       | Validate struct fields using tags           |
| `ValidateStructArrayRules`   | `(slice any) error`                     | Validate each element in a slice using tags |

**Error types:** `ValidationError` (field rule failures) and `SchemaError` (unknown fields, bad JSON, size limit).

### Quick Start

```go
import "github.com/DiepDao/go_pkg/validate"

type User struct {
    Name  string `json:"name"  validate:"required,notblank"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age"   validate:"required,min=18,max=100"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User

    // 1. Validate JSON schema + unmarshal
    if err := validate.ValidateRequestSchema(&user, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 2. Validate field rules
    if err := validate.ValidateStructRules(user); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // user is valid
}
```

### Array Payload Validation

```go
func bulkCreateUsers(w http.ResponseWriter, r *http.Request) {
    var users []User

    if err := validate.ValidateRequestSchemaArray(&users, r); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := validate.ValidateStructArrayRules(users); err != nil {
        // errors prefixed with "item[i]: ..."
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
}
```

### Struct-Only Validation (no HTTP)

```go
user := User{Name: "Alice", Email: "alice@example.com", Age: 25}
if err := validate.ValidateStructRules(user); err != nil {
    log.Fatal(err)
}
```

### Error Handling

```go
switch e := err.(type) {
case *validate.SchemaError:
    // unknown field, bad JSON, body too large
    fmt.Println(e.Field, e.Message)
case *validate.ValidationError:
    // field rule failures
    fmt.Println(e.Field, e.Message)
}
```

### Supported Validation Tags

| Tag         | Description                       | Example                           |
| ----------- | --------------------------------- | --------------------------------- |
| `required`  | Must be present and non-zero      | `validate:"required"`             |
| `notblank`  | String cannot be empty/whitespace | `validate:"notblank"`             |
| `email`     | Valid email format                | `validate:"email"`                |
| `min`       | Minimum value/length              | `validate:"min=18"`               |
| `max`       | Maximum value/length              | `validate:"max=100"`              |
| `len`       | Exact length                      | `validate:"len=5"`                |
| `oneof`     | Must be one of listed values      | `validate:"oneof=red blue green"` |
| `e164`      | E.164 phone number                | `validate:"e164"`                 |
| `omitempty` | Skip validation if empty          | `validate:"omitempty,email"`      |

### Security

- Request body capped at **1 MB** to prevent DoS
- Unknown/extra JSON fields are **rejected**
- Custom `notblank` rule prevents whitespace-only strings

---

## defaults Package

Set default values on struct fields using `mod:"default=..."` tags. Built on [go-playground/mold](https://github.com/go-playground/mold).

### Exported API

| Function                 | Signature                                       | Description                                  |
| ------------------------ | ----------------------------------------------- | -------------------------------------------- |
| `SetDefaults`            | `(data interface{}) error`                      | Apply defaults to a struct pointer           |
| `SetDefaultsWithContext` | `(data interface{}, ctx context.Context) error` | Apply defaults with a custom context         |
| `ApplyDefaults[T]`       | `(data *T) (*T, error)`                         | Generic helper that returns the pointer back |

**Error type:** `DefaultsError`

### Quick Start

```go
import "github.com/DiepDao/go_pkg/defaults"

type Config struct {
    Host    string `mod:"default=localhost"`
    Port    int    `mod:"default=8080"`
    Debug   bool   `mod:"default=false"`
    Timeout int    `mod:"default=30"`
}

func main() {
    cfg := &Config{}
    if err := defaults.SetDefaults(cfg); err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%+v\n", cfg)
    // {Host:localhost Port:8080 Debug:false Timeout:30}
}
```

### Partial Defaults

Fields that already have values are **not overwritten**:

```go
cfg := &Config{Host: "prod.example.com", Port: 443}
defaults.SetDefaults(cfg)
// Host stays "prod.example.com", Port stays 443
// Debug -> false, Timeout -> 30
```

### Generic Helper

```go
cfg, err := defaults.ApplyDefaults(&Config{})
// cfg is *Config with defaults applied
```

### Combined with validate

A common pattern is: unmarshal -> set defaults -> validate:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := validate.ValidateRequestSchema(&user, r); err != nil { ... }
    if err := defaults.SetDefaults(&user); err != nil { ... }
    if err := validate.ValidateStructRules(user); err != nil { ... }
    // proceed
}
```

### Supported Types

`string`, `int`, `int32`, `int64`, `uint`, `uint32`, `uint64`, `float32`, `float64`, `bool`

---

## array Package

Generic, type-safe slice utilities inspired by JavaScript's `Array.filter()`, `Array.map()`, and `Array.some()`. Zero dependencies, pure Go generics.

### Exported API

| Function          | Signature                                                | Description                        |
| ----------------- | -------------------------------------------------------- | ---------------------------------- |
| `Filter`          | `[T any](in []T, predicate func(T) bool) []T`            | Keep elements matching predicate   |
| `FilterIndexed`   | `[T any](in []T, predicate func(T, int, []T) bool) []T`  | Filter with index + array access   |
| `Map`             | `[T, R any](in []T, mapper func(T) R) []R`               | Transform each element             |
| `MapIndexed`      | `[T, R any](in []T, mapper func(T, int, []T) R) []R`     | Map with index + array access      |
| `Includes`        | `[T any](in []T, predicate func(T) bool) bool`           | True if any element matches        |
| `IncludesIndexed` | `[T any](in []T, predicate func(T, int, []T) bool) bool` | Includes with index + array access |

### Filter

Returns a new slice with elements that pass the predicate. Original slice is never modified.

```go
import "github.com/DiepDao/go_pkg/array"

numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

evens := array.Filter(numbers, func(n int) bool {
    return n%2 == 0
})
// [2, 4, 6, 8, 10]
```

With structs:

```go
type User struct {
    Name     string
    IsActive bool
}

users := []User{
    {Name: "Alice", IsActive: true},
    {Name: "Bob", IsActive: false},
    {Name: "Charlie", IsActive: true},
}

active := array.Filter(users, func(u User) bool {
    return u.IsActive
})
// [{Alice true}, {Charlie true}]
```

### FilterIndexed

Same as `Filter` but the callback also receives the index and original slice:

```go
// Keep elements at even indices
result := array.FilterIndexed(nums, func(v int, i int, arr []int) bool {
    return i%2 == 0
})
```

### Map

Transform each element. Input and output types can differ.

```go
names := array.Map(users, func(u User) string {
    return u.Name
})
// ["Alice", "Bob", "Charlie"]

doubled := array.Map([]int{1, 2, 3}, func(n int) int {
    return n * 2
})
// [2, 4, 6]
```

Struct-to-struct transformation:

```go
type UserDTO struct {
    FullName string
}

dtos := array.Map(users, func(u User) UserDTO {
    return UserDTO{FullName: u.Name}
})
```

### MapIndexed

Same as `Map` but the callback also receives the index and original slice:

```go
numbered := array.MapIndexed(words, func(v string, i int, arr []string) string {
    return fmt.Sprintf("%d. %s", i+1, v)
})
// ["1. apple", "2. banana", "3. cherry"]
```

### Includes

Returns `true` if **any** element matches. Short-circuits on first match.

```go
hasAdmin := array.Includes(users, func(u User) bool {
    return u.Role == "admin"
})

if array.Includes(orders, func(o Order) bool { return o.Status == "pending" }) {
    fmt.Println("Pending orders exist")
}
```

### IncludesIndexed

Same as `Includes` but the callback also receives the index and original slice:

```go
hasExpensiveAtOddIndex := array.IncludesIndexed(products, func(p Product, i int, arr []Product) bool {
    return i%2 == 1 && p.Price > 100
})
```

### Chaining Filter + Map

```go
// Filter expensive products, then extract names
names := array.Map(
    array.Filter(products, func(p Product) bool { return p.Price > 50 }),
    func(p Product) string { return p.Name },
)
```

### Edge Cases

| Scenario               | Filter               | Map          | Includes        |
| ---------------------- | -------------------- | ------------ | --------------- |
| `nil` predicate/mapper | Returns shallow copy | Panics       | Returns `false` |
| Empty slice            | Returns `[]`         | Returns `[]` | Returns `false` |
| No matches             | Returns `[]`         | N/A          | Returns `false` |

### Performance

- **Filter** pre-allocates with `cap(in)` for fewer re-allocations
- **Map** pre-allocates exact output size
- **Includes** exits on first match
- All operations are **immutable** -- input slices are never modified

---

## Development

### Running Tests

```bash
# All packages (detailed output)
./test-all.sh

# Quick run
./test-quick.sh

# Individual packages
go test -v ./validate
go test -v ./defaults
go test -v ./array

# With benchmarks
go test -v -bench=. ./array
```

### Test Coverage

```bash
./test-coverage.sh          # terminal report
./test-coverage.sh --html   # HTML report in coverage/
```

### Watch Mode

Auto-run tests on file changes (requires `inotify-tools` on Linux or `fswatch` on macOS):

```bash
./test-watch.sh
```

### Publishing a New Version

```bash
git tag v1.0.20
git push origin v1.0.20
```

Users then update with:

```bash
go get github.com/DiepDao/go_pkg@v1.0.20
```

---

## License

See [LICENSE](LICENSE) for details.

# Changelog

## [TBA]

### Added
- validate: `ValidateRequestSchemaArray` to validate JSON arrays of objects against struct `json` tags and unmarshal into slices
- validate: `ValidateStructArrayRules` to validate each array/slice element using tags and aggregate errors with item indices
- Documentation: README updated with array-of-objects examples, error handling notes, and migration snippet
- Tests: Unit tests covering array schema and per-item validation in `validate/sanitize_test.go`

## [1.0.14]
## [13-09-2025]

### Added
- **New `defaults` package**: Provides easy default value setting for struct fields using struct tags
- Support for setting default values using the `default` tag with mold library integration
- `SetDefaults()` function for applying defaults to struct pointers
- `SetDefaultsWithContext()` function for advanced default setting with custom context
- `ApplyDefaults()` generic helper function for functional programming patterns
- Custom `DefaultsError` type for structured error handling
- Comprehensive unit tests for the defaults package
- Documentation and examples for the defaults package in README.md

### Updated
- **validate package unit tests**: Updated all tests to use new function names (`ValidateStructRules`, `ValidateRequestSchema`)
- Added proper testing of custom error types (`ValidationError`, `SchemaError`) using `errors.As()`
- Added `TestBackwardCompatibility()` to ensure deprecated functions still work
- Improved test function names to match the new API while maintaining comprehensive coverage


### Dependencies
- Added `github.com/go-playground/mold/v4` dependency for default value transformation

## [1.0.13]
## [20-05-2025]

### Added
- check schema new flag to ensure no blank data

## [1.0.9]
## [20-05-2025]

### Changed
- error returns lowercase to field required field on validation

## [1.0.7]

## [20-05-2025]

### Changed
- library clean up


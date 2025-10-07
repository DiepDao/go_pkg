# Testing Guide

Complete guide to testing in this Go package repository.

## Quick Start

```bash
# Run all tests with detailed failure analysis
./test-all.sh

# Quick test run
./test-quick.sh

# Generate coverage report
./test-coverage.sh

# Watch mode (auto-run on file changes)
./test-watch.sh
```

## Available Test Scripts

| Script | Purpose | When to Use |
|--------|---------|-------------|
| `test-all.sh` | Comprehensive test runner with detailed failure analysis | When tests fail and you need to debug |
| `test-quick.sh` | Quick test runner with minimal output | During development for rapid feedback |
| `test-coverage.sh` | Generate coverage reports | Before committing or during code review |
| `test-watch.sh` | Auto-run tests on file changes | During active TDD development |

## Detailed Script Usage

### test-all.sh - Comprehensive Test Runner

**Best for:** Debugging test failures

```bash
./test-all.sh
```

**Output includes:**
- ✅ Package-by-package test results
- ❌ Failed test function names
- 📍 Exact file and line numbers
- 📝 Error messages and reasons
- 📊 Summary statistics

**Example failure output:**
```
❌ Test: TestValidateRequestSchema
   Location: sanitize_test.go:287
   Reason:
     Expected SchemaError type, got validate.SchemaError
```

### test-quick.sh - Quick Test Runner

**Best for:** Rapid development feedback

```bash
./test-quick.sh
```

**Features:**
- Fast execution
- Minimal output
- Suggests detailed analysis on failure

### test-coverage.sh - Coverage Report

**Best for:** Code quality checks

```bash
# Basic report
./test-coverage.sh

# HTML report (opens in browser)
./test-coverage.sh --html
```

**Coverage thresholds:**
- 🟢 **High:** ≥80% (excellent)
- 🟡 **Medium:** 60-79% (good)
- 🔴 **Low:** <60% (needs improvement)

**Example output:**
```
Package                    Coverage
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
array                     100.0%
defaults                  72.0%
validate                  85.5%
```

### test-watch.sh - Watch Mode

**Best for:** TDD workflow

```bash
./test-watch.sh
```

**Requirements:**
```bash
# Ubuntu/Debian
sudo apt-get install inotify-tools
```

**Features:**
- Watches all `.go` files
- Auto-runs tests on save
- Shows timestamp of each run
- Press Ctrl+C to stop

## Testing Individual Packages

```bash
# Test a specific package
go test -v ./array
go test -v ./validate
go test -v ./defaults

# Run with benchmarks
go test -v -bench=. ./array

# Run with coverage
go test -v -cover ./array

# Run specific test
go test -v -run TestFilter_Integers ./array
```

## Development Workflows

### Workflow 1: Test-Driven Development (TDD)

```bash
# 1. Start watch mode
./test-watch.sh

# 2. Write failing test
# 3. Implement feature
# 4. See tests pass automatically
# 5. Refactor with confidence
```

### Workflow 2: Feature Development

```bash
# 1. Quick check before starting
./test-quick.sh

# 2. Develop feature
# 3. Run tests frequently
./test-quick.sh

# 4. When done, run comprehensive tests
./test-all.sh

# 5. Check coverage
./test-coverage.sh
```

### Workflow 3: Bug Fixing

```bash
# 1. Run comprehensive tests to see failure details
./test-all.sh

# 2. Identify failing test and location
# 3. Fix the bug
# 4. Verify fix
./test-all.sh

# 5. Ensure coverage didn't drop
./test-coverage.sh
```

### Workflow 4: Code Review / PR

```bash
# 1. Run all tests
./test-all.sh

# 2. Generate coverage report
./test-coverage.sh --html

# 3. Review coverage/coverage.html in browser
# 4. Identify uncovered code
# 5. Add tests for uncovered areas
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test Suite

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Make scripts executable
        run: chmod +x *.sh
      
      - name: Run Tests
        run: ./test-all.sh
      
      - name: Generate Coverage
        run: ./test-coverage.sh --html
      
      - name: Upload Coverage Report
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage/
      
      - name: Check Coverage Threshold
        run: |
          # Fail if any package has <60% coverage
          ./test-coverage.sh || exit 1
```

### GitLab CI Example

```yaml
test:
  image: golang:1.23
  script:
    - chmod +x *.sh
    - ./test-all.sh
    - ./test-coverage.sh --html
  artifacts:
    paths:
      - coverage/
    expire_in: 1 week
  coverage: '/total:.*?(\d+\.\d+)%/'
```

## Troubleshooting

### Tests Pass Locally But Fail in CI

**Possible causes:**
- Different Go versions
- Missing dependencies
- Environment-specific issues

**Solution:**
```bash
# Check Go version matches CI
go version

# Ensure all dependencies are present
go mod download
go mod verify
```

### Coverage Report Shows 0%

**Possible causes:**
- Tests are failing
- No test files in package

**Solution:**
```bash
# Check if tests pass first
./test-all.sh

# Verify test files exist
find . -name "*_test.go"
```

### Watch Mode Not Working

**Possible causes:**
- `inotify-tools` not installed
- File system doesn't support inotify

**Solution:**
```bash
# Install inotify-tools
sudo apt-get install inotify-tools

# Or use manual test runs
./test-quick.sh
```

## Best Practices

### 1. Run Tests Before Committing
```bash
./test-all.sh && git commit
```

### 2. Maintain High Coverage
- Aim for ≥80% coverage on critical packages
- Use `./test-coverage.sh --html` to find gaps

### 3. Write Descriptive Test Names
```go
// Good
func TestFilter_EmptySlice(t *testing.T) { ... }

// Bad
func TestFilter1(t *testing.T) { ... }
```

### 4. Test Edge Cases
- Empty inputs
- Nil values
- Boundary conditions
- Error conditions

### 5. Use Table-Driven Tests
```go
func TestFilter(t *testing.T) {
    tests := []struct {
        name     string
        input    []int
        expected []int
    }{
        {"empty", []int{}, []int{}},
        {"single", []int{1}, []int{1}},
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

## Package-Specific Notes

### array Package
- 100% test coverage ✅
- Includes benchmarks for performance testing
- Tests cover all generic type combinations

### defaults Package
- 72% test coverage
- Tests cover all default value types
- Includes error handling tests

### validate Package
- Tests currently have some failures (type assertion issues)
- Comprehensive validation scenarios
- Tests for both success and failure cases

## Summary

This repository includes a complete testing infrastructure:

- ✅ **4 test scripts** for different use cases
- 📊 **Coverage reporting** with HTML visualization
- 🔄 **Watch mode** for TDD workflow
- 🎨 **Color-coded output** for easy reading
- 📈 **Detailed failure analysis** with line numbers
- 🚀 **CI/CD ready** with proper exit codes

All scripts automatically discover packages and require no configuration!

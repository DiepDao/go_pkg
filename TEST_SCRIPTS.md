# Test Runner Scripts

This repository includes automated test runner scripts that make it easy to test all packages.

## Scripts

### `test-all.sh` - Comprehensive Test Runner

**Purpose:** Runs all package tests with detailed failure analysis and reporting.

**Usage:**
```bash
./test-all.sh
```

**Features:**
- ✅ Auto-discovers all packages containing test files
- 📦 Tests each package individually
- 🎨 Color-coded output (green for pass, red for fail)
- 📊 Detailed statistics (packages tested, tests passed/failed)
- 🔍 **Detailed failure reporting:**
  - Exact test function name that failed
  - File name and line number where the failure occurred
  - Error message and reason for failure
  - Context from the test output

**Example Output:**

When all tests pass:
```
========================================
  Go Package Test Runner
========================================

🔍 Discovering packages...

Found 3 package(s) with tests:
  • ./array
  • ./defaults
  • ./validate

========================================
  Running Tests
========================================

📦 Testing package: array
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
=== RUN   TestFilter_Integers
--- PASS: TestFilter_Integers (0.00s)
...
✅ PASSED

========================================
  Test Summary
========================================

Packages:
  Total:   3
  Passed:  3

Tests:
  Total:   33
  Passed:  33

✅ ALL TESTS PASSED
```

When tests fail:
```
📦 Testing package: validate
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
...
❌ FAILED

Analyzing failures...

  ❌ Test: TestValidateRequestSchema
     Location: sanitize_test.go:287
     Reason:
       Expected SchemaError type, got validate.SchemaError

========================================
  Test Summary
========================================

Packages:
  Total:   3
  Passed:  2
  Failed:  1

Tests:
  Total:   33
  Passed:  29
  Failed:  4

Failed Packages:
  ❌ validate

Failed Tests Details:
  ❌ validate
     Test: TestValidateRequestSchema
     Location: sanitize_test.go:287

❌ TESTS FAILED
```

**Exit Codes:**
- `0` - All tests passed
- `1` - One or more tests failed

---

### `test-quick.sh` - Quick Test Runner

**Purpose:** Runs all package tests quickly with minimal output.

**Usage:**
```bash
./test-quick.sh
```

**Features:**
- ⚡ Fast execution
- 📝 Minimal output (only test results)
- ✅ Simple pass/fail indication
- 💡 Suggests running `test-all.sh` for detailed analysis on failure

**Example Output:**

When all tests pass:
```
Running all package tests...

[test output]

✅ ALL TESTS PASSED
```

When tests fail:
```
Running all package tests...

[test output]

❌ SOME TESTS FAILED
Run ./test-all.sh for detailed failure analysis
```

**Exit Codes:**
- `0` - All tests passed
- `1` - One or more tests failed

---

### `test-coverage.sh` - Coverage Report Generator

**Purpose:** Generates detailed test coverage reports for all packages.

**Usage:**
```bash
# Basic coverage report
./test-coverage.sh

# Generate HTML coverage report
./test-coverage.sh --html
```

**Features:**
- 📊 Coverage percentage for each package
- 🎨 Color-coded results:
  - 🟢 Green: ≥80% coverage (high)
  - 🟡 Yellow: 60-79% coverage (medium)
  - 🔴 Red: <60% coverage (low)
- 📁 Saves coverage profiles in `coverage/` directory
- 🌐 HTML report generation for detailed visualization
- 📈 Overall statistics and breakdown

**Example Output:**
```
========================================
  Test Coverage Report
========================================

Generating coverage for 3 package(s)...

📦 array
   Coverage: 100.0%

📦 defaults
   Coverage: 72.0%

📦 validate
   Coverage: 85.5%

========================================
  Coverage Summary
========================================

Package                    Coverage
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
array                     100.0%
defaults                  72.0%
validate                  85.5%

Coverage files saved in: coverage/
Run with --html flag to generate HTML report

Overall Statistics:
  Packages tested: 3
  High coverage (≥80%): 2
  Medium coverage (60-79%): 1
```

**HTML Report:**
When using `--html` flag, opens a detailed HTML report showing:
- Line-by-line coverage
- Uncovered code blocks
- Coverage heatmap

**Exit Codes:**
- `0` - All packages tested successfully
- `1` - One or more packages failed tests

---

### `test-watch.sh` - Watch Mode Test Runner

**Purpose:** Automatically runs tests when Go files change (for development workflow).

**Usage:**
```bash
./test-watch.sh
```

**Features:**
- 👀 Watches all `.go` files in the repository
- 🔄 Automatically runs tests on file changes
- ⚡ Uses quick test mode for fast feedback
- 🕐 Shows timestamp of each test run
- ⌨️ Press Ctrl+C to stop

**Requirements:**
```bash
# Ubuntu/Debian
sudo apt-get install inotify-tools

# macOS (requires modification to use fswatch)
brew install fswatch
```

**Example Output:**
```
========================================
  Test Watch Mode
========================================

Watching for file changes in Go files...
Press Ctrl+C to stop

========================================
  Running Tests (14:25:30)
========================================

[test output]

✅ ALL TESTS PASSED

Watching for changes...
```

**Use Case:**
Perfect for TDD (Test-Driven Development) workflow - make changes and see results immediately.

---

## Workflow Recommendations

### During Active Development (TDD)
Use watch mode for instant feedback:
```bash
./test-watch.sh
```
Make changes, save files, and see test results automatically.

### Quick Validation
Use `test-quick.sh` for rapid feedback:
```bash
./test-quick.sh
```

### When Tests Fail
Use `test-all.sh` to get detailed failure information:
```bash
./test-all.sh
```
Shows exact function names, line numbers, and error reasons.

### Before Committing
Run comprehensive tests and check coverage:
```bash
./test-all.sh && ./test-coverage.sh
```

### Code Review / Quality Check
Generate HTML coverage report:
```bash
./test-coverage.sh --html
```
Open `coverage/coverage.html` in your browser to review uncovered code.

### CI/CD Integration
All scripts can be used in CI/CD pipelines. They exit with appropriate codes:

```yaml
# Example GitHub Actions
name: Test

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
      
      - name: Run Tests
        run: ./test-all.sh
      
      - name: Generate Coverage
        run: ./test-coverage.sh --html
      
      - name: Upload Coverage
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage/
```

---

## Package-Specific Testing

You can also test individual packages:

```bash
# Test a specific package
go test -v ./array
go test -v ./validate
go test -v ./defaults

# Run with benchmarks
go test -v -bench=. ./array

# Run with coverage
go test -v -cover ./...
```

---

## Troubleshooting

### Script Permission Denied
Make the scripts executable:
```bash
chmod +x test-all.sh test-quick.sh
```

### No Tests Found
Ensure your test files follow Go conventions:
- Files must end with `_test.go`
- Test functions must start with `Test`
- Files must be in a package directory

### Color Output Not Working
Some terminals may not support ANSI color codes. The scripts will still work, but without colors.

---

## Script Maintenance

The test scripts automatically discover packages, so:
- ✅ No updates needed when adding new packages
- ✅ No updates needed when adding new tests
- ✅ Works with any Go project structure

Just ensure your test files follow standard Go naming conventions (`*_test.go`).

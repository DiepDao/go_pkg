#!/bin/bash

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Counters
TOTAL_PACKAGES=0
PASSED_PACKAGES=0
FAILED_PACKAGES=0
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Arrays to store results
declare -a FAILED_PACKAGE_NAMES
declare -a FAILED_TEST_DETAILS

echo -e "${BOLD}${CYAN}========================================${NC}"
echo -e "${BOLD}${CYAN}  Go Package Test Runner${NC}"
echo -e "${BOLD}${CYAN}========================================${NC}\n"

# Find all packages with test files
echo -e "${BLUE}🔍 Discovering packages...${NC}\n"

# Get all directories containing *_test.go files
PACKAGES=$(find . -name "*_test.go" -type f | sed 's|/[^/]*$||' | sort -u | grep -v "^\.$")

if [ -z "$PACKAGES" ]; then
    echo -e "${RED}❌ No test files found!${NC}"
    exit 1
fi

# Count packages
TOTAL_PACKAGES=$(echo "$PACKAGES" | wc -l)
echo -e "${CYAN}Found ${BOLD}$TOTAL_PACKAGES${NC}${CYAN} package(s) with tests:${NC}"
echo "$PACKAGES" | while read pkg; do
    echo -e "  ${BLUE}•${NC} $pkg"
done
echo ""

echo -e "${BOLD}${CYAN}========================================${NC}"
echo -e "${BOLD}${CYAN}  Running Tests${NC}"
echo -e "${BOLD}${CYAN}========================================${NC}\n"

# Create temporary directory for test outputs
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Run tests for each package
while IFS= read -r package; do
    if [ -z "$package" ]; then
        continue
    fi
    
    PACKAGE_NAME=$(echo "$package" | sed 's|^\./||')
    echo -e "${BOLD}${BLUE}📦 Testing package: ${CYAN}$PACKAGE_NAME${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    
    # Run tests and capture output
    TEST_OUTPUT="$TEMP_DIR/${PACKAGE_NAME//\//_}.txt"
    
    # Run go test with verbose output and capture exit code properly
    set +e  # Temporarily disable exit on error
    go test -v "$package" 2>&1 | tee "$TEST_OUTPUT"
    TEST_EXIT_CODE=${PIPESTATUS[0]}
    set -e  # Re-enable exit on error
    
    if [ $TEST_EXIT_CODE -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED${NC}\n"
        PASSED_PACKAGES=$((PASSED_PACKAGES + 1))
        
        # Count passed tests
        PACKAGE_PASSED=$(grep -c "^--- PASS:" "$TEST_OUTPUT" 2>/dev/null || echo "0")
        PACKAGE_PASSED=$(echo "$PACKAGE_PASSED" | tr -d '\n' | tr -d ' ')
        PASSED_TESTS=$((PASSED_TESTS + PACKAGE_PASSED))
        TOTAL_TESTS=$((TOTAL_TESTS + PACKAGE_PASSED))
    else
        echo -e "${RED}❌ FAILED${NC}\n"
        FAILED_PACKAGES=$((FAILED_PACKAGES + 1))
        FAILED_PACKAGE_NAMES+=("$PACKAGE_NAME")
        
        # Parse failures
        echo -e "${YELLOW}Analyzing failures...${NC}\n"
        
        # Extract failed test names
        FAILED_TESTS_IN_PKG=$(grep "^--- FAIL:" "$TEST_OUTPUT" | awk '{print $3}')
        
        if [ -n "$FAILED_TESTS_IN_PKG" ] && [ "$FAILED_TESTS_IN_PKG" != "" ]; then
            while IFS= read -r test_name; do
                if [ -z "$test_name" ]; then
                    continue
                fi
                
                FAILED_TESTS=$((FAILED_TESTS + 1))
                TOTAL_TESTS=$((TOTAL_TESTS + 1))
                
                echo -e "${RED}  ❌ Test: ${BOLD}$test_name${NC}"
                
                # Extract error details for this specific test
                # Look for the test failure section
                ERROR_DETAILS=$(awk "
                    /^=== RUN   $test_name\$/ { capture=1; next }
                    /^--- FAIL: $test_name/ { if (capture) exit }
                    /^=== RUN/ { if (capture && \$3 != \"$test_name\") exit }
                    capture && /Error Trace:|Error:|Test:/ { print }
                    capture && /^\s+[a-zA-Z_][a-zA-Z0-9_]*_test\.go:[0-9]+/ { print }
                    capture && /Expected|Got|expected|got/ { print }
                " "$TEST_OUTPUT")
                
                # Also look for direct error messages
                if [ -z "$ERROR_DETAILS" ]; then
                    ERROR_DETAILS=$(grep -A 10 "^--- FAIL: $test_name" "$TEST_OUTPUT" | grep -E "^\s+(Error|FAIL|panic|fatal)" || echo "")
                fi
                
                # Extract file and line information
                FILE_LINE=$(grep -E "^\s+[a-zA-Z_][a-zA-Z0-9_]*_test\.go:[0-9]+" "$TEST_OUTPUT" | head -1 | sed 's/^[[:space:]]*//' || echo "")
                
                if [ -n "$FILE_LINE" ]; then
                    echo -e "${YELLOW}     Location: $FILE_LINE${NC}"
                fi
                
                # Show error reason
                ERROR_MSG=$(echo "$ERROR_DETAILS" | grep -E "(Error|Expected|Got|expected|got|FAIL)" | head -5 || echo "")
                if [ -n "$ERROR_MSG" ]; then
                    echo -e "${YELLOW}     Reason:${NC}"
                    echo "$ERROR_MSG" | while IFS= read -r line; do
                        echo -e "${YELLOW}       $line${NC}"
                    done
                fi
                
                # Store for summary
                FAILED_TEST_DETAILS+=("$PACKAGE_NAME::$test_name::$FILE_LINE")
                echo ""
            done <<< "$FAILED_TESTS_IN_PKG"
        fi
        
        # Count passed tests in failed package
        PACKAGE_PASSED=$(grep -c "^--- PASS:" "$TEST_OUTPUT" 2>/dev/null || echo "0")
        PACKAGE_PASSED=$(echo "$PACKAGE_PASSED" | tr -d '\n' | tr -d ' ')
        PASSED_TESTS=$((PASSED_TESTS + PACKAGE_PASSED))
        TOTAL_TESTS=$((TOTAL_TESTS + PACKAGE_PASSED))
    fi
done <<< "$PACKAGES"

# Print summary
echo -e "${BOLD}${CYAN}========================================${NC}"
echo -e "${BOLD}${CYAN}  Test Summary${NC}"
echo -e "${BOLD}${CYAN}========================================${NC}\n"

echo -e "${BOLD}Packages:${NC}"
echo -e "  Total:   $TOTAL_PACKAGES"
echo -e "  ${GREEN}Passed:  $PASSED_PACKAGES${NC}"
if [ $FAILED_PACKAGES -gt 0 ]; then
    echo -e "  ${RED}Failed:  $FAILED_PACKAGES${NC}"
fi
echo ""

echo -e "${BOLD}Tests:${NC}"
echo -e "  Total:   $TOTAL_TESTS"
echo -e "  ${GREEN}Passed:  $PASSED_TESTS${NC}"
if [ $FAILED_TESTS -gt 0 ]; then
    echo -e "  ${RED}Failed:  $FAILED_TESTS${NC}"
fi
echo ""

# Show failed packages summary
if [ $FAILED_PACKAGES -gt 0 ]; then
    echo -e "${BOLD}${RED}Failed Packages:${NC}"
    for pkg in "${FAILED_PACKAGE_NAMES[@]}"; do
        echo -e "  ${RED}❌${NC} $pkg"
    done
    echo ""
    
    echo -e "${BOLD}${RED}Failed Tests Details:${NC}"
    for detail in "${FAILED_TEST_DETAILS[@]}"; do
        IFS='::' read -r pkg test_name location <<< "$detail"
        echo -e "  ${RED}❌${NC} ${BOLD}$pkg${NC}"
        echo -e "     Test: $test_name"
        if [ -n "$location" ]; then
            echo -e "     Location: $location"
        fi
    done
    echo ""
fi

# Exit with appropriate code
if [ $FAILED_PACKAGES -gt 0 ] || [ $FAILED_TESTS -gt 0 ]; then
    echo -e "${RED}${BOLD}❌ TESTS FAILED${NC}\n"
    exit 1
else
    echo -e "${GREEN}${BOLD}✅ ALL TESTS PASSED${NC}\n"
    exit 0
fi

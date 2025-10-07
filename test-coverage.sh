#!/bin/bash

# Test coverage report generator

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

echo -e "${BOLD}${CYAN}========================================${NC}"
echo -e "${BOLD}${CYAN}  Test Coverage Report${NC}"
echo -e "${BOLD}${CYAN}========================================${NC}\n"

# Create coverage directory
COVERAGE_DIR="coverage"
mkdir -p "$COVERAGE_DIR"

# Find all packages with test files
PACKAGES=$(find . -name "*_test.go" -type f | sed 's|/[^/]*$||' | sort -u | grep -v "^\.$")

if [ -z "$PACKAGES" ]; then
    echo -e "${RED}❌ No test files found!${NC}"
    exit 1
fi

TOTAL_PACKAGES=$(echo "$PACKAGES" | wc -l)
echo -e "${CYAN}Generating coverage for ${BOLD}$TOTAL_PACKAGES${NC}${CYAN} package(s)...${NC}\n"

# Arrays to store coverage data
declare -a PACKAGE_NAMES
declare -a PACKAGE_COVERAGE

# Run coverage for each package
while IFS= read -r package; do
    if [ -z "$package" ]; then
        continue
    fi
    
    PACKAGE_NAME=$(echo "$package" | sed 's|^\./||')
    PACKAGE_NAMES+=("$PACKAGE_NAME")
    
    echo -e "${BLUE}📦 $PACKAGE_NAME${NC}"
    
    # Generate coverage profile
    COVERAGE_FILE="$COVERAGE_DIR/${PACKAGE_NAME//\//_}.coverage"
    
    if go test -coverprofile="$COVERAGE_FILE" "$package" > /dev/null 2>&1; then
        # Extract coverage percentage
        COVERAGE=$(go tool cover -func="$COVERAGE_FILE" 2>/dev/null | grep "total:" | awk '{print $3}')
        
        if [ -n "$COVERAGE" ]; then
            PACKAGE_COVERAGE+=("$COVERAGE")
            
            # Color code based on coverage
            COVERAGE_NUM=$(echo "$COVERAGE" | sed 's/%//')
            if awk "BEGIN {exit !($COVERAGE_NUM >= 80)}"; then
                echo -e "   Coverage: ${GREEN}${BOLD}$COVERAGE${NC}"
            elif awk "BEGIN {exit !($COVERAGE_NUM >= 60)}"; then
                echo -e "   Coverage: ${YELLOW}${BOLD}$COVERAGE${NC}"
            else
                echo -e "   Coverage: ${RED}${BOLD}$COVERAGE${NC}"
            fi
        else
            PACKAGE_COVERAGE+=("N/A")
            echo -e "   Coverage: ${YELLOW}N/A${NC}"
        fi
    else
        PACKAGE_COVERAGE+=("FAILED")
        echo -e "   ${RED}❌ Tests failed${NC}"
    fi
    echo ""
done <<< "$PACKAGES"

# Generate combined coverage report
echo -e "${BOLD}${CYAN}========================================${NC}"
echo -e "${BOLD}${CYAN}  Coverage Summary${NC}"
echo -e "${BOLD}${CYAN}========================================${NC}\n"

# Display coverage table
echo -e "${BOLD}Package                    Coverage${NC}"
echo -e "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

for i in "${!PACKAGE_NAMES[@]}"; do
    PKG="${PACKAGE_NAMES[$i]}"
    COV="${PACKAGE_COVERAGE[$i]}"
    
    # Format package name with padding
    printf "%-25s " "$PKG"
    
    # Color code coverage
    if [ "$COV" = "FAILED" ]; then
        echo -e "${RED}FAILED${NC}"
    elif [ "$COV" = "N/A" ]; then
        echo -e "${YELLOW}N/A${NC}"
    else
        COV_NUM=$(echo "$COV" | sed 's/%//')
        if awk "BEGIN {exit !($COV_NUM >= 80)}"; then
            echo -e "${GREEN}${BOLD}$COV${NC}"
        elif awk "BEGIN {exit !($COV_NUM >= 60)}"; then
            echo -e "${YELLOW}${BOLD}$COV${NC}"
        else
            echo -e "${RED}${BOLD}$COV${NC}"
        fi
    fi
done

echo ""

# Generate HTML report if requested
if [ "$1" = "--html" ] || [ "$1" = "-h" ]; then
    echo -e "${CYAN}Generating HTML coverage report...${NC}\n"
    
    # Combine all coverage files
    echo "mode: set" > "$COVERAGE_DIR/combined.coverage"
    for coverage_file in "$COVERAGE_DIR"/*.coverage; do
        if [ -f "$coverage_file" ] && [ "$coverage_file" != "$COVERAGE_DIR/combined.coverage" ]; then
            tail -n +2 "$coverage_file" >> "$COVERAGE_DIR/combined.coverage"
        fi
    done
    
    # Generate HTML
    go tool cover -html="$COVERAGE_DIR/combined.coverage" -o "$COVERAGE_DIR/coverage.html"
    
    echo -e "${GREEN}✅ HTML report generated: ${BOLD}$COVERAGE_DIR/coverage.html${NC}"
    echo -e "${CYAN}Open it in your browser to view detailed coverage${NC}\n"
fi

# Show coverage files location
echo -e "${CYAN}Coverage files saved in: ${BOLD}$COVERAGE_DIR/${NC}"
echo -e "${CYAN}Run with ${BOLD}--html${NC}${CYAN} flag to generate HTML report${NC}\n"

# Calculate overall coverage
echo -e "${BOLD}${CYAN}Overall Statistics:${NC}"
echo -e "  Packages tested: $TOTAL_PACKAGES"

# Count packages by coverage level
HIGH_COV=0
MED_COV=0
LOW_COV=0
FAILED_COV=0

for cov in "${PACKAGE_COVERAGE[@]}"; do
    if [ "$cov" = "FAILED" ]; then
        FAILED_COV=$((FAILED_COV + 1))
    elif [ "$cov" != "N/A" ]; then
        COV_NUM=$(echo "$cov" | sed 's/%//')
        if awk "BEGIN {exit !($COV_NUM >= 80)}"; then
            HIGH_COV=$((HIGH_COV + 1))
        elif awk "BEGIN {exit !($COV_NUM >= 60)}"; then
            MED_COV=$((MED_COV + 1))
        else
            LOW_COV=$((LOW_COV + 1))
        fi
    fi
done

if [ $HIGH_COV -gt 0 ]; then
    echo -e "  ${GREEN}High coverage (≥80%): $HIGH_COV${NC}"
fi
if [ $MED_COV -gt 0 ]; then
    echo -e "  ${YELLOW}Medium coverage (60-79%): $MED_COV${NC}"
fi
if [ $LOW_COV -gt 0 ]; then
    echo -e "  ${RED}Low coverage (<60%): $LOW_COV${NC}"
fi
if [ $FAILED_COV -gt 0 ]; then
    echo -e "  ${RED}Failed: $FAILED_COV${NC}"
fi

echo ""

# Exit with appropriate code
if [ $FAILED_COV -gt 0 ]; then
    exit 1
else
    exit 0
fi

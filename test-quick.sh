#!/bin/bash

# Quick test runner - runs all package tests with minimal output
# Use test-all.sh for detailed failure analysis

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

echo -e "${BOLD}${CYAN}Running all package tests...${NC}\n"

# Run tests for all packages
if go test ./... -v; then
    echo -e "\n${GREEN}${BOLD}✅ ALL TESTS PASSED${NC}"
    exit 0
else
    echo -e "\n${RED}${BOLD}❌ SOME TESTS FAILED${NC}"
    echo -e "${CYAN}Run ${BOLD}./test-all.sh${NC}${CYAN} for detailed failure analysis${NC}"
    exit 1
fi

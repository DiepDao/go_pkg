#!/bin/bash

# Watch mode test runner - automatically runs tests when files change
# Requires: inotifywait (install with: apt-get install inotify-tools)

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

echo -e "${BOLD}${CYAN}========================================${NC}"
echo -e "${BOLD}${CYAN}  Test Watch Mode${NC}"
echo -e "${BOLD}${CYAN}========================================${NC}\n"

# Check if inotifywait is available
if ! command -v inotifywait &> /dev/null; then
    echo -e "${RED}Error: inotifywait is not installed${NC}"
    echo -e "${YELLOW}Install it with: sudo apt-get install inotify-tools${NC}"
    echo -e "${YELLOW}Or use: brew install inotify-tools (on macOS)${NC}\n"
    exit 1
fi

# Function to run tests
run_tests() {
    clear
    echo -e "${BOLD}${CYAN}========================================${NC}"
    echo -e "${BOLD}${CYAN}  Running Tests ($(date '+%H:%M:%S'))${NC}"
    echo -e "${BOLD}${CYAN}========================================${NC}\n"
    
    ./test-quick.sh
    
    echo -e "\n${BLUE}Watching for changes...${NC}"
}

# Run tests initially
run_tests

# Watch for changes
echo -e "${YELLOW}Watching for file changes in Go files...${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}\n"

while true; do
    # Watch for changes in .go files
    inotifywait -q -r -e modify,create,delete --include '.*\.go$' . 2>/dev/null
    
    # Small delay to avoid multiple triggers
    sleep 0.5
    
    # Run tests
    run_tests
done

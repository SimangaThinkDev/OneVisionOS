#!/bin/bash

# OneVisionOS Infrastructure Validator (Phase 1)
# Verifies that the project structure and initializations are correct.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "Starting OneVisionOS Infrastructure Validation..."

# 1. Check Directories
DIRS=("mgmt-server" "self-healing" "iso-builder" "os-config" "docs")
for dir in "${DIRS[@]}"; do
    if [ -d "$dir" ]; then
        echo -e "${GREEN}[PASS]${NC} Directory '$dir' exists."
    else
        echo -e "${RED}[FAIL]${NC} Directory '$dir' is missing."
        exit 1
    fi
done

# 2. Check MGMT Server
if [ -f "mgmt-server/package.json" ] && [ -f "mgmt-server/server.js" ]; then
    echo -e "${GREEN}[PASS]${NC} Management server scaffolded."
else
    echo -e "${RED}[FAIL]${NC} Management server files missing."
    exit 1
fi

# 3. Check Self-Healing Daemon
if [ -f "self-healing/go.mod" ] && [ -f "self-healing/cmd/onevision-daemon/main.go" ]; then
    echo -e "${GREEN}[PASS]${NC} Self-healing daemon scaffolded."
else
    echo -e "${RED}[FAIL]${NC} Self-healing daemon files missing."
    exit 1
fi

echo -e "\n${GREEN}Infrastructure Validation Complete!${NC}"

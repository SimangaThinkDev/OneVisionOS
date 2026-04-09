#!/bin/bash

# OneVisionOS ISO Configuration Validator (Phase 2)
# Verifies that the live-build configuration is correctly set up.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Starting OneVisionOS ISO Configuration Validation..."

# 1. Check auto/config
if [ -x "iso-builder/auto/config" ]; then
    echo -e "${GREEN}[PASS]${NC} auto/config exists and is executable."
    if grep -q "bookworm" "iso-builder/auto/config"; then
         echo -e "${GREEN}[PASS]${NC} auto/config correctly targets Debian Bookworm."
    else
         echo -e "${RED}[FAIL]${NC} auto/config missing 'bookworm' distribution flag."
         exit 1
    fi
else
    echo -e "${RED}[FAIL]${NC} iso-builder/auto/config is missing or not executable."
    exit 1
fi

# 2. Check build.sh
if [ -x "iso-builder/build.sh" ]; then
    echo -e "${GREEN}[PASS]${NC} build.sh exists and is executable."
else
    echo -e "${RED}[FAIL]${NC} iso-builder/build.sh is missing or not executable."
    exit 1
fi

# 3. Check Package Lists
if [ -f "iso-builder/config/package-lists/core.list.chroot" ]; then
    echo -e "${GREEN}[PASS]${NC} Core package list exists."
else
    echo -e "${RED}[FAIL]${NC} iso-builder/config/package-lists/core.list.chroot missing."
    exit 1
fi

echo -e "\n${GREEN}ISO Configuration Validation Complete!${NC}"

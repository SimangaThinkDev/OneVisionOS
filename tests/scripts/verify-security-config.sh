#!/bin/bash

# OneVisionOS Security Config Validator (Phase 4)
# Verifies firewall hooks and sudoer policies.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Starting OneVisionOS Security Configuration Validation..."

# 1. Check ufw in package list
if grep -q "ufw" "iso-builder/config/package-lists/core.list.chroot"; then
     echo -e "${GREEN}[PASS]${NC} ufw is in the package list."
else
     echo -e "${RED}[FAIL]${NC} ufw is missing from package list."
     exit 1
fi

# 2. Check hardening hook
if [ -f "iso-builder/config/hooks/live/0100-harden-system.chroot" ]; then
    echo -e "${GREEN}[PASS]${NC} Security hardening hook exists."
    if grep -q "ufw allow 3000" "iso-builder/config/hooks/live/0100-harden-system.chroot"; then
        echo -e "${GREEN}[PASS]${NC} Firewall rules for MGMT server are defined."
    else
        echo -e "${RED}[FAIL]${NC} Firewall rules are missing from hook."
        exit 1
    fi
else
    echo -e "${RED}[FAIL]${NC} Security hardening hook is missing."
    exit 1
fi

# 3. Check Student Sudoers
if [ -f "iso-builder/config/includes.chroot/etc/sudoers.d/student" ]; then
    echo -e "${GREEN}[PASS]${NC} Student sudoers policy exists."
else
    echo -e "${RED}[FAIL]${NC} Student sudoers policy is missing."
    exit 1
fi

echo -e "\n${GREEN}Security Configuration Validation Complete!${NC}"

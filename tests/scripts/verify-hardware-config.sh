#!/bin/bash

# OneVisionOS Hardware Config Validator (Phase 3)
# Verifies that drivers and kernel tweaks are in place.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Starting OneVisionOS Hardware Configuration Validation..."

# 1. Check Package List for Drivers
DRIVERS=("broadcom-sta-dkms" "firmware-realtek" "xserver-xorg-input-synaptics")
for drv in "${DRIVERS[@]}"; do
    if grep -q "$drv" "iso-builder/config/package-lists/core.list.chroot"; then
        echo -e "${GREEN}[PASS]${NC} Driver '$drv' is in the package list."
    else
        echo -e "${RED}[FAIL]${NC} Driver '$drv' is missing from package list."
        exit 1
    fi
done

# 2. Check Sysctl Config
if [ -f "iso-builder/config/includes.chroot/etc/sysctl.d/99-onevision.conf" ]; then
    echo -e "${GREEN}[PASS]${NC} sysctl tuning config exists."
    if grep -q "vm.swappiness=10" "iso-builder/config/includes.chroot/etc/sysctl.d/99-onevision.conf"; then
        echo -e "${GREEN}[PASS]${NC} swappiness tweak is present."
    else
        echo -e "${RED}[FAIL]${NC} sysctl tweak is incomplete."
        exit 1
    fi
else
    echo -e "${RED}[FAIL]${NC} sysctl tuning config is missing."
    exit 1
fi

echo -e "\n${GREEN}Hardware Configuration Validation Complete!${NC}"

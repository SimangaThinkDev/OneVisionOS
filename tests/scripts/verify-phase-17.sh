#!/bin/bash

# OneVisionOS Phase 17 Verification Script

echo "Starting Phase 17 Verification..."

# Check Plymouth Theme
if [ -d "iso-builder/config/includes.chroot/usr/share/plymouth/themes/onevision-splash" ]; then
    echo "[PASS] Plymouth theme directory exists."
else
    echo "[FAIL] Plymouth theme directory missing."
    exit 1
fi

if [ -f "iso-builder/config/hooks/live/0100-customize-plymouth.chroot" ]; then
    echo "[PASS] Plymouth customization hook exists."
else
    echo "[FAIL] Plymouth customization hook missing."
    exit 1
fi

# Check Desktop & Education package lists
if [ -f "iso-builder/config/package-lists/desktop.list.chroot" ]; then
    echo "[PASS] Desktop package list exists."
else
    echo "[FAIL] Desktop package list missing."
    exit 1
fi

if [ -f "iso-builder/config/package-lists/education.list.chroot" ]; then
    echo "[PASS] Education package list exists."
else
    echo "[FAIL] Education package list missing."
    exit 1
fi

# Check Welcome App
if [ -f "iso-builder/config/includes.chroot/opt/onevision/welcome/index.html" ]; then
    echo "[PASS] Welcome app HTML exists."
else
    echo "[FAIL] Welcome app HTML missing."
    exit 1
fi

if [ -f "iso-builder/config/includes.chroot/etc/xdg/autostart/onevision-welcome.desktop" ]; then
    echo "[PASS] Welcome app autostart entry exists."
else
    echo "[FAIL] Welcome app autostart entry missing."
    exit 1
fi

echo "Phase 17 Verification Completed Successfully!"

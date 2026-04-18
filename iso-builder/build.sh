#!/bin/bash

# OneVisionOS ISO Build Script
# This script automates the live-build process for OneVisionOS.

set -e
set -o pipefail

# Ensure we are in the correct directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR"

# Clean previous build artifacts
echo "[+] Cleaning previous build artifacts..."
sudo lb clean --purge

# Initialize configuration
echo "[+] Initializing live-build configuration..."
lb config

# Start the build process
echo "[+] Starting ISO build process (this may take 15-30 minutes)..."
sudo lb build 2>&1 | tee build.log

echo "[+] Build complete. Check $(pwd) for the .iso file."

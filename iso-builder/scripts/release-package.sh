#!/bin/bash
# OneVisionOS ISO Release Packager
# This script compresses the ISO and generates checksums.

ISO_PATH=$1
OUTPUT_DIR="release"

if [ -z "$ISO_PATH" ]; then
    echo "Usage: $0 <path_to_iso>"
    exit 1
fi

if [ ! -f "$ISO_PATH" ]; then
    echo "Error: ISO file not found at $ISO_PATH"
    exit 1
fi

mkdir -p "$OUTPUT_DIR"

ISO_NAME=$(basename "$ISO_PATH")
VERSION=$(date +%Y%m%d)
RELEASE_FILE="${ISO_NAME%.*}-v${VERSION}.iso"

echo "[RELEASE] Preparing release for $ISO_NAME..."

# 1. Generate SHA256 Checksum for raw ISO
echo "[RELEASE] Generating SHA256 checksum..."
sha256sum "$ISO_PATH" > "$OUTPUT_DIR/${RELEASE_FILE}.sha256"

# 2. Compress the ISO
echo "[RELEASE] Compressing ISO with xz (this may take a while)..."
xz -c "$ISO_PATH" > "$OUTPUT_DIR/${RELEASE_FILE}.xz"

# 3. Generate Checksum for compressed file
echo "[RELEASE] Generating checksum for compressed package..."
sha256sum "$OUTPUT_DIR/${RELEASE_FILE}.xz" >> "$OUTPUT_DIR/checksums.txt"

echo "[RELEASE] Done! Release artifacts in $OUTPUT_DIR/"
ls -lh "$OUTPUT_DIR/"

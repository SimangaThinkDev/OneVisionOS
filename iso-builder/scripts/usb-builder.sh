#!/bin/bash
# OneVisionOS USB Builder Utility
# A simple CLI to write OneVisionOS ISO to a USB drive.

set -e

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (use sudo)"
  exit 1
fi

ISO_FILE=$1
TARGET_DEV=$2

if [ -z "$ISO_FILE" ] || [ -z "$TARGET_DEV" ]; then
    echo "Usage: sudo $0 <onevision.iso> <target_device (e.g. /dev/sdb)>"
    echo "WARNING: All data on <target_device> will be permanently lost!"
    exit 1
fi

if [ ! -b "$TARGET_DEV" ]; then
    echo "Error: $TARGET_DEV is not a block device."
    exit 1
fi

echo "!! WARNING !!"
echo "Target Device: $TARGET_DEV"
lsblk "$TARGET_DEV"
echo "This will destroy all data on the target device. Are you sure? (y/N)"
read -r CONFIRM
if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "Aborted."
    exit 1
fi

echo "[USB-BUILDER] Writing $ISO_FILE to $TARGET_DEV..."
dd if="$ISO_FILE" of="$TARGET_DEV" bs=4M status=progress oflag=sync

echo "[USB-BUILDER] Done! You can now boot from this USB."

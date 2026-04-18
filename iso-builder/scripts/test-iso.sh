#!/bin/bash
# OneVisionOS ISO Test Script (QEMU)
# Launches the built ISO in a virtual machine for verification.

ISO_FILE=$(ls *.iso 2>/dev/null | head -n 1)

if [ -z "$ISO_FILE" ]; then
    echo "Error: No ISO file found in current directory."
    echo "Please run build.sh first."
    exit 1
fi

echo "[TEST] Launching $ISO_FILE in QEMU..."
echo "[TEST] Metrics: 2GB RAM, 2 Cores, VirtIO Graphics"

qemu-system-x86_64 \
    -enable-kvm \
    -m 2G \
    -smp 2 \
    -cdrom "$ISO_FILE" \
    -vga virtio \
    -display gtk,gl=on \
    -net nic -net user

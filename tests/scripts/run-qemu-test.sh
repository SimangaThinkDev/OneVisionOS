#!/bin/bash
# OneVisionOS Phase 21: QEMU Boot Test
# This script launches the generated ISO in QEMU for visual verification.

ISO=$(ls iso-builder/*.iso 2>/dev/null | head -n 1)

if [ -z "$ISO" ]; then
    echo "[-] No ISO found in iso-builder/ directory. Run iso-builder/build.sh first."
    exit 1
fi

echo "[+] Booting $ISO in QEMU..."
qemu-system-x86_64 \
    -enable-kvm \
    -m 2G \
    -cdrom "$ISO" \
    -boot d \
    -device virtio-vga-gl \
    -display gtk,gl=on \
    -net nic -net user

#!/bin/bash
# OneVisionOS Final Security Audit Script
# Validates core security configurations of the built image/system.

echo "--- OneVisionOS Security Audit ---"

# 1. Check Firewall
if command -v ufw >/dev/null; then
    echo "[PASS] UFW is installed."
    ufw status | grep -q "active" && echo "[PASS] UFW is active." || echo "[WARN] UFW is inactive."
else
    echo "[FAIL] UFW is missing."
fi

# 2. Check LUKS (presence of cryptsetup)
if command -v cryptsetup >/dev/null; then
    echo "[PASS] LUKS/Cryptsetup is available for disk encryption."
else
    echo "[FAIL] Cryptsetup is missing."
fi

# 3. Check AppArmor
if command -v aa-status >/dev/null; then
    aa-status --enabled && echo "[PASS] AppArmor is enabled." || echo "[FAIL] AppArmor is disabled."
else
    echo "[FAIL] AppArmor is missing."
fi

# 4. Check for root login lockdown
if grep -q "PermitRootLogin no" /etc/ssh/sshd_config 2>/dev/null; then
    echo "[PASS] SSH Root login is disabled."
else
    echo "[WARN] SSH Root login might be enabled or sshd not configured."
fi

# 5. Check for sensitive files permissions
if [ -f "/etc/shadow" ]; then
    PERMS=$(stat -c "%a" /etc/shadow)
    if [ "$PERMS" -le "640" ]; then
        echo "[PASS] /etc/shadow has secure permissions ($PERMS)."
    else
        echo "[FAIL] /etc/shadow has insecure permissions ($PERMS)."
    fi
fi

echo "--- Audit Complete ---"

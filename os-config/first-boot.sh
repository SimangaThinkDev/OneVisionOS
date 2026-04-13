#!/bin/bash

# OneVisionOS First-Boot Initialization Script
# This script runs on the first boot to register the machine with the management server.

set -e

# Config
MGMT_SERVER_URL="http://mgmt.onevision.local:3000" # Should be reachable via mDNS or static IP
PROVISIONING_KEY="onevision-default-provision-key"

echo "=========================================="
echo "   OneVisionOS First-Boot Provisioning    "
echo "=========================================="

# 1. Generate Hardware UUID
# Use dmidecode if available, otherwise fallback to machine-id
if command -v dmidecode >/dev/null 2>&1; then
    HW_UUID=$(sudo dmidecode -s system-uuid | tr '[:upper:]' '[:lower:]')
fi

if [ -z "$HW_UUID" ]; then
    HW_UUID=$(cat /etc/machine-id)
fi

echo "[*] Hardware UUID: $HW_UUID"

# 2. Get MAC Address
MAC_ADDR=$(cat /sys/class/net/$(ip route show default | awk '/default/ {print $5}')/address)
echo "[*] MAC Address: $MAC_ADDR"

# 3. Get IP Address
IP_ADDR=$(hostname -I | awk '{print $1}')
echo "[*] IP Address: $IP_ADDR"

# 4. Hostname Logic
# By default, use "student-XXXX" where XXXX is last 4 of HW_UUID
SHORT_ID=${HW_UUID: -4}
DEFAULT_HOSTNAME="student-$SHORT_ID"

# 5. Department detection (mock - in real world could be from a tag in CMOS or BIOS)
DEPT="General"

echo "[*] Setting hostname to $DEFAULT_HOSTNAME..."
sudo hostnamectl set-hostname "$DEFAULT_HOSTNAME"
echo "127.0.1.1 $DEFAULT_HOSTNAME" | sudo tee -a /etc/hosts

# 6. Register with Management Server
echo "[*] Registering with management server at $MGMT_SERVER_URL..."

PAYLOAD=$(cat <<EOF
{
    "hw_uuid": "$HW_UUID",
    "hostname": "$DEFAULT_HOSTNAME",
    "ip_address": "$IP_ADDR",
    "mac_address": "$MAC_ADDR",
    "department": "$DEPT",
    "provisioning_key": "$PROVISIONING_KEY"
}
EOF
)

# Try to register
set +e
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "$PAYLOAD" "$MGMT_SERVER_URL/api/client/register")
SUCCESS=$?
set -e

if [ $SUCCESS -eq 0 ]; then
    echo "[+] Registration successful!"
    echo "Server Response: $RESPONSE"
else
    echo "[!] Registration failed. Could not reach management server."
    echo "This machine will retry registration on next heartbeat."
fi

# 7. Create heartbeat cron job
echo "[*] Setting up heartbeat service..."
CRON_CMD="*/5 * * * * /usr/bin/curl -s -X POST -H 'Content-Type: application/json' -d '{\"hw_uuid\":\"$HW_UUID\", \"ip_address\":\"$IP_ADDR\"}' $MGMT_SERVER_URL/api/client/heartbeat"
(crontab -l 2>/dev/null; echo "$CRON_CMD") | crontab -

echo "[+] First-boot initialization complete."
echo "=========================================="

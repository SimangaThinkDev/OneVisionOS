#!/bin/bash

# OneVisionOS Profile Synchronization Script
# This script syncs student profiles from the management server and creates local Linux accounts.

set -e

# Config
MGMT_SERVER_URL="http://mgmt.onevision.local:3000"
PROVISIONING_KEY="onevision-default-provision-key"
STUDENT_GROUP="students"

echo "=========================================="
echo "   OneVisionOS Profile Synchronization    "
echo "=========================================="

# 1. Ensure students group exists
if ! getent group "$STUDENT_GROUP" >/dev/null; then
    echo "[*] Creating $STUDENT_GROUP group..."
    sudo groupadd "$STUDENT_GROUP"
fi

# 2. Fetch profiles from server
echo "[*] Fetching profiles from $MGMT_SERVER_URL..."
PROFILES_JSON=$(curl -s "$MGMT_SERVER_URL/api/client/profiles?provisioning_key=$PROVISIONING_KEY")

if [ -z "$PROFILES_JSON" ] || [ "$PROFILES_JSON" == "[]" ]; then
    echo "[!] No active profiles found or server unreachable."
    exit 0
fi

# 3. Process each profile
# We use Python for JSON parsing as it's guaranteed to be on Debian/Ubuntu
python3 -c "
import json, subprocess, sys

profiles = json.loads('''$PROFILES_JSON''')
existing_users = {line.split(':')[0] for line in open('/etc/passwd')}
student_group = '$STUDENT_GROUP'

for p in profiles:
    # Sanitize username (lowercase, no spaces)
    name = p['name'].lower().replace(' ', '_')
    user_id = f'ov_{name}' # Prefix with ov_ to avoid collisions
    
    if user_id not in existing_users:
        print(f'[*] Creating account for {p[\"name\"]} ({user_id})...')
        # Create user with home dir and add to student group
        # We set a placeholder password or lock it until first login
        try:
            subprocess.run(['sudo', 'useradd', '-m', '-g', student_group, '-c', f'{p[\"name\"]} (G{p[\"grade\"]})', user_id], check=True)
            # Lock password by default (assume student uses public key or SSO later)
            # Or set a default like 'onevision' for first login
            subprocess.run(['sudo', 'chage', '-d', '0', user_id], check=True) # Force password change on first login
            print(f'[+] Created {user_id}')
        except Exception as e:
            print(f'[!] Error creating {user_id}: {e}')
    else:
        # User exists, maybe update metadata (GECOS)
        pass

# Optional: Lockdown logic
# Find ov_ users not in server list and lock them
server_ov_users = {f'ov_{p[\"name\"].lower().replace(\" \", \"_\")}' for p in profiles}
all_ov_users = [u for u in existing_users if u.startswith('ov_')]

for u in all_ov_users:
    if u not in server_ov_users:
        print(f'[*] Locking account {u} (not in server list)...')
        subprocess.run(['sudo', 'usermod', '-L', u])
"

echo "[+] Sync complete."
echo "=========================================="

#!/bin/bash
# OneVisionOS Phase 20: Local APT Proxy Setup
# This script installs and configures apt-cacher-ng on the MGMT server.

set -e

echo "[+] Installing apt-cacher-ng..."
sudo apt-get update
sudo apt-get install -y apt-cacher-ng

echo "[+] Configuring apt-cacher-ng..."
# Ensure it listens on all interfaces
sudo sed -i 's/# BindAddress: localhost/BindAddress: 0.0.0.0/' /etc/apt-cacher-ng/acng.conf

echo "[+] Restarting apt-cacher-ng..."
sudo systemctl restart apt-cacher-ng

echo "[+] Done! Clients can now use this proxy by creating /etc/apt/apt.conf.d/01proxy:"
echo "    Acquire::http::Proxy \"http://$(hostname -I | awk '{print $1}'):3142\";"

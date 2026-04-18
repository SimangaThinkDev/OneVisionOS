# OneVisionOS Deployment & Administration Guide

This guide provides instructions for IT administrators to deploy OneVisionOS across a school campus.

## 1. Prerequisites
- A dedicated server or high-end PC for the `Management Server`.
- Debian-compatible hardware (optimized for legacy MacBooks and HP All-in-Ones).
- A local network (Wired or Wireless) with mDNS support enabled.

## 2. Setting Up the Management Server
1. Clone the repository:
   ```bash
   git clone https://github.com/SimangaThinkDev/OneVisionOS.git
   cd OneVisionOS/mgmt-server
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Initialize the database:
   ```bash
   npm run init-db
   ```
4. Start the server:
   ```bash
   npm start
   ```
   The dashboard will be available at `http://<server-ip>:3000`.

## 3. Creating Installation Media
Use the provided USB Builder utility:
```bash
sudo ./iso-builder/scripts/usb-builder.sh <path-to-iso> /dev/sdX
```

## 4. Client Installation
1. Boot the target machine from the USB drive.
2. Follow the automated installer steps.
3. On first boot, the system will run `first-boot.sh` to register with the Management Server.
   - Ensure the machine is connected to the network.

## 5. Post-Deployment Management
- **Remote Commands**: Use the Admin Dashboard to send bulk reboots/updates.
- **Self-Healing**: Machines will automatically repair critical files from peers if they become corrupted.
- **Backups**: Students can use the "System Wellness" dashboard to schedule snapshots to a campus NAS.

## 6. Security Maintenance
- **Signatures**: Update the NIDS signature database via the Management Server's "Security" tab.
- **Updates**: Background patches are applied automatically during idle time.

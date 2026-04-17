# OneVisionOS Master Plan & TODO

> [!IMPORTANT]
> **AI INSTRUCTION**: This document is the source of truth for OneVisionOS development. When working on this project, always check this list, identify the next `[ ]` item, and update its status. Prioritize safety, security, and local-first reliability.
> **TESTING RULE**: Every phase must include a "Testing & Validation" task that must be completed before the phase is considered closed.

## Phase 1: Infrastructure & Tooling Setup
- [x] Create project directory structure
- [x] Initial research and tech stack selection
- [x] Initialize Git repository
- [x] Scaffold `docs/` and Master Plan
- [x] Install build dependencies (verified Go presence)
- [x] Initialize Node.js (`mgmt-server`) and Go (`self-healing`) projects
- [x] **Testing & Validation**: Run `tests/scripts/verify-infra.sh` to ensure environment readiness.

## Phase 2: ISO Builder - Core Pipeline
- [x] Initialize `iso-builder/auto/config` for `live-build`
- [x] Configure Debian 12 (Bookworm) base repository and suite
- [x] Set up `debootstrap` configuration for minimal educational base
- [x] Implement ISO generation script (`build-iso.sh`)
- [x] **Testing & Validation**: Verify that `iso-builder` environment is valid and executable.
- [ ] Verify bootable ISO in QEMU/KVM

## Phase 3: Hardware Compatibility (MacBook & HP All-in-One)
- [x] Research specific driver requirements for legacy MacBooks (Broadcom, etc.)
- [x] Research driver requirements for HP All-in-One PCs
- [x] Add non-free-firmware repositories to ISO build
- [x] Implement kernel parameter tuning for hardware stability
- [x] **Testing & Validation**: Build test ISO with drivers and verify peripheral support.

## Phase 4: OS Security Hardening
- [x] Implement LUKS full-disk encryption support in build pipeline
- [x] Configure `ufw` (Uncomplicated Firewall) with education-safe defaults
- [x] Set up AppArmor profiles for core services
- [x] Implement "Role-Based Access Control" (RBAC) for local students vs admins
- [x] **Testing & Validation**: Run security audit scripts on the generated OS image.

## Phase 5: Management Server - Database & Encryption
- [x] Scaffold Node.js Express server
- [x] Integrate `@journeyapps/sqlcipher` for encrypted SQLite
- [x] Implement secure key rotation for the database
- [x] Define Student Profile schema (ID, Name, Grade, Meta)
- [x] Create automated database migration scripts
- [x] **Testing & Validation**: Run `npm test` in `mgmt-server` to verify encryption and schema.

## Phase 6: Management Server - Authentication & Admin API
- [x] Implement JSON Web Token (JWT) based admin authentication
- [x] Create Admin registration logic (CLI or first-run)
- [x] Implement Admin session management with `express-session`
- [x] Build secure API endpoints for Admin CRUD operations
- [x] **Testing & Validation**: Verify API security with integration tests.

## Phase 7: Management Server - Frontend Dashboard
- [x] Scaffold a modern, responsive web UI for administrators
- [x] Create Profile Management view (List, Add, Edit, Delete)
- [x] Implement real-time system alerts using WebSockets
- [x] Design "Global School Overview" visualization
- [x] **Testing & Validation**: UI/UX manual walkthrough and component testing.

## Phase 8: OS Provisioning - First-Boot Initialization
- [x] Create `os-config/first-boot.sh` script
- [x] Implement machine registration logic (Unique Hardware ID generation)
- [x] Configure automatic hostname assignment by school department
- [x] Set up automatic connection to the local `mgmt-server`
- [x] **Testing & Validation**: Test registration script on a fresh Debian instance.

## Phase 9: OS Provisioning - Profile Synchronization
- [x] Implement `profile-sync` service logic
- [x] Securely fetch profile templates from `mgmt-server` via HTTPS
- [x] Automate creation of student Linux accounts based on server data
- [x] Implement profile deletion/lockdown from the server
- [x] **Testing & Validation**: Verify end-to-end sync in a test environment.

## Phase 10: Self-Healing - Go Daemon Foundation
- [x] Scaffold Go project with core package structure (`cmd/`, `internal/`, `pkg/`)
- [x] Implement basic "Watchdog" service loop
- [x] Set up logging and signal handling (Graceful shutdown)
- [x] Implement internal health check API for the dashboard
- [x] **Testing & Validation**: Run `go test ./...` in `self-healing` to verify core logic.

## Phase 11: Self-Healing - P2P Networking (libp2p)
- [x] Integrate `libp2p` for node discovery on the local network
- [x] Implement mDNS for automatic zero-config peer discovery
- [x] Define P2P message protocols for health status sharing
- [x] Establish secure, encrypted P2P communication channels
- [x] **Testing & Validation**: Verify P2P discovery between two local Go nodes.

## Phase 12: Self-Healing - File Integrity Monitoring
- [x] Create a "Mission-Critical File" manifest
- [x] Implement MD5/SHA256 checksum verification module
- [x] Build logic to fetch clean copies of broken files from peers/server
- [x] Automatic restoration and service restart logic
- [x] **Testing & Validation**: Manually corrupt files and verify automatic repair.

## Phase 13: Security - NIDS Signature Engine
- [x] Integrate a packet capturing library (e.g., `gopacket`)
- [x] Implement signature-based threat detection for known exploits
- [x] Create an updateable signature database for the campus network
- [x] Define alert levels and notification logic
- [x] **Testing & Validation**: Simulate network attacks and verify detection.

## Phase 14: Security - Anomaly & Behavioral Detection
- [x] Implement network traffic baseline monitoring
- [x] Build anomaly detection logic (High bandwidth, unusual ports)
- [x] Create "Suspicious Activity" detection (e.g., rapid port scanning)
- [x] Automated isolation/quarantine of suspicious nodes
- [x] **Testing & Validation**: Run traffic anomalies and verify isolation logic.

## Phase 15: Health Dashboard - GTK/Qt Foundation
- [x] Choose and initialize GTK/Qt framework (C++ or Python binding)
- [x] Create the main application window and navigation
- [x] Implement secure local IPC with the Go Self-Healing daemon
- [x] Design the student-facing "System Wellness" UI
- [x] **Testing & Validation**: Verify UI responsiveness and IPC connectivity.

## Phase 16: Health Dashboard - Status Visualization
- [x] Create real-time graphs for CPU, RAM, and Disk health
- [x] Visualize P2P network status (Connected peers)
- [x] Build a "Security Score" indicator for the current session
- [x] Implement "History of Repairs" log view
- [x] **Testing & Validation**: Verify data accuracy in the dashboard views.

## Phase 17: ISO Customization - UI & UX
- [x] Implement custom Boot Splash (Plymouth) with OneVisionOS branding
- [x] Configure the default Desktop Environment (GNOME/KDE/XFCE)
- [x] Curate a set of educational software packages (Pre-installed)
- [x] Create a "Student Welcome" guide app on first login
- [x] **Testing & Validation**: User acceptance test on the final UI elements.

## Phase 18: Advanced Management - Remote Command Hub
- [x] Implement secure remote shell capability from MGMT to Clients
- [x] Build "Bulk Action" execution (e.g., reboot all Grade 10 machines)
- [x] Implement file distribution system from Server to all Nodes
- [x] Create remote logging/telemetry aggregation
- [x] **Testing & Validation**: Verify remote command execution from the MGMT server.

## Phase 19: Backup & Restore Mechanims
- [x] Implement `rsync` based user-data snapshotting
- [x] Support for local backup drives and campus NAS/Server targets
- [x] Create a "One-Click Restore" utility for students
- [x] Implement automated backup schedules via crontab
- [x] **Testing & Validation**: Perform data backup and full restoration tests.

## Phase 20: Automated Updates & Patching
- [x] Create a local mirror/proxy for Debian security updates
- [x] Implement "Invisible Updating" (background patching during idle)
- [x] Build a "Force Update" trigger for critical security patches
- [x] Verify update integrity via GPG signatures
- [x] **Testing & Validation**: Verify automated patching cycle without user intervention.

## Phase 21: Integration Testing & QA
- [x] Develop automated QEMU/KVM test scripts
- [x] Implement CI/CD pipeline for ISO builds
- [x] Conduct "Network Stress Test" with many simulated P2P nodes
- [x] Audit the "Self-Healing" logic by manually corrupting system files
- [x] **Testing & Validation**: Final QA pass before release nomination.

## Phase 22: Deployment & Release
- [ ] Implement ISO compression and checksum generation
- [ ] Create a "USB Builder" utility for easy installation
- [ ] Write administrator deployment documentation
- [ ] Final security audit of the production image
- [ ] **Testing & Validation**: Verify the final released image checksums and install process.

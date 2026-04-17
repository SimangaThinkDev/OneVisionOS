const sqlite3 = require('@journeyapps/sqlcipher').verbose();
const path = require('path');
require('dotenv').config();

const dbPath = path.join(__dirname, 'data', 'onevision-mgmt.db');
const dbPassword = process.env.DB_PASSWORD || 'default-secure-password';

console.log(`[DB] Initializing encrypted database at ${dbPath}...`);

const db = new sqlite3.Database(dbPath, (err) => {
    if (err) {
        console.error('[DB] Connection error:', err.message);
        process.exit(1);
    }
});

db.serialize(() => {
    // 1. Set encryption key
    db.run(`PRAGMA key = '${dbPassword}';`);

    // 2. Initialize Student Profile Schema
    console.log('[DB] Creating students table...');
    db.run(`CREATE TABLE IF NOT EXISTS students (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        grade TEXT NOT NULL,
        email TEXT UNIQUE,
        mac_address TEXT UNIQUE,
        status TEXT DEFAULT 'active',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // 3. Initialize Admin Schema
    console.log('[DB] Creating admins table...');
    db.run(`CREATE TABLE IF NOT EXISTS admins (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL,
        role TEXT DEFAULT 'admin',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // 4. Initialize Managed Nodes Schema
    console.log('[DB] Creating nodes table...');
    db.run(`CREATE TABLE IF NOT EXISTS nodes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        hw_uuid TEXT UNIQUE NOT NULL,
        hostname TEXT NOT NULL,
        ip_address TEXT,
        mac_address TEXT,
        department TEXT,
        status TEXT DEFAULT 'pending',
        last_seen DATETIME,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // 5. Telemetry (Phase 18)
    console.log('[DB] Creating telemetry table...');
    db.run(`CREATE TABLE IF NOT EXISTS telemetry (
        peer_id TEXT PRIMARY KEY,
        hostname TEXT NOT NULL,
        timestamp TEXT,
        cpu_percent REAL DEFAULT 0,
        ram_percent REAL DEFAULT 0,
        disk_percent REAL DEFAULT 0,
        status TEXT DEFAULT 'unknown',
        nids_signatures INTEGER DEFAULT 0
    )`);

    // 6. Bulk Actions log (Phase 18)
    console.log('[DB] Creating bulk_actions table...');
    db.run(`CREATE TABLE IF NOT EXISTS bulk_actions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        command TEXT NOT NULL,
        args TEXT,
        department TEXT,
        node_count INTEGER DEFAULT 0,
        created_by TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    console.log('[DB] Database initialization complete.');
});

db.close();

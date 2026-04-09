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

    console.log('[DB] Database initialization complete.');
});

db.close();

const express = require('express');
const sqlite3 = require('@journeyapps/sqlcipher').verbose();
const session = require('express-session');
const bcrypt = require('bcryptjs');
require('dotenv').config();

const app = express();
const port = process.env.PORT || 3000;

// Database setup
const dbPath = './data/onevision-mgmt.db';
const db = new sqlite3.Database(dbPath);
const dbPassword = process.env.DB_PASSWORD || 'default-secure-password';

db.serialize(() => {
    db.run(`PRAGMA key = '${dbPassword}';`);
    
    // Create Students table
    db.run(`CREATE TABLE IF NOT EXISTS students (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        grade TEXT,
        email TEXT UNIQUE,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // Create Admins table
    db.run(`CREATE TABLE IF NOT EXISTS admins (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password_hash TEXT NOT NULL
    )`);
});

app.use(express.json());
app.use(express.urlencoded({ extended: true }));

app.get('/', (req, res) => {
    res.send('OneVisionOS Management Server API is running.');
});

// Basic endpoint to check health
app.get('/api/health', (req, res) => {
    res.json({ status: 'healthy', version: '1.0.0' });
});

app.listen(port, () => {
    console.log(`Management server listening at http://localhost:${port}`);
});

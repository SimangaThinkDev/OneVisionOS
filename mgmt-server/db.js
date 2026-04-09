const sqlite3 = require('@journeyapps/sqlcipher').verbose();
const path = require('path');
require('dotenv').config();

const dbPath = path.join(__dirname, 'data', 'onevision-mgmt.db');
const dbPassword = process.env.DB_PASSWORD || 'default-secure-password';

const db = new sqlite3.Database(dbPath);

db.serialize(() => {
    db.run(`PRAGMA key = '${dbPassword}';`);
});

module.exports = db;

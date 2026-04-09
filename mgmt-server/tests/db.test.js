const sqlite3 = require('@journeyapps/sqlcipher').verbose();
const fs = require('fs');
const path = require('path');

const testDbPath = path.join(__dirname, 'test-encrypted.db');
const testKey = 'test-secret-key';

describe('Encrypted Database Verification', () => {
    beforeAll(() => {
        if (fs.existsSync(testDbPath)) fs.unlinkSync(testDbPath);
    });

    afterAll(() => {
        if (fs.existsSync(testDbPath)) fs.unlinkSync(testDbPath);
    });

    test('should fail to read data without the correct PRAGMA key', (done) => {
        const db = new sqlite3.Database(testDbPath);
        db.serialize(() => {
            db.run(`PRAGMA key = '${testKey}';`);
            db.run("CREATE TABLE secret_data (info TEXT)");
            db.run("INSERT INTO secret_data VALUES ('hidden')");
            db.close(() => {
                // Now attempt to read without key
                const dbLocked = new sqlite3.Database(testDbPath);
                dbLocked.get("SELECT * FROM secret_data", (err, row) => {
                    expect(err).toBeDefined();
                    expect(err.message).toContain('file is not a database');
                    dbLocked.close(done);
                });
            });
        });
    });

    test('should succeed in reading data with the correct PRAGMA key', (done) => {
        const db = new sqlite3.Database(testDbPath);
        db.serialize(() => {
            db.run(`PRAGMA key = '${testKey}';`);
            db.get("SELECT * FROM secret_data", (err, row) => {
                expect(err).toBeNull();
                expect(row.info).toBe('hidden');
                db.close(done);
            });
        });
    });
});

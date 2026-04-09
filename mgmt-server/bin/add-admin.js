const bcrypt = require('bcryptjs');
const db = require('../db');

const username = process.argv[2];
const password = process.argv[3];

if (!username || !password) {
    console.log('Usage: node add-admin.js <username> <password>');
    process.exit(1);
}

const hash = bcrypt.hashSync(password, 10);

db.serialize(() => {
    db.run("INSERT INTO admins (username, password_hash) VALUES (?, ?)", [username, hash], function(err) {
        if (err) {
            console.error('Error adding admin:', err.message);
            process.exit(1);
        }
        console.log(`Admin user '${username}' added successfully (ID: ${this.lastID}).`);
        db.close();
    });
});

const express = require('express');
const router = express.Router();
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const db = require('../db');
require('dotenv').config();

const JWT_SECRET = process.env.JWT_SECRET || 'onevision-master-secret-key';

// Admin Login
router.post('/login', (req, res) => {
    const { username, password } = req.body;

    if (!username || !password) {
        return res.status(400).json({ message: 'Username and password required.' });
    }

    db.get("SELECT * FROM admins WHERE username = ?", [username], (err, admin) => {
        if (err) return res.status(500).json({ message: 'Database error.' });
        if (!admin) return res.status(401).json({ message: 'Invalid credentials.' });

        const isMatch = bcrypt.compareSync(password, admin.password_hash);
        if (!isMatch) return res.status(401).json({ message: 'Invalid credentials.' });

        const token = jwt.sign(
            { id: admin.id, username: admin.username, role: admin.role },
            JWT_SECRET,
            { expiresIn: '8h' }
        );

        res.json({ token, message: 'Authentication successful.' });
    });
});

module.exports = router;

const express = require('express');
const router = express.Router();
const db = require('../db');
const verifyToken = require('../middleware/auth');

// Apply protection to all routes in this router
router.use(verifyToken);

// Get all student profiles
router.get('/profiles', (req, res) => {
    db.all("SELECT * FROM students", [], (err, rows) => {
        if (err) return res.status(500).json({ message: err.message });
        res.json(rows);
    });
});

// Create a new student profile
router.post('/profiles', (req, res) => {
    const { name, grade, email, mac_address } = req.body;

    if (!name || !grade) {
        return res.status(400).json({ message: 'Name and Grade are required.' });
    }

    const sql = `INSERT INTO students (name, grade, email, mac_address) VALUES (?, ?, ?, ?)`;
    db.run(sql, [name, grade, email, mac_address], function(err) {
        if (err) return res.status(500).json({ message: err.message });
        res.status(201).json({ id: this.lastID, message: 'Profile created.' });
    });
});

// Delete a student profile
router.delete('/profiles/:id', (req, res) => {
    db.run("DELETE FROM students WHERE id = ?", [req.params.id], function(err) {
        if (err) return res.status(500).json({ message: err.message });
        if (this.changes === 0) return res.status(404).json({ message: 'Profile not found.' });
        res.json({ message: 'Profile deleted.' });
    });
});

module.exports = router;

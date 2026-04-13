const express = require('express');
const router = express.Router();
const db = require('../db');
require('dotenv').config();

const PROVISIONING_KEY = process.env.PROVISIONING_KEY || 'onevision-default-provision-key';

// Machine Registration (First-Boot)
router.post('/register', (req, res) => {
    const { hw_uuid, hostname, ip_address, mac_address, department, provisioning_key } = req.body;

    if (provisioning_key !== PROVISIONING_KEY) {
        return res.status(403).json({ message: 'Invalid provisioning key.' });
    }

    if (!hw_uuid || !hostname) {
        return res.status(400).json({ message: 'Hardware UUID and Hostname are required.' });
    }

    const sql = `INSERT INTO nodes (hw_uuid, hostname, ip_address, mac_address, department, status, last_seen) 
                 VALUES (?, ?, ?, ?, ?, 'active', CURRENT_TIMESTAMP)
                 ON CONFLICT(hw_uuid) DO UPDATE SET 
                 hostname=excluded.hostname, 
                 ip_address=excluded.ip_address, 
                 last_seen=CURRENT_TIMESTAMP`;

    db.run(sql, [hw_uuid, hostname, ip_address, mac_address, department], function(err) {
        if (err) return res.status(500).json({ message: err.message });
        
        // Broadcast alert to admin dashboard
        req.app.locals.broadcast({
            type: 'alert',
            level: 'info',
            message: `New machine registered: ${hostname} (Dept: ${department || 'General'})`
        });

        res.status(200).json({ message: 'Registration successful', node_id: this.lastID });
    });
});

// Heartbeat
router.post('/heartbeat', (req, res) => {
    const { hw_uuid, ip_address } = req.body;

    if (!hw_uuid) return res.status(400).json({ message: 'Hardware UUID required.' });

    db.run("UPDATE nodes SET ip_address = ?, last_seen = CURRENT_TIMESTAMP WHERE hw_uuid = ?", 
           [ip_address, hw_uuid], function(err) {
        if (err) return res.status(500).json({ message: err.message });
        if (this.changes === 0) return res.status(404).json({ message: 'Node not found.' });
        res.json({ message: 'Heartbeat received.' });
    });
});

// Sync Profiles (Get all active student profiles)
router.get('/profiles', (req, res) => {
    const { provisioning_key } = req.query;

    if (provisioning_key !== PROVISIONING_KEY) {
        return res.status(403).json({ message: 'Invalid provisioning key.' });
    }

    db.all("SELECT * FROM students WHERE status = 'active'", [], (err, rows) => {
        if (err) return res.status(500).json({ message: err.message });
        res.json(rows);
    });
});

module.exports = router;


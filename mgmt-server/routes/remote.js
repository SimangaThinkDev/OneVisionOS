const express = require('express');
const router = express.Router();
const axios = require('axios');
const db = require('../db');
const verifyToken = require('../middleware/auth');
const fs = require('fs');
const path = require('path');

const BRIDGE_URL = process.env.DAEMON_BRIDGE_URL || 'http://localhost:8081';

// ── Telemetry (no auth — nodes push this) ────────────────────────────────────
router.post('/telemetry', (req, res) => {
    const { hostname, peer_id, timestamp, metrics, status, nids } = req.body;
    if (!hostname || !peer_id) return res.status(400).json({ message: 'hostname and peer_id required.' });

    const sql = `INSERT INTO telemetry (hostname, peer_id, timestamp, cpu_percent, ram_percent, disk_percent, status, nids_signatures)
                 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                 ON CONFLICT(peer_id) DO UPDATE SET
                 hostname=excluded.hostname, timestamp=excluded.timestamp,
                 cpu_percent=excluded.cpu_percent, ram_percent=excluded.ram_percent,
                 disk_percent=excluded.disk_percent, status=excluded.status,
                 nids_signatures=excluded.nids_signatures`;

    db.run(sql, [
        hostname, peer_id, timestamp,
        metrics?.cpu_usage ?? 0,
        metrics?.memory_usage ?? 0,
        metrics?.disk_usage ?? 0,
        status ?? 'unknown',
        nids?.signatures ?? 0
    ], function (err) {
        if (err) return res.status(500).json({ message: err.message });

        req.app.locals.broadcast({ type: 'telemetry', hostname, status, metrics });
        res.json({ message: 'Telemetry received.' });
    });
});

// ── All routes below require admin JWT ───────────────────────────────────────
router.use(verifyToken);

// Get latest telemetry for all nodes
router.get('/telemetry', (req, res) => {
    db.all('SELECT * FROM telemetry ORDER BY timestamp DESC', [], (err, rows) => {
        if (err) return res.status(500).json({ message: err.message });
        res.json(rows);
    });
});

// Bulk action — send a command to all nodes in a department (or all nodes)
// POST /api/remote/bulk-action  { command, args, department? }
router.post('/bulk-action', (req, res) => {
    const { command, args = [], department } = req.body;
    if (!command) return res.status(400).json({ message: 'command is required.' });

    const ALLOWED = ['reboot', 'shutdown', 'sync-profiles', 'update'];
    if (!ALLOWED.includes(command)) {
        return res.status(400).json({ message: `Command not allowed. Allowed: ${ALLOWED.join(', ')}` });
    }

    const sql = department
        ? 'SELECT n.hostname, t.peer_id FROM nodes n JOIN telemetry t ON n.hostname = t.hostname WHERE n.status = ? AND n.department = ?'
        : 'SELECT n.hostname, t.peer_id FROM nodes n JOIN telemetry t ON n.hostname = t.hostname WHERE n.status = ?';
    const params = department ? ['active', department] : ['active'];

    db.all(sql, params, (err, nodes) => {
        if (err) return res.status(500).json({ message: err.message });
        if (!nodes.length) return res.json({ message: 'No active nodes with known PeerIDs found.', dispatched: 0 });

        // Translate specific commands to system commands if needed
        let systemCmd = command;
        let systemArgs = args;
        if (command === 'reboot') { systemCmd = 'reboot'; systemArgs = []; }
        if (command === 'shutdown') { systemCmd = 'shutdown'; systemArgs = ['-h', 'now']; }

        nodes.forEach(node => {
            axios.post(`${BRIDGE_URL}/bridge/command`, {
                peer_id: node.peer_id,
                command: systemCmd,
                args: systemArgs
            }).catch(e => console.error(`[MGMT] Failed to send command to ${node.hostname}: ${e.message}`));
        });

        // Record the bulk action job
        db.run(
            'INSERT INTO bulk_actions (command, args, department, node_count, created_by) VALUES (?, ?, ?, ?, ?)',
            [command, JSON.stringify(args), department ?? 'all', nodes.length, req.user.username],
            function (jobErr) {
                const jobId = this?.lastID;
                req.app.locals.broadcast({
                    type: 'alert', level: 'info',
                    message: `Bulk action '${command}' dispatched to ${nodes.length} node(s) (job #${jobId})`
                });
                res.json({ message: 'Bulk action dispatched.', job_id: jobId, dispatched: nodes.length, nodes: nodes.map(n => n.hostname) });
            }
        );
    });
});

// File distribution — push a file to all nodes (or a department)
// POST /api/remote/distribute  { dest_path, department?, content_b64 }
router.post('/distribute', (req, res) => {
    const { dest_path, department, content_b64 } = req.body;
    if (!dest_path || !content_b64) return res.status(400).json({ message: 'dest_path and content_b64 required.' });

    // Prevent path traversal
    const normalized = path.normalize(dest_path);
    if (normalized.startsWith('..')) return res.status(400).json({ message: 'Invalid dest_path.' });

    const sql = department
        ? 'SELECT n.hostname, t.peer_id FROM nodes n JOIN telemetry t ON n.hostname = t.hostname WHERE n.status = ? AND n.department = ?'
        : 'SELECT n.hostname, t.peer_id FROM nodes n JOIN telemetry t ON n.hostname = t.hostname WHERE n.status = ?';
    const params = department ? ['active', department] : ['active'];

    db.all(sql, params, (err, nodes) => {
        if (err) return res.status(500).json({ message: err.message });
        if (!nodes.length) return res.json({ message: 'No active nodes with known PeerIDs found.', dispatched: 0 });

        nodes.forEach(node => {
            axios.post(`${BRIDGE_URL}/bridge/distribute`, {
                peer_id: node.peer_id,
                dest_path: dest_path,
                content_b64: content_b64,
                mode: 0o644
            }).catch(e => console.error(`[MGMT] Failed to distribute file to ${node.hostname}: ${e.message}`));
        });

        db.run(
            'INSERT INTO bulk_actions (command, args, department, node_count, created_by) VALUES (?, ?, ?, ?, ?)',
            ['distribute-file', JSON.stringify({ dest_path }), department ?? 'all', nodes.length, req.user.username],
            function () {
                req.app.locals.broadcast({
                    type: 'alert', level: 'info',
                    message: `File distribution to '${dest_path}' queued for ${nodes.length} node(s)`
                });
                res.json({ message: 'File distribution queued.', dispatched: nodes.length });
            }
        );
    });
});

module.exports = router;

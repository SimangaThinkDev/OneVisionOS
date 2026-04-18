const request = require('supertest');
const express = require('express');
const db = require('../db');
const remoteRoutes = require('../routes/remote');
const authRoutes = require('../routes/auth');
const bcrypt = require('bcryptjs');
const axios = require('axios');

jest.mock('axios');

const app = express();
app.use(express.json());
app.locals.broadcast = jest.fn();
app.use('/api/remote', remoteRoutes);
app.use('/api/auth', authRoutes);

const testAdmin = { username: 'remoteadmin', password: 'password123' };
let authToken = '';

describe('Remote Command Hub API', () => {
    beforeAll((done) => {
        db.serialize(() => {
            // Ensure schema exists
            db.run(`CREATE TABLE IF NOT EXISTS admins (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE NOT NULL, password_hash TEXT NOT NULL, role TEXT DEFAULT 'admin', created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`);
            db.run(`CREATE TABLE IF NOT EXISTS telemetry (peer_id TEXT PRIMARY KEY, hostname TEXT NOT NULL, timestamp TEXT, cpu_percent REAL DEFAULT 0, ram_percent REAL DEFAULT 0, disk_percent REAL DEFAULT 0, status TEXT DEFAULT 'unknown', nids_signatures INTEGER DEFAULT 0)`);
            db.run(`CREATE TABLE IF NOT EXISTS nodes (id INTEGER PRIMARY KEY AUTOINCREMENT, hw_uuid TEXT UNIQUE NOT NULL, hostname TEXT NOT NULL, department TEXT, status TEXT DEFAULT 'pending', created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`);
            db.run(`CREATE TABLE IF NOT EXISTS bulk_actions (id INTEGER PRIMARY KEY AUTOINCREMENT, command TEXT NOT NULL, args TEXT, department TEXT, node_count INTEGER DEFAULT 0, created_by TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`);

            db.run("DELETE FROM admins");
            db.run("DELETE FROM telemetry");
            db.run("DELETE FROM nodes");
            db.run("DELETE FROM bulk_actions");
            
            const hash = bcrypt.hashSync(testAdmin.password, 10);
            db.run("INSERT INTO admins (username, password_hash) VALUES (?, ?)", [testAdmin.username, hash]);
            
            // Insert a test node and telemetry
            db.run("INSERT INTO nodes (hostname, hw_uuid, status, department) VALUES (?, ?, ?, ?)", ['test-node-1', 'test-uuid-1', 'active', 'CS']);
            db.run("INSERT INTO telemetry (hostname, peer_id, timestamp, status) VALUES (?, ?, ?, ?)", 
                ['test-node-1', 'peer-1', Date.now().toString(), 'online'], done);
        });
    });

    beforeEach(() => {
        jest.clearAllMocks();
    });

    test('should allow nodes to push telemetry without auth', async () => {
        const res = await request(app)
            .post('/api/remote/telemetry')
            .send({
                hostname: 'test-node-1',
                peer_id: 'peer-1',
                timestamp: Date.now(),
                metrics: { cpu_usage: 10, memory_usage: 20, disk_usage: 30 },
                status: 'online'
            });
        
        expect(res.statusCode).toEqual(200);
        expect(res.body.message).toBe('Telemetry received.');
        expect(app.locals.broadcast).toHaveBeenCalledWith(expect.objectContaining({ type: 'telemetry' }));
    });

    test('should require auth for GET /telemetry', async () => {
        const res = await request(app).get('/api/remote/telemetry');
        expect(res.statusCode).toEqual(403);
    });

    test('should fetch telemetry with valid token', async () => {
        // Login first
        const loginRes = await request(app).post('/api/auth/login').send(testAdmin);
        authToken = loginRes.body.token;

        const res = await request(app)
            .get('/api/remote/telemetry')
            .set('Authorization', `Bearer ${authToken}`);
        
        expect(res.statusCode).toEqual(200);
        expect(res.body.length).toBeGreaterThan(0);
        expect(res.body[0].hostname).toBe('test-node-1');
    });

    test('should dispatch bulk action with valid token', async () => {
        axios.post.mockResolvedValue({ data: { success: true } });

        const res = await request(app)
            .post('/api/remote/bulk-action')
            .set('Authorization', `Bearer ${authToken}`)
            .send({
                command: 'reboot',
                department: 'CS'
            });
        
        expect(res.statusCode).toEqual(200);
        expect(res.body.dispatched).toBe(1);
        expect(axios.post).toHaveBeenCalled();
        expect(app.locals.broadcast).toHaveBeenCalledWith(expect.objectContaining({ type: 'alert' }));
    });

    test('should block unauthorized commands in bulk-action', async () => {
        const res = await request(app)
            .post('/api/remote/bulk-action')
            .set('Authorization', `Bearer ${authToken}`)
            .send({
                command: 'rm -rf /'
            });
        
        expect(res.statusCode).toEqual(400);
        expect(res.body.message).toContain('Command not allowed');
    });

    test('should queue file distribution', async () => {
        axios.post.mockResolvedValue({ data: { success: true } });

        const res = await request(app)
            .post('/api/remote/distribute')
            .set('Authorization', `Bearer ${authToken}`)
            .send({
                dest_path: '/tmp/test.txt',
                content_b64: 'SGVsbG8gd29ybGQ=' // Hello world
            });
        
        expect(res.statusCode).toEqual(200);
        expect(res.body.dispatched).toBe(1);
        expect(axios.post).toHaveBeenCalled();
    });
});

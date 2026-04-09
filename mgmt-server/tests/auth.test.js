const request = require('supertest');
const express = require('express');
const db = require('../db');
const authRoutes = require('../routes/auth');
const adminRoutes = require('../routes/admin');
const bcrypt = require('bcryptjs');

const app = express();
app.use(express.json());
app.use('/api/auth', authRoutes);
app.use('/api/admin', adminRoutes);

const testAdmin = { username: 'testadmin', password: 'password123' };

jest.setTimeout(10000);

describe('Admin Authentication & API Integration', () => {
    beforeAll((done) => {
        // Clear and setup test schema + admin
        db.serialize(() => {
            db.run("DELETE FROM admins");
            db.run("DELETE FROM students");
            const hash = bcrypt.hashSync(testAdmin.password, 10);
            db.run("INSERT INTO admins (username, password_hash) VALUES (?, ?)", [testAdmin.username, hash], done);
        });
    });

    let authToken = '';

    test('should login successfully and return a token', async () => {
        const res = await request(app)
            .post('/api/auth/login')
            .send(testAdmin);
        
        expect(res.statusCode).toEqual(200);
        expect(res.body).toHaveProperty('token');
        authToken = res.body.token;
    });

    test('should deny access to profiles without a token', async () => {
        const res = await request(app).get('/api/admin/profiles');
        if (res.statusCode !== 403) console.log('DEBUG 500:', res.body);
        expect(res.statusCode).toEqual(403);
    });

    test('should create a student profile with a valid token', async () => {
        const res = await request(app)
            .post('/api/admin/profiles')
            .set('Authorization', `Bearer ${authToken}`)
            .send({
                name: 'Innocent Student',
                grade: '12th',
                email: 'innocent@school.edu',
                mac_address: '00:11:22:33:44:55'
            });
        
        expect(res.statusCode).toEqual(201);
        expect(res.body.message).toBe('Profile created.');
    });

    test('should fetch student profiles with a valid token', async () => {
        const res = await request(app)
            .get('/api/admin/profiles')
            .set('Authorization', `Bearer ${authToken}`);
        
        expect(res.statusCode).toEqual(200);
        expect(res.body.length).toBeGreaterThan(0);
        expect(res.body[0].name).toBe('Innocent Student');
    });
});

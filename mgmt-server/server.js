const express = require('express');
const db = require('./db');
const http = require('http');
const WebSocket = require('ws');
const path = require('path');
require('dotenv').config();

const authRoutes = require('./routes/auth');
const adminRoutes = require('./routes/admin');
const clientRoutes = require('./routes/client');

const app = express();

const server = http.createServer(app);
const wss = new WebSocket.Server({ server });

const port = process.env.PORT || 3000;

app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(express.static(path.join(__dirname, 'public')));

// Routes
app.use('/api/auth', authRoutes);
app.use('/api/admin', adminRoutes);
app.use('/api/client', clientRoutes);


// WebSocket connection
wss.on('connection', (ws) => {
    console.log('New client connected');
    ws.on('message', (message) => {
        console.log(`Received message: ${message}`);
    });
    ws.send(JSON.stringify({ type: 'welcome', message: 'Connected to OneVisionOS Management Server' }));
});

// Broadcast helper
app.set('wss', wss);
app.locals.broadcast = (data) => {
    wss.clients.forEach((client) => {
        if (client.readyState === WebSocket.OPEN) {
            client.send(JSON.stringify(data));
        }
    });
};

app.get('/api/health', (req, res) => {
    res.json({ status: 'healthy', version: '1.0.0' });
});

// Serve the admin dashboard for all other routes to support SPA
app.get('*path', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'index.html'));
});




server.listen(port, () => {
    console.log(`Management server listening at http://localhost:${port}`);
});


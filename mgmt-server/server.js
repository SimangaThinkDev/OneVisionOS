const express = require('express');
const db = require('./db');
require('dotenv').config();

const authRoutes = require('./routes/auth');
const adminRoutes = require('./routes/admin');

const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Routes
app.use('/api/auth', authRoutes);
app.use('/api/admin', adminRoutes);

app.get('/', (req, res) => {
    res.send('OneVisionOS Management Server API is running.');
});

app.get('/api/health', (req, res) => {
    res.json({ status: 'healthy', version: '1.0.0' });
});

app.listen(port, () => {
    console.log(`Management server listening at http://localhost:${port}`);
});

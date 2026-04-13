document.addEventListener('DOMContentLoaded', () => {
    // Initialize Lucide Icons
    lucide.createIcons();

    // DOM Elements
    const navItems = document.querySelectorAll('nav a');
    const sections = {
        dashboard: document.getElementById('dashboard-view'),
        profiles: document.getElementById('profiles-view'),
        nodes: document.getElementById('nodes-view')
    };
    const pageTitle = document.getElementById('page-title');
    const profileModal = document.getElementById('profile-modal');
    const profileForm = document.getElementById('profile-form');
    const alertsContainer = document.getElementById('alerts-container');
    const systemLog = document.getElementById('system-log');
    const gradeChart = document.getElementById('grade-chart');

    // State
    let currentView = 'dashboard';
    let profiles = [];

    // Navigation Logic
    navItems.forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            const target = item.id.replace('nav-', '');
            if (sections[target]) {
                switchView(target);
            } else {
                addAlert('info', `Section "${target}" is under development.`);
            }
        });
    });

    function switchView(view) {
        // Hide all sections
        Object.values(sections).forEach(s => s.style.display = 'none');
        navItems.forEach(i => i.classList.remove('active'));

        // Show target section
        sections[view].style.display = 'block';
        document.getElementById(`nav-${view}`).classList.add('active');
        pageTitle.textContent = view.charAt(0).toUpperCase() + view.slice(1);
        currentView = view;

        if (view === 'profiles') {
            fetchProfiles();
        } else if (view === 'nodes') {
            fetchNodes();
        }
    }

    // Modal Logic
    document.getElementById('btn-add-profile').addEventListener('click', () => {
        profileForm.reset();
        document.getElementById('profile-id').value = '';
        document.getElementById('modal-title').textContent = 'New Student Profile';
        profileModal.style.display = 'flex';
    });

    document.getElementById('btn-cancel').addEventListener('click', () => {
        profileModal.style.display = 'none';
    });

    // Fetch Nodes
    async function fetchNodes() {
        try {
            const response = await fetch('/api/admin/nodes', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
            });
            if (!response.ok) throw new Error('Failed to fetch nodes');
            const nodes = await response.json();
            renderNodes(nodes);
            document.getElementById('stat-nodes').textContent = nodes.filter(n => n.status === 'active').length;
        } catch (error) {
            addAlert('critical', error.message);
        }
    }

    function renderNodes(nodes) {
        const tbody = document.querySelector('#nodes-table tbody');
        if (nodes.length === 0) {
            tbody.innerHTML = '<tr><td colspan="5" style="text-align: center; color: var(--text-muted);">No managed nodes connected yet.</td></tr>';
            return;
        }
        tbody.innerHTML = nodes.map(n => `
            <tr>
                <td>${n.hostname}</td>
                <td>${n.ip_address || '-'}</td>
                <td><span style="color: ${n.status === 'active' ? 'var(--accent-color)' : 'var(--text-muted)'}">${n.status}</span></td>
                <td>${new Date(n.last_seen).toLocaleString()}</td>
                <td>
                    <button class="btn-icon" title="View Details"><i data-lucide="info"></i></button>
                </td>
            </tr>
        `).join('');
        lucide.createIcons();
    }

    // Fetch Profiles

    async function fetchProfiles() {
        try {
            const response = await fetch('/api/admin/profiles', {
                headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
            });
            if (!response.ok) throw new Error('Failed to fetch profiles');
            profiles = await response.json();
            renderProfiles();
            updateStats();
            renderGradeChart();
        } catch (error) {
            addAlert('critical', error.message);
        }
    }

    function renderProfiles() {
        const tbody = document.querySelector('#profiles-table tbody');
        tbody.innerHTML = profiles.map(p => `
            <tr>
                <td>${p.id}</td>
                <td>${p.name}</td>
                <td>${p.grade}</td>
                <td>${p.email || '-'}</td>
                <td><code>${p.mac_address || '-'}</code></td>
                <td>
                    <button class="btn-icon" onclick="deleteProfile(${p.id})" style="background:none; border:none; cursor:pointer; color:var(--danger-color);">
                        <i data-lucide="trash-2"></i>
                    </button>
                </td>
            </tr>
        `).join('');
        lucide.createIcons();
    }

    function renderGradeChart() {
        const counts = {};
        profiles.forEach(p => {
            counts[p.grade] = (counts[p.grade] || 0) + 1;
        });

        const maxCount = Math.max(...Object.values(counts), 1);
        gradeChart.innerHTML = Object.entries(counts).sort((a,b) => a[0]-b[0]).map(([grade, count]) => `
            <div style="flex: 1; display: flex; flex-direction: column; align-items: center; gap: 0.5rem;">
                <div style="width: 100%; background: var(--primary-color); height: ${(count/maxCount)*100}%; border-radius: 4px 4px 0 0; position: relative;" title="${count} students">
                    <span style="position: absolute; top: -20px; width: 100%; text-align: center; font-size: 0.75rem;">${count}</span>
                </div>
                <span style="font-size: 0.75rem; color: var(--text-muted);">G${grade}</span>
            </div>
        `).join('');
    }

    function updateStats() {
        const statProfiles = document.getElementById('stat-profiles');
        if (statProfiles) statProfiles.textContent = profiles.length;
    }

    function logSystemMessage(message) {
        const div = document.createElement('div');
        div.style.padding = '0.5rem 0';
        div.style.borderBottom = '1px solid rgba(255,255,255,0.05)';
        div.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
        systemLog.prepend(div);
    }


    // Add Profile
    profileForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const data = {
            name: document.getElementById('profile-name').value,
            grade: document.getElementById('profile-grade').value,
            email: document.getElementById('profile-email').value,
            mac_address: document.getElementById('profile-mac').value
        };

        try {
            const response = await fetch('/api/admin/profiles', {
                method: 'POST',
                headers: { 
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                },
                body: JSON.stringify(data)
            });

            if (!response.ok) throw new Error('Failed to save profile');
            
            addAlert('success', 'Profile saved successfully.');
            profileModal.style.display = 'none';
            fetchProfiles();
        } catch (error) {
            addAlert('critical', error.message);
        }
    });

    // Delete Profile (Global function for easy onclick)
    window.deleteProfile = async (id) => {
        if (!confirm('Are you sure you want to delete this profile?')) return;

        try {
            const response = await fetch(`/api/admin/profiles/${id}`, {
                method: 'DELETE',
                headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
            });

            if (!response.ok) throw new Error('Failed to delete profile');
            
            addAlert('info', 'Profile deleted.');
            fetchProfiles();
        } catch (error) {
            addAlert('critical', error.message);
        }
    };

    // WebSocket for Real-time Alerts
    function setupWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const socket = new WebSocket(`${protocol}//${window.location.host}`);

        socket.onopen = () => {
            console.log('Connected to Management Server WebSocket');
        };

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === 'alert') {
                addAlert(data.level || 'info', data.message);
                logSystemMessage(data.message);
            } else if (data.type === 'welcome') {
                console.log('Server:', data.message);
                logSystemMessage('Connection established with management server.');
            }
        };


        socket.onclose = () => {
            console.log('WebSocket connection closed. Retrying in 5s...');
            setTimeout(setupWebSocket, 5000);
        };
    }

    function addAlert(level, message) {
        const div = document.createElement('div');
        div.className = `alert-item ${level}`;
        div.innerHTML = `
            <div style="font-weight: 700; margin-bottom: 0.25rem;">${level.toUpperCase()}</div>
            <div>${message}</div>
        `;
        alertsContainer.prepend(div);
        setTimeout(() => {
            div.style.opacity = '0';
            div.style.transform = 'translateX(20px)';
            setTimeout(() => div.remove(), 300);
        }, 8000);
    }

    // Authentication Logic
    const loginOverlay = document.getElementById('login-overlay');
    const loginForm = document.getElementById('login-form');

    if (localStorage.getItem('token')) {
        loginOverlay.style.display = 'none';
        setupWebSocket();
        fetchProfiles();
    }

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const username = document.getElementById('login-username').value;
        const password = document.getElementById('login-password').value;

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });

            const data = await response.json();
            if (response.ok) {
                localStorage.setItem('token', data.token);
                loginOverlay.style.display = 'none';
                addAlert('success', 'Logged in successfully.');
                setupWebSocket();
                fetchProfiles();
            } else {
                throw new Error(data.message || 'Login failed');
            }
        } catch (error) {
            alert(error.message);
        }
    });

    // Initialize
    // Don't call fetchProfiles here, it's called after login
});


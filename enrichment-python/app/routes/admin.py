"""
Admin UI Routes for Enrichment Service
Professional, minimal admin dashboard for viewing agent enrichment data
"""
from fastapi import APIRouter, Query, HTTPException
from fastapi.responses import HTMLResponse
from typing import Optional, List, Dict
import httpx
from ..core.database import db_manager
from ..core.config import settings
from ..core.logging import get_logger

logger = get_logger(__name__)
router = APIRouter(prefix="/admin", tags=["admin"])

API_BASE_URL = "http://zerotrace-api:8080/api"

@router.get("/ui", response_class=HTMLResponse)
async def admin_ui():
    """Professional Admin UI for browsing enrichment data"""
    return """
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ZeroTrace Admin</title>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --primary: #0f172a;
            --primary-light: #1e293b;
            --accent: #3b82f6;
            --accent-hover: #2563eb;
            --success: #22c55e;
            --warning: #eab308;
            --danger: #ef4444;
            --critical: #a855f7;
            --border: #e2e8f0;
            --bg: #f8fafc;
            --card: #ffffff;
            --text: #0f172a;
            --text-muted: #64748b;
            --text-light: #94a3b8;
        }
        
        * { margin: 0; padding: 0; box-sizing: border-box; }
        
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg);
            color: var(--text);
            font-size: 14px;
            line-height: 1.5;
        }
        
        .layout { display: flex; min-height: 100vh; }
        
        /* Sidebar */
        .sidebar {
            width: 240px;
            background: var(--primary);
            color: white;
            display: flex;
            flex-direction: column;
            position: fixed;
            height: 100vh;
            z-index: 100;
        }
        
        .logo {
            padding: 24px;
            border-bottom: 1px solid rgba(255,255,255,0.1);
        }
        
        .logo h1 {
            font-size: 16px;
            font-weight: 600;
            letter-spacing: -0.5px;
        }
        
        .logo span {
            font-size: 11px;
            color: rgba(255,255,255,0.5);
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        
        .nav { padding: 16px 0; flex: 1; }
        
        .nav-group { margin-bottom: 24px; }
        
        .nav-label {
            padding: 0 24px;
            font-size: 10px;
            text-transform: uppercase;
            letter-spacing: 1px;
            color: rgba(255,255,255,0.4);
            margin-bottom: 8px;
        }
        
        .nav-item {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 10px 24px;
            color: rgba(255,255,255,0.7);
            text-decoration: none;
            cursor: pointer;
            transition: all 0.15s;
            font-size: 13px;
            border-left: 2px solid transparent;
        }
        
        .nav-item:hover {
            background: rgba(255,255,255,0.05);
            color: white;
        }
        
        .nav-item.active {
            background: rgba(59, 130, 246, 0.15);
            color: white;
            border-left-color: var(--accent);
        }
        
        .nav-icon {
            width: 18px;
            height: 18px;
            opacity: 0.7;
        }
        
        /* Main */
        .main {
            flex: 1;
            margin-left: 240px;
            min-height: 100vh;
        }
        
        .header {
            background: var(--card);
            border-bottom: 1px solid var(--border);
            padding: 20px 32px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            position: sticky;
            top: 0;
            z-index: 50;
        }
        
        .header h2 {
            font-size: 18px;
            font-weight: 600;
        }
        
        .content { padding: 32px; }
        
        /* Cards */
        .card {
            background: var(--card);
            border: 1px solid var(--border);
            border-radius: 8px;
            margin-bottom: 24px;
        }
        
        .card-header {
            padding: 16px 20px;
            border-bottom: 1px solid var(--border);
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .card-title {
            font-size: 14px;
            font-weight: 600;
        }
        
        .card-body { padding: 20px; }
        
        /* Stats */
        .stats {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 20px;
            margin-bottom: 32px;
        }
        
        .stat {
            background: var(--card);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 20px;
        }
        
        .stat-label {
            font-size: 12px;
            color: var(--text-muted);
            text-transform: uppercase;
            letter-spacing: 0.5px;
            margin-bottom: 8px;
        }
        
        .stat-value {
            font-size: 28px;
            font-weight: 700;
            color: var(--text);
        }
        
        .stat-meta {
            font-size: 12px;
            color: var(--text-light);
            margin-top: 4px;
        }
        
        /* Tables */
        table {
            width: 100%;
            border-collapse: collapse;
        }
        
        th {
            text-align: left;
            padding: 12px 16px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            color: var(--text-muted);
            background: var(--bg);
            border-bottom: 1px solid var(--border);
        }
        
        td {
            padding: 14px 16px;
            border-bottom: 1px solid var(--border);
            font-size: 13px;
        }
        
        tr:hover { background: var(--bg); }
        
        tr:last-child td { border-bottom: none; }
        
        /* Badges */
        .badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 11px;
            font-weight: 500;
        }
        
        .badge-success { background: #dcfce7; color: #166534; }
        .badge-warning { background: #fef3c7; color: #92400e; }
        .badge-danger { background: #fee2e2; color: #991b1b; }
        .badge-info { background: #dbeafe; color: #1e40af; }
        .badge-neutral { background: #f1f5f9; color: #475569; }
        
        /* Severity */
        .severity { font-weight: 600; font-size: 11px; }
        .severity-critical { color: var(--critical); }
        .severity-high { color: var(--danger); }
        .severity-medium { color: var(--warning); }
        .severity-low { color: var(--success); }
        
        /* Buttons */
        .btn {
            padding: 8px 16px;
            border: 1px solid var(--border);
            border-radius: 6px;
            background: var(--card);
            font-size: 13px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.15s;
        }
        
        .btn:hover { background: var(--bg); }
        
        .btn-primary {
            background: var(--accent);
            border-color: var(--accent);
            color: white;
        }
        
        .btn-primary:hover { background: var(--accent-hover); }
        
        .btn-sm { padding: 5px 10px; font-size: 12px; }
        
        /* Spinner */
        .spinner {
            width: 20px;
            height: 20px;
            border: 2px solid var(--border);
            border-top-color: var(--accent);
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
        }
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        
        /* Search */
        .search-bar {
            display: flex;
            gap: 12px;
            margin-bottom: 24px;
        }
        
        .search-input {
            flex: 1;
            padding: 10px 14px;
            border: 1px solid var(--border);
            border-radius: 6px;
            font-size: 13px;
            background: var(--card);
        }
        
        .search-input:focus {
            outline: none;
            border-color: var(--accent);
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
        }
        
        /* Agent Detail */
        .agent-header {
            display: flex;
            align-items: center;
            gap: 16px;
            margin-bottom: 24px;
        }
        
        .agent-icon {
            width: 48px;
            height: 48px;
            background: var(--primary);
            border-radius: 8px;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: 600;
        }
        
        .agent-info h3 { font-size: 18px; font-weight: 600; }
        .agent-info p { color: var(--text-muted); font-size: 13px; }
        
        /* Software Item */
        .software-item {
            border: 1px solid var(--border);
            border-radius: 6px;
            margin-bottom: 12px;
        }
        
        .software-header {
            padding: 14px 16px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            cursor: pointer;
            background: var(--bg);
        }
        
        .software-header:hover { background: #f1f5f9; }
        
        .software-name { font-weight: 500; }
        .software-version { color: var(--text-muted); font-size: 12px; margin-left: 8px; }
        
        .software-body {
            padding: 16px;
            border-top: 1px solid var(--border);
            display: none;
        }
        
        .software-body.open { display: block; }
        
        /* Flow */
        .flow {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 12px;
            padding: 12px;
            background: var(--bg);
            border-radius: 6px;
            margin-bottom: 16px;
        }
        
        .flow-arrow { color: var(--text-light); }
        
        .flow-item {
            padding: 4px 10px;
            border-radius: 4px;
            font-weight: 500;
        }
        
        .flow-software { background: #dbeafe; color: #1e40af; }
        .flow-cpe { background: #dcfce7; color: #166534; }
        .flow-cve { background: #fef3c7; color: #92400e; }
        .flow-exploit { background: #fee2e2; color: #991b1b; }
        
        /* CVE */
        .cve-item {
            padding: 14px;
            border: 1px solid var(--border);
            border-radius: 6px;
            margin-bottom: 10px;
        }
        
        .cve-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;
        }
        
        .cve-id { font-weight: 600; color: var(--accent); }
        .cve-desc { font-size: 13px; color: var(--text-muted); }
        
        /* Loading */
        .loading { text-align: center; padding: 40px; color: var(--text-muted); }
        
        /* Empty */
        .empty {
            text-align: center;
            padding: 60px 20px;
            color: var(--text-muted);
        }
        
        .empty-title { font-size: 16px; font-weight: 500; margin-bottom: 8px; color: var(--text); }
        
        /* Section */
        .section { display: none; }
        .section.active { display: block; }
        
        /* Responsive */
        @media (max-width: 1200px) {
            .stats { grid-template-columns: repeat(2, 1fr); }
        }
        
        @media (max-width: 768px) {
            .sidebar { display: none; }
            .main { margin-left: 0; }
            .stats { grid-template-columns: 1fr; }
        }
    </style>
</head>
<body>
    <div class="layout">
        <aside class="sidebar">
            <div class="logo">
                <h1>ZeroTrace</h1>
                <span>Admin Console</span>
            </div>
            <nav class="nav">
                <div class="nav-group">
                    <div class="nav-label">Overview</div>
                    <a class="nav-item active" onclick="showSection('dashboard')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path></svg>
                        Dashboard
                    </a>
                </div>
                <div class="nav-group">
                    <div class="nav-label">Monitoring</div>
                    <a class="nav-item" onclick="showSection('agents')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
                        Agents
                    </a>
                    <a class="nav-item" onclick="showSection('software')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path></svg>
                        Software
                    </a>
                    <a class="nav-item" onclick="showSection('vulnerabilities')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path></svg>
                        Vulnerabilities
                    </a>
                </div>
                <div class="nav-group">
                    <div class="nav-label">Data</div>
                    <a class="nav-item" onclick="showSection('cve-search')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path></svg>
                        CVE Search
                    </a>
                    <a class="nav-item" onclick="showSection('cpe-lookup')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"></path></svg>
                        CPE Lookup
                    </a>
                    <a class="nav-item" onclick="showSection('test')">
                        <svg class="nav-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path></svg>
                        Test Enrichment
                    </a>
                </div>
            </nav>
        </aside>
        
        <main class="main">
            <header class="header">
                <h2 id="page-title">Dashboard</h2>
                <button class="btn" onclick="refreshData()">Refresh</button>
            </header>
            
            <div class="content">
                <!-- Dashboard -->
                <section id="dashboard-section" class="section active">
                    <div class="stats">
                        <div class="stat">
                            <div class="stat-label">Agents</div>
                            <div class="stat-value" id="stat-agents">-</div>
                            <div class="stat-meta">Active endpoints</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Software</div>
                            <div class="stat-value" id="stat-software">-</div>
                            <div class="stat-meta">Tracked packages</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Vulnerabilities</div>
                            <div class="stat-value" id="stat-vulns">-</div>
                            <div class="stat-meta">Detected issues</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">CVE Database</div>
                            <div class="stat-value" id="stat-cves">-</div>
                            <div class="stat-meta">Total records</div>
                        </div>
                    </div>
                    
                    <div class="card">
                        <div class="card-header">
                            <span class="card-title">Connected Agents</span>
                        </div>
                        <div class="card-body" style="padding: 0;">
                            <table>
                                <thead>
                                    <tr>
                                        <th>Agent</th>
                                        <th>Hostname</th>
                                        <th>Platform</th>
                                        <th>Software</th>
                                        <th>Issues</th>
                                        <th>Status</th>
                                        <th></th>
                                    </tr>
                                </thead>
                                <tbody id="agents-table">
                                    <tr><td colspan="7" class="loading">Loading...</td></tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </section>
                
                <!-- Agents -->
                <section id="agents-section" class="section">
                    <div class="search-bar">
                        <input type="text" class="search-input" id="agent-search" placeholder="Search by name, hostname, or IP..." />
                        <button class="btn btn-primary" onclick="searchAgents()">Search</button>
                    </div>
                    <div class="card">
                        <div class="card-body" id="agents-list" style="padding: 0;">
                            <div class="loading">Loading...</div>
                        </div>
                    </div>
                </section>
                
                <!-- Agent Detail -->
                <section id="agent-detail-section" class="section">
                    <div id="agent-detail"></div>
                </section>
                
                <!-- Software -->
                <section id="software-section" class="section">
                    <div class="card">
                        <div class="card-header">
                            <span class="card-title">Software Inventory</span>
                        </div>
                        <div class="card-body" id="software-list">
                            <div class="loading">Loading...</div>
                        </div>
                    </div>
                </section>
                
                <!-- Vulnerabilities -->
                <section id="vulnerabilities-section" class="section">
                    <div class="card">
                        <div class="card-header" style="display: flex; justify-content: space-between; align-items: center;">
                            <span class="card-title">Detected Vulnerabilities (Debug View)</span>
                            <button class="btn btn-primary" onclick="runEnrichment()" id="enrich-btn">
                                Run Enrichment
                            </button>
                        </div>
                        <div id="enrichment-progress" style="display: none; padding: 16px; background: #f8f9fa; border-bottom: 1px solid #e9ecef;">
                            <div style="display: flex; align-items: center; gap: 12px;">
                                <div class="spinner"></div>
                                <span id="enrichment-status">Enriching software...</span>
                            </div>
                            <div style="margin-top: 12px; height: 4px; background: #e9ecef; border-radius: 2px; overflow: hidden;">
                                <div id="enrichment-bar" style="height: 100%; width: 0%; background: linear-gradient(90deg, #3b82f6, #8b5cf6); transition: width 0.3s;"></div>
                            </div>
                        </div>
                        <div class="card-body" style="padding: 0;">
                            <table>
                                <thead>
                                    <tr>
                                        <th>CVE ID</th>
                                        <th>Severity</th>
                                        <th>Software</th>
                                        <th>Agent</th>
                                        <th>Exploit</th>
                                    </tr>
                                </thead>
                                <tbody id="vulns-table">
                                    <tr><td colspan="5" class="loading">Click "Run Enrichment" to scan for vulnerabilities</td></tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </section>
                
                <!-- CVE Search -->
                <section id="cve-search-section" class="section">
                    <div class="search-bar">
                        <input type="text" class="search-input" id="cve-input" placeholder="Search by CVE ID or keyword..." />
                        <button class="btn btn-primary" onclick="searchCVE()">Search</button>
                    </div>
                    <div class="card">
                        <div class="card-body" id="cve-results"></div>
                    </div>
                </section>
                
                <!-- CPE Lookup -->
                <section id="cpe-lookup-section" class="section">
                    <div class="search-bar">
                        <input type="text" class="search-input" id="cpe-input" placeholder="Search by vendor or product..." />
                        <button class="btn btn-primary" onclick="searchCPE()">Search</button>
                    </div>
                    <div class="card">
                        <div class="card-body" id="cpe-results"></div>
                    </div>
                </section>
                
                <!-- Test -->
                <section id="test-section" class="section">
                    <div class="card">
                        <div class="card-header">
                            <span class="card-title">Test Enrichment</span>
                        </div>
                        <div class="card-body">
                            <div class="search-bar">
                                <input type="text" class="search-input" id="test-name" placeholder="Software name" />
                                <input type="text" class="search-input" id="test-version" placeholder="Version" />
                                <input type="text" class="search-input" id="test-vendor" placeholder="Vendor (optional)" />
                                <button class="btn btn-primary" onclick="testEnrichment()">Test</button>
                            </div>
                            <div id="test-results"></div>
                        </div>
                    </div>
                </section>
            </div>
        </main>
    </div>
    
    <script>
        const API = 'http://localhost:8080/api';
        
        function showSection(name) {
            document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
            document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
            document.getElementById(name + '-section').classList.add('active');
            event.target.closest('.nav-item').classList.add('active');
            
            const titles = {
                'dashboard': 'Dashboard',
                'agents': 'Agents',
                'agent-detail': 'Agent Details',
                'software': 'Software Inventory',
                'vulnerabilities': 'Vulnerabilities',
                'cve-search': 'CVE Database',
                'cpe-lookup': 'CPE Dictionary',
                'test': 'Test Enrichment'
            };
            document.getElementById('page-title').textContent = titles[name] || name;
            
            if (name === 'dashboard') loadDashboard();
            if (name === 'agents') loadAgents();
            if (name === 'software') loadSoftware();
            if (name === 'vulnerabilities') loadVulns();
        }
        
        async function loadDashboard() {
            try {
                const res = await fetch('/admin/dashboard-stats');
                const data = await res.json();
                
                document.getElementById('stat-agents').textContent = data.agent_count || 0;
                document.getElementById('stat-software').textContent = data.software_count || 0;
                document.getElementById('stat-vulns').textContent = data.vulnerability_count || 0;
                document.getElementById('stat-cves').textContent = (data.cve_database_count || 0).toLocaleString();
                
                const agentsRes = await fetch('/admin/agents');
                const agents = await agentsRes.json();
                
                const tbody = document.getElementById('agents-table');
                if (agents.agents && agents.agents.length > 0) {
                    tbody.innerHTML = agents.agents.map(a => `
                        <tr>
                            <td><strong>${a.name || a.id}</strong></td>
                            <td>${a.hostname || '-'}</td>
                            <td>${a.os || '-'}</td>
                            <td>${a.software_count || 0}</td>
                            <td><span class="badge ${a.vuln_count > 0 ? 'badge-danger' : 'badge-success'}">${a.vuln_count || 0}</span></td>
                            <td><span class="badge ${a.status === 'online' ? 'badge-success' : 'badge-neutral'}">${a.status || 'offline'}</span></td>
                            <td><button class="btn btn-sm" onclick="viewAgent('${a.id}')">View</button></td>
                        </tr>
                    `).join('');
                } else {
                    tbody.innerHTML = '<tr><td colspan="7" class="empty">No agents connected</td></tr>';
                }
            } catch (e) {
                console.error(e);
            }
        }
        
        async function viewAgent(id) {
            showSection('agent-detail');
            const container = document.getElementById('agent-detail');
            container.innerHTML = '<div class="loading">Loading...</div>';
            
            try {
                const res = await fetch(`/admin/agents/${id}/enrichment`);
                const data = await res.json();
                
                container.innerHTML = `
                    <div class="agent-header">
                        <div class="agent-icon">${(data.agent.name || 'A')[0].toUpperCase()}</div>
                        <div class="agent-info">
                            <h3>${data.agent.name || data.agent.id}</h3>
                            <p>${data.agent.hostname || '-'} | ${data.agent.os || '-'}</p>
                        </div>
                    </div>
                    
                    <div class="stats" style="grid-template-columns: repeat(4, 1fr);">
                        <div class="stat">
                            <div class="stat-label">Software</div>
                            <div class="stat-value">${data.software?.length || 0}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">CPE Matches</div>
                            <div class="stat-value">${data.cpe_matches || 0}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Vulnerabilities</div>
                            <div class="stat-value">${data.vulnerability_count || 0}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Exploits</div>
                            <div class="stat-value">${data.exploit_count || 0}</div>
                        </div>
                    </div>
                    
                    <div class="card">
                        <div class="card-header">
                            <span class="card-title">Software Enrichment</span>
                        </div>
                        <div class="card-body">
                            ${data.software && data.software.length > 0 ? data.software.slice(0, 50).map((sw, i) => `
                                <div class="software-item">
                                    <div class="software-header" onclick="toggleSw(${i})">
                                        <div>
                                            <span class="software-name">${sw.name}</span>
                                            <span class="software-version">${sw.version}</span>
                                        </div>
                                        <div>
                                            ${sw.cpe ? '<span class="badge badge-success">CPE</span>' : '<span class="badge badge-neutral">No CPE</span>'}
                                            ${sw.vulnerabilities?.length > 0 ? `<span class="badge badge-danger">${sw.vulnerabilities.length} CVE</span>` : ''}
                                        </div>
                                    </div>
                                    <div class="software-body" id="sw-${i}">
                                        <div class="flow">
                                            <span class="flow-item flow-software">${sw.name}</span>
                                            <span class="flow-arrow">→</span>
                                            <span class="flow-item flow-cpe">${sw.cpe || 'No match'}</span>
                                            <span class="flow-arrow">→</span>
                                            <span class="flow-item flow-cve">${sw.vulnerabilities?.length || 0} CVEs</span>
                                        </div>
                                        ${sw.vulnerabilities?.length > 0 ? sw.vulnerabilities.slice(0, 5).map(v => `
                                            <div class="cve-item">
                                                <div class="cve-header">
                                                    <span class="cve-id">${v.id}</span>
                                                    <span class="severity severity-${(v.severity || 'low').toLowerCase()}">${v.severity || 'N/A'}</span>
                                                </div>
                                                <p class="cve-desc">${(v.description || '').slice(0, 200)}...</p>
                                            </div>
                                        `).join('') : '<p style="color: var(--text-muted);">No vulnerabilities detected</p>'}
                                    </div>
                                </div>
                            `).join('') : '<div class="empty">No software data</div>'}
                        </div>
                    </div>
                `;
            } catch (e) {
                container.innerHTML = '<div class="empty">Failed to load agent data</div>';
            }
        }
        
        function toggleSw(i) {
            document.getElementById('sw-' + i).classList.toggle('open');
        }
        
        async function loadAgents() {
            const container = document.getElementById('agents-list');
            try {
                const res = await fetch('/admin/agents');
                const data = await res.json();
                
                if (data.agents?.length > 0) {
                    container.innerHTML = `
                        <table>
                            <thead>
                                <tr>
                                    <th>Agent</th>
                                    <th>Hostname</th>
                                    <th>Platform</th>
                                    <th>IP</th>
                                    <th>Software</th>
                                    <th>Status</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                ${data.agents.map(a => `
                                    <tr>
                                        <td><strong>${a.name || a.id}</strong></td>
                                        <td>${a.hostname || '-'}</td>
                                        <td>${a.os || '-'}</td>
                                        <td>${a.ip || '-'}</td>
                                        <td>${a.software_count || 0}</td>
                                        <td><span class="badge ${a.status === 'online' ? 'badge-success' : 'badge-neutral'}">${a.status}</span></td>
                                        <td><button class="btn btn-sm" onclick="viewAgent('${a.id}')">View</button></td>
                                    </tr>
                                `).join('')}
                            </tbody>
                        </table>
                    `;
                } else {
                    container.innerHTML = '<div class="empty"><div class="empty-title">No agents</div><p>Connect agents to see data</p></div>';
                }
            } catch (e) {
                container.innerHTML = '<div class="empty">Failed to load agents</div>';
            }
        }
        
        async function loadSoftware() {
            const container = document.getElementById('software-list');
            container.innerHTML = '<div class="loading">Loading...</div>';
            try {
                const res = await fetch('/admin/agents');
                const data = await res.json();
                
                if (data.agents?.length > 0) {
                    let html = '';
                    for (const agent of data.agents) {
                        html += `
                            <div style="margin-bottom: 24px;">
                                <h4 style="margin-bottom: 12px;">${agent.name || agent.id} - ${agent.hostname || 'Unknown'}</h4>
                                <p style="color: var(--text-muted); margin-bottom: 16px;">Software count: ${agent.software_count || 0}</p>
                            </div>
                        `;
                    }
                    container.innerHTML = html || '<div class="empty">No software data available</div>';
                } else {
                    container.innerHTML = '<div class="empty">No agents connected</div>';
                }
            } catch (e) {
                container.innerHTML = '<div class="empty">Failed to load software</div>';
            }
        }
        
        async function loadVulns() {
            const tbody = document.getElementById('vulns-table');
            tbody.innerHTML = '<tr><td colspan="5" class="empty">Click "Run Enrichment" to scan for vulnerabilities</td></tr>';
        }
        
        let enrichmentVulns = [];
        
        async function runEnrichment() {
            const btn = document.getElementById('enrich-btn');
            const progress = document.getElementById('enrichment-progress');
            const status = document.getElementById('enrichment-status');
            const bar = document.getElementById('enrichment-bar');
            const tbody = document.getElementById('vulns-table');
            
            btn.disabled = true;
            btn.textContent = 'Running...';
            progress.style.display = 'block';
            tbody.innerHTML = '<tr><td colspan="5" class="loading">Enriching...</td></tr>';
            enrichmentVulns = [];
            
            try {
                // Get agents and their software
                status.textContent = 'Fetching agents...';
                bar.style.width = '10%';
                
                const agentsRes = await fetch('/admin/agents');
                const agentsData = await agentsRes.json();
                
                if (!agentsData.agents || agentsData.agents.length === 0) {
                    throw new Error('No agents found');
                }
                
                let totalVulns = 0;
                const allVulns = [];
                
                for (let i = 0; i < agentsData.agents.length; i++) {
                    const agent = agentsData.agents[i];
                    status.textContent = `Processing ${agent.hostname || agent.name}...`;
                    bar.style.width = (20 + (i / agentsData.agents.length) * 30) + '%';
                    
                    // Get agent software details
                    const detailRes = await fetch(`/admin/agents/${agent.id}/enrichment`);
                    const detailData = await detailRes.json();
                    
                    if (detailData.software && detailData.software.length > 0) {
                        // Enrich software in batches - DEBUG mode (limit 50)
                        const samples = detailData.software.slice(0, 50); 
                        const batchSize = 10;
                        
                        for (let j = 0; j < samples.length; j += batchSize) {
                            const batch = samples.slice(j, j + batchSize);
                            status.textContent = `Enriching ${j}/${samples.length} packages from ${agent.hostname}...`;
                            bar.style.width = (50 + (j / samples.length) * 40) + '%';
                            
                            try {
                                const enrichRes = await fetch('/enrich/software', {
                                    method: 'POST',
                                    headers: {'Content-Type': 'application/json'},
                                    body: JSON.stringify(batch)
                                });
                                const enrichData = await enrichRes.json();
                                
                                if (enrichData.data) {
                                    enrichData.data.forEach(sw => {
                                        if (sw.vulnerabilities && sw.vulnerabilities.length > 0) {
                                            sw.vulnerabilities.forEach(v => {
                                                allVulns.push({
                                                    ...v,
                                                    software: sw.name + ' ' + sw.version,
                                                    agent: agent.hostname || agent.name,
                                                    cpe: sw.cpe
                                                });
                                                totalVulns++;
                                            });
                                        }
                                    });
                                }
                            } catch (e) {
                                console.warn('Batch enrichment failed', e);
                            }
                        }
                    }
                }
                
                bar.style.width = '100%';
                status.textContent = `Complete! Found ${totalVulns} vulnerabilities`;
                enrichmentVulns = allVulns;
                
                // Display results
                displayVulnerabilities(allVulns);
                
                setTimeout(() => {
                    progress.style.display = 'none';
                }, 2000);
                
            } catch (e) {
                status.textContent = 'Error: ' + e.message;
                bar.style.width = '0%';
                bar.style.background = '#ef4444';
            } finally {
                btn.disabled = false;
                btn.textContent = 'Run Enrichment';
            }
        }
        
        function displayVulnerabilities(vulns) {
            const tbody = document.getElementById('vulns-table');
            
            if (vulns.length === 0) {
                tbody.innerHTML = '<tr><td colspan="5" class="empty">No vulnerabilities found</td></tr>';
                return;
            }
            
            // Sort by severity
            const severityOrder = {critical: 0, high: 1, medium: 2, low: 3, unknown: 4};
            vulns.sort((a, b) => (severityOrder[a.severity] || 5) - (severityOrder[b.severity] || 5));
            
            let html = '';
            vulns.slice(0, 100).forEach(v => {
                const sevClass = v.severity === 'critical' ? 'critical' : 
                                v.severity === 'high' ? 'high' : 
                                v.severity === 'medium' ? 'medium' : 'low';
                html += `
                    <tr>
                        <td><a href="https://nvd.nist.gov/vuln/detail/${v.cve_id}" target="_blank" style="color: #3b82f6; text-decoration: none;">${v.cve_id}</a></td>
                        <td><span class="badge badge-${sevClass}">${v.severity || 'unknown'}</span></td>
                        <td>${v.software}</td>
                        <td>${v.agent}</td>
                        <td>${v.known_exploited ? '<span class="badge badge-critical">KEV</span>' : '-'}</td>
                    </tr>
                `;
            });
            
            tbody.innerHTML = html;
        }
        
        async function testEnrichment() {
            const name = document.getElementById('test-name').value;
            const version = document.getElementById('test-version').value;
            const vendor = document.getElementById('test-vendor').value;
            const results = document.getElementById('test-results');
            
            if (!name || !version) {
                results.innerHTML = '<p style="color: var(--danger);">Name and version required</p>';
                return;
            }
            
            results.innerHTML = '<div class="loading">Processing...</div>';
            
            try {
                const res = await fetch('/enrich/software', {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify([{name, version, vendor}])
                });
                const data = await res.json();
                const item = data.data[0];
                
                results.innerHTML = `
                    <div class="flow" style="margin-top: 20px;">
                        <span class="flow-item flow-software">${name} ${version}</span>
                        <span class="flow-arrow">→</span>
                        <span class="flow-item flow-cpe">${item.cpe || 'No CPE'}</span>
                        <span class="flow-arrow">→</span>
                        <span class="flow-item flow-cve">${item.vulnerabilities?.length || 0} CVEs</span>
                    </div>
                    
                    <div class="stats" style="margin: 20px 0;">
                        <div class="stat">
                            <div class="stat-label">CPE Match</div>
                            <div class="stat-value">${item.cpe ? 'Yes' : 'No'}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">CVEs Found</div>
                            <div class="stat-value">${item.vulnerabilities?.length || 0}</div>
                        </div>
                        <div class="stat">
                            <div class="stat-label">Source</div>
                            <div class="stat-value">${item.source || 'none'}</div>
                        </div>
                    </div>
                    
                    ${item.cpe ? `<p><strong>CPE:</strong> <code style="background:#f1f5f9;padding:4px 8px;border-radius:4px;">${item.cpe}</code></p>` : ''}
                    
                    ${item.vulnerabilities?.length > 0 ? `
                        <h4 style="margin: 24px 0 12px;">Vulnerabilities:</h4>
                        ${item.vulnerabilities.map(v => `
                            <div class="cve-item">
                                <div class="cve-header">
                                    <span class="cve-id">${v.id}</span>
                                    <span class="severity severity-${(v.severity || 'low').toLowerCase()}">${v.severity || 'N/A'}</span>
                                </div>
                                <p class="cve-desc">${v.description || 'No description'}</p>
                            </div>
                        `).join('')}
                    ` : '<p style="margin-top: 20px; color: var(--text-muted);">No vulnerabilities found</p>'}
                `;
            } catch (e) {
                results.innerHTML = '<p style="color: var(--danger);">Test failed</p>';
            }
        }
        
        async function searchCVE() {
            const q = document.getElementById('cve-input').value;
            const results = document.getElementById('cve-results');
            
            if (!q || q.length < 3) {
                results.innerHTML = '<p style="color: var(--text-muted);">Enter at least 3 characters</p>';
                return;
            }
            
            results.innerHTML = '<div class="loading">Searching...</div>';
            
            try {
                const res = await fetch(`/admin/cve/search?q=${encodeURIComponent(q)}`);
                const data = await res.json();
                
                if (data.results?.length > 0) {
                    results.innerHTML = data.results.map(cve => `
                        <div class="cve-item">
                            <div class="cve-header">
                                <span class="cve-id">${cve.id}</span>
                                <span class="severity severity-${(cve.severity || 'low').toLowerCase()}">${cve.severity || 'N/A'}</span>
                            </div>
                            <p class="cve-desc">${cve.description}</p>
                        </div>
                    `).join('');
                } else {
                    results.innerHTML = '<div class="empty">No results found</div>';
                }
            } catch (e) {
                results.innerHTML = '<p style="color: var(--danger);">Search failed</p>';
            }
        }
        
        async function searchCPE() {
            const q = document.getElementById('cpe-input').value;
            const results = document.getElementById('cpe-results');
            
            if (!q || q.length < 2) {
                results.innerHTML = '<p style="color: var(--text-muted);">Enter at least 2 characters</p>';
                return;
            }
            
            results.innerHTML = '<div class="loading">Searching...</div>';
            
            try {
                const res = await fetch(`/admin/cpe/search?q=${encodeURIComponent(q)}`);
                const data = await res.json();
                
                if (data.results?.length > 0) {
                    results.innerHTML = data.results.map(cpe => `
                        <div class="cve-item">
                            <code style="word-break:break-all;">${cpe.cpe}</code>
                            ${cpe.score ? `<span class="badge badge-info" style="margin-left:8px;">Score: ${cpe.score}</span>` : ''}
                        </div>
                    `).join('');
                } else {
                    results.innerHTML = '<div class="empty">No results found</div>';
                }
            } catch (e) {
                results.innerHTML = '<p style="color: var(--danger);">Search failed</p>';
            }
        }
        
        function refreshData() {
            const active = document.querySelector('.section.active');
            if (active) {
                const id = active.id.replace('-section', '');
                if (id === 'dashboard') loadDashboard();
                if (id === 'agents') loadAgents();
            }
        }
        
        document.querySelectorAll('.search-input').forEach(input => {
            input.addEventListener('keypress', e => {
                if (e.key === 'Enter') {
                    const btn = input.parentElement.querySelector('button');
                    if (btn) btn.click();
                }
            });
        });
        
        loadDashboard();
    </script>
</body>
</html>
    """


@router.get("/dashboard-stats")
async def dashboard_stats():
    """Get dashboard statistics"""
    try:
        cve_result = await db_manager.fetch_one("SELECT count(*) as count FROM cves")
        cve_count = cve_result['count'] if cve_result else 0
        
        agent_count = 0
        software_count = 0
        vuln_count = 0
        
        try:
            async with httpx.AsyncClient(timeout=5.0) as client:
                res = await client.get(f"{API_BASE_URL}/agents/", follow_redirects=True)
                if res.status_code == 200:
                    agents = res.json().get('data', [])
                    agent_count = len(agents)
                    for agent in agents:
                        # Software is in metadata.dependencies
                        metadata = agent.get('metadata', {})
                        deps = metadata.get('dependencies', [])
                        if isinstance(deps, list):
                            software_count += len(deps)
                        vuln_count += metadata.get('total_vulnerabilities', 0)
        except:
            pass
        
        return {
            "agent_count": agent_count,
            "software_count": software_count,
            "vulnerability_count": vuln_count,
            "exploit_count": 0,
            "cve_database_count": cve_count
        }
    except Exception as e:
        logger.error(f"Dashboard stats error: {e}")
        return {"agent_count": 0, "software_count": 0, "vulnerability_count": 0, "exploit_count": 0, "cve_database_count": 0}


@router.get("/agents")
async def list_agents():
    """List all agents"""
    try:
        agents = []
        
        try:
            async with httpx.AsyncClient(timeout=5.0) as client:
                res = await client.get(f"{API_BASE_URL}/agents/", follow_redirects=True)
                if res.status_code == 200:
                    for agent in res.json().get('data', []):
                        # Software is stored in metadata.dependencies
                        metadata = agent.get('metadata', {})
                        deps = metadata.get('dependencies', [])
                        agents.append({
                            "id": agent.get('id', ''),
                            "name": agent.get('name', agent.get('id', '')),
                            "hostname": agent.get('hostname', ''),
                            "os": agent.get('os', ''),
                            "ip": agent.get('ip_address', ''),
                            "software_count": len(deps) if isinstance(deps, list) else 0,
                            "vuln_count": metadata.get('total_vulnerabilities', 0),
                            "status": agent.get('status', 'offline'),
                            "last_seen": agent.get('last_seen', '')
                        })
        except Exception as e:
            logger.warning(f"Failed to fetch agents: {e}")
        
        return {"agents": agents}
    except Exception as e:
        logger.error(f"List agents error: {e}")
        return {"agents": []}


@router.get("/agents/{agent_id}/enrichment")
async def get_agent_enrichment(agent_id: str):
    """Get enrichment data for agent"""
    try:
        agent_data = {"id": agent_id, "name": agent_id, "hostname": "", "os": "", "last_seen": ""}
        software_list = []
        
        try:
            async with httpx.AsyncClient(timeout=5.0) as client:
                res = await client.get(f"{API_BASE_URL}/agents/{agent_id}/", follow_redirects=True)
                if res.status_code == 200:
                    agent = res.json().get('data', {})
                    metadata = agent.get('metadata', {})
                    agent_data = {
                        "id": agent.get('id', agent_id),
                        "name": agent.get('name', agent_id),
                        "hostname": agent.get('hostname', ''),
                        "os": agent.get('os', ''),
                        "last_seen": agent.get('last_seen', '')
                    }
                    
                    # Software is in metadata.dependencies
                    raw_software = metadata.get('dependencies', [])
                    if isinstance(raw_software, list):
                        for sw in raw_software[:100]:
                            software_list.append({
                                "name": sw.get('name', ''),
                                "version": sw.get('version', ''),
                                "vendor": sw.get('vendor', sw.get('description', '')),
                                "cpe": None,
                                "vulnerabilities": [],
                                "exploit_count": 0
                            })
                    
                    vuln_count = metadata.get('total_vulnerabilities', 0)
        except Exception as e:
            logger.warning(f"Failed to fetch agent: {e}")
            vuln_count = 0
        
        return {
            "agent": agent_data,
            "software": software_list,
            "cpe_matches": 0,
            "vulnerability_count": vuln_count,
            "exploit_count": 0
        }
    except Exception as e:
        logger.error(f"Agent enrichment error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/stats")
async def get_stats():
    """Get database stats"""
    try:
        cve_result = await db_manager.fetch_one("SELECT count(*) as count FROM cves")
        cve_count = cve_result['count'] if cve_result else 0
        
        size_result = await db_manager.fetch_one("SELECT pg_size_pretty(pg_total_relation_size('cves')) as size")
        db_size = size_result['size'] if size_result else 'N/A'
        
        return {"cve_count": cve_count, "cpe_count": 167575, "db_size": db_size, "cache_size": 0}
    except Exception as e:
        logger.error(f"Stats error: {e}")
        return {"cve_count": 0, "cpe_count": 0, "db_size": "N/A", "cache_size": 0}


@router.get("/cve/search")
async def search_cve(q: str = Query(..., min_length=3)):
    """Search CVE database"""
    try:
        query = """
            SELECT id, description, 
                   data->'metrics'->'cvssMetricV31'->0->'cvssData'->>'baseSeverity' as severity
            FROM cves
            WHERE id ILIKE $1 OR description ILIKE $1
            LIMIT 20
        """
        results = await db_manager.fetch_all(query, f"%{q}%")
        
        return {
            "count": len(results),
            "results": [
                {
                    "id": r.get('id', ''),
                    "description": (r.get('description', '') or '')[:300] + "...",
                    "severity": r.get('severity') or 'UNKNOWN'
                }
                for r in results
            ]
        }
    except Exception as e:
        logger.error(f"CVE search error: {e}")
        return {"count": 0, "results": []}


@router.get("/cpe/search")
async def search_cpe(q: str = Query(..., min_length=2)):
    """Search CPE dictionary"""
    try:
        import valkey
        
        rdb = valkey.Valkey(host='zerotrace-valkey', port=6379, db=8, decode_responses=True, socket_timeout=5.0)
        results = rdb.zrevrange(f"s:{q.lower()}", 0, 19, withscores=True)
        
        return {
            "count": len(results),
            "results": [{"cpe": cpe, "score": int(score)} for cpe, score in results]
        }
    except Exception as e:
        logger.error(f"CPE search error: {e}")
        return {"count": 0, "results": []}

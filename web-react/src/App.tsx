import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Scans from './pages/Scans';
import Vulnerabilities from './pages/Vulnerabilities';
import Topology from './pages/Topology';
import Agents from './pages/Agents';
import Login from './pages/Login';
import { AuthProvider } from './contexts/AuthContext';
import './styles/terminal-theme.css'; // Terminal-inspired dark theme

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      retry: 1,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Router>
        <AuthProvider>
          <div className="min-h-screen">
            <Routes>
              <Route path="/login" element={<Login />} />
              <Route path="/" element={<Layout />}>
                <Route index element={<Dashboard />} />
                <Route path="agents" element={<Agents />} />
                <Route path="vulnerabilities" element={<Vulnerabilities />} />
                <Route path="issues" element={<Issues />} />
                <Route path="scans" element={<Scans />} />
                <Route path="reports" element={<Reports />} />
                <Route path="topology" element={<Topology />} />
                <Route path="settings" element={<Settings />} />
              </Route>
            </Routes>
            <Toaster position="top-right" />
          </div>
        </AuthProvider>
      </Router>
    </QueryClientProvider>
  );
}

// Placeholder components for new routes
const Issues = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-3xl font-bold text-gold text-glow">SECURITY ISSUES</h1>
      <p className="text-text-secondary">SECURITY ISSUES AND REMEDIATION TRACKING</p>
    </div>
    <div className="card card-terminal">
      <p className="text-text-secondary">ISSUES MANAGEMENT INTERFACE COMING SOON...</p>
    </div>
  </div>
);

const Reports = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-3xl font-bold text-gold text-glow">SECURITY REPORTS</h1>
      <p className="text-text-secondary">SECURITY REPORTS AND ANALYTICS</p>
    </div>
    <div className="card card-terminal">
      <p className="text-text-secondary">REPORTS INTERFACE COMING SOON...</p>
    </div>
  </div>
);

const Settings = () => (
  <div className="space-y-6">
    <div>
      <h1 className="text-3xl font-bold text-gold text-glow">SYSTEM SETTINGS</h1>
      <p className="text-text-secondary">SYSTEM CONFIGURATION AND PREFERENCES</p>
    </div>
    <div className="card card-terminal">
      <p className="text-text-secondary">SETTINGS INTERFACE COMING SOON...</p>
    </div>
  </div>
);

export default App;

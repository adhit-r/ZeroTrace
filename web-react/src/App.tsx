import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import LayoutMinimal from './components/LayoutMinimal';
import DashboardMinimal from './pages/DashboardMinimal';
import './styles/zerotrace-theme.css'; // ZeroTrace neubrutalist theme
import './styles/neobrutal.css'; // Neobrutal design system

// Lazy load other pages to prevent import errors from blocking app startup
const Layout = React.lazy(() => import('./components/Layout'));
const Dashboard = React.lazy(() => import('./pages/Dashboard'));
const Scans = React.lazy(() => import('./pages/Scans'));
const Vulnerabilities = React.lazy(() => import('./pages/Vulnerabilities'));
const Topology = React.lazy(() => import('./pages/Topology'));
const Agents = React.lazy(() => import('./pages/Agents'));
const AssetDetail = React.lazy(() => import('./pages/AssetDetail'));
const SecurityDashboard = React.lazy(() => import('./pages/SecurityDashboard'));
const Settings = React.lazy(() => import('./pages/Settings'));
const Profile = React.lazy(() => import('./pages/Profile'));
const OrganizationProfile = React.lazy(() => import('./pages/OrganizationProfile'));
const AIAnalytics = React.lazy(() => import('./pages/AIAnalytics'));
const RiskHeatmaps = React.lazy(() => import('./pages/RiskHeatmaps'));
const SecurityMaturity = React.lazy(() => import('./pages/SecurityMaturity'));
const Compliance = React.lazy(() => import('./pages/Compliance'));
const TechStack = React.lazy(() => import('./pages/TechStack'));
const ScanProcessing = React.lazy(() => import('./pages/ScanProcessing'));
const VulnerabilityAnalysis = React.lazy(() => import('./pages/VulnerabilityAnalysis'));
const ComplianceReports = React.lazy(() => import('./pages/ComplianceReports'));
const ScannerDetails = React.lazy(() => import('./pages/ScannerDetails'));

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
        <div className="min-h-screen">
          <Suspense fallback={<div className="flex items-center justify-center h-screen">Loading...</div>}>
            <Routes>
              <Route path="/" element={<LayoutMinimal><DashboardMinimal /></LayoutMinimal>} />
              <Route path="/agents" element={<LayoutMinimal><Agents /></LayoutMinimal>} />
              <Route path="/agents/:id" element={<LayoutMinimal><AssetDetail /></LayoutMinimal>} />
              <Route path="/vulnerabilities" element={<LayoutMinimal><Vulnerabilities /></LayoutMinimal>} />
              <Route path="/vulnerabilities/:id/analysis" element={<LayoutMinimal><VulnerabilityAnalysis /></LayoutMinimal>} />
              <Route path="/security" element={<LayoutMinimal><SecurityDashboard /></LayoutMinimal>} />
              <Route path="/scans" element={<LayoutMinimal><Scans /></LayoutMinimal>} />
              <Route path="/scan-processing" element={<LayoutMinimal><ScanProcessing /></LayoutMinimal>} />
              <Route path="/topology" element={<LayoutMinimal><Topology /></LayoutMinimal>} />
              <Route path="/settings" element={<LayoutMinimal><Settings /></LayoutMinimal>} />
              <Route path="/profile" element={<LayoutMinimal><Profile /></LayoutMinimal>} />
              <Route path="/organization-profile" element={<LayoutMinimal><OrganizationProfile /></LayoutMinimal>} />
              <Route path="/ai-analytics" element={<LayoutMinimal><AIAnalytics /></LayoutMinimal>} />
              <Route path="/risk-heatmaps" element={<LayoutMinimal><RiskHeatmaps /></LayoutMinimal>} />
              <Route path="/security-maturity" element={<LayoutMinimal><SecurityMaturity /></LayoutMinimal>} />
              <Route path="/compliance" element={<LayoutMinimal><Compliance /></LayoutMinimal>} />
              <Route path="/compliance/reports" element={<LayoutMinimal><ComplianceReports /></LayoutMinimal>} />
              <Route path="/tech-stack" element={<LayoutMinimal><TechStack /></LayoutMinimal>} />
              <Route path="/scanner-details" element={<LayoutMinimal><ScannerDetails /></LayoutMinimal>} />
            </Routes>
          </Suspense>
        </div>
      </Router>
    </QueryClientProvider>
  );
}


export default App;
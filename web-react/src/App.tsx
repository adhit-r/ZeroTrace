import React, { Suspense } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Toaster } from 'react-hot-toast';
import LayoutMinimal from './components/LayoutMinimal';
import DashboardMinimal from './pages/DashboardMinimal';
import LoadingSpinner from './components/LoadingSpinner';
import { ErrorBoundary } from './components/ErrorBoundary';
import './styles/zerotrace-theme.css'; // ZeroTrace neubrutalist theme

// Lazy load pages with code splitting
// Group related pages for better chunk organization
// Note: Layout and Dashboard are loaded on-demand via routes

// Dashboard & Analytics pages (heavy with charts)
const SecurityDashboard = React.lazy(() => import('./pages/SecurityDashboard'));
const AIAnalytics = React.lazy(() => import('./pages/AIAnalytics'));
const RiskHeatmaps = React.lazy(() => import('./pages/RiskHeatmaps'));
const SecurityMaturity = React.lazy(() => import('./pages/SecurityMaturity'));

// Vulnerability pages
const Vulnerabilities = React.lazy(() => import('./pages/Vulnerabilities'));
const VulnerabilityAnalysis = React.lazy(() => import('./pages/VulnerabilityAnalysis'));

// Agent & Asset pages
const Agents = React.lazy(() => import('./pages/Agents'));
const AssetDetail = React.lazy(() => import('./pages/AssetDetail'));
const ScannerDetails = React.lazy(() => import('./pages/ScannerDetails'));

// Applications page
const Applications = React.lazy(() => import('./pages/Applications'));
const ApplicationAnalysis = React.lazy(() => import('./pages/ApplicationAnalysis'));

// Scan pages
// const Scans = React.lazy(() => import('./pages/Scans'));
const NetworkScanner = React.lazy(() => import('./pages/NetworkScanner'));
const ScanProcessing = React.lazy(() => import('./pages/ScanProcessing'));

// Topology pages (heavy with d3/reactflow)
const Topology = React.lazy(() => import('./pages/Topology'));
const AttackPaths = React.lazy(() => import('./pages/AttackPaths'));

// Compliance pages
const Compliance = React.lazy(() => import('./pages/Compliance'));
const ComplianceReports = React.lazy(() => import('./pages/ComplianceReports'));

// Settings & Profile pages
const Settings = React.lazy(() => import('./pages/Settings'));
const Profile = React.lazy(() => import('./pages/Profile'));
const OrganizationProfile = React.lazy(() => import('./pages/OrganizationProfile'));
const TechStack = React.lazy(() => import('./pages/TechStack'));

// Optimized React Query configuration
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000, // 5 minutes
      gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
      retry: 1,
      refetchOnWindowFocus: false, // Reduce unnecessary refetches
      refetchOnReconnect: true,
    },
    mutations: {
      retry: 1,
    },
  },
});

function App() {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <Router>
          <div className="min-h-screen">
            <Toaster
              position="top-right"
              toastOptions={{
                duration: 4000,
                style: {
                  background: '#fff',
                  color: '#000',
                  border: '3px solid #000',
                  borderRadius: '0',
                  fontFamily: 'inherit',
                  fontWeight: 'bold',
                },
                success: {
                  iconTheme: {
                    primary: '#10b981',
                    secondary: '#fff',
                  },
                },
                error: {
                  iconTheme: {
                    primary: '#ef4444',
                    secondary: '#fff',
                  },
                },
              }}
            />
            <ErrorBoundary>
              <Suspense fallback={<LoadingSpinner />}>
                <Routes>
                  <Route path="/" element={<LayoutMinimal><DashboardMinimal /></LayoutMinimal>} />
                  <Route path="/agents" element={<LayoutMinimal><Agents /></LayoutMinimal>} />
                  <Route path="/agents/:id" element={<LayoutMinimal><AssetDetail /></LayoutMinimal>} />
                  <Route path="/vulnerabilities" element={<LayoutMinimal><Vulnerabilities /></LayoutMinimal>} />
                  <Route path="/vulnerabilities/:id/analysis" element={<LayoutMinimal><VulnerabilityAnalysis /></LayoutMinimal>} />
                  <Route path="/apps" element={<LayoutMinimal><Applications /></LayoutMinimal>} />
                  <Route path="/applications" element={<LayoutMinimal><Applications /></LayoutMinimal>} />
                  <Route path="/application-analysis" element={<LayoutMinimal><ApplicationAnalysis /></LayoutMinimal>} />
                  <Route path="/security" element={<LayoutMinimal><SecurityDashboard /></LayoutMinimal>} />
                  <Route path="/scans" element={<LayoutMinimal><NetworkScanner /></LayoutMinimal>} />
                  <Route path="/network-scanner" element={<LayoutMinimal><NetworkScanner /></LayoutMinimal>} />
                  <Route path="/scan-processing" element={<LayoutMinimal><ScanProcessing /></LayoutMinimal>} />
                  <Route path="/topology" element={<LayoutMinimal><Topology /></LayoutMinimal>} />
                  <Route path="/network-topology" element={<LayoutMinimal><Topology /></LayoutMinimal>} />
                  <Route path="/attack-paths" element={<LayoutMinimal><AttackPaths /></LayoutMinimal>} />
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
            </ErrorBoundary>
          </div>
        </Router>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}


export default App;
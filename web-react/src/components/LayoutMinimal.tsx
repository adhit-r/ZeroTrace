import React from 'react';
import { Link, useLocation } from 'react-router-dom';

interface LayoutProps {
  children: React.ReactNode;
}

export default function LayoutMinimal({ children }: LayoutProps) {
  const location = useLocation();

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  const getLinkClassName = (path: string) => {
    const baseClasses = "block py-2 px-3 rounded font-bold mb-2 transition-all border-3 border-black shadow-neubrutalist-sm";
    if (isActive(path)) {
      return `${baseClasses} bg-orange-500 text-white shadow-neubrutalist-md`;
    }
    return `${baseClasses} bg-white text-black hover:bg-gray-100 hover:shadow-neubrutalist-md hover:translate-x-0.5 hover:translate-y-0.5`;
  };

  return (
    <div className="min-h-screen bg-gray-100 flex">
      {/* Sidebar */}
      <div className="w-64 bg-white border-r-3 border-black">
        <div className="p-4 border-b-3 border-black">
          <h1 className="text-xl font-black text-orange-500 uppercase">ZeroTrace</h1>
        </div>
        <nav className="p-4 space-y-1">
          {/* Main */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Main</p>
            <Link to="/" className={getLinkClassName('/')}>Dashboard</Link>
            <Link to="/agents" className={getLinkClassName('/agents')}>Agents</Link>
            <Link to="/vulnerabilities" className={getLinkClassName('/vulnerabilities')}>Vulnerabilities</Link>
            <Link to="/apps" className={getLinkClassName('/apps')}>Applications</Link>
          </div>

          {/* Scanning */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Scanning</p>
            <Link to="/scans" className={getLinkClassName('/scans')}>Network Scanner</Link>
            <Link to="/topology" className={getLinkClassName('/topology')}>Network Topology</Link>
            <Link to="/scanner-details" className={getLinkClassName('/scanner-details')}>Scanner Details</Link>
            <Link to="/scan-processing" className={getLinkClassName('/scan-processing')}>Scan Processing</Link>
          </div>

          {/* Security & Analytics */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Security</p>
            <Link to="/security" className={getLinkClassName('/security')}>Security Dashboard</Link>
            <Link to="/ai-analytics" className={getLinkClassName('/ai-analytics')}>AI Analytics</Link>
            <Link to="/risk-heatmaps" className={getLinkClassName('/risk-heatmaps')}>Risk Heatmaps</Link>
            <Link to="/security-maturity" className={getLinkClassName('/security-maturity')}>Security Maturity</Link>
          </div>

          {/* Compliance */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Compliance</p>
            <Link to="/compliance" className={getLinkClassName('/compliance')}>Compliance</Link>
            <Link to="/compliance/reports" className={getLinkClassName('/compliance/reports')}>Compliance Reports</Link>
          </div>

          {/* Settings */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Settings</p>
            <Link to="/settings" className={getLinkClassName('/settings')}>Settings</Link>
            <Link to="/profile" className={getLinkClassName('/profile')}>Profile</Link>
            <Link to="/organization-profile" className={getLinkClassName('/organization-profile')}>Organization</Link>
            <Link to="/tech-stack" className={getLinkClassName('/tech-stack')}>Tech Stack</Link>
          </div>
        </nav>
      </div>

      {/* Main content */}
      <div className="flex-1 flex flex-col">
        <div className="bg-white border-b-3 border-black p-4">
          <p className="font-bold">Enterprise Security Management</p>
        </div>
        <div className="flex-1 overflow-auto p-8">
          {children}
        </div>
      </div>
    </div>
  );
}

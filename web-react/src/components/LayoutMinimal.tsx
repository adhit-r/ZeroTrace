import React from 'react';

interface LayoutProps {
  children: React.ReactNode;
}

export default function LayoutMinimal({ children }: LayoutProps) {
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
            <a href="/" className="block py-2 px-3 rounded bg-orange-500 text-white font-bold mb-2">Dashboard</a>
            <a href="/agents" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Agents</a>
            <a href="/vulnerabilities" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Vulnerabilities</a>
            <a href="/apps" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Applications</a>
          </div>

          {/* Scanning */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Scanning</p>
            <a href="/scans" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Network Scanner</a>
            <a href="/topology" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Network Topology</a>
            <a href="/scanner-details" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Scanner Details</a>
            <a href="/scan-processing" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Scan Processing</a>
          </div>

          {/* Security & Analytics */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Security</p>
            <a href="/security" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Security Dashboard</a>
            <a href="/ai-analytics" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">AI Analytics</a>
            <a href="/risk-heatmaps" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Risk Heatmaps</a>
            <a href="/security-maturity" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Security Maturity</a>
          </div>

          {/* Compliance */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Compliance</p>
            <a href="/compliance" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Compliance</a>
            <a href="/compliance/reports" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Compliance Reports</a>
          </div>

          {/* Settings */}
          <div className="mb-4">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2 px-2">Settings</p>
            <a href="/settings" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Settings</a>
            <a href="/profile" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Profile</a>
            <a href="/organization-profile" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Organization</a>
            <a href="/tech-stack" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Tech Stack</a>
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

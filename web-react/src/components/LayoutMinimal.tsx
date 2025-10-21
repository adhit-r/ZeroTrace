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
        <nav className="p-4">
          <a href="/" className="block py-2 px-3 rounded bg-orange-500 text-white font-bold mb-2">Dashboard</a>
          <a href="/agents" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Agents</a>
          <a href="/vulnerabilities" className="block py-2 px-3 rounded hover:bg-gray-100 font-bold mb-2">Vulnerabilities</a>
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

import React from 'react';

export default function Dashboard() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-4">Dashboard</h1>
      <p className="text-gray-600">Welcome to ZeroTrace vulnerability detection platform.</p>
      
      <div className="grid grid-cols-4 gap-4 mt-8">
        <div className="bg-white p-6 rounded border-3 border-black">
          <div className="text-lg font-bold">0</div>
          <div className="text-sm text-gray-600">Total Assets</div>
        </div>
        <div className="bg-white p-6 rounded border-3 border-black">
          <div className="text-lg font-bold">0</div>
          <div className="text-sm text-gray-600">Vulnerabilities</div>
        </div>
        <div className="bg-white p-6 rounded border-3 border-black">
          <div className="text-lg font-bold">0</div>
          <div className="text-sm text-gray-600">Critical Issues</div>
        </div>
        <div className="bg-white p-6 rounded border-3 border-black">
          <div className="text-lg font-bold">-</div>
          <div className="text-sm text-gray-600">Compliance</div>
        </div>
      </div>
    </div>
  );
}

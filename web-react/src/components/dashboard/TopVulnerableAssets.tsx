import React from 'react';
import { Eye, ArrowRight } from 'lucide-react';

interface Asset {
  name: string;
  vulnerabilities: number;
  critical: number;
}

interface TopVulnerableAssetsProps {
  assets: Asset[];
}

const TopVulnerableAssets: React.FC<TopVulnerableAssetsProps> = ({ assets }) => {
  if (!assets || assets.length === 0) {
    return (
      <div>
        <h2 className="text-xl font-black text-black uppercase mb-4">Top Vulnerable Assets</h2>
        <p className="text-gray-500">No vulnerable assets found.</p>
      </div>
    );
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-xl font-black text-black uppercase">Top Vulnerable Assets</h2>
        <button className="flex items-center gap-2 px-4 py-2 bg-gray-100 text-black font-bold uppercase tracking-wide rounded-lg border-3 border-black shadow-neo-brutal-small hover:shadow-neo-brutal-small-hover transition-all">
          View All <ArrowRight className="h-4 w-4" />
        </button>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b-3 border-black">
              <th className="text-left p-3 font-black uppercase">Asset Name</th>
              <th className="text-left p-3 font-black uppercase">Total Vulns</th>
              <th className="text-left p-3 font-black uppercase">Critical</th>
              <th className="text-left p-3 font-black uppercase">Actions</th>
            </tr>
          </thead>
          <tbody>
            {assets.map((asset, index) => (
              <tr key={index} className="border-b-2 border-gray-200">
                <td className="p-3 font-bold">{asset.name}</td>
                <td className="p-3">
                  <span className="px-2 py-1 bg-yellow-400 text-black font-bold rounded-md border-2 border-black">
                    {asset.vulnerabilities}
                  </span>
                </td>
                <td className="p-3">
                  <span className="px-2 py-1 bg-red-500 text-white font-bold rounded-md border-2 border-black">
                    {asset.critical}
                  </span>
                </td>
                <td className="p-3">
                  <button className="p-2 bg-white rounded-lg border-3 border-black shadow-neo-brutal-small">
                    <Eye className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default TopVulnerableAssets;


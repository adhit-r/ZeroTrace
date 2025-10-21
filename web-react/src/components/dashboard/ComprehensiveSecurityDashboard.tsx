import React from 'react';
import { 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  Activity, 
  BarChart3, 
  TrendingUp, 
  TrendingDown,
  Target,
  Zap
} from 'lucide-react';

const ComprehensiveSecurityDashboard: React.FC = () => {
  return (
    <div className="space-y-6">
      {/* Security Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Overall Security Score */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-green-100 rounded border-2 border-black">
              <Shield className="h-6 w-6 text-green-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-bold text-green-600">85%</div>
              <div className="text-sm text-gray-600 uppercase tracking-wider">Security Score</div>
            </div>
          </div>
          <div className="text-sm text-green-600 font-bold">
            <TrendingUp className="h-4 w-4 inline-block mr-1" />
            +5% from last week
          </div>
        </div>

        {/* Critical Vulnerabilities */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-red-100 rounded border-2 border-black">
              <AlertTriangle className="h-6 w-6 text-red-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-bold text-red-600">12</div>
              <div className="text-sm text-gray-600 uppercase tracking-wider">Critical Issues</div>
            </div>
          </div>
          <div className="text-sm text-red-600 font-bold">
            <TrendingDown className="h-4 w-4 inline-block mr-1" />
            -3 from last week
          </div>
        </div>

        {/* Compliance Status */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-blue-100 rounded border-2 border-black">
              <CheckCircle className="h-6 w-6 text-blue-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-bold text-blue-600">92%</div>
              <div className="text-sm text-gray-600 uppercase tracking-wider">Compliance</div>
            </div>
          </div>
          <div className="text-sm text-blue-600 font-bold">
            <TrendingUp className="h-4 w-4 inline-block mr-1" />
            +2% from last month
          </div>
        </div>

        {/* Active Threats */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
          <div className="flex items-center justify-between mb-4">
            <div className="p-3 bg-orange-100 rounded border-2 border-black">
              <Target className="h-6 w-6 text-orange-600" />
            </div>
            <div className="text-right">
              <div className="text-3xl font-bold text-orange-600">3</div>
              <div className="text-sm text-gray-600 uppercase tracking-wider">Active Threats</div>
            </div>
          </div>
          <div className="text-sm text-orange-600 font-bold">
            <Activity className="h-4 w-4 inline-block mr-1" />
            Real-time monitoring
          </div>
        </div>
      </div>

      {/* Security Metrics Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Vulnerability Trends */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Vulnerability Trends</h3>
          <div className="h-64 flex items-center justify-center bg-gray-50 border-2 border-gray-300 rounded">
            <div className="text-center">
              <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-600 font-bold">Vulnerability Chart</p>
              <p className="text-sm text-gray-500">Real-time vulnerability tracking</p>
            </div>
          </div>
        </div>

        {/* Security Categories */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Security Categories</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-3 bg-gray-50 border-2 border-gray-300 rounded">
              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-blue-600" />
                <span className="font-bold">Network Security</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-sm text-green-600 font-bold">Good</span>
                <CheckCircle className="h-4 w-4 text-green-600" />
              </div>
            </div>
            <div className="flex items-center justify-between p-3 bg-gray-50 border-2 border-gray-300 rounded">
              <div className="flex items-center gap-3">
                <AlertTriangle className="h-5 w-5 text-red-600" />
                <span className="font-bold">System Vulnerabilities</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-sm text-red-600 font-bold">Critical</span>
                <AlertTriangle className="h-4 w-4 text-red-600" />
              </div>
            </div>
            <div className="flex items-center justify-between p-3 bg-gray-50 border-2 border-gray-300 rounded">
              <div className="flex items-center gap-3">
                <Zap className="h-5 w-5 text-orange-600" />
                <span className="font-bold">API Security</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-sm text-orange-600 font-bold">Warning</span>
                <Activity className="h-4 w-4 text-orange-600" />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Security Events */}
      <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <h3 className="text-xl font-black text-black uppercase mb-4">Recent Security Events</h3>
        <div className="space-y-3">
          <div className="flex items-center gap-3 p-3 bg-red-50 border-2 border-red-300 rounded">
            <AlertTriangle className="h-5 w-5 text-red-600" />
            <div className="flex-1">
              <p className="font-bold text-red-800">Critical vulnerability detected</p>
              <p className="text-sm text-red-600">CVE-2024-1234 in Apache HTTP Server</p>
            </div>
            <span className="text-xs text-red-600 font-bold">2 min ago</span>
          </div>
          <div className="flex items-center gap-3 p-3 bg-orange-50 border-2 border-orange-300 rounded">
            <Activity className="h-5 w-5 text-orange-600" />
            <div className="flex-1">
              <p className="font-bold text-orange-800">Suspicious network activity</p>
              <p className="text-sm text-orange-600">Multiple failed login attempts detected</p>
            </div>
            <span className="text-xs text-orange-600 font-bold">15 min ago</span>
          </div>
          <div className="flex items-center gap-3 p-3 bg-green-50 border-2 border-green-300 rounded">
            <CheckCircle className="h-5 w-5 text-green-600" />
            <div className="flex-1">
              <p className="font-bold text-green-800">Security scan completed</p>
              <p className="text-sm text-green-600">All systems scanned successfully</p>
            </div>
            <span className="text-xs text-green-600 font-bold">1 hour ago</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ComprehensiveSecurityDashboard;

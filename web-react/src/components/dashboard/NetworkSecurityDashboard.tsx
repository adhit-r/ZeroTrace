import React from 'react';
import { 
  Network, 
  Shield, 
  Activity, 
  AlertTriangle, 
  CheckCircle, 
  Eye, 
  Zap,
  Lock,
  Wifi,
  Server,
  BarChart3
} from 'lucide-react';

const NetworkSecurityDashboard: React.FC = () => {
  const networkMetrics = [
    {
      name: 'Firewall Status',
      status: 'active',
      value: '100%',
      color: 'green',
      icon: Shield
    },
    {
      name: 'VPN Connections',
      status: 'active',
      value: '3/5',
      color: 'green',
      icon: Lock
    },
    {
      name: 'Network Traffic',
      status: 'normal',
      value: '2.4 GB/s',
      color: 'blue',
      icon: Activity
    },
    {
      name: 'Blocked Threats',
      status: 'active',
      value: '127',
      color: 'red',
      icon: AlertTriangle
    }
  ];

  const securityEvents = [
    {
      time: '2 min ago',
      type: 'blocked',
      source: import.meta.env.VITE_EXAMPLE_IP_1 || '192.168.1.100',
      description: 'Suspicious port scan detected',
      severity: 'high',
      color: 'red'
    },
    {
      time: '5 min ago',
      type: 'allowed',
      source: '10.0.0.50',
      description: 'VPN connection established',
      severity: 'info',
      color: 'green'
    },
    {
      time: '12 min ago',
      type: 'blocked',
      source: '203.0.113.45',
      description: 'Malicious payload detected',
      severity: 'critical',
      color: 'red'
    },
    {
      time: '18 min ago',
      type: 'monitored',
      source: import.meta.env.VITE_EXAMPLE_IP_2 || '192.168.1.200',
      description: 'High bandwidth usage',
      severity: 'medium',
      color: 'orange'
    }
  ];

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'text-red-600 bg-red-100 border-red-300';
      case 'high':
        return 'text-red-600 bg-red-100 border-red-300';
      case 'medium':
        return 'text-orange-600 bg-orange-100 border-orange-300';
      case 'low':
        return 'text-yellow-600 bg-yellow-100 border-yellow-300';
      default:
        return 'text-green-600 bg-green-100 border-green-300';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'text-green-600 bg-green-100 border-green-300';
      case 'normal':
        return 'text-blue-600 bg-blue-100 border-blue-300';
      case 'warning':
        return 'text-orange-600 bg-orange-100 border-orange-300';
      default:
        return 'text-gray-600 bg-gray-100 border-gray-300';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <h2 className="text-2xl font-black text-black uppercase mb-2">Network Security</h2>
        <p className="text-gray-600 font-bold">Real-time network monitoring and threat detection</p>
      </div>

      {/* Network Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {networkMetrics.map((metric, index) => {
          const IconComponent = metric.icon;
          return (
            <div key={index} className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal hover:shadow-neo-brutal-hover transition-all">
              <div className="flex items-center justify-between mb-4">
                <div className="p-3 bg-gray-100 rounded border-2 border-black">
                  <IconComponent className="h-6 w-6 text-gray-600" />
                </div>
                <div className="text-right">
                  <div className={`text-2xl font-bold ${
                    metric.color === 'green' ? 'text-green-600' : 
                    metric.color === 'red' ? 'text-red-600' : 
                    'text-blue-600'
                  }`}>
                    {metric.value}
                  </div>
                  <div className="text-sm text-gray-600 uppercase tracking-wider">Current</div>
                </div>
              </div>
              <h3 className="text-lg font-black text-black uppercase mb-2">{metric.name}</h3>
              <div className={`px-3 py-1 rounded border-2 font-bold uppercase text-xs ${getStatusColor(metric.status)}`}>
                {metric.status}
              </div>
            </div>
          );
        })}
      </div>

      {/* Network Overview */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Network Topology */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Network Topology</h3>
          <div className="h-64 flex items-center justify-center bg-gray-50 border-2 border-gray-300 rounded">
            <div className="text-center">
              <Network className="h-12 w-12 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-600 font-bold">Network Diagram</p>
              <p className="text-sm text-gray-500">Real-time network visualization</p>
            </div>
          </div>
        </div>

        {/* Security Status */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Security Status</h3>
          <div className="space-y-4">
            <div className="flex items-center justify-between p-3 bg-green-50 border-2 border-green-300 rounded">
              <div className="flex items-center gap-3">
                <CheckCircle className="h-5 w-5 text-green-600" />
                <span className="font-bold text-green-800">Firewall Active</span>
              </div>
              <span className="text-sm text-green-600 font-bold">100%</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-green-50 border-2 border-green-300 rounded">
              <div className="flex items-center gap-3">
                <Lock className="h-5 w-5 text-green-600" />
                <span className="font-bold text-green-800">VPN Secure</span>
              </div>
              <span className="text-sm text-green-600 font-bold">3/5</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-orange-50 border-2 border-orange-300 rounded">
              <div className="flex items-center gap-3">
                <AlertTriangle className="h-5 w-5 text-orange-600" />
                <span className="font-bold text-orange-800">Traffic Monitoring</span>
              </div>
              <span className="text-sm text-orange-600 font-bold">Active</span>
            </div>
            <div className="flex items-center justify-between p-3 bg-red-50 border-2 border-red-300 rounded">
              <div className="flex items-center gap-3">
                <Shield className="h-5 w-5 text-red-600" />
                <span className="font-bold text-red-800">Threats Blocked</span>
              </div>
              <span className="text-sm text-red-600 font-bold">127</span>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Security Events */}
      <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
        <h3 className="text-xl font-black text-black uppercase mb-4">Recent Security Events</h3>
        <div className="space-y-3">
          {securityEvents.map((event, index) => (
            <div key={index} className={`flex items-center gap-3 p-3 border-2 rounded ${
              event.color === 'red' ? 'bg-red-50 border-red-300' :
              event.color === 'green' ? 'bg-green-50 border-green-300' :
              'bg-orange-50 border-orange-300'
            }`}>
              <div className={`p-2 rounded border-2 ${
                event.color === 'red' ? 'bg-red-100 border-red-300' :
                event.color === 'green' ? 'bg-green-100 border-green-300' :
                'bg-orange-100 border-orange-300'
              }`}>
                {event.type === 'blocked' ? (
                  <AlertTriangle className="h-4 w-4 text-red-600" />
                ) : event.type === 'allowed' ? (
                  <CheckCircle className="h-4 w-4 text-green-600" />
                ) : (
                  <Eye className="h-4 w-4 text-orange-600" />
                )}
              </div>
              <div className="flex-1">
                <p className="font-bold text-gray-800">{event.description}</p>
                <p className="text-sm text-gray-600">Source: {event.source}</p>
              </div>
              <div className="text-right">
                <div className={`px-2 py-1 rounded border font-bold text-xs ${getSeverityColor(event.severity)}`}>
                  {event.severity}
                </div>
                <p className="text-xs text-gray-500 mt-1">{event.time}</p>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Network Performance */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Traffic Analysis */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Traffic Analysis</h3>
          <div className="h-48 flex items-center justify-center bg-gray-50 border-2 border-gray-300 rounded">
            <div className="text-center">
              <BarChart3 className="h-10 w-10 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-600 font-bold">Traffic Chart</p>
              <p className="text-sm text-gray-500">Network traffic visualization</p>
            </div>
          </div>
        </div>

        {/* Security Recommendations */}
        <div className="p-6 bg-white border-3 border-black rounded-lg shadow-neo-brutal">
          <h3 className="text-xl font-black text-black uppercase mb-4">Security Recommendations</h3>
          <div className="space-y-3">
            <div className="p-3 bg-blue-50 border-2 border-blue-300 rounded">
              <div className="flex items-start gap-3">
                <Zap className="h-5 w-5 text-blue-600 mt-0.5" />
                <div>
                  <p className="font-bold text-blue-800">Update Firewall Rules</p>
                  <p className="text-sm text-blue-600">Review and update firewall configurations</p>
                </div>
              </div>
            </div>
            <div className="p-3 bg-orange-50 border-2 border-orange-300 rounded">
              <div className="flex items-start gap-3">
                <AlertTriangle className="h-5 w-5 text-orange-600 mt-0.5" />
                <div>
                  <p className="font-bold text-orange-800">Monitor VPN Usage</p>
                  <p className="text-sm text-orange-600">Check for unusual VPN connection patterns</p>
                </div>
              </div>
            </div>
            <div className="p-3 bg-green-50 border-2 border-green-300 rounded">
              <div className="flex items-start gap-3">
                <CheckCircle className="h-5 w-5 text-green-600 mt-0.5" />
                <div>
                  <p className="font-bold text-green-800">Network Segmentation</p>
                  <p className="text-sm text-green-600">Current segmentation is optimal</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NetworkSecurityDashboard;

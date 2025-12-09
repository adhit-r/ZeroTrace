import React from 'react';
import { X, Package, AlertTriangle, Shield, Server } from 'lucide-react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import type { Application } from '@/services/applicationService';

interface AppDetailModalProps {
  application: Application | null;
  onClose: () => void;
}

const AppDetailModal: React.FC<AppDetailModalProps> = ({ application, onClose }) => {
  if (!application) return null;

  const getRiskBadgeColor = (risk: string) => {
    switch (risk) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-300';
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low':
        return 'bg-blue-100 text-blue-800 border-blue-300';
      default:
        return 'bg-green-100 text-green-800 border-green-300';
    }
  };

  const getSeverityBadgeColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'critical':
        return 'bg-red-100 text-red-800 border-red-300';
      case 'high':
        return 'bg-orange-100 text-orange-800 border-orange-300';
      case 'medium':
        return 'bg-yellow-100 text-yellow-800 border-yellow-300';
      case 'low':
        return 'bg-blue-100 text-blue-800 border-blue-300';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-300';
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
      <Card className="bg-white border-4 border-black rounded-lg max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="p-6 border-b-4 border-black flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Package className="h-6 w-6 text-orange-600" />
            <div>
              <h2 className="text-2xl font-black text-black uppercase">{application.name}</h2>
              <p className="text-sm text-gray-600">{application.vendor}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          >
            <X className="h-5 w-5 text-black" />
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Basic Info */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-xs text-gray-600 uppercase font-bold mb-1">Version</p>
              <p className="text-lg font-bold text-black">{application.version}</p>
            </div>
            <div>
              <p className="text-xs text-gray-600 uppercase font-bold mb-1">Risk Level</p>
              <Badge className={getRiskBadgeColor(application.riskLevel)}>
                {application.riskLevel.toUpperCase()}
              </Badge>
            </div>
            <div>
              <p className="text-xs text-gray-600 uppercase font-bold mb-1">Status</p>
              <Badge
                className={
                  application.status === 'vulnerable'
                    ? 'bg-red-100 text-red-800 border-red-300'
                    : 'bg-green-100 text-green-800 border-green-300'
                }
              >
                {application.status.toUpperCase()}
              </Badge>
            </div>
            <div>
              <p className="text-xs text-gray-600 uppercase font-bold mb-1">Vulnerabilities</p>
              <p className="text-lg font-black text-red-600">{application.vulnerabilities}</p>
            </div>
          </div>

          {/* Agent Info */}
          <div className="flex items-center gap-2 p-4 bg-gray-50 border-2 border-gray-300 rounded-lg">
            <Server className="h-5 w-5 text-blue-600" />
            <div>
              <p className="text-xs text-gray-600 uppercase font-bold">Installed on</p>
              <p className="font-bold text-black">{application.agentName}</p>
            </div>
          </div>

          {/* Classification & Path */}
          {(application.classification || application.path) && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {application.classification && (
                <div>
                  <p className="text-xs text-gray-600 uppercase font-bold mb-1">Classification</p>
                  <p className="font-bold text-black uppercase">{application.classification}</p>
                </div>
              )}
              {application.path && (
                <div>
                  <p className="text-xs text-gray-600 uppercase font-bold mb-1">Installation Path</p>
                  <p className="font-bold text-black text-sm break-all">{application.path}</p>
                </div>
              )}
            </div>
          )}

          {/* Vulnerabilities */}
          {application.vulnerabilityDetails && application.vulnerabilityDetails.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-4">
                <AlertTriangle className="h-5 w-5 text-red-600" />
                <h3 className="text-xl font-black text-black uppercase">
                  Vulnerabilities ({application.vulnerabilityDetails.length})
                </h3>
              </div>
              <div className="space-y-3">
                {application.vulnerabilityDetails.map((vuln, idx) => (
                  <Card
                    key={idx}
                    className="p-4 bg-red-50 border-2 border-red-300 rounded-lg"
                  >
                    <div className="flex items-start justify-between mb-2">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <h4 className="font-bold text-black">{vuln.title}</h4>
                          {vuln.cve && (
                            <Badge className="bg-gray-100 text-gray-800 border-gray-300">
                              {vuln.cve}
                            </Badge>
                          )}
                          <Badge className={getSeverityBadgeColor(vuln.severity)}>
                            {vuln.severity.toUpperCase()}
                          </Badge>
                        </div>
                        {vuln.description && (
                          <p className="text-sm text-gray-700 mt-2">{vuln.description}</p>
                        )}
                      </div>
                      {vuln.score && (
                        <div className="text-right">
                          <p className="text-lg font-black text-red-600">{vuln.score.toFixed(1)}</p>
                          <p className="text-xs text-gray-600">CVSS</p>
                        </div>
                      )}
                    </div>
                  </Card>
                ))}
              </div>
            </div>
          )}

          {/* No Vulnerabilities */}
          {application.vulnerabilities === 0 && (
            <div className="p-6 bg-green-50 border-2 border-green-300 rounded-lg text-center">
              <Shield className="h-12 w-12 text-green-600 mx-auto mb-2" />
              <p className="font-bold text-green-800">No Known Vulnerabilities</p>
              <p className="text-sm text-green-700 mt-1">This application appears to be secure.</p>
            </div>
          )}
        </div>
      </Card>
    </div>
  );
};

export default AppDetailModal;


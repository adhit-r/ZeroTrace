import React, { useEffect, useState } from 'react';
import { complianceService, ComplianceScore, ComplianceFinding, EvidenceItem } from '@/services/complianceService';

const FRAMEWORKS = ['SOC2', 'ISO27001', 'PCI DSS', 'HIPAA', 'SOX'];

export default function Compliance() {
  const [organizationId] = useState<string>('123e4567-e89b-12d3-a456-426614174000'); // TODO: fetch from profile/auth
  const [framework, setFramework] = useState<string>('SOC2');
  const [score, setScore] = useState<ComplianceScore | null>(null);
  const [findings, setFindings] = useState<ComplianceFinding[]>([]);
  const [evidence, setEvidence] = useState<EvidenceItem[]>([]);
  const [recommendations, setRecommendations] = useState<string[]>([]);
  const [executive, setExecutive] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const [sc, fi, ev, rec, ex] = await Promise.all([
          complianceService.getScore(organizationId),
          complianceService.getFindings(organizationId),
          complianceService.getEvidence(organizationId),
          complianceService.getRecommendations(organizationId),
          complianceService.getExecutiveSummary(organizationId)
        ]);
        setScore(sc);
        setFindings(fi);
        setEvidence(ev);
        setRecommendations(rec);
        setExecutive(ex);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [organizationId, framework]);

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Compliance</h1>
          <div className="flex items-center gap-2">
            <span className="text-sm font-bold">Framework:</span>
            <select value={framework} onChange={(e) => setFramework(e.target.value)} className="neobrutal-form-select">
              {FRAMEWORKS.map(f => <option key={f} value={f}>{f}</option>)}
            </select>
          </div>
        </div>

        {/* Score */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Score</h2>
          {score ? (
            <div className="flex items-center gap-6">
              <div className="text-5xl font-extrabold">{score.score.toFixed(1)}</div>
              <span className="neobrutal-badge">Level: {score.level}</span>
              <span className="neobrutal-badge">Framework: {score.framework}</span>
            </div>
          ) : <div className="text-gray-500">No score</div>}
        </div>

        {/* Findings */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Findings</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {findings.map(f => (
              <div key={f.id} className="neobrutal-section">
                <div className="text-sm font-bold">{f.title}</div>
                <div className="text-xs">Control: {f.controlId}</div>
                <div className="text-xs">Severity: {f.severity} • Status: {f.status}</div>
              </div>
            ))}
            {findings.length === 0 && <div className="text-gray-500">No findings</div>}
          </div>
        </div>

        {/* Evidence */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Evidence</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {evidence.map(e => (
              <div key={e.id} className="neobrutal-section">
                <div className="text-sm font-bold">{e.title}</div>
                <div className="text-xs">Status: {e.status}</div>
              </div>
            ))}
            {evidence.length === 0 && <div className="text-gray-500">No evidence</div>}
          </div>
        </div>

        {/* Recommendations */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Recommendations</h2>
          <ul className="list-disc pl-6">
            {recommendations.map((r, i) => (
              <li key={i} className="neobrutal-text">{r}</li>
            ))}
            {recommendations.length === 0 && <div className="text-gray-500">No recommendations</div>}
          </ul>
        </div>

        {/* Executive Summary */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Executive Summary</h2>
          {executive ? (
            <pre className="text-xs whitespace-pre-wrap">{JSON.stringify(executive, null, 2)}</pre>
          ) : <div className="text-gray-500">No executive summary</div>}
        </div>
      </div>
    </div>
  );
}

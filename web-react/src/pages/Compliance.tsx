import { useEffect, useState } from 'react';
import { complianceService } from '@/services/complianceService';
import type { ComplianceScore, ComplianceFinding, EvidenceItem } from '@/services/complianceService';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';

const FRAMEWORKS = ['SOC2', 'ISO27001', 'PCI DSS', 'HIPAA', 'SOX'];

export default function Compliance() {
  // Get organization ID from environment or fetch from API/auth context
  const [organizationId] = useState<string | null>(
    import.meta.env.VITE_ORGANIZATION_ID || null
  );
  const [framework, setFramework] = useState<string>('SOC2');
  const [score, setScore] = useState<ComplianceScore | null>(null);
  const [findings, setFindings] = useState<ComplianceFinding[]>([]);
  const [evidence, setEvidence] = useState<EvidenceItem[]>([]);
  const [recommendations, setRecommendations] = useState<string[]>([]);
  const [executive, setExecutive] = useState<any>(null);

  useEffect(() => {
    const load = async () => {
      if (!organizationId) return;
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
        // loaded
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
            <Select value={framework} onValueChange={setFramework}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select framework" />
              </SelectTrigger>
              <SelectContent>
                {FRAMEWORKS.map(f => (
                  <SelectItem key={f} value={f}>{f}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Score */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Score</CardTitle>
          </CardHeader>
          <CardContent>
            {score ? (
              <div className="flex items-center gap-6">
                <div className="text-5xl font-extrabold">{score.score.toFixed(1)}</div>
                <Badge>Level: {score.level}</Badge>
                <Badge variant="secondary">Framework: {score.framework}</Badge>
              </div>
            ) : <div className="text-gray-500">No score</div>}
          </CardContent>
        </Card>

        {/* Findings */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Findings</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {findings.map(f => (
                <div key={f.id} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                  <div className="text-sm font-bold">{f.title}</div>
                  <div className="text-xs">Control: {f.controlId}</div>
                  <div className="text-xs">Severity: {f.severity} • Status: {f.status}</div>
                </div>
              ))}
              {findings.length === 0 && <div className="text-gray-500">No findings</div>}
            </div>
          </CardContent>
        </Card>

        {/* Evidence */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Evidence</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {evidence.map(e => (
                <div key={e.id} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                  <div className="text-sm font-bold">{e.title}</div>
                  <div className="text-xs">Status: {e.status}</div>
                </div>
              ))}
              {evidence.length === 0 && <div className="text-gray-500">No evidence</div>}
            </div>
          </CardContent>
        </Card>

        {/* Recommendations */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Recommendations</CardTitle>
          </CardHeader>
          <CardContent>
            <ul className="list-disc pl-6">
              {recommendations.map((r, i) => (
                <li key={i} className="font-semibold">{r}</li>
              ))}
              {recommendations.length === 0 && <div className="text-gray-500">No recommendations</div>}
            </ul>
          </CardContent>
        </Card>

        {/* Executive Summary */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Executive Summary</CardTitle>
          </CardHeader>
          <CardContent>
            {executive ? (
              <pre className="text-xs whitespace-pre-wrap">{JSON.stringify(executive, null, 2)}</pre>
            ) : <div className="text-gray-500">No executive summary</div>}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

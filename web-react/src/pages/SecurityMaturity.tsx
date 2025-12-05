import React, { useEffect, useState } from 'react';
import { maturityService, MaturityScore, DimensionScore, ImprovementItem, PeerComparison } from '@/services/maturityService';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

export default function SecurityMaturity() {
  // Get organization ID from environment or fetch from API/auth context
  const [organizationId] = useState<string | null>(
    import.meta.env.VITE_ORGANIZATION_ID || null
  );
  const [score, setScore] = useState<MaturityScore | null>(null);
  const [benchmark, setBenchmark] = useState<PeerComparison | null>(null);
  const [roadmap, setRoadmap] = useState<ImprovementItem[]>([]);
  const [trends, setTrends] = useState<Array<{ date: string; score: number }>>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const [s, b, r, t] = await Promise.all([
          maturityService.getScore(organizationId),
          maturityService.getBenchmark(organizationId),
          maturityService.getRoadmap(organizationId),
          maturityService.getTrends(organizationId)
        ]);
        setScore(s);
        setBenchmark(b);
        setRoadmap(r);
        setTrends(t);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [organizationId]);

  if (loading) {
    return (
      <div className="p-6">
        <div className="animate-pulse h-8 w-40 bg-gray-200 mb-4"></div>
        <div className="animate-pulse h-64 bg-gray-200"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Security Maturity</h1>
          {score && <Badge>Level: {score.level}</Badge>}
        </div>

        {/* Overall Score */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Overall Score</CardTitle>
          </CardHeader>
          <CardContent>
            {score ? (
              <div className="flex items-center gap-6">
                <div className="text-5xl font-extrabold">{score.overallScore.toFixed(1)}</div>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 flex-1">
                  {score.dimensions.map((d: DimensionScore, i: number) => (
                    <div key={i} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <div className="text-xs font-bold uppercase">{d.name}</div>
                      <div className="text-2xl font-extrabold">{d.score.toFixed(1)}</div>
                      <div className="text-xs text-gray-600">Weight: {Math.round(d.weight * 100)}%</div>
                    </div>
                  ))}
                </div>
              </div>
            ) : <div className="text-gray-500">No data</div>}
          </CardContent>
        </Card>

        {/* Benchmark */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Industry Benchmark</CardTitle>
          </CardHeader>
          <CardContent>
            {benchmark ? (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4 text-center">
                  <div className="text-xs font-bold uppercase">Industry</div>
                  <div className="text-2xl font-extrabold">{benchmark.industry}</div>
                </div>
                <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4 text-center">
                  <div className="text-xs font-bold uppercase">Percentile</div>
                  <div className="text-2xl font-extrabold">{Math.round(benchmark.percentile * 100)}%</div>
                </div>
                <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4 text-center">
                  <div className="text-xs font-bold uppercase">Peers Above / Below</div>
                  <div className="text-2xl font-extrabold">{benchmark.peersAbove} / {benchmark.peersBelow}</div>
                </div>
              </div>
            ) : <div className="text-gray-500">No benchmark data</div>}
          </CardContent>
        </Card>

        {/* Trends */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Maturity Trends</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
              {trends.map((t, i) => (
                <div key={i} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4 text-center">
                  <div className="text-xs font-bold">{t.date}</div>
                  <div className="text-2xl font-extrabold">{t.score.toFixed(1)}</div>
                </div>
              ))}
              {trends.length === 0 && <div className="text-gray-500">No trend data</div>}
            </div>
          </CardContent>
        </Card>

        {/* Roadmap */}
        <Card className="p-6">
          <CardHeader>
            <CardTitle>Improvement Roadmap</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {roadmap.map((item: ImprovementItem) => (
                <div key={item.id} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                  <div className="text-sm font-bold">{item.title}</div>
                  <div className="text-xs text-gray-600">Priority: {item.priority} • Effort: {item.effort}</div>
                </div>
              ))}
              {roadmap.length === 0 && <div className="text-gray-500">No roadmap items</div>}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

import React, { useEffect, useState } from 'react';
import { maturityService, MaturityScore, DimensionScore, ImprovementItem, PeerComparison } from '@/services/maturityService';

export default function SecurityMaturity() {
  const [organizationId] = useState<string>('123e4567-e89b-12d3-a456-426614174000'); // TODO: fetch from profile/auth
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
          {score && <span className="neobrutal-badge">Level: {score.level}</span>}
        </div>

        {/* Overall Score */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Overall Score</h2>
          {score ? (
            <div className="flex items-center gap-6">
              <div className="text-5xl font-extrabold">{score.overallScore.toFixed(1)}</div>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 flex-1">
                {score.dimensions.map((d: DimensionScore, i: number) => (
                  <div key={i} className="neobrutal-section">
                    <div className="text-xs font-bold uppercase">{d.name}</div>
                    <div className="text-2xl font-extrabold">{d.score.toFixed(1)}</div>
                    <div className="text-xs text-gray-600">Weight: {Math.round(d.weight * 100)}%</div>
                  </div>
                ))}
              </div>
            </div>
          ) : <div className="text-gray-500">No data</div>}
        </div>

        {/* Benchmark */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Industry Benchmark</h2>
          {benchmark ? (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="neobrutal-section text-center">
                <div className="text-xs font-bold uppercase">Industry</div>
                <div className="text-2xl font-extrabold">{benchmark.industry}</div>
              </div>
              <div className="neobrutal-section text-center">
                <div className="text-xs font-bold uppercase">Percentile</div>
                <div className="text-2xl font-extrabold">{Math.round(benchmark.percentile * 100)}%</div>
              </div>
              <div className="neobrutal-section text-center">
                <div className="text-xs font-bold uppercase">Peers Above / Below</div>
                <div className="text-2xl font-extrabold">{benchmark.peersAbove} / {benchmark.peersBelow}</div>
              </div>
            </div>
          ) : <div className="text-gray-500">No benchmark data</div>}
        </div>

        {/* Trends */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Maturity Trends</h2>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            {trends.map((t, i) => (
              <div key={i} className="neobrutal-section text-center">
                <div className="text-xs font-bold">{t.date}</div>
                <div className="text-2xl font-extrabold">{t.score.toFixed(1)}</div>
              </div>
            ))}
            {trends.length === 0 && <div className="text-gray-500">No trend data</div>}
          </div>
        </div>

        {/* Roadmap */}
        <div className="neobrutal-card p-6">
          <h2 className="neobrutal-title mb-4">Improvement Roadmap</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {roadmap.map((item: ImprovementItem) => (
              <div key={item.id} className="neobrutal-section">
                <div className="text-sm font-bold">{item.title}</div>
                <div className="text-xs text-gray-600">Priority: {item.priority} • Effort: {item.effort}</div>
              </div>
            ))}
            {roadmap.length === 0 && <div className="text-gray-500">No roadmap items</div>}
          </div>
        </div>
      </div>
    </div>
  );
}

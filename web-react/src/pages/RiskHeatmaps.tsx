import React, { useEffect, useState } from 'react';
import { heatmapService, HeatmapData, HeatmapPoint, RiskDistributionBucket, TrendPoint } from '@/services/heatmapService';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';

function HeatmapGrid({ data }: { data: HeatmapData }) {
  const max = Math.max(1, ...data.data.map(d => d.value));
  return (
    <div className="overflow-auto">
      <div className="grid" style={{ gridTemplateColumns: `150px repeat(${data.xAxis.length}, minmax(60px, 1fr))` }}>
        <div></div>
        {data.xAxis.map(x => (
          <div key={x} className="p-2 text-xs font-bold border-b border-black">{x}</div>
        ))}
        {data.yAxis.map(y => (
          <React.Fragment key={y}>
            <div className="p-2 text-xs font-bold border-r border-black">{y}</div>
            {data.xAxis.map(x => {
              const point = data.data.find(p => p.x === x && p.y === y);
              const value = point?.value || 0;
              const intensity = value / max;
              const bg = `rgba(255, 0, 0, ${intensity})`;
              return (
                <div key={x + y} className="h-10 border border-black flex items-center justify-center" style={{ background: bg }}>
                  <span className="text-xs font-bold">{value}</span>
                </div>
              );
            })}
          </React.Fragment>
        ))}
      </div>
    </div>
  );
}

export default function RiskHeatmaps() {
  // Get organization ID from environment or fetch from API/auth context
  const [organizationId] = useState<string | null>(
    import.meta.env.VITE_ORGANIZATION_ID || null
  );
  const [heatmap, setHeatmap] = useState<HeatmapData | null>(null);
  const [hotspots, setHotspots] = useState<HeatmapPoint[]>([]);
  const [distribution, setDistribution] = useState<RiskDistributionBucket[]>([]);
  const [trends, setTrends] = useState<TrendPoint[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      try {
        const [hm, hs, rd, tr] = await Promise.all([
          heatmapService.getHeatmap(organizationId),
          heatmapService.getHotspots(organizationId),
          heatmapService.getRiskDistribution(organizationId),
          heatmapService.getTrends(organizationId)
        ]);
        setHeatmap(hm);
        setHotspots(hs);
        setDistribution(rd);
        setTrends(tr);
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
          <h1 className="text-3xl font-bold">Risk Heatmaps</h1>
          <Badge>Org: {organizationId.slice(0, 8)}...</Badge>
        </div>

        <Card className="p-4">
          <CardHeader>
            <CardTitle>Severity by Technology</CardTitle>
          </CardHeader>
          <CardContent>
            {heatmap ? (
              <HeatmapGrid data={heatmap} />
            ) : (
              <div className="text-gray-500">No data</div>
            )}
          </CardContent>
        </Card>

        <Card className="p-4">
          <CardHeader>
            <CardTitle>Hotspots</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {hotspots.map((h, i) => (
                <div key={i} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                  <div className="flex items-center justify-between">
                    <span className="font-bold">{h.y} / {h.x}</span>
                    <Badge variant="destructive">{h.value}</Badge>
                  </div>
                </div>
              ))}
              {hotspots.length === 0 && <div className="text-gray-500">No hotspots</div>}
            </div>
          </CardContent>
        </Card>

        <Card className="p-4">
          <CardHeader>
            <CardTitle>Risk Distribution</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {distribution.map((b, i) => (
                <div key={i} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4 text-center">
                  <div className="text-sm font-bold uppercase">{b.severity}</div>
                  <div className="text-2xl font-extrabold">{b.count}</div>
                </div>
              ))}
              {distribution.length === 0 && <div className="text-gray-500">No distribution data</div>}
            </div>
          </CardContent>
        </Card>

        <Card className="p-4">
          <CardHeader>
            <CardTitle>Trends</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
              {trends.map((t, i) => (
                <div key={i} className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                  <div className="text-xs font-bold">{t.date}</div>
                  <div className="text-xs">Crit: <span className="font-bold text-red-600">{t.critical}</span></div>
                  <div className="text-xs">High: <span className="font-bold text-orange-600">{t.high}</span></div>
                  <div className="text-xs">Med: <span className="font-bold text-yellow-600">{t.medium}</span></div>
                  <div className="text-xs">Low: <span className="font-bold text-green-600">{t.low}</span></div>
                </div>
              ))}
              {trends.length === 0 && <div className="text-gray-500">No trend data</div>}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

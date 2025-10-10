import React from 'react';
import { TrendingUp, TrendingDown, AlertTriangle, Shield, Clock, Target } from 'lucide-react';

interface KPIMetric {
  id: string;
  label: string;
  value: number | string;
  delta?: {
    value: number;
    trend: 'up' | 'down' | 'neutral';
    period: string;
  };
  sparkline?: number[];
  tooltip?: string;
  icon?: React.ReactNode;
  color?: 'default' | 'critical' | 'warning' | 'success';
}

interface KPIRibbonProps {
  metrics: KPIMetric[];
  loading?: boolean;
  className?: string;
}

const KPIRibbon: React.FC<KPIRibbonProps> = ({ metrics, loading = false, className = '' }) => {
  const getIcon = (id: string) => {
    switch (id) {
      case 'active_critical_cves':
        return <AlertTriangle className="w-5 h-5" />;
      case 'mean_time_to_remediate':
        return <Clock className="w-5 h-5" />;
      case 'compliance_percent':
        return <Shield className="w-5 h-5" />;
      case 'scan_coverage':
        return <Target className="w-5 h-5" />;
      case 'exploit_presence':
        return <AlertTriangle className="w-5 h-5" />;
      default:
        return <TrendingUp className="w-5 h-5" />;
    }
  };

  const getColorClasses = (color: string = 'default') => {
    switch (color) {
      case 'critical':
        return 'border-red-500 bg-red-50 text-red-900';
      case 'warning':
        return 'border-yellow-500 bg-yellow-50 text-yellow-900';
      case 'success':
        return 'border-green-500 bg-green-50 text-green-900';
      default:
        return 'border-black bg-white text-black';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'up':
        return <TrendingUp className="w-4 h-4 text-red-600" />;
      case 'down':
        return <TrendingDown className="w-4 h-4 text-green-600" />;
      default:
        return <div className="w-4 h-4" />;
    }
  };

  if (loading) {
    return (
      <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6 gap-4 ${className}`}>
        {Array.from({ length: 6 }).map((_, i) => (
          <div key={i} className="animate-pulse">
            <div className="h-24 bg-gray-200 rounded border-3 border-black shadow-lg"></div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-6 gap-4 ${className}`}>
      {metrics.map((metric) => (
        <div
          key={metric.id}
          className={`
            relative p-4 rounded border-3 border-black bg-white shadow-lg hover:shadow-xl
            transition-all duration-150 ease-in-out hover:translate-x-1 hover:translate-y-1
            hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)]
            ${getColorClasses(metric.color)}
          `}
          title={metric.tooltip}
        >
          {/* Header */}
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center gap-2">
              {metric.icon || getIcon(metric.id)}
              <span className="text-xs font-bold uppercase tracking-wider text-gray-600">
                {metric.label}
              </span>
            </div>
            {metric.delta && (
              <div className="flex items-center gap-1">
                {getTrendIcon(metric.delta.trend)}
                <span className={`text-xs font-bold ${
                  metric.delta.trend === 'up' ? 'text-red-600' : 
                  metric.delta.trend === 'down' ? 'text-green-600' : 
                  'text-gray-600'
                }`}>
                  {Math.abs(metric.delta.value)}%
                </span>
              </div>
            )}
          </div>

          {/* Value */}
          <div className="text-2xl font-bold mb-1">
            {typeof metric.value === 'number' ? metric.value.toLocaleString() : metric.value}
          </div>

          {/* Sparkline */}
          {metric.sparkline && metric.sparkline.length > 0 && (
            <div className="h-8 w-full">
              <svg className="w-full h-full" viewBox="0 0 100 32">
                <polyline
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  points={metric.sparkline.map((value, index) => 
                    `${(index / (metric.sparkline!.length - 1)) * 100},${32 - (value / Math.max(...metric.sparkline!)) * 32}`
                  ).join(' ')}
                />
              </svg>
            </div>
          )}

          {/* Delta period */}
          {metric.delta && (
            <div className="text-xs text-gray-500 mt-1">
              vs {metric.delta.period}
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default KPIRibbon;

import React, { useState, useEffect } from 'react';
import { aiService } from '@/services/aiService';

interface AIAnalysis {
  vulnerabilityId: string;
  analysisTimestamp: string;
  exploitIntelligence: {
    exploitAvailability: boolean;
    exploitSources: string[];
    exploitComplexity: string;
    exploitLikelihood: number;
    cisaKev: boolean;
    confidenceScore: number;
  };
  predictiveAnalysis: {
    riskTrend: string;
    exploitProbability: number;
    patchAvailability: boolean;
    patchTimeline: string;
    businessImpact: string;
    remediationUrgency: string;
  };
  remediationGuidance: {
    immediateActions: string[];
    longTermActions: string[];
    complexityScore: number;
    estimatedEffort: string;
    requiredSkills: string[];
  };
  businessContext: {
    affectedSystems: string[];
    businessCriticality: string;
    complianceImpact: string[];
    financialRisk: number;
    reputationRisk: string;
  };
  aiConfidence: {
    overallConfidence: number;
    dataQuality: number;
    modelAccuracy: number;
    uncertaintyFactors: string[];
  };
}

interface VulnerabilityTrend {
  date: string;
  critical: number;
  high: number;
  medium: number;
  low: number;
  total: number;
}

interface ExploitIntelligence {
  cveId: string;
  exploitAvailability: boolean;
  exploitSources: string[];
  exploitComplexity: string;
  exploitLikelihood: number;
  cisaKev: boolean;
  confidenceScore: number;
  lastUpdated: string;
}

function ProgressBar({ value }: { value: number }) {
  return (
    <div className="neobrutal-progress">
      <div className="neobrutal-progress-bar" style={{ width: `${Math.max(0, Math.min(100, value))}%` }} />
    </div>
  );
}

export default function AIAnalytics() {
  const [analysis, setAnalysis] = useState<AIAnalysis | null>(null);
  const [trends, setTrends] = useState<VulnerabilityTrend[]>([]);
  const [exploitIntel, setExploitIntel] = useState<ExploitIntelligence | null>(null);
  const [loading, setLoading] = useState(false);
  const [selectedVulnerability, setSelectedVulnerability] = useState<string>('');
  const [activeTab, setActiveTab] = useState<'trends' | 'analysis' | 'exploit' | 'predictive'>('trends');

  useEffect(() => {
    loadTrends();
  }, []);

  const loadTrends = async () => {
    setLoading(true);
    try {
      const data = await aiService.getVulnerabilityTrends();
      setTrends(data);
    } catch (error) {
      console.error('Failed to load vulnerability trends:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadAnalysis = async (vulnerabilityId: string) => {
    if (!vulnerabilityId) return;
    setLoading(true);
    try {
      const data = await aiService.getComprehensiveAnalysis(vulnerabilityId);
      setAnalysis(data);
    } catch (error) {
      console.error('Failed to load AI analysis:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadExploitIntelligence = async (cveId: string) => {
    setLoading(true);
    try {
      const data = await aiService.getExploitIntelligence(cveId);
      setExploitIntel(data);
    } catch (error) {
      console.error('Failed to load exploit intelligence:', error);
    } finally {
      setLoading(false);
    }
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.8) return 'text-green-600';
    if (confidence >= 0.6) return 'text-yellow-600';
    return 'text-red-600';
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-7xl mx-auto">
          <div className="animate-pulse">
            <div className="h-8 bg-gray-200 rounded w-1/4 mb-6"></div>
            <div className="grid gap-6">
              <div className="h-64 bg-gray-200 rounded"></div>
              <div className="h-64 bg-gray-200 rounded"></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">AI Analytics Dashboard</h1>
          <div className="flex gap-4">
            <button onClick={loadTrends} className="neobrutal-button secondary px-4 py-2">Refresh Trends</button>
            <button onClick={() => loadAnalysis(selectedVulnerability)} disabled={!selectedVulnerability} className="neobrutal-button primary px-4 py-2">
              Analyze Selected
            </button>
          </div>
        </div>

        {/* Tabs */}
        <div className="neobrutal-nav flex gap-2 p-2">
          {(
            [
              { key: 'trends', label: 'Vulnerability Trends' },
              { key: 'analysis', label: 'AI Analysis' },
              { key: 'exploit', label: 'Exploit Intelligence' },
              { key: 'predictive', label: 'Predictive Analysis' },
            ] as const
          ).map(tab => (
            <button
              key={tab.key}
              className={`neobrutal-nav-item ${activeTab === tab.key ? 'active' : ''}`}
              onClick={() => setActiveTab(tab.key)}
            >
              {tab.label}
            </button>
          ))}
        </div>

        {/* Trends */}
        {activeTab === 'trends' && (
          <div className="neobrutal-card p-4 space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-2">
              <div className="neobrutal-section">
                <h3 className="neobrutal-subtitle mb-2">Critical Vulnerabilities</h3>
                <div className="text-3xl font-bold text-red-600">
                  {trends.reduce((sum, trend) => sum + trend.critical, 0)}
                </div>
                <div className="text-sm text-gray-600">Last 30 days</div>
              </div>
              <div className="neobrutal-section">
                <h3 className="neobrutal-subtitle mb-2">High Vulnerabilities</h3>
                <div className="text-3xl font-bold text-orange-600">
                  {trends.reduce((sum, trend) => sum + trend.high, 0)}
                </div>
                <div className="text-sm text-gray-600">Last 30 days</div>
              </div>
              <div className="neobrutal-section">
                <h3 className="neobrutal-subtitle mb-2">Medium Vulnerabilities</h3>
                <div className="text-3xl font-bold text-yellow-600">
                  {trends.reduce((sum, trend) => sum + trend.medium, 0)}
                </div>
                <div className="text-sm text-gray-600">Last 30 days</div>
              </div>
              <div className="neobrutal-section">
                <h3 className="neobrutal-subtitle mb-2">Total Vulnerabilities</h3>
                <div className="text-3xl font-bold text-gray-900">
                  {trends.reduce((sum, trend) => sum + trend.total, 0)}
                </div>
                <div className="text-sm text-gray-600">Last 30 days</div>
              </div>
            </div>
            <div className="neobrutal-container">
              <h3 className="neobrutal-subtitle mb-4">Vulnerability Trend Chart</h3>
              <div className="h-64 bg-gray-100 flex items-center justify-center">
                <p className="text-gray-500">Chart visualization would go here</p>
              </div>
            </div>
          </div>
        )}

        {/* Analysis */}
        {activeTab === 'analysis' && (
          analysis ? (
            <div className="space-y-6">
              <div className="neobrutal-card p-4 space-y-4">
                <h2 className="neobrutal-title">Exploit Intelligence</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Exploit Availability</h4>
                    <span className={`neobrutal-badge ${analysis.exploitIntelligence.exploitAvailability ? 'danger' : 'success'}`}>
                      {analysis.exploitIntelligence.exploitAvailability ? 'Available' : 'Not Available'}
                    </span>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Exploit Likelihood</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.exploitIntelligence.exploitLikelihood * 100} />
                      <span className="font-bold">{(analysis.exploitIntelligence.exploitLikelihood * 100).toFixed(1)}%</span>
                    </div>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">CISA KEV</h4>
                    <span className={`neobrutal-badge ${analysis.exploitIntelligence.cisaKev ? 'danger' : 'secondary'}`}>
                      {analysis.exploitIntelligence.cisaKev ? 'Listed' : 'Not Listed'}
                    </span>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Confidence Score</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.exploitIntelligence.confidenceScore * 100} />
                      <span className={`font-bold ${getConfidenceColor(analysis.exploitIntelligence.confidenceScore)}`}>
                        {(analysis.exploitIntelligence.confidenceScore * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>
                </div>
                <div className="neobrutal-section">
                  <h4 className="neobrutal-subtitle mb-2">Exploit Sources</h4>
                  <div className="flex flex-wrap gap-2">
                    {analysis.exploitIntelligence.exploitSources.map((source, index) => (
                      <span key={index} className="neobrutal-badge">{source}</span>
                    ))}
                  </div>
                </div>
              </div>

              <div className="neobrutal-card p-4 space-y-4">
                <h2 className="neobrutal-title">Predictive Analysis</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Risk Trend</h4>
                    <span className={`neobrutal-badge ${analysis.predictiveAnalysis.riskTrend === 'increasing' ? 'danger' : 'success'}`}>
                      {analysis.predictiveAnalysis.riskTrend}
                    </span>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Exploit Probability</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.predictiveAnalysis.exploitProbability * 100} />
                      <span className="font-bold">{(analysis.predictiveAnalysis.exploitProbability * 100).toFixed(1)}%</span>
                    </div>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Patch Availability</h4>
                    <span className={`neobrutal-badge ${analysis.predictiveAnalysis.patchAvailability ? 'success' : 'warning'}`}>
                      {analysis.predictiveAnalysis.patchAvailability ? 'Available' : 'Not Available'}
                    </span>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Remediation Urgency</h4>
                    <span className={`neobrutal-badge ${analysis.predictiveAnalysis.remediationUrgency === 'high' ? 'danger' : 'secondary'}`}>
                      {analysis.predictiveAnalysis.remediationUrgency}
                    </span>
                  </div>
                </div>
              </div>

              <div className="neobrutal-card p-4 space-y-4">
                <h2 className="neobrutal-title">AI-Generated Remediation Plan</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Immediate Actions</h4>
                    <ul className="space-y-1">
                      {analysis.remediationGuidance.immediateActions.map((action, index) => (
                        <li key={index} className="flex items-start gap-2">
                          <span className="text-red-500 font-bold">•</span>
                          <span className="text-sm">{action}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Long-term Actions</h4>
                    <ul className="space-y-1">
                      {analysis.remediationGuidance.longTermActions.map((action, index) => (
                        <li key={index} className="flex items-start gap-2">
                          <span className="text-blue-500 font-bold">•</span>
                          <span className="text-sm">{action}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Complexity Score</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.remediationGuidance.complexityScore * 100} />
                      <span className="font-bold">{(analysis.remediationGuidance.complexityScore * 100).toFixed(1)}%</span>
                    </div>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Estimated Effort</h4>
                    <p className="font-bold">{analysis.remediationGuidance.estimatedEffort}</p>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Required Skills</h4>
                    <div className="flex flex-wrap gap-1">
                      {analysis.remediationGuidance.requiredSkills.map((skill, index) => (
                        <span key={index} className="neobrutal-badge text-xs">{skill}</span>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              <div className="neobrutal-card p-4 space-y-4">
                <h2 className="neobrutal-title">Business Impact Analysis</h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Affected Systems</h4>
                    <div className="flex flex-wrap gap-2">
                      {analysis.businessContext.affectedSystems.map((system, index) => (
                        <span key={index} className="neobrutal-badge">{system}</span>
                      ))}
                    </div>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Business Criticality</h4>
                    <span className={`neobrutal-badge ${analysis.businessContext.businessCriticality === 'high' ? 'danger' : 'secondary'}`}>
                      {analysis.businessContext.businessCriticality}
                    </span>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Financial Risk</h4>
                    <p className="text-2xl font-bold text-red-600">${analysis.businessContext.financialRisk.toLocaleString()}</p>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Reputation Risk</h4>
                    <span className={`neobrutal-badge ${analysis.businessContext.reputationRisk === 'high' ? 'danger' : 'secondary'}`}>
                      {analysis.businessContext.reputationRisk}
                    </span>
                  </div>
                </div>
              </div>

              <div className="neobrutal-card p-4 space-y-4">
                <h2 className="neobrutal-title">AI Confidence Metrics</h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Overall Confidence</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.aiConfidence.overallConfidence * 100} />
                      <span className={`font-bold ${getConfidenceColor(analysis.aiConfidence.overallConfidence)}`}>
                        {(analysis.aiConfidence.overallConfidence * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Data Quality</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.aiConfidence.dataQuality * 100} />
                      <span className={`font-bold ${getConfidenceColor(analysis.aiConfidence.dataQuality)}`}>
                        {(analysis.aiConfidence.dataQuality * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>
                  <div className="neobrutal-section">
                    <h4 className="neobrutal-subtitle mb-2">Model Accuracy</h4>
                    <div className="flex items-center gap-2">
                      <ProgressBar value={analysis.aiConfidence.modelAccuracy * 100} />
                      <span className={`font-bold ${getConfidenceColor(analysis.aiConfidence.modelAccuracy)}`}>
                        {(analysis.aiConfidence.modelAccuracy * 100).toFixed(1)}%
                      </span>
                    </div>
                  </div>
                </div>
                <div className="neobrutal-section">
                  <h4 className="neobrutal-subtitle mb-2">Uncertainty Factors</h4>
                  <div className="flex flex-wrap gap-2">
                    {analysis.aiConfidence.uncertaintyFactors.map((factor, index) => (
                      <span key={index} className="neobrutal-badge">{factor}</span>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className="neobrutal-card">
              <div className="text-center p-12">
                <h3 className="neobrutal-subtitle mb-4">No Analysis Selected</h3>
                <p className="text-gray-600 mb-4">Select a vulnerability to perform AI-powered analysis</p>
                <button onClick={() => loadAnalysis('sample-vuln-1')} className="neobrutal-button primary px-4 py-2">Analyze Sample Vulnerability</button>
              </div>
            </div>
          )
        )}

        {/* Exploit Intelligence */}
        {activeTab === 'exploit' && (
          exploitIntel ? (
            <div className="neobrutal-card p-4 space-y-4">
              <h2 className="neobrutal-title">Exploit Intelligence Report</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="neobrutal-section">
                  <h4 className="neobrutal-subtitle mb-2">CVE ID</h4>
                  <p className="font-mono font-bold">{exploitIntel.cveId}</p>
                </div>
                <div className="neobrutal-section">
                  <h4 className="neobrutal-subtitle mb-2">Exploit Availability</h4>
                  <span className={`neobrutal-badge ${exploitIntel.exploitAvailability ? 'danger' : 'success'}`}>
                    {exploitIntel.exploitAvailability ? 'Available' : 'Not Available'}
                  </span>
                </div>
                <div className="neobrutal-section">
                  <h4 className="neobrutal-subtitle mb-2">Exploit Complexity</h4>
                  <span className={`neobrutal-badge ${exploitIntel.exploitComplexity === 'low' ? 'success' : exploitIntel.exploitComplexity === 'medium' ? 'warning' : 'danger'}`}>
                    {exploitIntel.exploitComplexity}
                  </span>
                </div>
                <div className="neobrutal-section">
                  <h4 className="neobrutal-subtitle mb-2">CISA KEV</h4>
                  <span className={`neobrutal-badge ${exploitIntel.cisaKev ? 'danger' : 'secondary'}`}>
                    {exploitIntel.cisaKev ? 'Listed' : 'Not Listed'}
                  </span>
                </div>
              </div>
              <div className="neobrutal-section">
                <h4 className="neobrutal-subtitle mb-2">Exploit Sources</h4>
                <div className="flex flex-wrap gap-2">
                  {exploitIntel.exploitSources.map((source, index) => (
                    <span key={index} className="neobrutal-badge">{source}</span>
                  ))}
                </div>
              </div>
            </div>
          ) : (
            <div className="neobrutal-card">
              <div className="text-center p-12">
                <h3 className="neobrutal-subtitle mb-4">No Exploit Intelligence</h3>
                <p className="text-gray-600 mb-4">Enter a CVE ID to get exploit intelligence</p>
                <button onClick={() => loadExploitIntelligence('CVE-2024-1234')} className="neobrutal-button primary px-4 py-2">Get Sample Intelligence</button>
              </div>
            </div>
          )
        )}

        {/* Predictive */}
        {activeTab === 'predictive' && (
          <div className="neobrutal-card p-8 text-center">
            <h3 className="neobrutal-subtitle mb-4">Predictive Analysis Coming Soon</h3>
            <p className="text-gray-600">AI-powered predictive analysis will be available in future updates</p>
          </div>
        )}
      </div>
    </div>
  );
}

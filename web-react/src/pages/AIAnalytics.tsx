import React, { useState, useEffect } from 'react';
import { aiService } from '@/services/aiService';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';

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

import { Progress } from '@/components/ui/progress';

function ProgressBar({ value }: { value: number }) {
  return <Progress value={value} />;
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
            <Button variant="secondary" onClick={loadTrends}>Refresh Trends</Button>
            <Button variant="default" onClick={() => loadAnalysis(selectedVulnerability)} disabled={!selectedVulnerability}>
              Analyze Selected
            </Button>
          </div>
        </div>

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as typeof activeTab)}>
          <TabsList>
            <TabsTrigger value="trends">Vulnerability Trends</TabsTrigger>
            <TabsTrigger value="analysis">AI Analysis</TabsTrigger>
            <TabsTrigger value="exploit">Exploit Intelligence</TabsTrigger>
            <TabsTrigger value="predictive">Predictive Analysis</TabsTrigger>
          </TabsList>

          <TabsContent value="trends">
            <Card className="p-4 space-y-4">
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-2">
                  <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                    <h3 className="font-bold uppercase tracking-wide mb-2">Critical Vulnerabilities</h3>
                    <div className="text-3xl font-bold text-red-600">
                      {trends.reduce((sum, trend) => sum + trend.critical, 0)}
                    </div>
                    <div className="text-sm text-gray-600">Last 30 days</div>
                  </div>
                  <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                    <h3 className="font-bold uppercase tracking-wide mb-2">High Vulnerabilities</h3>
                    <div className="text-3xl font-bold text-orange-600">
                      {trends.reduce((sum, trend) => sum + trend.high, 0)}
                    </div>
                    <div className="text-sm text-gray-600">Last 30 days</div>
                  </div>
                  <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                    <h3 className="font-bold uppercase tracking-wide mb-2">Medium Vulnerabilities</h3>
                    <div className="text-3xl font-bold text-yellow-600">
                      {trends.reduce((sum, trend) => sum + trend.medium, 0)}
                    </div>
                    <div className="text-sm text-gray-600">Last 30 days</div>
                  </div>
                  <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                    <h3 className="font-bold uppercase tracking-wide mb-2">Total Vulnerabilities</h3>
                    <div className="text-3xl font-bold text-gray-900">
                      {trends.reduce((sum, trend) => sum + trend.total, 0)}
                    </div>
                    <div className="text-sm text-gray-600">Last 30 days</div>
                  </div>
                </div>
                <div className="border-4 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                  <h3 className="font-bold uppercase tracking-wide mb-4">Vulnerability Trend Chart</h3>
                  <div className="h-64 bg-gray-100 flex items-center justify-center">
                    <p className="text-gray-500">Chart visualization would go here</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="analysis">
            {analysis ? (
              <div className="space-y-6">
                <Card className="p-4 space-y-4">
                  <CardHeader>
                    <CardTitle>Exploit Intelligence</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Availability</h4>
                        <Badge variant={analysis.exploitIntelligence.exploitAvailability ? 'destructive' : 'success'}>
                          {analysis.exploitIntelligence.exploitAvailability ? 'Available' : 'Not Available'}
                        </Badge>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Likelihood</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.exploitIntelligence.exploitLikelihood * 100} />
                          <span className="font-bold">{(analysis.exploitIntelligence.exploitLikelihood * 100).toFixed(1)}%</span>
                        </div>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">CISA KEV</h4>
                        <Badge variant={analysis.exploitIntelligence.cisaKev ? 'destructive' : 'secondary'}>
                          {analysis.exploitIntelligence.cisaKev ? 'Listed' : 'Not Listed'}
                        </Badge>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Confidence Score</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.exploitIntelligence.confidenceScore * 100} />
                          <span className={`font-bold ${getConfidenceColor(analysis.exploitIntelligence.confidenceScore)}`}>
                            {(analysis.exploitIntelligence.confidenceScore * 100).toFixed(1)}%
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Sources</h4>
                      <div className="flex flex-wrap gap-2">
                        {analysis.exploitIntelligence.exploitSources.map((source, index) => (
                          <Badge key={index}>{source}</Badge>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card className="p-4 space-y-4">
                  <CardHeader>
                    <CardTitle>Predictive Analysis</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Risk Trend</h4>
                        <Badge variant={analysis.predictiveAnalysis.riskTrend === 'increasing' ? 'destructive' : 'success'}>
                          {analysis.predictiveAnalysis.riskTrend}
                        </Badge>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Probability</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.predictiveAnalysis.exploitProbability * 100} />
                          <span className="font-bold">{(analysis.predictiveAnalysis.exploitProbability * 100).toFixed(1)}%</span>
                        </div>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Patch Availability</h4>
                        <Badge variant={analysis.predictiveAnalysis.patchAvailability ? 'success' : 'warning'}>
                          {analysis.predictiveAnalysis.patchAvailability ? 'Available' : 'Not Available'}
                        </Badge>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Remediation Urgency</h4>
                        <Badge variant={analysis.predictiveAnalysis.remediationUrgency === 'high' ? 'destructive' : 'secondary'}>
                          {analysis.predictiveAnalysis.remediationUrgency}
                        </Badge>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card className="p-4 space-y-4">
                  <CardHeader>
                    <CardTitle>AI-Generated Remediation Plan</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Immediate Actions</h4>
                        <ul className="space-y-1">
                          {analysis.remediationGuidance.immediateActions.map((action, index) => (
                            <li key={index} className="flex items-start gap-2">
                              <span className="text-red-500 font-bold">•</span>
                              <span className="text-sm">{action}</span>
                            </li>
                          ))}
                        </ul>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Long-term Actions</h4>
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
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Complexity Score</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.remediationGuidance.complexityScore * 100} />
                          <span className="font-bold">{(analysis.remediationGuidance.complexityScore * 100).toFixed(1)}%</span>
                        </div>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Estimated Effort</h4>
                        <p className="font-bold">{analysis.remediationGuidance.estimatedEffort}</p>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Required Skills</h4>
                        <div className="flex flex-wrap gap-1">
                          {analysis.remediationGuidance.requiredSkills.map((skill, index) => (
                            <Badge key={index} className="text-xs">{skill}</Badge>
                          ))}
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card className="p-4 space-y-4">
                  <CardHeader>
                    <CardTitle>Business Impact Analysis</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Affected Systems</h4>
                        <div className="flex flex-wrap gap-2">
                          {analysis.businessContext.affectedSystems.map((system, index) => (
                            <Badge key={index}>{system}</Badge>
                          ))}
                        </div>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Business Criticality</h4>
                        <Badge variant={analysis.businessContext.businessCriticality === 'high' ? 'destructive' : 'secondary'}>
                          {analysis.businessContext.businessCriticality}
                        </Badge>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Financial Risk</h4>
                        <p className="text-2xl font-bold text-red-600">${analysis.businessContext.financialRisk.toLocaleString()}</p>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Reputation Risk</h4>
                        <Badge variant={analysis.businessContext.reputationRisk === 'high' ? 'destructive' : 'secondary'}>
                          {analysis.businessContext.reputationRisk}
                        </Badge>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <Card className="p-4 space-y-4">
                  <CardHeader>
                    <CardTitle>AI Confidence Metrics</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Overall Confidence</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.aiConfidence.overallConfidence * 100} />
                          <span className={`font-bold ${getConfidenceColor(analysis.aiConfidence.overallConfidence)}`}>
                            {(analysis.aiConfidence.overallConfidence * 100).toFixed(1)}%
                          </span>
                        </div>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Data Quality</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.aiConfidence.dataQuality * 100} />
                          <span className={`font-bold ${getConfidenceColor(analysis.aiConfidence.dataQuality)}`}>
                            {(analysis.aiConfidence.dataQuality * 100).toFixed(1)}%
                          </span>
                        </div>
                      </div>
                      <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                        <h4 className="font-bold uppercase tracking-wide mb-2">Model Accuracy</h4>
                        <div className="flex items-center gap-2">
                          <ProgressBar value={analysis.aiConfidence.modelAccuracy * 100} />
                          <span className={`font-bold ${getConfidenceColor(analysis.aiConfidence.modelAccuracy)}`}>
                            {(analysis.aiConfidence.modelAccuracy * 100).toFixed(1)}%
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <h4 className="font-bold uppercase tracking-wide mb-2">Uncertainty Factors</h4>
                      <div className="flex flex-wrap gap-2">
                        {analysis.aiConfidence.uncertaintyFactors.map((factor, index) => (
                          <Badge key={index}>{factor}</Badge>
                        ))}
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            ) : (
              <Card>
                <CardContent className="text-center p-12">
                  <h3 className="font-bold uppercase tracking-wide mb-4">No Analysis Selected</h3>
                  <p className="text-gray-600 mb-4">Select a vulnerability to perform AI-powered analysis</p>
                  <Button variant="default" onClick={() => loadAnalysis('sample-vuln-1')}>Analyze Sample Vulnerability</Button>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          <TabsContent value="exploit">
            {exploitIntel ? (
              <Card className="p-4 space-y-4">
                <CardHeader>
                  <CardTitle>Exploit Intelligence Report</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <h4 className="font-bold uppercase tracking-wide mb-2">CVE ID</h4>
                      <p className="font-mono font-bold">{exploitIntel.cveId}</p>
                    </div>
                    <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Availability</h4>
                      <Badge variant={exploitIntel.exploitAvailability ? 'destructive' : 'success'}>
                        {exploitIntel.exploitAvailability ? 'Available' : 'Not Available'}
                      </Badge>
                    </div>
                    <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Complexity</h4>
                      <Badge variant={exploitIntel.exploitComplexity === 'low' ? 'success' : exploitIntel.exploitComplexity === 'medium' ? 'warning' : 'destructive'}>
                        {exploitIntel.exploitComplexity}
                      </Badge>
                    </div>
                    <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                      <h4 className="font-bold uppercase tracking-wide mb-2">CISA KEV</h4>
                      <Badge variant={exploitIntel.cisaKev ? 'destructive' : 'secondary'}>
                        {exploitIntel.cisaKev ? 'Listed' : 'Not Listed'}
                      </Badge>
                    </div>
                  </div>
                  <div className="border-2 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] bg-white p-4">
                    <h4 className="font-bold uppercase tracking-wide mb-2">Exploit Sources</h4>
                    <div className="flex flex-wrap gap-2">
                      {exploitIntel.exploitSources.map((source, index) => (
                        <Badge key={index}>{source}</Badge>
                      ))}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ) : (
              <Card>
                <CardContent className="text-center p-12">
                  <h3 className="font-bold uppercase tracking-wide mb-4">No Exploit Intelligence</h3>
                  <p className="text-gray-600 mb-4">Enter a CVE ID to get exploit intelligence</p>
                  <Button variant="default" onClick={() => loadExploitIntelligence('CVE-2024-1234')}>Get Sample Intelligence</Button>
                </CardContent>
              </Card>
            )}
          </TabsContent>

          <TabsContent value="predictive">
            <Card className="p-8 text-center">
              <CardContent>
                <h3 className="font-bold uppercase tracking-wide mb-4">Predictive Analysis Coming Soon</h3>
                <p className="text-gray-600">AI-powered predictive analysis will be available in future updates</p>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}

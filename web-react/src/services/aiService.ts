import { api } from './api';

export interface AIAnalysis {
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

export interface VulnerabilityTrend {
  date: string;
  critical: number;
  high: number;
  medium: number;
  low: number;
  total: number;
}

export interface ExploitIntelligence {
  cveId: string;
  exploitAvailability: boolean;
  exploitSources: string[];
  exploitComplexity: string;
  exploitLikelihood: number;
  cisaKev: boolean;
  confidenceScore: number;
  lastUpdated: string;
}

export interface PredictiveAnalysis {
  vulnerabilityId: string;
  riskTrend: string;
  exploitProbability: number;
  patchAvailability: boolean;
  patchTimeline: string;
  businessImpact: string;
  remediationUrgency: string;
  confidenceScore: number;
}

export interface RemediationPlan {
  planId: string;
  title: string;
  description: string;
  severity: string;
  estimatedTotalTime: string;
  complexityLevel: string;
  steps: Array<{
    stepNumber: number;
    title: string;
    description: string;
    estimatedTime: string;
    requiredSkills: string[];
    priority: string;
  }>;
}

export interface BulkAnalysis {
  analysisId: string;
  vulnerabilityIds: string[];
  analysisType: string;
  results: AIAnalysis[];
  summary: {
    totalAnalyzed: number;
    highRiskCount: number;
    mediumRiskCount: number;
    lowRiskCount: number;
    averageConfidence: number;
  };
}

export const aiService = {
  async getComprehensiveAnalysis(vulnerabilityId: string): Promise<AIAnalysis> {
    try {
      const response = await api.get(`/api/ai-analysis/vulnerabilities/${vulnerabilityId}/comprehensive`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch comprehensive analysis:', error);
      throw error;
    }
  },

  async getVulnerabilityTrends(): Promise<VulnerabilityTrend[]> {
    try {
      const response = await api.get('/api/ai-analysis/vulnerabilities/trends');
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch vulnerability trends:', error);
      return [];
    }
  },

  async getExploitIntelligence(cveId: string): Promise<ExploitIntelligence> {
    try {
      const response = await api.get(`/api/ai-analysis/exploit-intelligence/${cveId}`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch exploit intelligence:', error);
      throw error;
    }
  },

  async getPredictiveAnalysis(vulnerabilityId: string): Promise<PredictiveAnalysis> {
    try {
      const response = await api.get(`/api/ai-analysis/vulnerabilities/${vulnerabilityId}/predictive`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch predictive analysis:', error);
      throw error;
    }
  },

  async getRemediationPlan(vulnerabilityId: string): Promise<RemediationPlan> {
    try {
      const response = await api.get(`/api/ai-analysis/vulnerabilities/${vulnerabilityId}/remediation-plan`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch remediation plan:', error);
      throw error;
    }
  },

  async getBulkAnalysis(vulnerabilityIds: string[], analysisType: string = 'comprehensive'): Promise<BulkAnalysis> {
    try {
      const response = await api.post('/api/ai-analysis/bulk-analysis', {
        vulnerability_ids: vulnerabilityIds,
        analysis_type: analysisType
      });
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch bulk analysis:', error);
      throw error;
    }
  }
};

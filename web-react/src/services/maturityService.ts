import { api } from './api';

export interface DimensionScore {
  name: string;
  score: number;
  weight: number;
}

export interface MaturityScore {
  overallScore: number;
  level: string;
  dimensions: DimensionScore[];
  benchmarkPercentile?: number;
}

export interface PeerComparison {
  industry: string;
  percentile: number;
  peersAbove: number;
  peersBelow: number;
}

export interface ImprovementItem {
  id: string;
  title: string;
  priority: string;
  effort: string;
}

export const maturityService = {
  async getScore(organizationId: string): Promise<MaturityScore> {
    const resp = await api.get(`/api/maturity/organizations/${organizationId}/score`);
    return resp.data.data;
  },

  async getBenchmark(organizationId: string): Promise<PeerComparison> {
    const resp = await api.get(`/api/maturity/organizations/${organizationId}/benchmark`);
    return resp.data.data;
  },

  async getRoadmap(organizationId: string): Promise<ImprovementItem[]> {
    const resp = await api.get(`/api/maturity/organizations/${organizationId}/roadmap`);
    return resp.data.data;
  },

  async getTrends(organizationId: string): Promise<Array<{ date: string; score: number }>> {
    const resp = await api.get(`/api/maturity/organizations/${organizationId}/trends`);
    return resp.data.data;
  },

  async getDimensions(organizationId: string): Promise<DimensionScore[]> {
    const resp = await api.get(`/api/maturity/organizations/${organizationId}/dimensions`);
    return resp.data.data;
  }
};

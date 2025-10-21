export interface HeatmapPoint {
  x: string;
  y: string;
  value: number;
}

export interface HeatmapData {
  type: string;
  title: string;
  xAxis: string[];
  yAxis: string[];
  data: HeatmapPoint[];
  generatedAt: string;
}

export interface RiskDistributionBucket {
  severity: string;
  count: number;
}

export interface TrendPoint {
  date: string;
  critical: number;
  high: number;
  medium: number;
  low: number;
}

import { api } from './api';

export const heatmapService = {
  async getHeatmap(organizationId: string): Promise<HeatmapData> {
    const resp = await api.get(`/api/heatmaps/organizations/${organizationId}`);
    return resp.data.data;
  },

  async getHotspots(organizationId: string): Promise<HeatmapPoint[]> {
    const resp = await api.get(`/api/heatmaps/organizations/${organizationId}/hotspots`);
    return resp.data.data;
  },

  async getRiskDistribution(organizationId: string): Promise<RiskDistributionBucket[]> {
    const resp = await api.get(`/api/heatmaps/organizations/${organizationId}/risk-distribution`);
    return resp.data.data;
  },

  async getTrends(organizationId: string): Promise<TrendPoint[]> {
    const resp = await api.get(`/api/heatmaps/organizations/${organizationId}/trends`);
    return resp.data.data;
  },

  async getRecommendations(organizationId: string): Promise<string[]> {
    const resp = await api.get(`/api/heatmaps/organizations/${organizationId}/recommendations`);
    return resp.data.data || [];
  }
};

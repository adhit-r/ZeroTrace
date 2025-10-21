import { api } from './api';

export interface ComplianceScore {
  framework: string;
  score: number;
  level: string;
}

export interface ComplianceFinding {
  id: string;
  controlId: string;
  severity: string;
  title: string;
  status: string;
}

export interface EvidenceItem {
  id: string;
  title: string;
  status: string;
}

export const complianceService = {
  async getReport(organizationId: string, framework: string, type = 'full', period = 'quarterly') {
    const resp = await api.get(`/api/compliance/organizations/${organizationId}/report`, {
      params: { framework, type, period }
    });
    return resp.data.data;
  },

  async getScore(organizationId: string) {
    const resp = await api.get(`/api/compliance/organizations/${organizationId}/score`);
    return resp.data.data as ComplianceScore;
  },

  async getFindings(organizationId: string) {
    const resp = await api.get(`/api/compliance/organizations/${organizationId}/findings`);
    return resp.data.data as ComplianceFinding[];
  },

  async getRecommendations(organizationId: string) {
    const resp = await api.get(`/api/compliance/organizations/${organizationId}/recommendations`);
    return resp.data.data as string[];
  },

  async getEvidence(organizationId: string) {
    const resp = await api.get(`/api/compliance/organizations/${organizationId}/evidence`);
    return resp.data.data as EvidenceItem[];
  },

  async getExecutiveSummary(organizationId: string) {
    const resp = await api.get(`/api/compliance/organizations/${organizationId}/executive-summary`);
    return resp.data.data;
  }
};

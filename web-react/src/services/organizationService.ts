import { api } from './api';

export interface OrganizationProfile {
  id: string;
  name: string;
  industry: string;
  size: string;
  riskTolerance: string;
  techStack: string[];
  complianceFrameworks: string[];
  securityPolicies: string[];
  riskWeights: {
    confidentiality: number;
    integrity: number;
    availability: number;
    compliance: number;
  };
}

export interface TechStackRelevance {
  technology: string;
  relevance: number;
  confidence: number;
  riskFactors: string[];
  recommendations: string[];
}

export interface IndustryRiskWeights {
  confidentiality: number;
  integrity: number;
  availability: number;
  compliance: number;
}

export const organizationService = {
  async getProfile(): Promise<OrganizationProfile | null> {
    try {
      const response = await api.get('/api/organizations/profile');
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch organization profile:', error);
      return null;
    }
  },

  async updateProfile(profile: OrganizationProfile): Promise<void> {
    try {
      await api.put('/api/organizations/profile', profile);
    } catch (error) {
      console.error('Failed to update organization profile:', error);
      throw error;
    }
  },

  async createProfile(profile: Omit<OrganizationProfile, 'id'>): Promise<OrganizationProfile> {
    try {
      const response = await api.post('/api/organizations/profile', profile);
      return response.data.data;
    } catch (error) {
      console.error('Failed to create organization profile:', error);
      throw error;
    }
  },

  async getTechStackRelevance(organizationId: string): Promise<TechStackRelevance[]> {
    try {
      const response = await api.get(`/api/organizations/${organizationId}/tech-stack/relevance`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch tech stack relevance:', error);
      return [];
    }
  },

  async getIndustryRiskWeights(organizationId: string): Promise<IndustryRiskWeights> {
    try {
      const response = await api.get(`/api/organizations/${organizationId}/risk-weights`);
      return response.data.data;
    } catch (error) {
      console.error('Failed to fetch industry risk weights:', error);
      return {
        confidentiality: 0.3,
        integrity: 0.3,
        availability: 0.2,
        compliance: 0.2
      };
    }
  }
};

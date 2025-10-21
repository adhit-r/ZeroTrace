/**
 * ZeroTrace API Response Types
 * Comprehensive type definitions for all API endpoints
 */

export interface APIResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: APIError;
  timestamp: string;
}

export interface APIError {
  code: string;
  message: string;
  details?: string;
}

// Dashboard API Types
export interface DashboardOverviewResponse {
  assets: {
    total: number;
    vulnerable: number;
    critical: number;
    high: number;
    medium: number;
    low: number;
    lastScan: string;
  };
  vulnerabilities: {
    total: number;
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  agents: {
    total: number;
    online: number;
  };
}

// Branch API Types
export interface Branch {
  id: string;
  name: string;
  location: string;
  type: 'headquarters' | 'branch' | 'datacenter' | 'cloud';
  status: 'active' | 'inactive' | 'maintenance';
  metrics: {
    totalAssets: number;
    criticalVulns: number;
    complianceScore: number;
    lastScan: string;
  };
  coordinates: { lat: number; lng: number };
  children?: Branch[];
}

export interface BranchListResponse {
  branches: Branch[];
  total: number;
  page: number;
  pageSize: number;
}

// Asset API Types
export interface Asset {
  id: string;
  hostname: string;
  ip: string;
  branch: string;
  location: string;
  owner: string;
  businessCriticality: 'critical' | 'high' | 'medium' | 'low';
  tags: string[];
  lastSeen: string;
  agentStatus: 'online' | 'offline' | 'maintenance';
  vulnerabilities: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  complianceScore: number;
  riskScore: number;
  suggestedFixes: number;
  metadata: {
    os: string;
    architecture: string;
    kernel: string;
    uptime: string;
    memory: string;
    cpu: string;
    disk: string;
  };
}

export interface AssetListResponse {
  assets: Asset[];
  total: number;
  page: number;
  pageSize: number;
  filters: {
    criticality: string[];
    status: string[];
    riskLevel: string[];
  };
}

// Vulnerability API Types
export interface Vulnerability {
  id: string;
  cve: string;
  title: string;
  description: string;
  severity: 'critical' | 'high' | 'medium' | 'low';
  cvss: number;
  exploitability: 'exploitable' | 'poc' | 'theoretical' | 'unknown';
  status: 'open' | 'patched' | 'ignored' | 'in_progress';
  published: string;
  suggestedFixes: string[];
  references: string[];
  affectedAssets: number;
  firstSeen: string;
  lastSeen: string;
}

export interface VulnerabilityListResponse {
  vulnerabilities: Vulnerability[];
  total: number;
  page: number;
  pageSize: number;
  filters: {
    severity: string[];
    status: string[];
    exploitability: string[];
  };
}

// Agent API Types
export interface Agent {
  id: string;
  name: string;
  version: string;
  hostname: string;
  os: string;
  status: 'online' | 'offline' | 'maintenance';
  lastSeen: string;
  cpuUsage: number;
  memoryUsage: number;
  ipAddress: string;
  metadata: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface AgentListResponse {
  agents: Agent[];
  total: number;
  online: number;
  offline: number;
}

// Scan API Types
export interface ScanResult {
  id: string;
  agentId: string;
  scanType: string;
  status: 'completed' | 'failed' | 'in_progress';
  startTime: string;
  endTime: string;
  vulnerabilitiesFound: number;
  assetsScanned: number;
  duration: string;
  metadata: Record<string, any>;
}

export interface ScanHistoryResponse {
  scans: ScanResult[];
  total: number;
  page: number;
  pageSize: number;
}

// Compliance API Types
export interface ComplianceFramework {
  id: string;
  name: string;
  version: string;
  description: string;
  controls: ComplianceControl[];
}

export interface ComplianceControl {
  id: string;
  name: string;
  description: string;
  category: string;
  status: 'compliant' | 'non_compliant' | 'not_applicable' | 'not_assessed';
  evidence: string[];
  lastAssessed: string;
}

export interface ComplianceScorecard {
  framework: string;
  overallScore: number;
  controls: {
    total: number;
    compliant: number;
    nonCompliant: number;
    notApplicable: number;
    notAssessed: number;
  };
  categories: Array<{
    name: string;
    score: number;
    controls: number;
  }>;
  lastUpdated: string;
}

// Risk API Types
export interface RiskAssessment {
  assetId: string;
  riskScore: number;
  factors: Array<{
    name: string;
    weight: number;
    score: number;
    description: string;
  }>;
  lastCalculated: string;
  trend: 'increasing' | 'decreasing' | 'stable';
}

// Patch API Types
export interface PatchRecommendation {
  id: string;
  vulnerabilityId: string;
  assetId: string;
  title: string;
  description: string;
  confidence: 'high' | 'medium' | 'low';
  rollbackSteps: string[];
  testingSteps: string[];
  estimatedTime: string;
  prerequisites: string[];
  status: 'pending' | 'testing' | 'deployed' | 'failed' | 'rolled_back';
}

export interface PatchQueueResponse {
  patches: PatchRecommendation[];
  total: number;
  pending: number;
  testing: number;
  deployed: number;
  failed: number;
}

// Alert API Types
export interface Alert {
  id: string;
  type: 'critical' | 'high' | 'medium' | 'low';
  title: string;
  description: string;
  timestamp: string;
  branch: string;
  asset: string;
  vulnerabilityId?: string;
  status: 'new' | 'acknowledged' | 'investigating' | 'resolved' | 'dismissed';
  assignedTo?: string;
  priority: 'urgent' | 'high' | 'medium' | 'low';
}

export interface AlertListResponse {
  alerts: Alert[];
  total: number;
  new: number;
  acknowledged: number;
  investigating: number;
  resolved: number;
}

// Report API Types
export interface Report {
  id: string;
  name: string;
  type: 'executive' | 'technical' | 'compliance' | 'custom';
  format: 'pdf' | 'excel' | 'csv' | 'json';
  status: 'generating' | 'completed' | 'failed';
  createdAt: string;
  downloadUrl?: string;
  parameters: Record<string, any>;
}

// User and Role API Types
export interface User {
  id: string;
  email: string;
  name: string;
  role: UserRole;
  companyId: string;
  branchId?: string;
  status: 'active' | 'inactive' | 'suspended';
  lastLogin: string;
  createdAt: string;
}

export type UserRole = 'global_ciso' | 'branch_ciso' | 'branch_it_manager' | 'security_analyst' | 'patch_engineer';

export interface UserPermissions {
  view: boolean;
  triage: boolean;
  remediate: boolean;
  export: boolean;
  scheduleScan: boolean;
  manageAgents: boolean;
  manageUsers: boolean;
  viewReports: boolean;
}

// Integration API Types
export interface Integration {
  id: string;
  name: string;
  type: 'jira' | 'servicenow' | 'slack' | 'teams' | 'aws' | 'azure' | 'gcp';
  status: 'active' | 'inactive' | 'error';
  configuration: Record<string, any>;
  lastSync: string;
  syncStatus: 'success' | 'failed' | 'in_progress';
}

// Analytics API Types
export interface AnalyticsData {
  timeRange: {
    start: string;
    end: string;
  };
  metrics: {
    totalAssets: number;
    totalVulnerabilities: number;
    criticalVulnerabilities: number;
    complianceScore: number;
    meanTimeToRemediate: number;
    scanCoverage: number;
  };
  trends: Array<{
    date: string;
    vulnerabilities: number;
    assets: number;
    compliance: number;
  }>;
  topRiskyAssets: Array<{
    assetId: string;
    hostname: string;
    riskScore: number;
    criticalVulns: number;
  }>;
  vulnerabilityDistribution: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
}

// Export API Types
export interface ExportRequest {
  format: 'csv' | 'excel' | 'pdf' | 'json';
  filters: {
    branches?: string[];
    assets?: string[];
    vulnerabilities?: string[];
    dateRange?: {
      start: string;
      end: string;
    };
  };
  fields: string[];
  includeMetadata: boolean;
}

export interface ExportResponse {
  jobId: string;
  status: 'queued' | 'processing' | 'completed' | 'failed';
  downloadUrl?: string;
  expiresAt: string;
}

// Webhook API Types
export interface Webhook {
  id: string;
  url: string;
  events: string[];
  secret: string;
  status: 'active' | 'inactive';
  lastTriggered: string;
  createdAt: string;
}

export interface WebhookPayload {
  event: string;
  timestamp: string;
  data: any;
  signature: string;
}

// Search API Types
export interface SearchRequest {
  query: string;
  filters: {
    type?: string[];
    severity?: string[];
    status?: string[];
    branch?: string[];
    dateRange?: {
      start: string;
      end: string;
    };
  };
  page: number;
  pageSize: number;
  sortBy: string;
  sortOrder: 'asc' | 'desc';
}

export interface SearchResponse {
  results: Array<{
    id: string;
    type: 'asset' | 'vulnerability' | 'branch' | 'agent';
    title: string;
    description: string;
    score: number;
    metadata: Record<string, any>;
  }>;
  total: number;
  page: number;
  pageSize: number;
  facets: Record<string, Array<{
    value: string;
    count: number;
  }>>;
}

// Bulk Action API Types
export interface BulkActionRequest {
  action: 'patch' | 'ignore' | 'assign' | 'export' | 'scan';
  targetIds: string[];
  parameters: Record<string, any>;
}

export interface BulkActionResponse {
  jobId: string;
  status: 'queued' | 'processing' | 'completed' | 'failed';
  results: Array<{
    id: string;
    status: 'success' | 'failed' | 'skipped';
    message?: string;
  }>;
}

// Configuration API Types
export interface SystemConfiguration {
  scanning: {
    defaultInterval: number;
    maxConcurrentScans: number;
    retentionPeriod: number;
  };
  notifications: {
    email: boolean;
    slack: boolean;
    webhook: boolean;
  };
  compliance: {
    frameworks: string[];
    assessmentInterval: number;
  };
  security: {
    passwordPolicy: Record<string, any>;
    sessionTimeout: number;
    mfaRequired: boolean;
  };
}

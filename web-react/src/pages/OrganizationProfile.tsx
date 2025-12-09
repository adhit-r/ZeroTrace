import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
// import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { organizationService } from '@/services/organizationService';

interface OrganizationProfile {
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

const INDUSTRIES = [
  'Technology', 'Healthcare', 'Finance', 'Government', 'Education',
  'Manufacturing', 'Retail', 'Energy', 'Telecommunications', 'Other'
];

const COMPANY_SIZES = [
  'Startup (1-10)', 'Small (11-50)', 'Medium (51-200)', 'Large (201-1000)', 'Enterprise (1000+)'
];

const RISK_TOLERANCE_LEVELS = [
  'Very Low', 'Low', 'Medium', 'High', 'Very High'
];

const COMPLIANCE_FRAMEWORKS = [
  'SOC2', 'ISO27001', 'PCI DSS', 'HIPAA', 'SOX', 'FedRAMP', 'FISMA', 'FERPA', 'GDPR'
];

const SECURITY_POLICIES = [
  'Password Policy', 'Access Control', 'Data Encryption', 'Incident Response',
  'Security Training', 'Vulnerability Management', 'Patch Management', 'Backup Policy'
];

const TECH_STACK_OPTIONS = [
  'AWS', 'Azure', 'GCP', 'Docker', 'Kubernetes', 'Terraform', 'Ansible',
  'Python', 'Java', 'Node.js', 'React', 'Angular', 'Vue.js', 'Go', 'Rust',
  'PostgreSQL', 'MySQL', 'MongoDB', 'Redis', 'Elasticsearch'
];

export default function OrganizationProfile() {
  const [profile, setProfile] = useState<OrganizationProfile>({
    id: '',
    name: '',
    industry: '',
    size: '',
    riskTolerance: '',
    techStack: [],
    complianceFrameworks: [],
    securityPolicies: [],
    riskWeights: {
      confidentiality: 0.3,
      integrity: 0.3,
      availability: 0.2,
      compliance: 0.2
    }
  });
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    loadProfile();
  }, []);

  const loadProfile = async () => {
    setLoading(true);
    try {
      const data = await organizationService.getProfile();
      if (data) {
        setProfile(data);
      }
    } catch (error) {
      console.error('Failed to load organization profile:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await organizationService.updateProfile(profile);
      // Show success message
    } catch (error) {
      console.error('Failed to save organization profile:', error);
    } finally {
      setSaving(false);
    }
  };

  const handleTechStackChange = (tech: string, checked: boolean) => {
    if (checked) {
      setProfile(prev => ({
        ...prev,
        techStack: [...prev.techStack, tech]
      }));
    } else {
      setProfile(prev => ({
        ...prev,
        techStack: prev.techStack.filter(t => t !== tech)
      }));
    }
  };

  const handleComplianceChange = (framework: string, checked: boolean) => {
    if (checked) {
      setProfile(prev => ({
        ...prev,
        complianceFrameworks: [...prev.complianceFrameworks, framework]
      }));
    } else {
      setProfile(prev => ({
        ...prev,
        complianceFrameworks: prev.complianceFrameworks.filter(f => f !== framework)
      }));
    }
  };

  const handlePolicyChange = (policy: string, checked: boolean) => {
    if (checked) {
      setProfile(prev => ({
        ...prev,
        securityPolicies: [...prev.securityPolicies, policy]
      }));
    } else {
      setProfile(prev => ({
        ...prev,
        securityPolicies: prev.securityPolicies.filter(p => p !== policy)
      }));
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 p-6">
        <div className="max-w-4xl mx-auto">
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
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Organization Profile</h1>
          <Button
            onClick={handleSave}
            disabled={saving}
          >
            {saving ? 'Saving...' : 'Save Profile'}
          </Button>
        </div>

        {/* Basic Information */}
        <Card>
          <CardHeader>
            <CardTitle className="text-xl font-bold">Basic Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label htmlFor="name">Organization Name</Label>
                <Input
                  id="name"
                  value={profile.name}
                  onChange={(e) => setProfile(prev => ({ ...prev, name: e.target.value }))}
                />
              </div>
              <div>
                <Label htmlFor="industry">Industry</Label>
                <Select
                  value={profile.industry}
                  onValueChange={(value) => setProfile(prev => ({ ...prev, industry: value }))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select industry" />
                  </SelectTrigger>
                  <SelectContent>
                    {INDUSTRIES.map(industry => (
                      <SelectItem key={industry} value={industry}>{industry}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label htmlFor="size">Company Size</Label>
                <Select
                  value={profile.size}
                  onValueChange={(value) => setProfile(prev => ({ ...prev, size: value }))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select company size" />
                  </SelectTrigger>
                  <SelectContent>
                    {COMPANY_SIZES.map(size => (
                      <SelectItem key={size} value={size}>{size}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div>
                <Label htmlFor="riskTolerance">Risk Tolerance</Label>
                <Select
                  value={profile.riskTolerance}
                  onValueChange={(value) => setProfile(prev => ({ ...prev, riskTolerance: value }))}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select risk tolerance" />
                  </SelectTrigger>
                  <SelectContent>
                    {RISK_TOLERANCE_LEVELS.map(level => (
                      <SelectItem key={level} value={level}>{level}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Tech Stack */}
        <Card>
          <CardHeader>
            <CardTitle className="text-xl font-bold">Technology Stack</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
              {TECH_STACK_OPTIONS.map(tech => (
                <label key={tech} className="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={profile.techStack.includes(tech)}
                    onChange={(e) => handleTechStackChange(tech, e.target.checked)}
                  />
                  <span className="text-sm font-medium">{tech}</span>
                </label>
              ))}
            </div>
            {profile.techStack.length > 0 && (
              <div className="mt-4">
                <Label>Selected Technologies:</Label>
                <div className="flex flex-wrap gap-2 mt-2">
                  {profile.techStack.map(tech => (
                    <Badge key={tech} variant="secondary">
                      {tech}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Compliance Frameworks */}
        <Card>
          <CardHeader>
            <CardTitle className="text-xl font-bold">Compliance Frameworks</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
              {COMPLIANCE_FRAMEWORKS.map(framework => (
                <label key={framework} className="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={profile.complianceFrameworks.includes(framework)}
                    onChange={(e) => handleComplianceChange(framework, e.target.checked)}
                  />
                  <span className="text-sm font-medium">{framework}</span>
                </label>
              ))}
            </div>
            {profile.complianceFrameworks.length > 0 && (
              <div className="mt-4">
                <Label>Selected Frameworks:</Label>
                <div className="flex flex-wrap gap-2 mt-2">
                  {profile.complianceFrameworks.map(framework => (
                    <Badge key={framework} variant="outline">
                      {framework}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Security Policies */}
        <Card>
          <CardHeader>
            <CardTitle className="text-xl font-bold">Security Policies</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {SECURITY_POLICIES.map(policy => (
                <label key={policy} className="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={profile.securityPolicies.includes(policy)}
                    onChange={(e) => handlePolicyChange(policy, e.target.checked)}
                  />
                  <span className="text-sm font-medium">{policy}</span>
                </label>
              ))}
            </div>
            {profile.securityPolicies.length > 0 && (
              <div className="mt-4">
                <Label>Selected Policies:</Label>
                <div className="flex flex-wrap gap-2 mt-2">
                  {profile.securityPolicies.map(policy => (
                    <Badge key={policy} variant="secondary">
                      {policy}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Risk Weights */}
        <Card>
          <CardHeader>
            <CardTitle className="text-xl font-bold">Risk Weight Configuration</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <Label htmlFor="confidentiality">Confidentiality Weight</Label>
                <Input
                  id="confidentiality"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                  value={profile.riskWeights.confidentiality}
                  onChange={(e) => setProfile(prev => ({
                    ...prev,
                    riskWeights: { ...prev.riskWeights, confidentiality: parseFloat(e.target.value) }
                  }))}
                />
              </div>
              <div>
                <Label htmlFor="integrity">Integrity Weight</Label>
                <Input
                  id="integrity"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                  value={profile.riskWeights.integrity}
                  onChange={(e) => setProfile(prev => ({
                    ...prev,
                    riskWeights: { ...prev.riskWeights, integrity: parseFloat(e.target.value) }
                  }))}
                />
              </div>
              <div>
                <Label htmlFor="availability">Availability Weight</Label>
                <Input
                  id="availability"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                  value={profile.riskWeights.availability}
                  onChange={(e) => setProfile(prev => ({
                    ...prev,
                    riskWeights: { ...prev.riskWeights, availability: parseFloat(e.target.value) }
                  }))}
                />
              </div>
              <div>
                <Label htmlFor="compliance">Compliance Weight</Label>
                <Input
                  id="compliance"
                  type="number"
                  min="0"
                  max="1"
                  step="0.1"
                  value={profile.riskWeights.compliance}
                  onChange={(e) => setProfile(prev => ({
                    ...prev,
                    riskWeights: { ...prev.riskWeights, compliance: parseFloat(e.target.value) }
                  }))}
                />
              </div>
            </div>
            <div className="p-4 bg-gray-100 rounded-lg">
              <p className="text-sm text-gray-600">
                <strong>Total Weight:</strong> {(profile.riskWeights.confidentiality +
                  profile.riskWeights.integrity +
                  profile.riskWeights.availability +
                  profile.riskWeights.compliance).toFixed(1)}
              </p>
              <p className="text-xs text-gray-500 mt-1">
                Weights should sum to 1.0 for optimal risk calculation
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

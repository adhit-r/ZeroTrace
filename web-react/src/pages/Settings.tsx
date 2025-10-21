import React, { useState, useEffect } from 'react';
import { Settings as SettingsIcon, Shield, Bell, User, Database, Building2, Target, Layers } from 'lucide-react';

interface OrganizationProfile {
  id?: string;
  organization_id: string;
  industry: string;
  risk_tolerance: 'CONSERVATIVE' | 'MODERATE' | 'AGGRESSIVE';
  tech_stack: {
    languages: string[];
    frameworks: string[];
    databases: string[];
    cloud_providers: string[];
    operating_systems: string[];
    containers: string[];
    dev_tools: string[];
    security_tools: string[];
  };
  compliance_frameworks: string[];
  security_policies: Record<string, any>;
  risk_weights: Record<string, any>;
}

const Settings: React.FC = () => {
  const [organizationProfile, setOrganizationProfile] = useState<OrganizationProfile | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [loading, setLoading] = useState(false);

  // Mock organization ID - in real app, this would come from auth context
  const organizationId = '123e4567-e89b-12d3-a456-426614174000';

  useEffect(() => {
    loadOrganizationProfile();
  }, []);

  const loadOrganizationProfile = async () => {
    try {
      setLoading(true);
      // Mock API call - replace with actual API call
      const mockProfile: OrganizationProfile = {
        organization_id: organizationId,
        industry: 'Technology',
        risk_tolerance: 'MODERATE',
        tech_stack: {
          languages: ['JavaScript', 'Python', 'Go'],
          frameworks: ['React', 'Django', 'Express.js'],
          databases: ['PostgreSQL', 'Redis'],
          cloud_providers: ['AWS', 'Google Cloud'],
          operating_systems: ['Linux', 'macOS'],
          containers: ['Docker', 'Kubernetes'],
          dev_tools: ['Git', 'Jenkins', 'VS Code'],
          security_tools: ['Nessus', 'OWASP ZAP']
        },
        compliance_frameworks: ['SOC2', 'ISO27001'],
        security_policies: {
          patch_management: {
            critical_patches: 'immediate',
            high_patches: '24_hours'
          }
        },
        risk_weights: {
          critical: 1.0,
          high: 0.8,
          medium: 0.6
        }
      };
      setOrganizationProfile(mockProfile);
    } catch (error) {
      console.error('Failed to load organization profile:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveProfile = async () => {
    if (!organizationProfile) return;
    
    try {
      setLoading(true);
      // Mock API call - replace with actual API call
      console.log('Saving organization profile:', organizationProfile);
      setIsEditing(false);
    } catch (error) {
      console.error('Failed to save organization profile:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleTechStackChange = (category: keyof OrganizationProfile['tech_stack'], value: string[]) => {
    if (!organizationProfile) return;
    setOrganizationProfile({
      ...organizationProfile,
      tech_stack: {
        ...organizationProfile.tech_stack,
        [category]: value
      }
    });
  };

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-4xl font-black text-black uppercase tracking-wider mb-2">SYSTEM SETTINGS</h1>
          <p className="text-lg text-gray-600 font-bold">CONFIGURE YOUR ZEROTRACE ENVIRONMENT</p>
        </div>

        {/* Organization Profile Section */}
        {organizationProfile && (
          <div className="mb-8">
            <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center">
                  <Building2 className="w-8 h-8 text-orange-500 mr-3" />
                  <div>
                    <h2 className="text-xl font-black text-black uppercase tracking-wider">ORGANIZATION PROFILE</h2>
                    <p className="text-sm text-gray-600">Configure your organization's security posture and technology stack</p>
                  </div>
                </div>
                <div className="flex gap-2">
                  {isEditing ? (
                    <>
                      <button
                        onClick={handleSaveProfile}
                        disabled={loading}
                        className="px-4 py-2 bg-green-500 text-white font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out disabled:opacity-50"
                      >
                        {loading ? 'SAVING...' : 'SAVE'}
                      </button>
                      <button
                        onClick={() => setIsEditing(false)}
                        className="px-4 py-2 bg-gray-500 text-white font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out"
                      >
                        CANCEL
                      </button>
                    </>
                  ) : (
                    <button
                      onClick={() => setIsEditing(true)}
                      className="px-4 py-2 bg-orange-500 text-white font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out"
                    >
                      EDIT PROFILE
                    </button>
                  )}
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {/* Industry & Risk Tolerance */}
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">INDUSTRY</label>
                    {isEditing ? (
                      <select
                        value={organizationProfile.industry}
                        onChange={(e) => setOrganizationProfile({...organizationProfile, industry: e.target.value})}
                        className="w-full px-3 py-2 border-3 border-black rounded font-bold"
                      >
                        <option value="Technology">Technology</option>
                        <option value="Healthcare">Healthcare</option>
                        <option value="Finance">Finance</option>
                        <option value="Government">Government</option>
                        <option value="Education">Education</option>
                        <option value="Manufacturing">Manufacturing</option>
                      </select>
                    ) : (
                      <div className="px-3 py-2 bg-gray-100 border-3 border-black rounded font-bold">{organizationProfile.industry}</div>
                    )}
                  </div>
                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">RISK TOLERANCE</label>
                    {isEditing ? (
                      <select
                        value={organizationProfile.risk_tolerance}
                        onChange={(e) => setOrganizationProfile({...organizationProfile, risk_tolerance: e.target.value as any})}
                        className="w-full px-3 py-2 border-3 border-black rounded font-bold"
                      >
                        <option value="CONSERVATIVE">Conservative</option>
                        <option value="MODERATE">Moderate</option>
                        <option value="AGGRESSIVE">Aggressive</option>
                      </select>
                    ) : (
                      <div className="px-3 py-2 bg-gray-100 border-3 border-black rounded font-bold">{organizationProfile.risk_tolerance}</div>
                    )}
                  </div>
                </div>

                {/* Technology Stack */}
                <div className="space-y-4">
                  <div className="flex items-center mb-3">
                    <Layers className="w-5 h-5 text-orange-500 mr-2" />
                    <h3 className="font-bold text-black uppercase">TECH STACK</h3>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">LANGUAGES</label>
                    {isEditing ? (
                      <input
                        type="text"
                        value={organizationProfile.tech_stack.languages.join(', ')}
                        onChange={(e) => handleTechStackChange('languages', e.target.value.split(',').map(s => s.trim()).filter(s => s))}
                        className="w-full px-3 py-2 border-3 border-black rounded font-bold"
                        placeholder="JavaScript, Python, Go"
                      />
                    ) : (
                      <div className="px-3 py-2 bg-gray-100 border-3 border-black rounded font-bold">
                        {organizationProfile.tech_stack.languages.join(', ')}
                      </div>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">FRAMEWORKS</label>
                    {isEditing ? (
                      <input
                        type="text"
                        value={organizationProfile.tech_stack.frameworks.join(', ')}
                        onChange={(e) => handleTechStackChange('frameworks', e.target.value.split(',').map(s => s.trim()).filter(s => s))}
                        className="w-full px-3 py-2 border-3 border-black rounded font-bold"
                        placeholder="React, Django, Express.js"
                      />
                    ) : (
                      <div className="px-3 py-2 bg-gray-100 border-3 border-black rounded font-bold">
                        {organizationProfile.tech_stack.frameworks.join(', ')}
                      </div>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">DATABASES</label>
                    {isEditing ? (
                      <input
                        type="text"
                        value={organizationProfile.tech_stack.databases.join(', ')}
                        onChange={(e) => handleTechStackChange('databases', e.target.value.split(',').map(s => s.trim()).filter(s => s))}
                        className="w-full px-3 py-2 border-3 border-black rounded font-bold"
                        placeholder="PostgreSQL, Redis"
                      />
                    ) : (
                      <div className="px-3 py-2 bg-gray-100 border-3 border-black rounded font-bold">
                        {organizationProfile.tech_stack.databases.join(', ')}
                      </div>
                    )}
                  </div>
                </div>

                {/* Compliance & Security */}
                <div className="space-y-4">
                  <div className="flex items-center mb-3">
                    <Target className="w-5 h-5 text-orange-500 mr-2" />
                    <h3 className="font-bold text-black uppercase">COMPLIANCE</h3>
                  </div>
                  
                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">FRAMEWORKS</label>
                    {isEditing ? (
                      <input
                        type="text"
                        value={organizationProfile.compliance_frameworks.join(', ')}
                        onChange={(e) => setOrganizationProfile({...organizationProfile, compliance_frameworks: e.target.value.split(',').map(s => s.trim()).filter(s => s)})}
                        className="w-full px-3 py-2 border-3 border-black rounded font-bold"
                        placeholder="SOC2, ISO27001"
                      />
                    ) : (
                      <div className="px-3 py-2 bg-gray-100 border-3 border-black rounded font-bold">
                        {organizationProfile.compliance_frameworks.join(', ')}
                      </div>
                    )}
                  </div>

                  <div>
                    <label className="block text-sm font-bold text-black uppercase mb-2">RISK WEIGHTS</label>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span className="font-bold">Critical:</span>
                        <span className="text-red-500 font-bold">{organizationProfile.risk_weights.critical || 1.0}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="font-bold">High:</span>
                        <span className="text-orange-500 font-bold">{organizationProfile.risk_weights.high || 0.8}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="font-bold">Medium:</span>
                        <span className="text-yellow-500 font-bold">{organizationProfile.risk_weights.medium || 0.6}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {/* Security Settings */}
          <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
            <div className="flex items-center mb-4">
              <Shield className="w-8 h-8 text-orange-500 mr-3" />
              <h2 className="text-xl font-black text-black uppercase tracking-wider">SECURITY</h2>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="font-bold text-black">Auto-scan enabled</span>
                <div className="w-12 h-6 bg-orange-500 border-2 border-black rounded-full relative">
                  <div className="w-5 h-5 bg-white border-2 border-black rounded-full absolute right-0.5 top-0.5"></div>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="font-bold text-black">Threat detection</span>
                <div className="w-12 h-6 bg-orange-500 border-2 border-black rounded-full relative">
                  <div className="w-5 h-5 bg-white border-2 border-black rounded-full absolute right-0.5 top-0.5"></div>
                </div>
              </div>
            </div>
          </div>

          {/* Notifications */}
          <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
            <div className="flex items-center mb-4">
              <Bell className="w-8 h-8 text-orange-500 mr-3" />
              <h2 className="text-xl font-black text-black uppercase tracking-wider">NOTIFICATIONS</h2>
            </div>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="font-bold text-black">Email alerts</span>
                <div className="w-12 h-6 bg-orange-500 border-2 border-black rounded-full relative">
                  <div className="w-5 h-5 bg-white border-2 border-black rounded-full absolute right-0.5 top-0.5"></div>
                </div>
              </div>
              <div className="flex items-center justify-between">
                <span className="font-bold text-black">Critical alerts</span>
                <div className="w-12 h-6 bg-orange-500 border-2 border-black rounded-full relative">
                  <div className="w-5 h-5 bg-white border-2 border-black rounded-full absolute right-0.5 top-0.5"></div>
                </div>
              </div>
            </div>
          </div>

          {/* User Management */}
          <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
            <div className="flex items-center mb-4">
              <User className="w-8 h-8 text-orange-500 mr-3" />
              <h2 className="text-xl font-black text-black uppercase tracking-wider">USERS</h2>
            </div>
            <div className="space-y-4">
              <div className="text-center">
                <div className="text-3xl font-black text-orange-500 mb-2">12</div>
                <div className="text-sm font-bold text-black uppercase">ACTIVE USERS</div>
              </div>
              <button className="w-full px-4 py-2 bg-orange-500 text-white font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out">
                MANAGE USERS
              </button>
            </div>
          </div>

          {/* Database Settings */}
          <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6">
            <div className="flex items-center mb-4">
              <Database className="w-8 h-8 text-orange-500 mr-3" />
              <h2 className="text-xl font-black text-black uppercase tracking-wider">DATABASE</h2>
            </div>
            <div className="space-y-4">
              <div className="text-center">
                <div className="text-2xl font-black text-green-500 mb-2">ONLINE</div>
                <div className="text-sm font-bold text-black uppercase">CONNECTION STATUS</div>
              </div>
              <button className="w-full px-4 py-2 bg-green-500 text-white font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out">
                TEST CONNECTION
              </button>
            </div>
          </div>

          {/* System Info */}
          <div className="bg-white border-3 border-black rounded-lg shadow-neubrutalist p-6 md:col-span-2">
            <div className="flex items-center mb-4">
              <SettingsIcon className="w-8 h-8 text-orange-500 mr-3" />
              <h2 className="text-xl font-black text-black uppercase tracking-wider">SYSTEM INFORMATION</h2>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="font-bold text-black">Version:</span>
                  <span className="text-orange-500 font-bold">v1.0.0</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-bold text-black">Uptime:</span>
                  <span className="text-green-500 font-bold">99.9%</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-bold text-black">Last Backup:</span>
                  <span className="text-blue-500 font-bold">2 hours ago</span>
                </div>
              </div>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="font-bold text-black">Agents:</span>
                  <span className="text-orange-500 font-bold">24 active</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-bold text-black">Scans:</span>
                  <span className="text-green-500 font-bold">1,247 completed</span>
                </div>
                <div className="flex justify-between">
                  <span className="font-bold text-black">Vulnerabilities:</span>
                  <span className="text-red-500 font-bold">156 found</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Settings;
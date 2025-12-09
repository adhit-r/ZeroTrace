import { api } from './api';

export interface Application {
  id: string;
  name: string;
  version: string;
  vendor: string;
  type: string;
  classification?: string;
  path?: string;
  agentId: string;
  agentName: string;
  vulnerabilities: number;
  riskLevel: 'critical' | 'high' | 'medium' | 'low' | 'safe';
  status: 'vulnerable' | 'safe' | 'unknown';
  vulnerabilityDetails?: VulnerabilityDetail[];
  installedAt?: string;
  lastUpdated?: string;
}

export interface VulnerabilityDetail {
  id: string;
  cve?: string;
  severity: string;
  title: string;
  description: string;
  score?: number;
}

export interface ApplicationStats {
  total: number;
  vulnerable: number;
  safe: number;
  byRiskLevel: {
    critical: number;
    high: number;
    medium: number;
    low: number;
    safe: number;
  };
  byClassification: Record<string, number>;
  byAgent: Record<string, number>;
}

// Helper function to infer vendor from application name and path
function inferVendorFromName(appName: string, appPath?: string): string {
  const nameLower = appName.toLowerCase();
  const pathLower = (appPath || '').toLowerCase();

  // Extract vendor from path if it's a macOS app
  if (pathLower.includes('/applications/')) {
    // For macOS apps in /Applications, try to infer from common patterns
    // Apple apps are usually in /System/Applications or have specific names
    if (pathLower.includes('/system/applications/') ||
      pathLower.includes('/system/library/')) {
      return 'Apple';
    }
  }

  // Common vendor mappings (expanded list)
  const vendorMap: Record<string, string> = {
    // Browsers
    'chrome': 'Google',
    'google chrome': 'Google',
    'chromium': 'Chromium',
    'firefox': 'Mozilla',
    'safari': 'Apple',
    'edge': 'Microsoft',
    'microsoft edge': 'Microsoft',
    'opera': 'Opera',
    'brave': 'Brave',
    'vivaldi': 'Vivaldi',

    // Adobe
    'adobe': 'Adobe',
    'photoshop': 'Adobe',
    'illustrator': 'Adobe',
    'premiere': 'Adobe',
    'after effects': 'Adobe',
    'acrobat': 'Adobe',
    'indesign': 'Adobe',
    'lightroom': 'Adobe',
    'xd': 'Adobe',

    // Apple
    'xcode': 'Apple',
    'finder': 'Apple',
    'mail': 'Apple',
    'messages': 'Apple',
    'calendar': 'Apple',
    'notes': 'Apple',
    // 'safari': 'Apple', // Duplicate
    'keynote': 'Apple',
    'pages': 'Apple',
    'numbers': 'Apple',
    'garageband': 'Apple',
    'logic': 'Apple',
    'final cut': 'Apple',
    'imovie': 'Apple',
    'quicktime': 'Apple',
    'itunes': 'Apple',
    'music': 'Apple',
    'tv': 'Apple',
    'podcasts': 'Apple',
    'books': 'Apple',
    'news': 'Apple',
    'stocks': 'Apple',
    'maps': 'Apple',
    'reminders': 'Apple',
    'contacts': 'Apple',
    'facetime': 'Apple',
    'photo booth': 'Apple',
    'preview': 'Apple',
    'textedit': 'Apple',
    'calculator': 'Apple',
    'system preferences': 'Apple',
    'system settings': 'Apple',

    // Microsoft
    // Microsoft
    'vscode': 'Microsoft',
    'visual studio code': 'Microsoft',
    'visual studio': 'Microsoft',
    'skype': 'Microsoft',
    'office': 'Microsoft',
    'word': 'Microsoft',
    'excel': 'Microsoft',
    'powerpoint': 'Microsoft',
    'outlook': 'Microsoft',
    'onenote': 'Microsoft',
    'teams': 'Microsoft',
    'azure': 'Microsoft',

    // Development Tools
    'intellij': 'JetBrains',
    'pycharm': 'JetBrains',
    'webstorm': 'JetBrains',
    'goland': 'JetBrains',
    'clion': 'JetBrains',
    'rider': 'JetBrains',
    'phpstorm': 'JetBrains',
    'rubymine': 'JetBrains',
    'docker': 'Docker',
    'docker desktop': 'Docker',
    'kubectl': 'Kubernetes',
    'kubernetes': 'Kubernetes',
    'git': 'Git',
    'github': 'GitHub',
    'gitlab': 'GitLab',
    'node': 'Node.js',
    'npm': 'Node.js',
    'yarn': 'Yarn',
    'go': 'Go Team',
    'golang': 'Go Team',
    'java': 'Oracle',
    'jdk': 'Oracle',
    'jre': 'Oracle',
    'maven': 'Apache',
    'gradle': 'Gradle',
    'sublime': 'Sublime HQ',
    'atom': 'GitHub',
    'vim': 'Bram Moolenaar',
    'emacs': 'GNU',

    // Media
    'vlc': 'VideoLAN',
    'spotify': 'Spotify',
    // Media
    // 'vlc': 'VideoLAN', // Duplicate
    // 'spotify': 'Spotify', // Duplicate
    'vlc media player': 'VideoLAN',
    'handbrake': 'HandBrake',
    'ffmpeg': 'FFmpeg',

    // Communication
    'slack': 'Slack',
    'zoom': 'Zoom',
    'discord': 'Discord',
    'telegram': 'Telegram',
    'whatsapp': 'WhatsApp',
    'messenger': 'Meta',
    'facebook': 'Meta',
    'instagram': 'Meta',

    // Utilities
    '7zip': '7-Zip',
    'notepad++': 'Notepad++',
    'the unarchiver': 'The Unarchiver',
    'cleanmymac': 'MacPaw',
    'little snitch': 'Objective Development',
    'bartender': 'Surtees Studios',
    'alfred': 'Alfred',
    'raycast': 'Raycast',
    '1password': '1Password',
    'lastpass': 'LastPass',
    'bitwarden': 'Bitwarden',

    // Other
    'dropbox': 'Dropbox',
    'google drive': 'Google',
    'onedrive': 'Microsoft',
    'icloud': 'Apple',
  };

  // Check for exact matches first
  for (const [key, vendor] of Object.entries(vendorMap)) {
    if (nameLower === key || nameLower.includes(key)) {
      return vendor;
    }
  }

  // Check for common patterns in name
  if (nameLower.includes('microsoft')) return 'Microsoft';
  if (nameLower.includes('google')) return 'Google';
  if (nameLower.includes('apple')) return 'Apple';
  if (nameLower.includes('adobe')) return 'Adobe';
  if (nameLower.includes('jetbrains')) return 'JetBrains';
  if (nameLower.includes('oracle')) return 'Oracle';
  if (nameLower.includes('apache')) return 'Apache';
  if (nameLower.includes('meta')) return 'Meta';
  if (nameLower.includes('facebook')) return 'Meta';
  if (nameLower.includes('amazon')) return 'Amazon';
  if (nameLower.includes('aws')) return 'Amazon';

  // For macOS apps, if path contains /Applications and no vendor found, try to infer from bundle structure
  if (pathLower.includes('/applications/') && !pathLower.includes('/system/')) {
    // Could be third-party app, but we don't have enough info
    // Return 'Unknown' for now
  }

  return 'Unknown';
}

export const applicationService = {
  async getAllApplications(): Promise<Application[]> {
    try {
      const response = await api.get('/api/agents/');
      const agents = response.data?.data || [];

      const allApps: Application[] = [];

      agents.forEach((agent: any) => {
        // Skip if agent is invalid
        if (!agent || !agent.id) return;

        // Try multiple possible metadata keys for dependencies
        const dependencies = agent.metadata?.dependencies ||
          agent.metadata?.installed_apps ||
          agent.metadata?.applications ||
          [];
        const vulnerabilities = agent.metadata?.vulnerabilities || [];
        const enrichedVulns = agent.metadata?.enriched_vulnerabilities || [];

        // Handle both array of objects and array of strings
        const depsArray = Array.isArray(dependencies) ? dependencies : [];

        // Skip if no dependencies
        if (depsArray.length === 0) return;

        depsArray.forEach((dep: any) => {
          // Skip if dep is null or not an object
          if (!dep || typeof dep !== 'object') return;

          // Extract app information with better vendor detection
          const appName = dep.name || dep.package_name || dep.Name || '';
          if (!appName) return; // Skip if no app name

          // Filter out system binaries and command-line tools that aren't real applications
          const systemBinaries = ['python', 'python3', 'python2', 'go', 'golang', 'node', 'npm', 'yarn', 'git', 'bash', 'zsh', 'sh', 'curl', 'wget', 'ssh', 'scp', 'rsync', 'tar', 'gzip', 'zip'];
          const nameLower = appName.toLowerCase();
          const appPath = dep.path || dep.Path || dep.install_path || '';
          const pathLower = (appPath || '').toLowerCase();

          // Only skip if it's a bare binary name without any app-like structure
          // If it has a path suggesting it's an actual app (like .app bundle or Program Files), keep it
          if (systemBinaries.some(binary => nameLower === binary || nameLower === `${binary}.exe`)) {
            // Skip bare system binaries that aren't in app directories
            if (!appPath ||
              (!pathLower.includes('/applications/') &&
                !pathLower.includes('\\program files\\') &&
                !pathLower.includes('.app') &&
                !pathLower.includes('/usr/local/bin/') &&
                !pathLower.includes('/opt/'))) {
              return; // Skip bare system binaries
            }
          }

          const appVersion = dep.version || dep.Version || 'unknown';
          const appType = dep.type || dep.Type || 'application';

          // Better vendor detection for macOS apps
          let vendor = dep.vendor || dep.Vendor || dep.publisher || dep.Publisher || '';

          // For macOS apps, try to extract vendor from CFBundleIdentifier or path
          if (!vendor || vendor === 'Unknown') {
            if (appType === 'macos_app' && appPath) {
              // Extract vendor from bundle identifier (e.g., com.apple.Safari -> Apple)
              const bundleId = dep.bundle_id || dep.BundleID || dep.cf_bundle_identifier || '';
              if (bundleId) {
                const parts = bundleId.split('.');
                if (parts.length >= 2) {
                  // Usually format is com.vendor.appname
                  const vendorPart = parts[1];
                  vendor = vendorPart.charAt(0).toUpperCase() + vendorPart.slice(1);
                }
              }

              // If still unknown, try to infer from app name and path
              if (!vendor || vendor === 'Unknown') {
                vendor = inferVendorFromName(appName, appPath);
              }
            } else {
              vendor = inferVendorFromName(appName, appPath);
            }
          }

          // Find vulnerabilities for this application (with null checks)
          const appVulns = (Array.isArray(vulnerabilities) ? vulnerabilities : []).filter((v: any) =>
            v && (
              v.package_name === appName ||
              v.package_name === dep.package_name ||
              v.name === appName ||
              (v.cve_id && dep.cve_ids && Array.isArray(dep.cve_ids) && dep.cve_ids.includes(v.cve_id))
            )
          );

          // Get enriched vulnerability details (with null checks)
          const vulnDetails: VulnerabilityDetail[] = (Array.isArray(enrichedVulns) ? enrichedVulns : [])
            .filter((v: any) =>
              v && (
                v.package_name === appName ||
                v.package_name === dep.package_name ||
                v.name === appName
              )
            )
            .map((v: any) => ({
              id: v.id || v.cve_id || '',
              cve: v.cve_id || '',
              severity: v.severity || 'unknown',
              title: v.title || v.name || 'Unknown',
              description: v.description || '',
              score: v.cvss_score || v.score || 0,
            }));

          // Determine risk level based on vulnerabilities
          let riskLevel: Application['riskLevel'] = 'safe';
          if (appVulns.length > 0 || vulnDetails.length > 0) {
            const allVulns = [...appVulns, ...vulnDetails];
            const hasCritical = allVulns.some((v: any) =>
              (v.severity && v.severity.toLowerCase() === 'critical') ||
              (v.score && v.score >= 9) ||
              (v.Severity && v.Severity.toLowerCase() === 'critical')
            );
            const hasHigh = allVulns.some((v: any) =>
              (v.severity && v.severity.toLowerCase() === 'high') ||
              (v.score && v.score >= 7 && v.score < 9) ||
              (v.Severity && v.Severity.toLowerCase() === 'high')
            );
            const hasMedium = allVulns.some((v: any) =>
              (v.severity && v.severity.toLowerCase() === 'medium') ||
              (v.score && v.score >= 4 && v.score < 7) ||
              (v.Severity && v.Severity.toLowerCase() === 'medium')
            );

            if (hasCritical) riskLevel = 'critical';
            else if (hasHigh) riskLevel = 'high';
            else if (hasMedium) riskLevel = 'medium';
            else riskLevel = 'low';
          }

          allApps.push({
            id: `${agent.id}-${appName}-${appVersion}`,
            name: appName || 'Unknown',
            version: appVersion,
            vendor: vendor || 'Unknown',
            type: appType,
            classification: dep.classification || dep.Classification || appType,
            path: appPath,
            agentId: agent.id,
            agentName: agent.name || agent.hostname || agent.Hostname || 'Unknown Agent',
            vulnerabilities: appVulns.length + vulnDetails.length,
            riskLevel,
            status: (appVulns.length > 0 || vulnDetails.length > 0) ? 'vulnerable' : 'safe',
            vulnerabilityDetails: vulnDetails,
            installedAt: dep.installed_at || dep.install_date || dep.InstallDate,
            lastUpdated: dep.last_updated || dep.update_date || dep.LastUpdated,
          });
        });
      });

      // Deduplicate by name+version+agent
      const uniqueApps = allApps.filter((app, index, self) =>
        index === self.findIndex((a) =>
          a.name === app.name &&
          a.version === app.version &&
          a.agentId === app.agentId
        )
      );

      return uniqueApps;
    } catch (error) {
      console.error('Failed to fetch applications:', error);
      return [];
    }
  },

  async getApplicationsByAgent(agentId: string): Promise<Application[]> {
    const allApps = await this.getAllApplications();
    return allApps.filter(app => app.agentId === agentId);
  },

  async getApplicationsByRiskLevel(riskLevel: Application['riskLevel']): Promise<Application[]> {
    const allApps = await this.getAllApplications();
    return allApps.filter(app => app.riskLevel === riskLevel);
  },

  async getApplicationsByClassification(classification: string): Promise<Application[]> {
    const allApps = await this.getAllApplications();
    return allApps.filter(app => app.classification === classification);
  },

  async getApplicationStats(): Promise<ApplicationStats> {
    const allApps = await this.getAllApplications();

    const stats: ApplicationStats = {
      total: allApps.length,
      vulnerable: allApps.filter(a => a.status === 'vulnerable').length,
      safe: allApps.filter(a => a.status === 'safe').length,
      byRiskLevel: {
        critical: allApps.filter(a => a.riskLevel === 'critical').length,
        high: allApps.filter(a => a.riskLevel === 'high').length,
        medium: allApps.filter(a => a.riskLevel === 'medium').length,
        low: allApps.filter(a => a.riskLevel === 'low').length,
        safe: allApps.filter(a => a.riskLevel === 'safe').length,
      },
      byClassification: {},
      byAgent: {},
    };

    // Count by classification
    allApps.forEach(app => {
      const classification = app.classification || 'unknown';
      stats.byClassification[classification] = (stats.byClassification[classification] || 0) + 1;
    });

    // Count by agent
    allApps.forEach(app => {
      stats.byAgent[app.agentId] = (stats.byAgent[app.agentId] || 0) + 1;
    });

    return stats;
  },

  async getVulnerableApplications(): Promise<Application[]> {
    const allApps = await this.getAllApplications();
    return allApps.filter(app => app.status === 'vulnerable');
  },
};


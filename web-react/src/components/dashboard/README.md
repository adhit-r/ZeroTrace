# ZeroTrace Dashboard Components

A comprehensive dashboard UI overhaul for ZeroTrace, a Qualys/Tenable alternative VMDR-style security product. This implementation provides enterprise-grade security management with role-based access control and hierarchical asset visibility.

## 🎯 Overview

ZeroTrace is designed for enterprise security teams to discover, prioritize, and remediate vulnerabilities across distributed branches and cloud/on-prem assets. The dashboard provides:

- **Metric-oriented, data-driven interface**
- **Hierarchical asset visibility**: Company → Branch → Assets
- **Role-aware views** for different user types
- **Multi-branch comparison and prioritization**
- **Quick remediation workflows**

## 👥 Target Users

- **Global CISO**: Enterprise-wide security oversight and risk management
- **Branch CISO**: Branch-level security management and compliance oversight  
- **Branch IT Manager**: Day-to-day asset management and vulnerability remediation
- **Security Analyst**: Threat investigation and vulnerability analysis
- **Patch Engineer**: Patch management and deployment coordination

## 🏗️ Architecture

### Core Components

1. **KPIRibbon** - Configurable metrics dashboard with sparklines and trends
2. **BranchSelector** - Hierarchical branch navigation with risk indicators
3. **AssetInventory** - Comprehensive asset management with filtering and bulk actions
4. **RiskHeatmap** - Geographic risk distribution visualization
5. **AssetDetail** - Detailed asset view with vulnerability management
6. **RoleBasedDashboard** - Context-aware dashboard layouts per user role
7. **GlobalDashboard** - Main dashboard orchestrator

### Design System Integration

All components are built using the ZeroTrace design system (`zerotrace-design-system.json`) with:

- **Neubrutalist design language** with bold borders and shadows
- **Consistent color palette** (black, white, orange accents)
- **Typography**: Space Grotesk font family
- **Interactive elements** with hover states and animations
- **Responsive layouts** for desktop, tablet, and mobile

## 🎨 Design System Mapping

### Color Tokens
```typescript
// Light theme colors
background: 'rgba(255, 255, 255, 1)'
foreground: 'rgba(0, 0, 0, 1)'
primary: 'rgba(0, 0, 0, 1)'
secondary: 'rgba(255, 107, 0, 1)'
destructive: 'rgba(239, 68, 68, 1)'
border: 'rgba(0, 0, 0, 1)'
```

### Component Styles
```typescript
// Button styles
button: {
  default: 'bg-black text-white border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]'
  secondary: 'bg-orange-500 text-white border-3 border-black'
  destructive: 'bg-red-500 text-white border-3 border-black'
}

// Card styles  
card: 'bg-white border-3 border-black shadow-[8px_8px_0px_0px_rgba(0,0,0,1)] rounded'

// Input styles
input: 'w-full h-11 border-3 border-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]'
```

## 📊 Key Features

### KPI Ribbon
- **Configurable metrics**: Active Critical CVEs, MTTR, Compliance %, Scan Coverage %
- **Trend indicators**: Sparklines and delta comparisons
- **Color-coded severity**: Critical (red), High (orange), Medium (yellow), Low (green)
- **Interactive tooltips** with detailed explanations

### Branch Selector
- **Hierarchical navigation**: Headquarters → Branches → Sub-locations
- **Risk indicators**: Visual risk scores and critical vulnerability counts
- **Search functionality**: Quick branch discovery
- **Role-based filtering**: Show only accessible branches

### Asset Inventory
- **Comprehensive filtering**: By criticality, status, risk level, tags
- **Bulk actions**: Patch, ignore, assign, export, scan
- **Sortable columns**: Risk score, vulnerabilities, last scan
- **Real-time updates**: Live asset status

### Risk Heatmap
- **Geographic visualization**: World map with risk indicators
- **Grid/Map toggle**: Alternative view modes
- **Risk filtering**: Filter by risk level (critical, high, medium, low)
- **Interactive markers**: Click to drill down to branch details

### Asset Detail
- **Comprehensive asset information**: System specs, network info, scan history
- **Vulnerability management**: Detailed CVE information with suggested fixes
- **Bulk vulnerability actions**: Patch, ignore, export
- **Network topology**: Interface and port information

## 🔐 Role-Based Access Control

### Global CISO
- **Enterprise-wide visibility**: All branches and assets
- **Risk heatmap**: Geographic risk distribution
- **Compliance scorecards**: Framework-based assessments
- **Executive reporting**: High-level metrics and trends

### Branch CISO  
- **Branch-level oversight**: Focused on specific location
- **Asset inventory**: Detailed asset management
- **Vulnerability trends**: Historical analysis
- **Compliance status**: Branch-specific compliance tracking

### Branch IT Manager
- **Operational focus**: Day-to-day asset management
- **Patch queue**: Prioritized remediation tasks
- **Scan status**: Agent health and scan scheduling
- **Ticket management**: Integration with ITSM tools

### Security Analyst
- **Threat investigation**: Deep-dive vulnerability analysis
- **Attack surface monitoring**: External exposure assessment
- **Threat intelligence**: External signals and trending CVEs
- **Investigation tools**: Forensic analysis capabilities

### Patch Engineer
- **Patch management**: Deployment coordination and scheduling
- **Testing workflows**: Patch validation and rollback planning
- **Deployment tracking**: Progress monitoring and reporting
- **Risk assessment**: Patch impact analysis

## 🚀 Getting Started

### Installation
```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

### Usage
```typescript
import { RoleBasedDashboard } from './components/dashboard/RoleBasedDashboard';

// Render dashboard with user role
<RoleBasedDashboard userRole="global_ciso" />
```

### Configuration
```typescript
// Configure KPI metrics
const kpiMetrics = [
  {
    id: 'active_critical_cves',
    label: 'Active Critical CVEs',
    value: 47,
    delta: { value: 12, trend: 'up', period: 'last 7 days' },
    sparkline: [45, 42, 38, 41, 44, 47, 43],
    color: 'critical'
  }
];

// Configure branch data
const branches = [
  {
    id: 'hq-nyc',
    name: 'Headquarters NYC',
    location: 'New York, NY',
    type: 'headquarters',
    status: 'active',
    metrics: {
      totalAssets: 450,
      criticalVulns: 12,
      complianceScore: 92,
      lastScan: '2025-01-09T10:30:00Z'
    }
  }
];
```

## 📱 Responsive Design

### Breakpoints
- **Mobile**: 640px and below
- **Tablet**: 768px - 1024px  
- **Desktop**: 1024px and above

### Mobile Optimizations
- **Progressive disclosure**: Key KPIs first, details on demand
- **Touch-friendly**: Larger tap targets and swipe gestures
- **Simplified navigation**: Collapsible menus and quick actions
- **Offline support**: Cached data for critical metrics

## 🔧 API Integration

### Endpoints
```typescript
// Dashboard overview
GET /api/dashboard/overview

// Branch management
GET /api/branches
GET /api/branches/{id}/assets

// Asset management  
GET /api/assets
GET /api/assets/{id}
POST /api/assets/bulk-action

// Vulnerability management
GET /api/vulnerabilities
GET /api/vulnerabilities/{id}
POST /api/vulnerabilities/bulk-action
```

### Data Models
```typescript
interface Asset {
  id: string;
  hostname: string;
  ip: string;
  branch: string;
  businessCriticality: 'critical' | 'high' | 'medium' | 'low';
  riskScore: number;
  vulnerabilities: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  complianceScore: number;
  agentStatus: 'online' | 'offline' | 'maintenance';
}
```

## 🧪 Testing

### Component Testing
```bash
# Run component tests
npm run test:components

# Run with coverage
npm run test:coverage
```

### Integration Testing
```bash
# Run integration tests
npm run test:integration

# Run E2E tests
npm run test:e2e
```

## 📈 Performance

### Optimization Strategies
- **Virtual scrolling**: For large asset lists
- **Lazy loading**: Component-level code splitting
- **Memoization**: React.memo for expensive components
- **Debounced search**: Optimized filtering and search
- **Caching**: API response caching with React Query

### Metrics
- **First Contentful Paint**: < 1.5s
- **Largest Contentful Paint**: < 2.5s
- **Cumulative Layout Shift**: < 0.1
- **Time to Interactive**: < 3.5s

## 🔒 Security

### Data Protection
- **Role-based access**: Granular permissions per user type
- **Data encryption**: TLS 1.3 for all API communications
- **Input validation**: XSS and injection prevention
- **CSRF protection**: Token-based request validation

### Privacy
- **Data minimization**: Only collect necessary information
- **Retention policies**: Automatic data cleanup
- **Audit logging**: Complete user action tracking
- **GDPR compliance**: Data subject rights support

## 🚀 Deployment

### Production Build
```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Variables
```bash
# API Configuration
VITE_API_BASE_URL=https://api.zerotrace.com
VITE_API_VERSION=v1

# Authentication
VITE_AUTH_PROVIDER=okta
VITE_AUTH_DOMAIN=zerotrace.okta.com

# Feature Flags
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_DEBUG=false
```

## 📚 Documentation

### Component Documentation
- **Storybook**: Interactive component documentation
- **Props API**: TypeScript interfaces and examples
- **Usage examples**: Real-world implementation patterns

### API Documentation
- **OpenAPI Spec**: Complete API specification
- **Postman Collection**: Ready-to-use API examples
- **SDK Examples**: JavaScript/TypeScript integration

## 🤝 Contributing

### Development Setup
```bash
# Clone repository
git clone https://github.com/zerotrace/dashboard.git

# Install dependencies
npm install

# Start development server
npm run dev

# Run tests
npm test
```

### Code Standards
- **TypeScript**: Strict type checking enabled
- **ESLint**: Airbnb configuration
- **Prettier**: Consistent code formatting
- **Husky**: Pre-commit hooks for quality gates

## 📄 License

MIT License - see LICENSE file for details.

## 🆘 Support

- **Documentation**: [docs.zerotrace.com](https://docs.zerotrace.com)
- **Issues**: [GitHub Issues](https://github.com/zerotrace/dashboard/issues)
- **Discord**: [ZeroTrace Community](https://discord.gg/zerotrace)
- **Email**: support@zerotrace.com

---

Built with ❤️ by the ZeroTrace team

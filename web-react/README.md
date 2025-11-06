# ZeroTrace Frontend

Modern React frontend for ZeroTrace vulnerability detection and management platform.

## Overview

ZeroTrace Frontend is a high-performance React application built with:

- **React 19.1.1** with TypeScript
- **Vite 7.1.3** for fast development and optimized builds
- **Tailwind CSS 3.4.17** with custom design system
- **shadcn/ui** components for consistent, accessible interfaces
- **React Router 7.8.2** for navigation
- **React Query** for data fetching and caching
- **Clerk** for authentication and multi-organization support

## Features

### Dashboard
- Real-time vulnerability monitoring
- Interactive charts and visualizations
- Risk heatmap visualization
- Security maturity scoring
- Compliance reporting

### Components
- Comprehensive dashboard components
- Vulnerability management interface
- Agent monitoring interface
- Compliance dashboard
- Security analytics interface

### User Interface
- Responsive design with Tailwind CSS
- Dark theme with terminal-inspired design
- Accessible components with shadcn/ui
- Real-time updates with WebSocket support

## Quick Start

### Prerequisites

- Node.js 24.8.0+ (or Bun)
- Bun (recommended) or npm
- Git

### Installation

```bash
# Clone repository
git clone https://github.com/adhit-r/ZeroTrace.git
cd ZeroTrace/web-react

# Install dependencies with Bun (recommended)
bun install

# Or with npm
npm install
```

### Configuration

```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your configuration
# Required: VITE_API_URL
# Optional: VITE_CLERK_PUBLISHABLE_KEY for Clerk authentication
```

### Development

```bash
# Start development server with Bun
bun run dev

# Or with npm
npm run dev

# Access at http://localhost:3000
```

### Building

```bash
# Build for production with Bun
bun run build

# Or with npm
npm run build

# Preview production build
bun run preview
```

## Configuration

### Environment Variables

See `.env.example` for all available configuration options.

#### Required

- `VITE_API_URL`: ZeroTrace API endpoint (default: http://localhost:8080)

#### Optional

- `VITE_CLERK_PUBLISHABLE_KEY`: Clerk publishable key for authentication
- `VITE_ENRICHMENT_URL`: Enrichment service URL
- `VITE_ENABLE_ANALYTICS`: Enable analytics (default: false)
- `VITE_ENABLE_DEBUG`: Enable debug mode (default: false)
- `VITE_ENV`: Environment (development/production)

### Configuration Files

- `.env.example`: Environment variable template
- `vite.config.ts`: Vite configuration
- `tailwind.config.js`: Tailwind CSS configuration
- `tsconfig.json`: TypeScript configuration
- `playwright.config.ts`: Playwright test configuration

## Project Structure

```
web-react/
├── src/
│   ├── components/          # React components
│   │   ├── dashboard/      # Dashboard components
│   │   ├── ui/             # shadcn/ui components
│   │   └── Layout.tsx      # Layout components
│   ├── pages/              # Route pages
│   │   ├── Dashboard.tsx
│   │   ├── Vulnerabilities.tsx
│   │   ├── Agents.tsx
│   │   └── ...
│   ├── services/           # API integration layer
│   │   ├── api.ts
│   │   ├── agentService.ts
│   │   ├── dashboardService.ts
│   │   └── ...
│   ├── contexts/          # React contexts
│   │   ├── AuthContext.tsx
│   │   └── ...
│   ├── styles/            # Global styles
│   │   ├── zerotrace-theme.css
│   │   ├── neobrutal.css
│   │   └── index.css
│   ├── types/             # TypeScript types
│   │   └── api.ts
│   ├── lib/               # Utility functions
│   │   └── utils.ts
│   ├── App.tsx            # Main application
│   └── main.tsx           # Entry point
├── tests/                  # Playwright tests
│   └── frontend-analysis.spec.ts
├── public/                # Static assets
├── dist/                  # Build output
├── package.json
├── vite.config.ts
├── tailwind.config.js
└── README.md
```

## Development

### Available Scripts

```bash
# Development
bun run dev          # Start development server
bun run build        # Build for production
bun run preview      # Preview production build

# Testing
bun run test         # Run tests (if configured)
npx playwright test  # Run Playwright E2E tests

# Linting
bun run lint         # Run ESLint
```

### Code Examples

**Example: Fetching Dashboard Data**
```typescript
import { useQuery } from '@tanstack/react-query';
import { dashboardService } from './services/dashboardService';

function Dashboard() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['dashboard', 'overview'],
    queryFn: () => dashboardService.getOverview(),
    refetchInterval: 30000 // Refresh every 30 seconds
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h1>Dashboard</h1>
      <p>Total Assets: {data.assets.total}</p>
      <p>Vulnerabilities: {data.vulnerabilities.total}</p>
    </div>
  );
}
```

**Example: Agent Service Integration**
```typescript
import { agentService } from './services/agentService';

// Get all agents
const agents = await agentService.getAgents();

// Get agent statistics
const stats = await agentService.getStats();

// Send agent heartbeat
await agentService.sendHeartbeat({
  agent_id: 'agent-123',
  organization_id: 'org-456',
  status: 'active',
  cpu_usage: 45.5,
  memory_usage: 62.3
});
```

**Example: Vulnerability Management**
```typescript
import { vulnerabilityService } from './services/vulnerabilityService';

// Get vulnerabilities with filters
const vulnerabilities = await vulnerabilityService.getVulnerabilities({
  severity: 'high',
  category: 'network',
  page: 1,
  page_size: 20
});

// Get vulnerability statistics
const stats = await vulnerabilityService.getStats();

// Export vulnerabilities
const exportData = await vulnerabilityService.export({
  format: 'csv',
  severity: 'high'
});
```

### Adding Dependencies

```bash
# Add dependency with Bun
bun add package-name

# Add dev dependency
bun add -d package-name

# Or with npm
npm install package-name
npm install -D package-name
```

### Code Style

- TypeScript for type safety
- ESLint for linting
- Prettier for formatting
- Functional components with hooks
- React Query for data fetching

## Testing

### Playwright E2E Tests

```bash
# Run all tests
npx playwright test

# Run tests in UI mode
npx playwright test --ui

# Run specific test
npx playwright test tests/frontend-analysis.spec.ts

# View test report
npx playwright show-report
```

### Test Coverage

- E2E tests with Playwright
- Component tests (if configured)
- Integration tests (if configured)

## Deployment

### Docker

```bash
# Build image
docker build -t zerotrace-web .

# Run container
docker run -p 3000:3000 --env-file .env zerotrace-web
```

### Docker Compose

```bash
# Start with docker-compose
docker-compose up web
```

### Production Build

```bash
# Build for production
bun run build

# Output will be in dist/ directory
# Serve with nginx or any static file server
```

### Nginx Configuration

See `nginx.conf` for production Nginx configuration.

## Performance

### Optimizations

- Code splitting with React.lazy
- Tree shaking and minification
- CSS optimization with Tailwind
- Image optimization
- Lazy loading components

### Performance Metrics

- **First Contentful Paint**: < 1.5s
- **Largest Contentful Paint**: < 2.5s
- **Time to Interactive**: < 3.5s
- **Cumulative Layout Shift**: < 0.1

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Documentation

- [API Documentation](../docs/api-v2-documentation.md)
- [Frontend Technology Analysis](../docs/frontend-technology-analysis.md)
- [Web Implementation](../docs/web-implementation.md)
- [Architecture Documentation](../docs/architecture.md)

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.

---

**Last Updated**: January 2025

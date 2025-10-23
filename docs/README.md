# ZeroTrace Documentation

Welcome to the ZeroTrace documentation. This comprehensive guide covers all aspects of the ZeroTrace security platform, from API usage to compliance frameworks.

## 📚 Documentation Index

### Core Documentation
- [API v2 Documentation](./api-v2-documentation.md) - Complete API reference
- [Scanner Modules Documentation](./scanner-modules-documentation.md) - Security scanner modules
- [Compliance Frameworks Documentation](./compliance-frameworks-documentation.md) - Compliance and regulatory frameworks
- [Database Schema v2](./database-schema-v2.md) - Database schema and migrations
- [Architecture Overview](./architecture.md) - System architecture and design

### Development Documentation
- [Development Setup](./development-setup.md) - Getting started with development
- [Performance Monitoring](./performance/monitoring-setup.md) - Performance monitoring setup
- [Network Discovery Implementation](./network-discovery-implementation.md) - Network scanning implementation
- [Scalable Data Processing](./scalable-data-processing.md) - Data processing architecture

### Component Documentation
- [Go Agent Documentation](./go-agent.md) - Agent implementation details
- [Go API Documentation](./go-api.md) - API server implementation
- [Python Enrichment Documentation](./python-enrichment.md) - Enrichment service details
- [Frontend Technology Analysis](./frontend-technology-analysis.md) - Frontend architecture
- [Web Implementation](./web-implementation.md) - Web interface implementation

### Security Documentation
- [Agent CVE Documentation](./agent-cve.md) - CVE detection and analysis
- [Monitoring Strategy](./monitoring-strategy.md) - Security monitoring approach
- [Organization Prioritization Design](./org-prioritization-design.md) - Risk prioritization

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Python 3.9+
- Node.js 18+
- PostgreSQL 13+
- Redis 6+

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/zerotrace/zerotrace.git
   cd zerotrace
   ```

2. **Set up the database**
   ```bash
   ./scripts/setup-database.sh
   ```

3. **Start the services**
   ```bash
   # Start all services
   docker-compose up -d
   
   # Or start individually
   # API Server
   cd api-go && go run cmd/api/main.go
   
   # Enrichment Service
   cd enrichment-python && python -m uvicorn app.main:app --port 5001
   
   # Frontend
   cd web-react && bun dev
   ```

4. **Deploy agents**
   ```bash
   # Download and run agent
   chmod +x zerotrace-agent
   ./zerotrace-agent
   ```

### First Steps

1. **Access the dashboard**: http://localhost:5173
2. **View API documentation**: http://localhost:8080/api/v2/docs
3. **Check agent status**: http://localhost:8080/api/agents
4. **Monitor vulnerabilities**: http://localhost:8080/api/vulnerabilities

## 🔧 API Reference

### Base URL
```
http://localhost:8080/api/v2
```

### Authentication
```bash
curl -H "Authorization: Bearer YOUR_API_KEY" \
     -H "Content-Type: application/json" \
     http://localhost:8080/api/v2/vulnerabilities
```

### Key Endpoints

#### Vulnerabilities
- `GET /api/v2/vulnerabilities` - List vulnerabilities
- `GET /api/v2/vulnerabilities/stats` - Vulnerability statistics
- `GET /api/v2/vulnerabilities/export` - Export vulnerabilities

#### Compliance
- `GET /api/v2/compliance/status` - Compliance status
- `GET /api/v2/compliance/gaps` - Compliance gaps

#### Network Scanning
- `POST /api/v2/network/scan` - Initiate network scan
- `GET /api/v2/network/scan/{id}/status` - Scan status
- `GET /api/v2/network/scan/{id}/results` - Scan results

#### Agents
- `GET /api/v2/agents/processing-status` - Agent processing status

## 🛡️ Security Categories

ZeroTrace supports comprehensive security scanning across 11 categories:

### 1. Network Security
- Port scanning and service detection
- SSL/TLS analysis
- Network topology mapping
- Protocol security assessment

### 2. Compliance
- CIS Benchmarks
- PCI-DSS compliance
- HIPAA compliance
- GDPR compliance
- SOC 2 compliance
- ISO 27001 compliance

### 3. System Vulnerabilities
- OS patch analysis
- Kernel vulnerability detection
- Driver issue identification
- End-of-life software detection

### 4. Authentication
- Password policy analysis
- Account security assessment
- Privilege escalation detection
- Multi-factor authentication analysis

### 5. Database Security
- Multi-database vulnerability scanning
- Access control assessment
- Encryption status verification
- Performance security analysis

### 6. API Security
- REST API security analysis
- GraphQL security assessment
- OWASP API Top 10 detection
- Shadow API discovery

### 7. Container Security
- Docker daemon configuration
- Container escape detection
- Kubernetes RBAC analysis
- IaC template scanning

### 8. AI/ML Security
- Model vulnerability assessment
- Training data security analysis
- LLM application security
- Adversarial attack detection

### 9. IoT/OT Security
- Device discovery
- Firmware analysis
- Protocol security assessment
- Industrial control system assessment

### 10. Privacy & Compliance
- PII detection
- GDPR/CCPA compliance
- Data retention analysis
- Consent mechanism verification

### 11. Web3 Security
- Smart contract vulnerability detection
- Wallet security analysis
- DApp security assessment
- DeFi protocol analysis

## 📊 Compliance Frameworks

### Supported Frameworks
- **CIS Benchmarks**: Security configuration guidelines
- **PCI-DSS**: Payment card industry standards
- **HIPAA**: Healthcare compliance
- **GDPR**: Data protection regulations
- **SOC 2**: Service organization controls
- **ISO 27001**: Information security management

### Compliance Features
- Automated compliance assessment
- Gap analysis and remediation
- Real-time compliance monitoring
- Compliance reporting and dashboards
- Automated remediation workflows

## 🔍 Scanner Modules

### Core Scanners
- **System Scanner**: OS and hardware information
- **Software Scanner**: Application discovery and analysis
- **Network Scanner**: Network security assessment
- **Configuration Scanner**: Compliance framework checks
- **System Vulnerability Scanner**: OS and kernel vulnerabilities
- **Authentication Scanner**: Access control assessment
- **Database Scanner**: Multi-database security
- **API Scanner**: API security assessment
- **Container Scanner**: Container and Kubernetes security
- **AI/ML Scanner**: AI and machine learning security
- **IoT/OT Scanner**: IoT and operational technology
- **Privacy Scanner**: Privacy and data protection
- **Web3 Scanner**: Web3 and blockchain security

### Scanner Configuration
```go
config := ScannerConfig{
    Timeout:       5 * time.Minute,
    ParallelScans: 10,
    MaxRetries:    3,
    LogLevel:      "info",
    EnableCaching: true,
}
```

## 🎨 Frontend Dashboards

### Comprehensive Security Dashboard
- Multi-category security overview
- Interactive category selection
- Compliance status monitoring
- Risk heatmap visualization
- Category distribution charts

### Individual Category Dashboards
- Category-specific metrics and KPIs
- Expandable sections (Overview, Trends, Actions)
- Real-time data visualization
- Interactive charts and graphs
- Filtering and search capabilities

### Network Security Dashboard
- Network topology visualization
- Host and service discovery
- SSL/TLS analysis
- Port scanning results
- Service vulnerability mapping

### Compliance Dashboard
- Framework selection and scoring
- Gap analysis and remediation
- Compliance trend monitoring
- Automated compliance workflows
- Compliance reporting

## 🧪 Testing

### Test Suite
```bash
# Run all tests
./scripts/run-tests.sh

# Run specific test suites
./scripts/run-tests.sh --go-only
./scripts/run-tests.sh --python-only
./scripts/run-tests.sh --frontend-only
./scripts/run-tests.sh --security-only
./scripts/run-tests.sh --load-tests
```

### Test Coverage
- **Unit Tests**: Individual component testing
- **Integration Tests**: Component interaction testing
- **Performance Tests**: Load and stress testing
- **Security Tests**: Security vulnerability testing
- **E2E Tests**: End-to-end workflow testing

## 📈 Performance Monitoring

### Metrics Collection
- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboards
- **Alerting**: Real-time alerts and notifications

### Key Metrics
- **API Performance**: Response times, throughput
- **Agent Performance**: Scan duration, resource usage
- **Enrichment Performance**: Processing latency, accuracy
- **Database Performance**: Query performance, connections

## 🔧 Configuration

### Environment Variables
```bash
# API Configuration
API_PORT=8080
API_HOST=0.0.0.0
API_KEY=your-api-key

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=zerotrace
DB_PASSWORD=your-password

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your-password

# Enrichment Service
ENRICHMENT_SERVICE_URL=http://localhost:5001
CVE_DATA_PATH=/data/cve_data.json
```

### Configuration Files
- `api-go/.env` - API server configuration
- `enrichment-python/.env` - Enrichment service configuration
- `web-react/.env` - Frontend configuration
- `docker-compose.yml` - Docker services configuration

## 🚀 Deployment

### Docker Deployment
```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d api
docker-compose up -d enrichment
docker-compose up -d frontend
```

### Production Deployment
```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Deploy to production
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods
kubectl get services
```

## 🔒 Security

### Security Features
- **Authentication**: JWT-based authentication
- **Authorization**: Role-based access control
- **Encryption**: Data encryption at rest and in transit
- **Audit Logging**: Comprehensive audit trails
- **Vulnerability Scanning**: Continuous security assessment

### Security Best Practices
- Regular security updates
- Secure configuration management
- Network segmentation
- Access control implementation
- Security monitoring and alerting



## 📄 License

ZeroTrace is licensed under the MIT License. See [LICENSE](../LICENSE) for details.

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

### Development Setup
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

### Code Style
- **Go**: Use `gofmt` and `golint`
- **Python**: Use `black` and `flake8`
- **JavaScript**: Use `prettier` and `eslint`


# Installation Guide

This guide will walk you through installing and setting up ZeroTrace for different environments.

## 📋 **Prerequisites**

### **System Requirements**
- **OS**: Linux, macOS, or Windows
- **CPU**: 2+ cores (4+ recommended)
- **Memory**: 8GB+ RAM (16GB+ recommended)
- **Storage**: 50GB+ available space
- **Network**: Internet connection for initial setup

### **Software Requirements**
- **Docker**: 20.10+ and Docker Compose 2.0+
- **Git**: 2.30+
- **Node.js**: 18+ (for development)
- **Go**: 1.21+ (for development)
- **Python**: 3.9+ (for development)

## 🚀 **Quick Installation**

### **1. Clone the Repository**
```bash
git clone https://github.com/radhi1991/ZeroTrace.git
cd ZeroTrace
```

### **2. Environment Setup**
```bash
# Copy environment files
cp api-go/env.example api-go/.env
cp agent-go/env.example agent-go/.env
cp enrichment-python/env.example enrichment-python/.env
```

### **3. Start Services**
```bash
# Start all services with Docker Compose
docker-compose up -d

# Or start individual services
docker-compose up -d postgres redis
docker-compose up -d api enrichment
docker-compose up -d frontend
```

### **4. Verify Installation**
```bash
# Check service status
docker-compose ps

# Test API health
curl http://localhost:8080/api/v1/health

# Test frontend
open http://localhost:3000
```

## 🔧 **Detailed Installation**

### **Development Environment**

#### **1. Backend Setup (Go API)**
```bash
cd api-go

# Install dependencies
go mod download

# Set up database
go run cmd/api/main.go migrate

# Start development server
go run cmd/api/main.go
```

#### **2. Enrichment Service (Python)**
```bash
cd enrichment-python

# Create virtual environment
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Start development server
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

#### **3. Frontend Setup (React)**
```bash
cd web-react

# Install dependencies
npm install

# Start development server
npm run dev
```

#### **4. Agent Setup (Go)**
```bash
cd agent-go

# Install dependencies
go mod download

# Build agent
go build -o zerotrace-agent cmd/agent/main.go

# Run agent
./zerotrace-agent
```

### **Production Environment**

#### **1. Docker Deployment**
```bash
# Build production images
docker-compose -f docker-compose.prod.yml build

# Start production services
docker-compose -f docker-compose.prod.yml up -d
```

#### **2. Kubernetes Deployment**
```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -n zerotrace
```

#### **3. Cloud Deployment**
```bash
# AWS ECS
aws ecs create-cluster --cluster-name zerotrace
aws ecs create-service --cluster zerotrace --service-name api --task-definition api

# Google Cloud Run
gcloud run deploy zerotrace-api --source api-go/
gcloud run deploy zerotrace-enrichment --source enrichment-python/
```

## ⚙️ **Configuration**

### **Environment Variables**

#### **API Configuration**
```bash
# api-go/.env
DATABASE_URL=postgresql://user:password@localhost:5432/zerotrace
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
API_PORT=8080
LOG_LEVEL=info
```

#### **Enrichment Configuration**
```bash
# enrichment-python/.env
DATABASE_URL=postgresql://user:password@localhost:5432/zerotrace
REDIS_URL=redis://localhost:6379
API_PORT=8000
LOG_LEVEL=info
NVD_API_KEY=your-nvd-api-key
```

#### **Agent Configuration**
```bash
# agent-go/.env
API_URL=http://localhost:8080
ENROLLMENT_TOKEN=your-enrollment-token
ORGANIZATION_ID=your-org-id
SCAN_INTERVAL=24h
HEARTBEAT_INTERVAL=5m
```

### **Database Setup**
```sql
-- Create database
CREATE DATABASE zerotrace;

-- Create user
CREATE USER zerotrace_user WITH PASSWORD 'secure_password';

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE zerotrace TO zerotrace_user;
```

## 🔍 **Verification**

### **Health Checks**
```bash
# API Health
curl http://localhost:8080/api/v1/health

# Enrichment Health
curl http://localhost:8000/health

# Frontend Health
curl http://localhost:3000

# Database Health
docker exec -it zerotrace-postgres psql -U zerotrace_user -d zerotrace -c "SELECT 1;"
```

### **Performance Tests**
```bash
# API Performance
ab -n 1000 -c 10 http://localhost:8080/api/v1/health

# Enrichment Performance
python tests/performance_test.py

# Agent Performance
go test ./agent-go/tests/ -bench=.
```

## 🚨 **Troubleshooting**

### **Common Issues**

#### **1. Database Connection Issues**
```bash
# Check database status
docker-compose logs postgres

# Reset database
docker-compose down
docker volume rm zerotrace_postgres_data
docker-compose up -d postgres
```

#### **2. Redis Connection Issues**
```bash
# Check Redis status
docker-compose logs redis

# Test Redis connection
docker exec -it zerotrace-redis redis-cli ping
```

#### **3. Port Conflicts**
```bash
# Check port usage
lsof -i :8080
lsof -i :8000
lsof -i :3000

# Change ports in docker-compose.yml
```

#### **4. Permission Issues**
```bash
# Fix file permissions
chmod +x scripts/*.sh
chmod 600 *.env

# Fix Docker permissions
sudo usermod -aG docker $USER
```

### **Logs and Debugging**
```bash
# View all logs
docker-compose logs

# View specific service logs
docker-compose logs api
docker-compose logs enrichment
docker-compose logs frontend

# Follow logs in real-time
docker-compose logs -f api
```

## 📚 **Next Steps**

After successful installation:

1. **Configure Monitoring**: Set up Prometheus and Grafana
2. **Set Up Alerting**: Configure AlertManager
3. **Deploy Agents**: Install agents on target systems
4. **Configure Organizations**: Set up company accounts
5. **Test Workflows**: Run end-to-end tests

## 📞 **Support**

If you encounter issues:

1. Check the [Troubleshooting Guide](Troubleshooting)
2. Search [GitHub Issues](https://github.com/radhi1991/ZeroTrace/issues)
3. Ask in [GitHub Discussions](https://github.com/radhi1991/ZeroTrace/discussions)
4. Review the [FAQ](FAQ)

---

**Last Updated**: January 2024  
**Version**: 1.0.0

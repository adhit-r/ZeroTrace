# Contributing Guidelines

Thank you for your interest in contributing to ZeroTrace! This document provides guidelines for contributing to the project.

## 🤝 **How to Contribute**

### **1. Fork and Clone**
```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/ZeroTrace.git
cd ZeroTrace

# Add upstream remote
git remote add upstream https://github.com/radhi1991/ZeroTrace.git
```

### **2. Create a Branch**
```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Or create a bug fix branch
git checkout -b fix/your-bug-fix-name
```

### **3. Make Changes**
- Follow the coding standards
- Write tests for new features
- Update documentation
- Ensure all tests pass

### **4. Commit Changes**
```bash
# Add your changes
git add .

# Commit with a descriptive message
git commit -m "feat: Add new feature description"

# Push to your fork
git push origin feature/your-feature-name
```

### **5. Create a Pull Request**
- Use the [Pull Request Template](.github/pull_request_template.md)
- Link to relevant issues
- Provide clear description of changes
- Include tests and documentation updates

## 📋 **Issue Reporting**

### **Before Creating an Issue**
1. Check existing issues and discussions
2. Search the documentation
3. Try to reproduce the issue
4. Gather relevant information

### **Issue Templates**
- [Bug Report](.github/ISSUE_TEMPLATE/bug_report.md)
- [Feature Request](.github/ISSUE_TEMPLATE/feature_request.md)
- [Performance Issue](.github/ISSUE_TEMPLATE/performance_issue.md)

### **Issue Guidelines**
- Use clear, descriptive titles
- Provide detailed descriptions
- Include steps to reproduce
- Add relevant logs and screenshots
- Specify environment details

## 🏗️ **Development Setup**

### **Prerequisites**
- Go 1.21+
- Python 3.9+
- Node.js 18+
- Docker & Docker Compose
- Git

### **Local Development**
```bash
# Clone repository
git clone https://github.com/radhi1991/ZeroTrace.git
cd ZeroTrace

# Set up environment
cp api-go/env.example api-go/.env
cp agent-go/env.example agent-go/.env
cp enrichment-python/env.example enrichment-python/.env

# Start services
docker-compose up -d postgres redis

# Start individual services
cd api-go && go run cmd/api/main.go
cd ../enrichment-python && uvicorn app.main:app --reload
cd ../web-react && npm run dev
cd ../agent-go && go run cmd/agent/main.go
```

### **Testing**
```bash
# Run all tests
go test ./...
python -m pytest enrichment-python/tests/
npm test

# Run specific tests
go test ./api-go/internal/handlers/ -v
python -m pytest enrichment-python/tests/test_enrichment.py -v
npm test -- --testNamePattern="API"
```

## 📝 **Coding Standards**

### **Go (Backend & Agent)**
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Write tests for all new code
- Use meaningful variable and function names
- Add comments for complex logic

```go
// Example Go code
func (s *Service) ProcessData(data []byte) (*Result, error) {
    if len(data) == 0 {
        return nil, errors.New("empty data")
    }
    
    result := &Result{
        ProcessedAt: time.Now(),
        Data:       data,
    }
    
    return result, nil
}
```

### **Python (Enrichment Service)**
- Use Black for formatting
- Follow PEP 8 style guide
- Write type hints
- Add docstrings for functions
- Use async/await for I/O operations

```python
# Example Python code
from typing import List, Optional
from dataclasses import dataclass

@dataclass
class EnrichmentResult:
    """Result of enrichment process."""
    app_id: str
    vulnerabilities: List[dict]
    processed_at: datetime

async def enrich_app_data(app_data: dict) -> Optional[EnrichmentResult]:
    """Enrich application data with CVE information."""
    if not app_data:
        return None
    
    vulnerabilities = await fetch_cve_data(app_data)
    
    return EnrichmentResult(
        app_id=app_data["id"],
        vulnerabilities=vulnerabilities,
        processed_at=datetime.utcnow()
    )
```

### **TypeScript/React (Frontend)**
- Use ESLint and Prettier
- Follow React best practices
- Use TypeScript for type safety
- Write unit tests with Jest
- Use functional components with hooks

```typescript
// Example React component
import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';

interface DashboardProps {
  organizationId: string;
}

export const Dashboard: React.FC<DashboardProps> = ({ organizationId }) => {
  const [selectedPeriod, setSelectedPeriod] = useState('7d');
  
  const { data, isLoading, error } = useQuery({
    queryKey: ['dashboard', organizationId, selectedPeriod],
    queryFn: () => fetchDashboardData(organizationId, selectedPeriod),
  });
  
  if (isLoading) return <LoadingSpinner />;
  if (error) return <ErrorMessage error={error} />;
  
  return (
    <div className="dashboard">
      <DashboardHeader period={selectedPeriod} onPeriodChange={setSelectedPeriod} />
      <DashboardMetrics data={data} />
    </div>
  );
};
```

## 🧪 **Testing Guidelines**

### **Test Coverage**
- Aim for 80%+ test coverage
- Write unit tests for all functions
- Write integration tests for APIs
- Write end-to-end tests for critical flows

### **Test Structure**
```bash
# Go tests
api-go/
├── internal/
│   ├── handlers/
│   │   ├── auth_test.go
│   │   └── scan_test.go
│   └── services/
│       └── auth_test.go

# Python tests
enrichment-python/
├── tests/
│   ├── test_enrichment.py
│   ├── test_api.py
│   └── conftest.py

# React tests
web-react/
├── src/
│   ├── components/
│   │   └── __tests__/
│   └── pages/
│       └── __tests__/
```

### **Running Tests**
```bash
# Go tests
go test ./... -v -cover

# Python tests
pytest enrichment-python/tests/ -v --cov=app

# React tests
npm test -- --coverage
```

## 📚 **Documentation**

### **Code Documentation**
- Add comments for complex logic
- Write clear function documentation
- Update README files
- Add inline examples

### **API Documentation**
- Document all API endpoints
- Include request/response examples
- Add error codes and messages
- Keep OpenAPI specs updated

### **User Documentation**
- Update installation guides
- Add troubleshooting sections
- Include configuration examples
- Maintain wiki pages

## 🔄 **Pull Request Process**

### **Before Submitting**
1. Ensure all tests pass
2. Update documentation
3. Follow coding standards
4. Add necessary tests
5. Self-review your changes

### **Pull Request Checklist**
- [ ] Code follows style guidelines
- [ ] Tests are added/updated
- [ ] Documentation is updated
- [ ] No breaking changes (or documented)
- [ ] Performance impact assessed
- [ ] Security implications considered

### **Review Process**
1. Automated checks must pass
2. Code review by maintainers
3. Address feedback and comments
4. Maintainer approval required
5. Merge to main branch

## 🏷️ **Commit Message Format**

Use conventional commit format:
```
type(scope): description

[optional body]

[optional footer]
```

### **Types**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build/tool changes

### **Examples**
```bash
feat(api): Add user authentication endpoint
fix(agent): Resolve memory leak in scanner
docs(readme): Update installation instructions
test(enrichment): Add CVE enrichment tests
```

## 🎯 **Areas for Contribution**

### **High Priority**
- Performance optimizations
- Security improvements
- Bug fixes
- Documentation updates

### **Medium Priority**
- New features
- UI/UX improvements
- Test coverage
- Monitoring enhancements

### **Low Priority**
- Code refactoring
- Style improvements
- Minor optimizations

## 📞 **Getting Help**

### **Resources**
- [GitHub Discussions](https://github.com/radhi1991/ZeroTrace/discussions)
- [Issue Tracker](https://github.com/radhi1991/ZeroTrace/issues)
- [Wiki](https://github.com/radhi1991/ZeroTrace/wiki)
- [Documentation](docs/)

### **Contact**
- Create a discussion for questions
- Use issues for bugs and features
- Join community conversations

## 🙏 **Recognition**

Contributors will be recognized in:
- GitHub contributors list
- Release notes
- Project documentation
- Community acknowledgments

---

**Thank you for contributing to ZeroTrace!** 🚀

**Last Updated**: January 2024  
**Version**: 1.0.0

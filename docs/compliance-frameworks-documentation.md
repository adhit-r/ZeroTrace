# ZeroTrace Compliance Frameworks Documentation

## Overview

ZeroTrace provides comprehensive compliance monitoring across multiple industry standards and regulatory frameworks. The platform supports automated compliance assessment, gap analysis, and remediation guidance.

## Supported Compliance Frameworks

### 1. CIS Benchmarks

**Purpose**: Center for Internet Security (CIS) Benchmarks provide security configuration guidelines for various technologies.

**Scope**:
- Operating Systems (Windows, Linux, macOS)
- Network Devices
- Cloud Platforms (AWS, Azure, GCP)
- Applications and Services

**Key Controls**:
- Access Control (CIS-1.1.1 to CIS-1.1.10)
- Authentication (CIS-2.1.1 to CIS-2.1.15)
- Authorization (CIS-3.1.1 to CIS-3.1.8)
- Data Protection (CIS-4.1.1 to CIS-4.1.12)
- Network Security (CIS-5.1.1 to CIS-5.1.20)

**Implementation**:
```go
scanner := NewConfigScanner()
result, err := scanner.Scan()
```

**Compliance Score Calculation**:
```
Score = (Passed Checks / Total Checks) * 100
```

**Example Output**:
```json
{
  "framework": "CIS",
  "score": 85.5,
  "total_checks": 100,
  "passed": 85,
  "failed": 15,
  "categories": {
    "access_control": {
      "score": 90.0,
      "passed": 9,
      "failed": 1
    },
    "authentication": {
      "score": 80.0,
      "passed": 12,
      "failed": 3
    }
  }
}
```

### 2. PCI-DSS (Payment Card Industry Data Security Standard)

**Purpose**: Security standards for organizations that handle credit card information.

**Requirements**:
- **Requirement 1**: Install and maintain a firewall configuration
- **Requirement 2**: Do not use vendor-supplied defaults
- **Requirement 3**: Protect stored cardholder data
- **Requirement 4**: Encrypt transmission of cardholder data
- **Requirement 5**: Use and regularly update anti-virus software
- **Requirement 6**: Develop and maintain secure systems
- **Requirement 7**: Restrict access by business need-to-know
- **Requirement 8**: Assign unique ID to each person with computer access
- **Requirement 9**: Restrict physical access to cardholder data
- **Requirement 10**: Track and monitor all access to network resources
- **Requirement 11**: Regularly test security systems and processes
- **Requirement 12**: Maintain a policy that addresses information security

**Implementation**:
```go
config := ComplianceConfig{
    Framework: "PCI-DSS",
    Version: "4.0",
    Level: "Level 1", // Based on transaction volume
}
scanner := NewConfigScannerWithConfig(config)
```

**Compliance Levels**:
- **Level 1**: 6M+ transactions/year
- **Level 2**: 1M-6M transactions/year
- **Level 3**: 20K-1M transactions/year
- **Level 4**: <20K transactions/year

### 3. HIPAA (Health Insurance Portability and Accountability Act)

**Purpose**: Security and privacy standards for healthcare organizations.

**Key Components**:
- **Administrative Safeguards**: Policies and procedures
- **Physical Safeguards**: Physical access controls
- **Technical Safeguards**: Technology-based protections

**Implementation**:
```go
config := ComplianceConfig{
    Framework: "HIPAA",
    EntityType: "Covered Entity", // or "Business Associate"
    RiskLevel: "High", // Low, Medium, High
}
scanner := NewConfigScannerWithConfig(config)
```

**Key Requirements**:
- Access Control (164.312(a))
- Audit Controls (164.312(b))
- Integrity (164.312(c))
- Person or Entity Authentication (164.312(d))
- Transmission Security (164.312(e))

### 4. GDPR (General Data Protection Regulation)

**Purpose**: Data protection and privacy regulation for EU citizens.

**Key Principles**:
- Lawfulness, fairness, and transparency
- Purpose limitation
- Data minimization
- Accuracy
- Storage limitation
- Integrity and confidentiality
- Accountability

**Implementation**:
```go
config := ComplianceConfig{
    Framework: "GDPR",
    DataController: true,
    DataProcessor: false,
    DataSubjectRights: true,
}
scanner := NewConfigScannerWithConfig(config)
```

**Key Articles**:
- **Article 5**: Principles relating to processing
- **Article 6**: Lawfulness of processing
- **Article 7**: Conditions for consent
- **Article 25**: Data protection by design and by default
- **Article 32**: Security of processing

### 5. SOC 2 (Service Organization Control 2)

**Purpose**: Security, availability, processing integrity, confidentiality, and privacy controls.

**Trust Service Criteria**:
- **CC6.1**: Logical and physical access security
- **CC6.2**: System access controls
- **CC6.3**: Data transmission and disposal
- **CC6.4**: System boundaries and data flow
- **CC6.5**: System processing integrity
- **CC6.6**: System availability
- **CC6.7**: System confidentiality
- **CC6.8**: System privacy

**Implementation**:
```go
config := ComplianceConfig{
    Framework: "SOC2",
    Type: "Type II", // Type I or Type II
    Criteria: []string{"Security", "Availability", "Confidentiality"},
}
scanner := NewConfigScannerWithConfig(config)
```

### 6. ISO 27001 (Information Security Management System)

**Purpose**: International standard for information security management.

**Key Domains**:
- **A.5**: Information security policies
- **A.6**: Organization of information security
- **A.7**: Human resource security
- **A.8**: Asset management
- **A.9**: Access control
- **A.10**: Cryptography
- **A.11**: Physical and environmental security
- **A.12**: Operations security
- **A.13**: Communications security
- **A.14**: System acquisition, development, and maintenance
- **A.15**: Supplier relationships
- **A.16**: Information security incident management
- **A.17**: Information security aspects of business continuity management
- **A.18**: Compliance

**Implementation**:
```go
config := ComplianceConfig{
    Framework: "ISO27001",
    Version: "2022",
    CertificationLevel: "Certified", // or "Implementing"
}
scanner := NewConfigScannerWithConfig(config)
```

## Compliance Assessment Process

### 1. Framework Selection

```go
frameworks := []string{"CIS", "PCI-DSS", "HIPAA", "GDPR", "SOC2", "ISO27001"}
config := ComplianceConfig{
    Frameworks: frameworks,
    AssessmentType: "comprehensive", // or "targeted"
}
```

### 2. Compliance Scanning

```go
scanner := NewConfigScanner()
result, err := scanner.Scan()
if err != nil {
    log.Fatal(err)
}

// Process compliance results
for _, check := range result.ComplianceChecks {
    fmt.Printf("Framework: %s\n", check.Framework)
    fmt.Printf("Requirement: %s\n", check.Requirement)
    fmt.Printf("Status: %s\n", check.Status)
    fmt.Printf("Gap: %s\n", check.Gap)
    fmt.Printf("Remediation: %s\n", check.Remediation)
}
```

### 3. Gap Analysis

```go
gaps := analyzeComplianceGaps(result)
for _, gap := range gaps {
    fmt.Printf("Framework: %s\n", gap.Framework)
    fmt.Printf("Requirement: %s\n", gap.Requirement)
    fmt.Printf("Priority: %s\n", gap.Priority)
    fmt.Printf("Effort: %s\n", gap.Effort)
    fmt.Printf("Timeline: %s\n", gap.Timeline)
}
```

## Compliance Scoring

### Score Calculation

```go
func calculateComplianceScore(checks []ComplianceCheck) float64 {
    totalChecks := len(checks)
    passedChecks := 0
    
    for _, check := range checks {
        if check.Status == "compliant" {
            passedChecks++
        }
    }
    
    return float64(passedChecks) / float64(totalChecks) * 100
}
```

### Weighted Scoring

```go
func calculateWeightedScore(checks []ComplianceCheck) float64 {
    totalWeight := 0.0
    weightedScore := 0.0
    
    for _, check := range checks {
        weight := getRequirementWeight(check.Requirement)
        totalWeight += weight
        
        if check.Status == "compliant" {
            weightedScore += weight
        }
    }
    
    return weightedScore / totalWeight * 100
}
```

## Compliance Reporting

### Executive Summary

```json
{
  "overall_score": 85.9,
  "frameworks": {
    "CIS": {
      "score": 85.5,
      "status": "compliant",
      "critical_gaps": 2,
      "high_gaps": 5
    },
    "PCI-DSS": {
      "score": 92.0,
      "status": "compliant",
      "critical_gaps": 0,
      "high_gaps": 1
    }
  },
  "recommendations": [
    "Implement multi-factor authentication",
    "Encrypt sensitive data at rest",
    "Regular security awareness training"
  ]
}
```

### Detailed Compliance Report

```json
{
  "framework": "CIS",
  "assessment_date": "2024-01-15T10:30:00Z",
  "next_assessment": "2024-04-15T10:30:00Z",
  "requirements": [
    {
      "id": "CIS-1.1.1",
      "title": "Ensure default passwords are changed",
      "status": "non_compliant",
      "priority": "critical",
      "description": "Default passwords are still in use",
      "remediation": "Change all default passwords immediately",
      "effort": "low",
      "timeline": "1 week",
      "evidence": [
        "Default admin password found on 3 systems",
        "Root password not changed on 2 servers"
      ]
    }
  ],
  "trends": {
    "score_change": 5.2,
    "improvement_areas": ["Access Control", "Authentication"],
    "regression_areas": ["Data Protection"]
  }
}
```

## Compliance Monitoring

### Real-time Monitoring

```go
func monitorCompliance() {
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        // Run compliance checks
        result, err := runComplianceChecks()
        if err != nil {
            log.Printf("Compliance check failed: %v", err)
            continue
        }
        
        // Check for compliance violations
        violations := detectComplianceViolations(result)
        if len(violations) > 0 {
            // Send alerts
            sendComplianceAlerts(violations)
        }
    }
}
```

### Compliance Alerts

```go
type ComplianceAlert struct {
    Framework    string    `json:"framework"`
    Requirement  string    `json:"requirement"`
    Severity     string    `json:"severity"`
    Description  string    `json:"description"`
    Timestamp    time.Time `json:"timestamp"`
    Remediation  string    `json:"remediation"`
}
```

## Compliance Automation

### Automated Remediation

```go
func autoRemediateCompliance(checks []ComplianceCheck) {
    for _, check := range checks {
        if check.AutoRemediable && check.Status == "non_compliant" {
            err := executeRemediation(check)
            if err != nil {
                log.Printf("Auto-remediation failed for %s: %v", 
                    check.Requirement, err)
            }
        }
    }
}
```

### Compliance Workflows

```go
type ComplianceWorkflow struct {
    ID          string                 `json:"id"`
    Framework   string                 `json:"framework"`
    Steps       []ComplianceStep       `json:"steps"`
    Triggers    []ComplianceTrigger    `json:"triggers"`
    Actions     []ComplianceAction     `json:"actions"`
}

type ComplianceStep struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Order       int       `json:"order"`
    Required    bool      `json:"required"`
}
```

## Compliance Integration

### API Integration

```go
func getComplianceStatus(framework string) (*ComplianceStatus, error) {
    url := fmt.Sprintf("/api/v2/compliance/status?framework=%s", framework)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var status ComplianceStatus
    err = json.NewDecoder(resp.Body).Decode(&status)
    return &status, err
}
```

### Database Integration

```sql
-- Compliance status table
CREATE TABLE compliance_status (
    id UUID PRIMARY KEY,
    framework VARCHAR(50) NOT NULL,
    score DECIMAL(5,2) NOT NULL,
    total_checks INTEGER NOT NULL,
    passed_checks INTEGER NOT NULL,
    failed_checks INTEGER NOT NULL,
    assessment_date TIMESTAMP NOT NULL,
    next_assessment TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Compliance gaps table
CREATE TABLE compliance_gaps (
    id UUID PRIMARY KEY,
    framework VARCHAR(50) NOT NULL,
    requirement VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    priority VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    gap_description TEXT NOT NULL,
    remediation TEXT NOT NULL,
    effort VARCHAR(20) NOT NULL,
    timeline VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Best Practices

### 1. Framework Selection
- Choose frameworks based on industry requirements
- Consider regulatory obligations
- Balance compliance with operational efficiency

### 2. Assessment Frequency
- Regular compliance assessments
- Continuous monitoring
- Automated compliance checking

### 3. Gap Management
- Prioritize critical gaps
- Implement remediation plans
- Track progress and improvements

### 4. Documentation
- Maintain compliance documentation
- Regular policy updates
- Training and awareness

### 5. Integration
- Integrate with existing systems
- Automate compliance processes
- Real-time monitoring and alerting

## Compliance Metrics

### Key Performance Indicators

```go
type ComplianceMetrics struct {
    OverallScore        float64 `json:"overall_score"`
    FrameworkScores     map[string]float64 `json:"framework_scores"`
    CriticalGaps        int     `json:"critical_gaps"`
    HighGaps           int     `json:"high_gaps"`
    MediumGaps         int     `json:"medium_gaps"`
    LowGaps            int     `json:"low_gaps"`
    RemediationRate    float64 `json:"remediation_rate"`
    ComplianceTrend    string  `json:"compliance_trend"`
}
```

### Compliance Dashboard

```json
{
  "metrics": {
    "overall_score": 85.9,
    "critical_gaps": 5,
    "high_gaps": 12,
    "medium_gaps": 25,
    "low_gaps": 8,
    "remediation_rate": 75.5,
    "compliance_trend": "improving"
  },
  "frameworks": {
    "CIS": {
      "score": 85.5,
      "status": "compliant",
      "last_assessment": "2024-01-15T10:30:00Z"
    },
    "PCI-DSS": {
      "score": 92.0,
      "status": "compliant",
      "last_assessment": "2024-01-15T10:30:00Z"
    }
  },
  "recommendations": [
    "Focus on critical gaps in CIS framework",
    "Implement automated compliance monitoring",
    "Regular security awareness training"
  ]
}
```

## Support

For compliance support and questions:
- **Documentation**: https://docs.zerotrace.com/compliance
- **Support Email**: compliance@zerotrace.com
- **GitHub Issues**: https://github.com/zerotrace/zerotrace/issues

# ZeroTrace Python Enrichment

## Overview
The ZeroTrace Python Enrichment service is responsible for advanced vulnerability analysis, data enrichment, and machine learning-based threat intelligence. It processes raw scan results from Go agents and enhances them with additional context, severity scoring, and predictive analytics.

## Architecture

### Service Structure
```
/enrichment-python
  /app
    /main.py              # FastAPI application entry point
    /api                  # API routes and handlers
      /v1
        /enrichment.py    # Enrichment endpoints
        /analysis.py      # Analysis endpoints
        /ml.py           # ML model endpoints
    /core                 # Core business logic
      /enrichment         # Data enrichment logic
        /cve_lookup.py    # CVE database lookups
        /severity.py      # Severity scoring
        /context.py       # Context enrichment
        /threat_intel.py  # Threat intelligence
      /analysis           # Analysis engines
        /static.py        # Static analysis
        /dynamic.py       # Dynamic analysis
        /behavioral.py    # Behavioral analysis
      /ml                 # Machine learning
        /models.py        # ML models
        /training.py      # Model training
        /prediction.py    # Prediction logic
        /features.py      # Feature engineering
    /services             # External service integrations
      /databases          # Database connections
        /postgres.py      # PostgreSQL client
        /redis.py         # Redis client
      /external           # External APIs
        /nvd.py          # NVD API client
        /github.py       # GitHub API client
        /virustotal.py   # VirusTotal API
    /models               # Data models
      /schemas.py         # Pydantic schemas
      /entities.py        # Database entities
      /dto.py            # Data transfer objects
    /utils                # Utilities
      /crypto.py         # Cryptographic functions
      /validation.py     # Validation utilities
      /logging.py        # Logging configuration
      /cache.py          # Caching utilities
  /tests                 # Test files
    /unit                # Unit tests
    /integration         # Integration tests
    /fixtures            # Test data
  /migrations            # Database migrations
  /scripts               # Utility scripts
    /train_models.py     # Model training script
    /update_cve_db.py    # CVE database update
  /config                # Configuration files
    /settings.py         # Application settings
    /logging.py          # Logging configuration
  /requirements.txt      # Python dependencies
  /pyproject.toml        # Poetry configuration
  Dockerfile
  docker-compose.yml
```

## Core Features

### 1. Vulnerability Enrichment
- CVE database lookups and correlation
- Severity scoring and prioritization
- False positive reduction
- Historical trend analysis
- Exploit availability checking

### 2. Threat Intelligence
- Real-time threat feeds integration
- Attack pattern recognition
- Risk scoring algorithms
- Predictive threat modeling
- Industry-specific threat context

### 3. Machine Learning Analysis
- Anomaly detection
- Risk prediction models
- Vulnerability clustering
- Trend forecasting
- Automated remediation suggestions

### 4. Data Processing
- Batch processing for large datasets
- Real-time streaming analysis
- Data normalization and cleaning
- Cross-reference analysis
- Performance optimization

## Implementation Details

### 1. FastAPI Application Setup

#### Main Application
```python
from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
import uvicorn

from app.core.config import settings
from app.api.v1 import enrichment, analysis, ml
from app.services.databases.postgres import get_db
from app.utils.logging import setup_logging

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    setup_logging()
    await init_database()
    await load_ml_models()
    yield
    # Shutdown
    await cleanup_resources()

app = FastAPI(
    title="ZeroTrace Enrichment API",
    description="Advanced vulnerability analysis and enrichment service",
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.ALLOWED_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(enrichment.router, prefix="/api/v1/enrichment", tags=["enrichment"])
app.include_router(analysis.router, prefix="/api/v1/analysis", tags=["analysis"])
app.include_router(ml.router, prefix="/api/v1/ml", tags=["ml"])

if __name__ == "__main__":
    uvicorn.run(
        "app.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=settings.DEBUG
    )
```

### 2. Enrichment Service

#### CVE Lookup Service
```python
from typing import List, Optional
import aiohttp
import asyncio
from app.models.schemas import Vulnerability, CVEInfo
from app.services.external.nvd import NVDClient
from app.utils.cache import cache

class CVELookupService:
    def __init__(self):
        self.nvd_client = NVDClient()
        self.cache_ttl = 3600  # 1 hour

    async def enrich_vulnerabilities(self, vulnerabilities: List[Vulnerability]) -> List[Vulnerability]:
        """Enrich vulnerabilities with CVE information"""
        enriched_vulns = []
        
        # Process in batches for efficiency
        batch_size = 50
        for i in range(0, len(vulnerabilities), batch_size):
            batch = vulnerabilities[i:i + batch_size]
            enriched_batch = await self._process_batch(batch)
            enriched_vulns.extend(enriched_batch)
            
        return enriched_vulns

    async def _process_batch(self, vulnerabilities: List[Vulnerability]) -> List[Vulnerability]:
        """Process a batch of vulnerabilities"""
        tasks = []
        for vuln in vulnerabilities:
            if vuln.cve_id:
                task = self._enrich_single_vulnerability(vuln)
                tasks.append(task)
            else:
                # Try to find CVE by package name and version
                task = self._find_cve_by_package(vuln)
                tasks.append(task)
        
        return await asyncio.gather(*tasks, return_exceptions=True)

    @cache(ttl=3600)
    async def _enrich_single_vulnerability(self, vulnerability: Vulnerability) -> Vulnerability:
        """Enrich a single vulnerability with CVE data"""
        try:
            cve_info = await self.nvd_client.get_cve(vulnerability.cve_id)
            if cve_info:
                vulnerability.cvss_score = cve_info.cvss_score
                vulnerability.severity = cve_info.severity
                vulnerability.description = cve_info.description
                vulnerability.references = cve_info.references
                vulnerability.published_date = cve_info.published_date
                vulnerability.last_modified_date = cve_info.last_modified_date
        except Exception as e:
            logger.error(f"Error enriching CVE {vulnerability.cve_id}: {e}")
        
        return vulnerability

    async def _find_cve_by_package(self, vulnerability: Vulnerability) -> Vulnerability:
        """Find CVE by package name and version"""
        try:
            cve_list = await self.nvd_client.search_by_package(
                vulnerability.package_name,
                vulnerability.package_version
            )
            if cve_list:
                # Use the most relevant CVE
                best_match = self._select_best_cve_match(cve_list, vulnerability)
                vulnerability.cve_id = best_match.cve_id
                vulnerability.cvss_score = best_match.cvss_score
                vulnerability.severity = best_match.severity
        except Exception as e:
            logger.error(f"Error finding CVE for package {vulnerability.package_name}: {e}")
        
        return vulnerability
```

#### Severity Scoring Service
```python
from typing import Dict, Any
import numpy as np
from app.models.schemas import Vulnerability, SeverityScore

class SeverityScoringService:
    def __init__(self):
        self.weights = {
            'cvss_score': 0.4,
            'exploit_available': 0.3,
            'attack_complexity': 0.2,
            'business_impact': 0.1
        }

    async def calculate_severity_score(self, vulnerability: Vulnerability) -> SeverityScore:
        """Calculate comprehensive severity score"""
        score = 0.0
        factors = {}

        # CVSS Score (0-10)
        if vulnerability.cvss_score:
            cvss_normalized = vulnerability.cvss_score / 10.0
            score += cvss_normalized * self.weights['cvss_score']
            factors['cvss_score'] = vulnerability.cvss_score

        # Exploit Availability
        exploit_score = await self._check_exploit_availability(vulnerability)
        score += exploit_score * self.weights['exploit_available']
        factors['exploit_available'] = exploit_score

        # Attack Complexity
        complexity_score = self._calculate_attack_complexity(vulnerability)
        score += complexity_score * self.weights['attack_complexity']
        factors['attack_complexity'] = complexity_score

        # Business Impact
        impact_score = self._calculate_business_impact(vulnerability)
        score += impact_score * self.weights['business_impact']
        factors['business_impact'] = impact_score

        # Determine severity level
        severity_level = self._determine_severity_level(score)

        return SeverityScore(
            overall_score=score,
            severity_level=severity_level,
            factors=factors,
            confidence=0.85
        )

    async def _check_exploit_availability(self, vulnerability: Vulnerability) -> float:
        """Check if exploit is available for the vulnerability"""
        # Check multiple sources
        sources = [
            self._check_exploit_db(vulnerability),
            self._check_github_exploits(vulnerability),
            self._check_metasploit(vulnerability)
        ]
        
        results = await asyncio.gather(*sources, return_exceptions=True)
        
        # Calculate exploit availability score
        exploit_count = sum(1 for result in results if result and isinstance(result, bool))
        return min(exploit_count / len(sources), 1.0)

    def _calculate_attack_complexity(self, vulnerability: Vulnerability) -> float:
        """Calculate attack complexity score"""
        complexity_factors = {
            'remote': 0.8,
            'local': 0.4,
            'physical': 0.2,
            'network': 0.9,
            'adjacent': 0.6
        }
        
        # Analyze vulnerability characteristics
        if 'remote' in vulnerability.description.lower():
            return complexity_factors['remote']
        elif 'local' in vulnerability.description.lower():
            return complexity_factors['local']
        
        return 0.5  # Default medium complexity

    def _calculate_business_impact(self, vulnerability: Vulnerability) -> float:
        """Calculate business impact score"""
        impact_keywords = {
            'data breach': 1.0,
            'rce': 0.9,
            'sql injection': 0.8,
            'xss': 0.6,
            'information disclosure': 0.5
        }
        
        description_lower = vulnerability.description.lower()
        max_impact = 0.0
        
        for keyword, impact in impact_keywords.items():
            if keyword in description_lower:
                max_impact = max(max_impact, impact)
        
        return max_impact

    def _determine_severity_level(self, score: float) -> str:
        """Determine severity level based on score"""
        if score >= 0.8:
            return "CRITICAL"
        elif score >= 0.6:
            return "HIGH"
        elif score >= 0.4:
            return "MEDIUM"
        elif score >= 0.2:
            return "LOW"
        else:
            return "INFO"
```

### 3. Machine Learning Models

#### Anomaly Detection Model
```python
import joblib
import numpy as np
from sklearn.ensemble import IsolationForest
from sklearn.preprocessing import StandardScaler
from app.models.schemas import Vulnerability, AnomalyScore

class AnomalyDetectionModel:
    def __init__(self):
        self.model = None
        self.scaler = StandardScaler()
        self.feature_columns = [
            'cvss_score', 'age_days', 'exploit_count',
            'affected_versions', 'patch_availability'
        ]

    async def load_model(self, model_path: str):
        """Load trained model"""
        try:
            self.model = joblib.load(model_path)
            logger.info(f"Loaded anomaly detection model from {model_path}")
        except Exception as e:
            logger.error(f"Error loading model: {e}")
            # Train a new model if loading fails
            await self.train_model()

    async def predict_anomaly(self, vulnerability: Vulnerability) -> AnomalyScore:
        """Predict if vulnerability is anomalous"""
        if not self.model:
            await self.load_model("models/anomaly_detection.pkl")

        # Extract features
        features = self._extract_features(vulnerability)
        
        # Scale features
        features_scaled = self.scaler.transform([features])
        
        # Predict
        prediction = self.model.predict(features_scaled)[0]
        score = self.model.score_samples(features_scaled)[0]
        
        return AnomalyScore(
            is_anomaly=prediction == -1,
            anomaly_score=score,
            confidence=0.85
        )

    def _extract_features(self, vulnerability: Vulnerability) -> np.ndarray:
        """Extract features from vulnerability"""
        features = []
        
        # CVSS Score
        features.append(vulnerability.cvss_score or 0.0)
        
        # Age in days
        age_days = (datetime.now() - vulnerability.published_date).days if vulnerability.published_date else 0
        features.append(age_days)
        
        # Exploit count (placeholder)
        features.append(vulnerability.exploit_count or 0)
        
        # Affected versions
        features.append(len(vulnerability.affected_versions) if vulnerability.affected_versions else 0)
        
        # Patch availability
        features.append(1.0 if vulnerability.patch_available else 0.0)
        
        return np.array(features)

    async def train_model(self, training_data: List[Vulnerability] = None):
        """Train the anomaly detection model"""
        if not training_data:
            # Load training data from database
            training_data = await self._load_training_data()

        # Extract features
        X = np.array([self._extract_features(vuln) for vuln in training_data])
        
        # Scale features
        X_scaled = self.scaler.fit_transform(X)
        
        # Train model
        self.model = IsolationForest(
            contamination=0.1,
            random_state=42,
            n_estimators=100
        )
        self.model.fit(X_scaled)
        
        # Save model
        joblib.dump(self.model, "models/anomaly_detection.pkl")
        logger.info("Trained and saved anomaly detection model")
```

#### Risk Prediction Model
```python
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
import pandas as pd

class RiskPredictionModel:
    def __init__(self):
        self.model = RandomForestClassifier(
            n_estimators=100,
            max_depth=10,
            random_state=42
        )
        self.feature_columns = [
            'cvss_score', 'severity_level', 'exploit_available',
            'attack_complexity', 'business_impact', 'age_days'
        ]

    async def predict_risk(self, vulnerability: Vulnerability) -> RiskPrediction:
        """Predict risk level for vulnerability"""
        features = self._extract_features(vulnerability)
        
        # Make prediction
        prediction = self.model.predict([features])[0]
        probabilities = self.model.predict_proba([features])[0]
        
        return RiskPrediction(
            risk_level=prediction,
            confidence=max(probabilities),
            probabilities={
                'low': probabilities[0],
                'medium': probabilities[1],
                'high': probabilities[2],
                'critical': probabilities[3]
            }
        )

    async def train_model(self, training_data: List[Vulnerability]):
        """Train the risk prediction model"""
        # Prepare training data
        X = []
        y = []
        
        for vuln in training_data:
            features = self._extract_features(vuln)
            X.append(features)
            y.append(vuln.actual_risk_level)  # Historical risk data
        
        X = np.array(X)
        y = np.array(y)
        
        # Split data
        X_train, X_test, y_train, y_test = train_test_split(
            X, y, test_size=0.2, random_state=42
        )
        
        # Train model
        self.model.fit(X_train, y_train)
        
        # Evaluate model
        y_pred = self.model.predict(X_test)
        report = classification_report(y_test, y_pred)
        logger.info(f"Model training report:\n{report}")
        
        # Save model
        joblib.dump(self.model, "models/risk_prediction.pkl")
```

### 4. API Endpoints

#### Enrichment Endpoints
```python
from fastapi import APIRouter, Depends, HTTPException
from typing import List
from app.models.schemas import (
    Vulnerability, EnrichedVulnerability, 
    EnrichmentRequest, EnrichmentResponse
)
from app.core.enrichment.cve_lookup import CVELookupService
from app.core.enrichment.severity import SeverityScoringService

router = APIRouter()

@router.post("/enrich", response_model=EnrichmentResponse)
async def enrich_vulnerabilities(
    request: EnrichmentRequest,
    cve_service: CVELookupService = Depends(),
    severity_service: SeverityScoringService = Depends()
):
    """Enrich vulnerabilities with additional data"""
    try:
        # Enrich with CVE data
        enriched_vulns = await cve_service.enrich_vulnerabilities(request.vulnerabilities)
        
        # Calculate severity scores
        for vuln in enriched_vulns:
            severity_score = await severity_service.calculate_severity_score(vuln)
            vuln.severity_score = severity_score
        
        return EnrichmentResponse(
            vulnerabilities=enriched_vulns,
            enrichment_time=datetime.now(),
            total_processed=len(enriched_vulns)
        )
    except Exception as e:
        logger.error(f"Error enriching vulnerabilities: {e}")
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/batch-enrich", response_model=EnrichmentResponse)
async def batch_enrich_vulnerabilities(
    request: EnrichmentRequest,
    background_tasks: BackgroundTasks,
    cve_service: CVELookupService = Depends()
):
    """Enrich vulnerabilities in background"""
    # Start background task
    task_id = str(uuid.uuid4())
    background_tasks.add_task(
        process_batch_enrichment,
        task_id,
        request.vulnerabilities,
        cve_service
    )
    
    return EnrichmentResponse(
        task_id=task_id,
        status="processing",
        message="Batch enrichment started"
    )

@router.get("/status/{task_id}")
async def get_enrichment_status(task_id: str):
    """Get status of background enrichment task"""
    status = await get_task_status(task_id)
    return status
```

#### Analysis Endpoints
```python
@router.post("/analyze", response_model=AnalysisResponse)
async def analyze_vulnerabilities(
    request: AnalysisRequest,
    anomaly_model: AnomalyDetectionModel = Depends(),
    risk_model: RiskPredictionModel = Depends()
):
    """Analyze vulnerabilities using ML models"""
    results = []
    
    for vuln in request.vulnerabilities:
        # Anomaly detection
        anomaly_score = await anomaly_model.predict_anomaly(vuln)
        
        # Risk prediction
        risk_prediction = await risk_model.predict_risk(vuln)
        
        results.append(AnalysisResult(
            vulnerability_id=vuln.id,
            anomaly_score=anomaly_score,
            risk_prediction=risk_prediction,
            analysis_timestamp=datetime.now()
        ))
    
    return AnalysisResponse(
        results=results,
        analysis_time=datetime.now(),
        total_analyzed=len(results)
    )
```

### 5. Data Models

#### Pydantic Schemas
```python
from pydantic import BaseModel, Field
from typing import List, Optional, Dict, Any
from datetime import datetime
from enum import Enum

class SeverityLevel(str, Enum):
    CRITICAL = "CRITICAL"
    HIGH = "HIGH"
    MEDIUM = "MEDIUM"
    LOW = "LOW"
    INFO = "INFO"

class Vulnerability(BaseModel):
    id: str
    cve_id: Optional[str] = None
    title: str
    description: str
    cvss_score: Optional[float] = None
    severity: Optional[SeverityLevel] = None
    package_name: Optional[str] = None
    package_version: Optional[str] = None
    published_date: Optional[datetime] = None
    references: List[str] = []
    affected_versions: List[str] = []
    patch_available: bool = False
    exploit_count: int = 0

class EnrichedVulnerability(Vulnerability):
    severity_score: Optional[SeverityScore] = None
    anomaly_score: Optional[AnomalyScore] = None
    risk_prediction: Optional[RiskPrediction] = None
    enrichment_metadata: Dict[str, Any] = {}

class SeverityScore(BaseModel):
    overall_score: float
    severity_level: SeverityLevel
    factors: Dict[str, float]
    confidence: float

class AnomalyScore(BaseModel):
    is_anomaly: bool
    anomaly_score: float
    confidence: float

class RiskPrediction(BaseModel):
    risk_level: str
    confidence: float
    probabilities: Dict[str, float]

class EnrichmentRequest(BaseModel):
    vulnerabilities: List[Vulnerability]
    options: Optional[Dict[str, Any]] = {}

class EnrichmentResponse(BaseModel):
    vulnerabilities: Optional[List[EnrichedVulnerability]] = None
    task_id: Optional[str] = None
    status: Optional[str] = None
    message: Optional[str] = None
    enrichment_time: Optional[datetime] = None
    total_processed: Optional[int] = None
```

### 6. Performance Optimizations

#### Async Processing
```python
import asyncio
from concurrent.futures import ThreadPoolExecutor

class AsyncEnrichmentService:
    def __init__(self):
        self.executor = ThreadPoolExecutor(max_workers=10)
        self.semaphore = asyncio.Semaphore(50)  # Limit concurrent requests

    async def enrich_batch_async(self, vulnerabilities: List[Vulnerability]) -> List[EnrichedVulnerability]:
        """Enrich vulnerabilities asynchronously"""
        async def enrich_single(vuln: Vulnerability) -> EnrichedVulnerability:
            async with self.semaphore:
                return await self._enrich_single_vulnerability(vuln)

        tasks = [enrich_single(vuln) for vuln in vulnerabilities]
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        # Filter out exceptions
        enriched_vulns = []
        for result in results:
            if isinstance(result, Exception):
                logger.error(f"Enrichment error: {result}")
            else:
                enriched_vulns.append(result)
        
        return enriched_vulns
```

#### Caching Strategy
```python
from functools import lru_cache
import redis.asyncio as redis

class CachedEnrichmentService:
    def __init__(self):
        self.redis_client = redis.Redis(host='localhost', port=6379, db=0)
        self.cache_ttl = 3600  # 1 hour

    async def get_cached_enrichment(self, cve_id: str) -> Optional[Dict]:
        """Get cached enrichment data"""
        cache_key = f"enrichment:{cve_id}"
        cached_data = await self.redis_client.get(cache_key)
        
        if cached_data:
            return json.loads(cached_data)
        return None

    async def cache_enrichment(self, cve_id: str, data: Dict):
        """Cache enrichment data"""
        cache_key = f"enrichment:{cve_id}"
        await self.redis_client.setex(
            cache_key,
            self.cache_ttl,
            json.dumps(data)
        )

    @lru_cache(maxsize=1000)
    def get_static_data(self, key: str) -> Any:
        """Get static data with in-memory cache"""
        # Static data that doesn't change often
        pass
```

## Configuration

### Environment Variables
```bash
# Application Configuration
ENRICHMENT_HOST=0.0.0.0
ENRICHMENT_PORT=8000
DEBUG=false

# Database Configuration
DATABASE_URL=postgresql://user:password@localhost/zerotrace
REDIS_URL=redis://localhost:6379/0

# External APIs
NVD_API_KEY=your-nvd-api-key
GITHUB_TOKEN=your-github-token
VIRUSTOTAL_API_KEY=your-virustotal-key

# ML Model Configuration
MODEL_PATH=/app/models
TRAINING_DATA_PATH=/app/data/training
MODEL_UPDATE_INTERVAL=86400

# Performance Configuration
MAX_CONCURRENT_REQUESTS=50
BATCH_SIZE=100
CACHE_TTL=3600
```

### Settings Configuration
```python
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # Application
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    DEBUG: bool = False
    
    # Database
    DATABASE_URL: str
    REDIS_URL: str
    
    # External APIs
    NVD_API_KEY: Optional[str] = None
    GITHUB_TOKEN: Optional[str] = None
    VIRUSTOTAL_API_KEY: Optional[str] = None
    
    # ML Models
    MODEL_PATH: str = "/app/models"
    TRAINING_DATA_PATH: str = "/app/data/training"
    MODEL_UPDATE_INTERVAL: int = 86400
    
    # Performance
    MAX_CONCURRENT_REQUESTS: int = 50
    BATCH_SIZE: int = 100
    CACHE_TTL: int = 3600
    
    class Config:
        env_file = ".env"

settings = Settings()
```

## Testing Strategy

### Unit Tests
```python
import pytest
from unittest.mock import AsyncMock, patch
from app.core.enrichment.cve_lookup import CVELookupService

@pytest.mark.asyncio
async def test_cve_lookup_enrichment():
    """Test CVE lookup enrichment"""
    service = CVELookupService()
    
    # Mock NVD client
    with patch.object(service.nvd_client, 'get_cve') as mock_get_cve:
        mock_get_cve.return_value = CVEInfo(
            cve_id="CVE-2021-1234",
            cvss_score=8.5,
            severity="HIGH",
            description="Test vulnerability"
        )
        
        vulnerability = Vulnerability(
            id="test-1",
            cve_id="CVE-2021-1234",
            title="Test Vulnerability"
        )
        
        enriched = await service._enrich_single_vulnerability(vulnerability)
        
        assert enriched.cvss_score == 8.5
        assert enriched.severity == "HIGH"
```

### Integration Tests
```python
@pytest.mark.asyncio
async def test_full_enrichment_pipeline():
    """Test full enrichment pipeline"""
    from app.api.v1.enrichment import enrich_vulnerabilities
    
    request = EnrichmentRequest(
        vulnerabilities=[
            Vulnerability(
                id="test-1",
                cve_id="CVE-2021-1234",
                title="Test Vulnerability"
            )
        ]
    )
    
    response = await enrich_vulnerabilities(request)
    
    assert response.total_processed == 1
    assert response.vulnerabilities[0].severity_score is not None
```

## Deployment

### Docker Configuration
```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create model directory
RUN mkdir -p /app/models

# Expose port
EXPOSE 8000

# Run application
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

### Local Development
```bash
# Install dependencies
pip install -r requirements.txt

# Run the service
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

# Run with Docker
podman build -t zerotrace-enrichment .
podman run -p 8000:8000 zerotrace-enrichment
```

import pytest
import asyncio
import json
from unittest.mock import Mock, patch, AsyncMock
from app.cve_enrichment import CVEEnrichmentService, CVEEnrichmentRequest, CVEEnrichmentResponse
from app.ai_matching.cpe_matcher import CPEMatcher
from app.version_matcher import CPEVersionMatcher


class TestCVEEnrichmentService:
    """Test cases for CVE enrichment service"""
    
    @pytest.fixture
    def enrichment_service(self):
        """Create a CVE enrichment service instance for testing"""
        return CVEEnrichmentService()
    
    @pytest.fixture
    def sample_request(self):
        """Create a sample enrichment request"""
        return CVEEnrichmentRequest(
            software_name="Google Chrome",
            version="120.0.6099.109",
            vendor="Google",
            path="/Applications/Google Chrome.app"
        )
    
    @pytest.fixture
    def sample_cve_data(self):
        """Create sample CVE data for testing"""
        return {
            "CVE-2024-1234": {
                "id": "CVE-2024-1234",
                "description": "Test vulnerability in Chrome",
                "severity": "high",
                "cvss_score": 7.5,
                "affected_versions": ["<120.0.6099.110"],
                "references": ["https://example.com/cve-2024-1234"]
            },
            "CVE-2024-1235": {
                "id": "CVE-2024-1235", 
                "description": "Another test vulnerability",
                "severity": "medium",
                "cvss_score": 5.0,
                "affected_versions": ["<120.0.6099.111"],
                "references": ["https://example.com/cve-2024-1235"]
            }
        }
    
    def test_enrichment_service_initialization(self, enrichment_service):
        """Test that enrichment service initializes correctly"""
        assert enrichment_service is not None
        assert hasattr(enrichment_service, 'cpe_matcher')
        assert hasattr(enrichment_service, 'version_matcher')
    
    @pytest.mark.asyncio
    async def test_enrich_software_success(self, enrichment_service, sample_request, sample_cve_data):
        """Test successful software enrichment"""
        with patch.object(enrichment_service, 'load_cve_data', return_value=sample_cve_data):
            with patch.object(enrichment_service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*"):
                with patch.object(enrichment_service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*", 0.95)):
                    
                    response = await enrichment_service.enrich_software(sample_request)
                    
                    assert response is not None
                    assert response.software_name == sample_request.software_name
                    assert response.version == sample_request.version
                    assert len(response.vulnerabilities) > 0
                    assert response.cpe is not None
                    assert response.confidence > 0
    
    @pytest.mark.asyncio
    async def test_enrich_software_no_cves_found(self, enrichment_service, sample_request):
        """Test enrichment when no CVEs are found"""
        with patch.object(enrichment_service, 'load_cve_data', return_value={}):
            with patch.object(enrichment_service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*"):
                with patch.object(enrichment_service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*", 0.95)):
                    
                    response = await enrichment_service.enrich_software(sample_request)
                    
                    assert response is not None
                    assert response.software_name == sample_request.software_name
                    assert len(response.vulnerabilities) == 0
                    assert response.cpe is not None
    
    @pytest.mark.asyncio
    async def test_enrich_software_cpe_not_found(self, enrichment_service, sample_request):
        """Test enrichment when CPE cannot be determined"""
        with patch.object(enrichment_service, 'load_cve_data', return_value={}):
            with patch.object(enrichment_service.cpe_matcher, 'get_cpe_for_software', return_value=None):
                with patch.object(enrichment_service.cpe_matcher, 'match_software_to_cpe', return_value=(None, 0.0)):
                    
                    response = await enrichment_service.enrich_software(sample_request)
                    
                    assert response is not None
                    assert response.software_name == sample_request.software_name
                    assert response.cpe is None
                    assert response.confidence == 0.0
    
    @pytest.mark.asyncio
    async def test_batch_enrichment(self, enrichment_service):
        """Test batch enrichment of multiple software items"""
        requests = [
            CVEEnrichmentRequest(
                software_name="Google Chrome",
                version="120.0.6099.109",
                vendor="Google",
                path="/Applications/Google Chrome.app"
            ),
            CVEEnrichmentRequest(
                software_name="Mozilla Firefox",
                version="121.0",
                vendor="Mozilla",
                path="/Applications/Firefox.app"
            )
        ]
        
        with patch.object(enrichment_service, 'load_cve_data', return_value={}):
            with patch.object(enrichment_service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:test:software:1.0:*:*:*:*:*:*:*"):
                with patch.object(enrichment_service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:test:software:1.0:*:*:*:*:*:*:*", 0.8)):
                    
                    responses = await enrichment_service.batch_enrich(requests)
                    
                    assert len(responses) == 2
                    assert all(isinstance(response, CVEEnrichmentResponse) for response in responses)
    
    @pytest.mark.asyncio
    async def test_enrichment_with_version_matching(self, enrichment_service, sample_request, sample_cve_data):
        """Test enrichment with version matching logic"""
        with patch.object(enrichment_service, 'load_cve_data', return_value=sample_cve_data):
            with patch.object(enrichment_service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*"):
                with patch.object(enrichment_service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*", 0.95)):
                    with patch.object(enrichment_service.version_matcher, 'match_cpe_version', return_value=True):
                        
                        response = await enrichment_service.enrich_software(sample_request)
                        
                        assert response is not None
                        assert len(response.vulnerabilities) > 0
                        
                        # Check that vulnerabilities are filtered by version
                        for vuln in response.vulnerabilities:
                            assert vuln.affected is True
    
    def test_load_cve_data_file_not_found(self, enrichment_service):
        """Test loading CVE data when file doesn't exist"""
        with patch('os.path.exists', return_value=False):
            data = enrichment_service.load_cve_data()
            assert data == {}
    
    def test_load_cve_data_invalid_json(self, enrichment_service):
        """Test loading CVE data with invalid JSON"""
        with patch('os.path.exists', return_value=True):
            with patch('builtins.open', mock_open(read_data="invalid json")):
                with patch('json.load', side_effect=json.JSONDecodeError("Invalid JSON", "", 0)):
                    data = enrichment_service.load_cve_data()
                    assert data == {}
    
    @pytest.mark.asyncio
    async def test_enrichment_performance(self, enrichment_service, sample_request):
        """Test enrichment performance with large dataset"""
        # Create large CVE dataset
        large_cve_data = {}
        for i in range(1000):
            large_cve_data[f"CVE-2024-{i:04d}"] = {
                "id": f"CVE-2024-{i:04d}",
                "description": f"Test vulnerability {i}",
                "severity": "medium",
                "cvss_score": 5.0,
                "affected_versions": ["<120.0.6099.110"],
                "references": [f"https://example.com/cve-2024-{i:04d}"]
            }
        
        with patch.object(enrichment_service, 'load_cve_data', return_value=large_cve_data):
            with patch.object(enrichment_service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*"):
                with patch.object(enrichment_service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*", 0.95)):
                    
                    import time
                    start_time = time.time()
                    response = await enrichment_service.enrich_software(sample_request)
                    end_time = time.time()
                    
                    # Should complete within reasonable time (5 seconds)
                    assert (end_time - start_time) < 5.0
                    assert response is not None


class TestCPEMatcher:
    """Test cases for CPE matching functionality"""
    
    @pytest.fixture
    def cpe_matcher(self):
        """Create a CPE matcher instance for testing"""
        return CPEMatcher()
    
    def test_get_cpe_for_software_exact_match(self, cpe_matcher):
        """Test CPE matching with exact software name match"""
        cpe = cpe_matcher.get_cpe_for_software("Google Chrome", "120.0.6099.109", "Google")
        assert cpe is not None
        assert "google" in cpe.lower()
        assert "chrome" in cpe.lower()
    
    def test_get_cpe_for_software_fuzzy_match(self, cpe_matcher):
        """Test CPE matching with fuzzy matching"""
        cpe = cpe_matcher.get_cpe_for_software("Chrome Browser", "120.0.6099.109", "Google")
        assert cpe is not None
    
    def test_get_cpe_for_software_no_match(self, cpe_matcher):
        """Test CPE matching when no match is found"""
        cpe = cpe_matcher.get_cpe_for_software("Unknown Software", "1.0", "Unknown Vendor")
        assert cpe is None
    
    def test_match_software_to_cpe_high_confidence(self, cpe_matcher):
        """Test software to CPE matching with high confidence"""
        cpe, confidence = cpe_matcher.match_software_to_cpe("Google Chrome", "120.0.6099.109", "Google")
        assert cpe is not None
        assert confidence > 0.8
    
    def test_match_software_to_cpe_low_confidence(self, cpe_matcher):
        """Test software to CPE matching with low confidence"""
        cpe, confidence = cpe_matcher.match_software_to_cpe("Unknown Software", "1.0", "Unknown Vendor")
        assert confidence < 0.5
    
    def test_cpe_similarity_calculation(self, cpe_matcher):
        """Test CPE similarity calculation"""
        similarity = cpe_matcher.calculate_similarity("Google Chrome", "Google Chrome Browser")
        assert 0 <= similarity <= 1
        assert similarity > 0.8  # Should be high for similar names


class TestCPEVersionMatcher:
    """Test cases for CPE version matching"""
    
    @pytest.fixture
    def version_matcher(self):
        """Create a version matcher instance for testing"""
        return CPEVersionMatcher()
    
    def test_match_cpe_version_exact_match(self, version_matcher):
        """Test version matching with exact version match"""
        result = version_matcher.match_cpe_version("120.0.6099.109", "120.0.6099.109")
        assert result is True
    
    def test_match_cpe_version_range_match(self, version_matcher):
        """Test version matching with version range"""
        result = version_matcher.match_cpe_version("120.0.6099.109", "<120.0.6099.110")
        assert result is True
    
    def test_match_cpe_version_no_match(self, version_matcher):
        """Test version matching when versions don't match"""
        result = version_matcher.match_cpe_version("120.0.6099.109", ">120.0.6099.110")
        assert result is False
    
    def test_parse_version_semver(self, version_matcher):
        """Test parsing semantic versions"""
        version = version_matcher.parse_version("1.2.3")
        assert version is not None
        assert version.major == 1
        assert version.minor == 2
        assert version.patch == 3
    
    def test_parse_version_calver(self, version_matcher):
        """Test parsing calendar versions"""
        version = version_matcher.parse_version("2024.01.15")
        assert version is not None
    
    def test_parse_version_pep440(self, version_matcher):
        """Test parsing PEP 440 versions"""
        version = version_matcher.parse_version("1.2.3a1")
        assert version is not None
    
    def test_compare_versions(self, version_matcher):
        """Test version comparison"""
        # Test various comparison scenarios
        assert version_matcher.compare_versions("1.2.3", "1.2.3") == 0
        assert version_matcher.compare_versions("1.2.3", "1.2.4") < 0
        assert version_matcher.compare_versions("1.2.4", "1.2.3") > 0
    
    def test_version_range_parsing(self, version_matcher):
        """Test parsing version ranges"""
        ranges = [
            "<1.2.3",
            "<=1.2.3", 
            ">1.2.3",
            ">=1.2.3",
            "~1.2.3",
            "^1.2.3",
            "1.2.3 - 1.2.5"
        ]
        
        for range_str in ranges:
            parsed = version_matcher.parse_version_range(range_str)
            assert parsed is not None


class TestIntegration:
    """Integration tests for the enrichment service"""
    
    @pytest.mark.asyncio
    async def test_full_enrichment_pipeline(self):
        """Test the complete enrichment pipeline"""
        service = CVEEnrichmentService()
        
        request = CVEEnrichmentRequest(
            software_name="Google Chrome",
            version="120.0.6099.109",
            vendor="Google",
            path="/Applications/Google Chrome.app"
        )
        
        # Mock external dependencies
        with patch.object(service, 'load_cve_data', return_value={}):
            with patch.object(service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*"):
                with patch.object(service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*", 0.95)):
                    
                    response = await service.enrich_software(request)
                    
                    assert response is not None
                    assert response.software_name == request.software_name
                    assert response.version == request.version
                    assert response.vendor == request.vendor
                    assert response.path == request.path
                    assert response.cpe is not None
                    assert response.confidence > 0
    
    @pytest.mark.asyncio
    async def test_batch_enrichment_pipeline(self):
        """Test batch enrichment pipeline"""
        service = CVEEnrichmentService()
        
        requests = [
            CVEEnrichmentRequest(
                software_name="Google Chrome",
                version="120.0.6099.109",
                vendor="Google",
                path="/Applications/Google Chrome.app"
            ),
            CVEEnrichmentRequest(
                software_name="Mozilla Firefox",
                version="121.0",
                vendor="Mozilla",
                path="/Applications/Firefox.app"
            )
        ]
        
        with patch.object(service, 'load_cve_data', return_value={}):
            with patch.object(service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:test:software:1.0:*:*:*:*:*:*:*"):
                with patch.object(service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:test:software:1.0:*:*:*:*:*:*:*", 0.8)):
                    
                    responses = await service.batch_enrich(requests)
                    
                    assert len(responses) == 2
                    assert all(isinstance(response, CVEEnrichmentResponse) for response in responses)
                    assert responses[0].software_name == requests[0].software_name
                    assert responses[1].software_name == requests[1].software_name


# Performance tests
class TestPerformance:
    """Performance tests for the enrichment service"""
    
    @pytest.mark.asyncio
    async def test_enrichment_performance_large_dataset(self):
        """Test enrichment performance with large CVE dataset"""
        service = CVEEnrichmentService()
        
        # Create large CVE dataset
        large_cve_data = {}
        for i in range(10000):
            large_cve_data[f"CVE-2024-{i:05d}"] = {
                "id": f"CVE-2024-{i:05d}",
                "description": f"Test vulnerability {i}",
                "severity": "medium",
                "cvss_score": 5.0,
                "affected_versions": ["<120.0.6099.110"],
                "references": [f"https://example.com/cve-2024-{i:05d}"]
            }
        
        request = CVEEnrichmentRequest(
            software_name="Google Chrome",
            version="120.0.6099.109",
            vendor="Google",
            path="/Applications/Google Chrome.app"
        )
        
        with patch.object(service, 'load_cve_data', return_value=large_cve_data):
            with patch.object(service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*"):
                with patch.object(service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:google:chrome:120.0.6099.109:*:*:*:*:*:*:*", 0.95)):
                    
                    import time
                    start_time = time.time()
                    response = await service.enrich_software(request)
                    end_time = time.time()
                    
                    # Should complete within reasonable time (10 seconds for large dataset)
                    assert (end_time - start_time) < 10.0
                    assert response is not None
    
    @pytest.mark.asyncio
    async def test_concurrent_enrichment(self):
        """Test concurrent enrichment requests"""
        service = CVEEnrichmentService()
        
        requests = [
            CVEEnrichmentRequest(
                software_name=f"Software {i}",
                version="1.0.0",
                vendor=f"Vendor {i}",
                path=f"/path/to/software{i}"
            ) for i in range(100)
        ]
        
        with patch.object(service, 'load_cve_data', return_value={}):
            with patch.object(service.cpe_matcher, 'get_cpe_for_software', return_value="cpe:2.3:a:test:software:1.0:*:*:*:*:*:*:*"):
                with patch.object(service.cpe_matcher, 'match_software_to_cpe', return_value=("cpe:2.3:a:test:software:1.0:*:*:*:*:*:*:*", 0.8)):
                    
                    import time
                    start_time = time.time()
                    responses = await service.batch_enrich(requests)
                    end_time = time.time()
                    
                    # Should complete within reasonable time (5 seconds for 100 concurrent requests)
                    assert (end_time - start_time) < 5.0
                    assert len(responses) == 100
                    assert all(isinstance(response, CVEEnrichmentResponse) for response in responses)


# Test configuration
@pytest.fixture(scope="session")
def event_loop():
    """Create an instance of the default event loop for the test session."""
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()


# Mock utilities
def mock_open(read_data):
    """Mock file open for testing"""
    from unittest.mock import mock_open
    return mock_open(read_data=read_data)

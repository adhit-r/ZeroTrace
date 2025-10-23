import pytest
import asyncio
import sys
import os
from pathlib import Path

# Add the parent directory to the Python path
sys.path.insert(0, str(Path(__file__).parent.parent))

# Configure asyncio for testing
@pytest.fixture(scope="session")
def event_loop():
    """Create an instance of the default event loop for the test session."""
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()


@pytest.fixture
def sample_cve_data():
    """Sample CVE data for testing"""
    return {
        "CVE-2024-1234": {
            "id": "CVE-2024-1234",
            "description": "Test vulnerability in Chrome",
            "severity": "high",
            "cvss_score": 7.5,
            "affected_versions": ["<120.0.6099.110"],
            "references": ["https://example.com/cve-2024-1234"],
            "cpe": "cpe:2.3:a:google:chrome:*:*:*:*:*:*:*:*"
        },
        "CVE-2024-1235": {
            "id": "CVE-2024-1235",
            "description": "Another test vulnerability",
            "severity": "medium",
            "cvss_score": 5.0,
            "affected_versions": ["<120.0.6099.111"],
            "references": ["https://example.com/cve-2024-1235"],
            "cpe": "cpe:2.3:a:google:chrome:*:*:*:*:*:*:*:*"
        }
    }


@pytest.fixture
def sample_software_data():
    """Sample software data for testing"""
    return [
        {
            "name": "Google Chrome",
            "version": "120.0.6099.109",
            "vendor": "Google",
            "path": "/Applications/Google Chrome.app"
        },
        {
            "name": "Mozilla Firefox",
            "version": "121.0",
            "vendor": "Mozilla",
            "path": "/Applications/Firefox.app"
        },
        {
            "name": "Microsoft Edge",
            "version": "120.0.2210.91",
            "vendor": "Microsoft",
            "path": "/Applications/Microsoft Edge.app"
        }
    ]


@pytest.fixture
def mock_cve_file(tmp_path):
    """Create a mock CVE data file for testing"""
    cve_file = tmp_path / "cve_data.json"
    cve_data = {
        "CVE-2024-1234": {
            "id": "CVE-2024-1234",
            "description": "Test vulnerability",
            "severity": "high",
            "cvss_score": 7.5,
            "affected_versions": ["<120.0.6099.110"],
            "references": ["https://example.com/cve-2024-1234"]
        }
    }
    
    import json
    with open(cve_file, 'w') as f:
        json.dump(cve_data, f)
    
    return str(cve_file)


# Test markers
def pytest_configure(config):
    """Configure pytest markers"""
    config.addinivalue_line(
        "markers", "slow: marks tests as slow (deselect with '-m \"not slow\"')"
    )
    config.addinivalue_line(
        "markers", "integration: marks tests as integration tests"
    )
    config.addinivalue_line(
        "markers", "performance: marks tests as performance tests"
    )


# Test collection
def pytest_collection_modifyitems(config, items):
    """Modify test collection to add markers"""
    for item in items:
        # Add slow marker to performance tests
        if "performance" in item.name or "test_enrichment_performance" in item.name:
            item.add_marker(pytest.mark.slow)
        
        # Add integration marker to integration tests
        if "integration" in item.name or "test_full_enrichment_pipeline" in item.name:
            item.add_marker(pytest.mark.integration)
        
        # Add performance marker to performance tests
        if "performance" in item.name:
            item.add_marker(pytest.mark.performance)

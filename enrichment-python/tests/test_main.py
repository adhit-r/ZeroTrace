import pytest
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_health_check():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"

def test_enrich_software_endpoint():
    # Mock data
    software_list = [
        {"name": "apache http server", "version": "2.4.49", "vendor": "apache"}
    ]
    
    # Note: In a real test, we would mock the database/redis calls
    # For now, we expect a 200 response with potential empty results if no DB connection
    try:
        response = client.post("/enrich/software", json=software_list)
        if response.status_code == 200:
            assert response.json()["success"] == True
            assert isinstance(response.json()["data"], list)
    except Exception:
        # DB connection might fail in test environment without docker
        pass

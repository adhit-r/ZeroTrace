import os
from typing import Optional, List, Dict, Any
try:
    from pydantic_settings import BaseSettings
    from pydantic import Field, AnyHttpUrl, validator
except ImportError:
    from pydantic import BaseSettings, Field, AnyHttpUrl, validator

class Settings(BaseSettings):
    """
    Unified Application Configuration
    Combines settings from all previous config files
    """
    # Application Settings
    app_name: str = "ZeroTrace Enrichment Service"
    environment: str = Field(default="production", env="ENVIRONMENT")
    debug: bool = Field(default=False, env="DEBUG")
    log_level: str = Field(default="INFO", env="LOG_LEVEL")
    
    # API Settings
    enrichment_host: str = Field(default="0.0.0.0", env="ENRICHMENT_HOST")
    enrichment_port: int = Field(default=8000, env="ENRICHMENT_PORT")
    cors_origins: str = Field(default="*", env="CORS_ORIGINS")
    enable_auth: bool = Field(default=False, env="ENABLE_AUTH")
    api_key: Optional[str] = Field(default=None, env="API_KEY")

    # Database Settings (PostgreSQL)
    database_url: Optional[str] = Field(default=None, env="DATABASE_URL")
    db_host: str = Field(default="postgres", env="DB_HOST")
    db_port: int = Field(default=5432, env="DB_PORT")
    db_name: str = Field(default="cve_db", env="DB_NAME")
    db_user: str = Field(default="postgres", env="DB_USER")
    db_password: str = Field(default="postgres", env="DB_PASSWORD")
    db_pool_min_size: int = Field(default=5, env="DB_POOL_MIN_SIZE")
    db_pool_max_size: int = Field(default=20, env="DB_POOL_MAX_SIZE")

    # Redis Settings (Cache & Vector Search)
    redis_url: Optional[str] = Field(default=None, env="REDIS_URL")
    redis_host: str = Field(default="redis", env="REDIS_HOST")
    redis_port: int = Field(default=6379, env="REDIS_PORT")
    redis_password: Optional[str] = Field(default=None, env="REDIS_PASSWORD")
    redis_db: int = Field(default=0, env="REDIS_DB")

    # Cache Settings
    cache_enabled: bool = Field(default=True, env="CACHE_ENABLED")
    cache_ttl: int = Field(default=3600, env="CACHE_TTL")
    cache_strategy: str = Field(default="lru", env="CACHE_STRATEGY")

    # Performance Settings
    uvloop_enabled: bool = Field(default=True, env="UVLOOP_ENABLED")
    orjson_enabled: bool = Field(default=True, env="ORJSON_ENABLED")
    max_concurrent_requests: int = Field(default=1000, env="MAX_CONCURRENT_REQUESTS")
    
    # AI/ML Settings
    enable_ai_matcher: bool = Field(default=True, env="ENABLE_AI_MATCHER")
    model_name: str = Field(default="all-MiniLM-L6-v2", env="MODEL_NAME")
    pgvector_enabled: bool = Field(default=True, env="PGVECTOR_ENABLED")
    
    # External APIs
    nvd_api_key: Optional[str] = Field(default=None, env="NVD_API_KEY")
    nvd_api_base_url: str = Field(default="https://services.nvd.nist.gov/rest/json/cves/2.0", env="NVD_API_BASE_URL")
    
    # CPE Guesser Integration
    cpe_guesser_url: str = Field(default="http://cpe-guesser:8000", env="CPE_GUESSER_URL")
    cpe_guesser_enabled: bool = Field(default=True, env="CPE_GUESSER_ENABLED")
    
    # Threat Intelligence Feeds
    cisa_kev_enabled: bool = Field(default=True, env="CISA_KEV_ENABLED")  # Enabled by default
    mitre_attack_enabled: bool = Field(default=False, env="MITRE_ATTACK_ENABLED")
    mitre_attack_data_url: str = Field(
        default="https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json",
        env="MITRE_ATTACK_DATA_URL"
    )
    alienvault_otx_enabled: bool = Field(default=False, env="ALIENVAULT_OTX_ENABLED")
    alienvault_otx_api_key: Optional[str] = Field(default=None, env="ALIENVAULT_OTX_API_KEY")
    opencve_enabled: bool = Field(default=False, env="OPENCVE_ENABLED")
    opencve_api_url: str = Field(default="http://localhost:8000", env="OPENCVE_API_URL")
    opencve_api_key: Optional[str] = Field(default=None, env="OPENCVE_API_KEY")

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"
        case_sensitive = False
        extra = "ignore"  # Ignore extra fields

    def get_database_dsn(self) -> str:
        if self.database_url:
            return self.database_url
        return f"postgresql://{self.db_user}:{self.db_password}@{self.db_host}:{self.db_port}/{self.db_name}"

    def get_redis_dsn(self) -> str:
        if self.redis_url:
            return self.redis_url
        auth = f":{self.redis_password}@" if self.redis_password else ""
        return f"redis://{auth}{self.redis_host}:{self.redis_port}/{self.redis_db}"
    
    def get_cors_origins_list(self) -> List[str]:
        if self.cors_origins == "*":
            return ["*"]
        return [origin.strip() for origin in self.cors_origins.split(",")]

settings = Settings()


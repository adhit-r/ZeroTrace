import logging
from typing import List, Dict, Optional, Tuple
from ..services.cpe_matcher import semantic_matcher
from ..core.config import settings

from ..core.logging import get_logger

logger = get_logger(__name__)

class HybridCPEMatcher:
    """
    Unified Hybrid CPE Matcher
    
    Legacy wrapper that redirects to the new enterprise SemanticCPEMatcher service.
    Kept for backward compatibility with existing code that imports this class.
    
    Strategies:
    1. Exact Match (L1): via SemanticCPEMatcher (which uses Redis/CPE Guesser logic internally)
    2. Semantic Match (L2): via SemanticCPEMatcher (pgvector)
    """
    
    def __init__(self):
        self.matcher = semantic_matcher
        self._initialized = False

    async def initialize(self):
        if not self._initialized:
            await self.matcher.initialize()
            self._initialized = True

    async def match_software_to_cpe(
        self,
        software_name: str,
        version: str = None,
        vendor: str = None
    ) -> List[Dict]:
        """
        Match software to CPE identifiers.
        """
        if not self._initialized:
            await self.initialize()

        # Build search query from parts
        query_parts = []
        if vendor: query_parts.append(vendor)
        query_parts.append(software_name)
        if version: query_parts.append(version)
        
        full_query = " ".join(query_parts)
        
        # Use new semantic search service
        # It handles both exact and semantic matching logic
        results = await self.matcher.search(full_query)
        
        # Adapt output format to legacy expectations if needed
        adapted_results = []
        for res in results:
            adapted_results.append({
                'cpe_string': res['cpe'],
                'confidence': self._score_to_confidence(res['score']),
                'combined_confidence': res['score'],
                'source': res['method']
            })
            
        return adapted_results

    def _score_to_confidence(self, score: float) -> str:
        if score > 0.85: return 'HIGH'
        if score > 0.65: return 'MEDIUM'
        return 'LOW'

# Global instance
hybrid_cpe_matcher = HybridCPEMatcher()

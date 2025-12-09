import asyncio
import numpy as np
from typing import List, Dict, Optional, Tuple
from ..core.config import settings
from ..core.logging import get_logger
from ..core.database import db_manager

logger = get_logger(__name__)

class SemanticCPEMatcher:
    """
    Enterprise CPE Matcher using Hybrid Search:
    1. Exact Match (L1): Fast Redis/Valkey lookup (via cpe-guesser)
    2. Semantic Match (L2): Vector Search using pgvector + SentenceTransformers
    
    This solves the OOM issue by offloading the index to the database.
    """
    
    def __init__(self):
        self.model = None
        self.model_name = settings.model_name
        self.device = "cpu" # Default to CPU for inference safety
        self._initialized = False

    async def initialize(self):
        """Lazy load the AI model to avoid startup blocking"""
        if self._initialized:
            return

        if not settings.enable_ai_matcher:
            logger.info("AI Matcher disabled via config")
            return

        try:
            # Import here to avoid heavy dependency if not used
            from sentence_transformers import SentenceTransformer
            import torch
            
            logger.info("Loading AI model for CPE matching", model=self.model_name)
            
            # Use GPU if available (for training mostly)
            if torch.cuda.is_available():
                self.device = "cuda"
            elif torch.backends.mps.is_available():
                self.device = "mps"
            
            # Load optimized model
            self.model = SentenceTransformer(self.model_name, device=self.device)
            self._initialized = True
            logger.info("AI model loaded successfully", device=self.device)
            
        except ImportError:
            logger.warning("sentence-transformers not installed. AI matching disabled.")
        except Exception as e:
            logger.error("Failed to load AI model", error=str(e))

    def encode(self, text: str) -> List[float]:
        """Convert text to vector embedding"""
        if not self.model:
            return []
        
        try:
            # Encode and convert numpy array to list
            embedding = self.model.encode(text, convert_to_numpy=True)
            return embedding.tolist()
        except Exception as e:
            logger.error("Encoding error", error=str(e))
            return []

    async def search(self, query: str, limit: int = 10) -> List[Dict]:
        """
        Perform semantic search for CPEs using pgvector
        """
        if not self._initialized:
            await self.initialize()

        if not self.model:
            return []

        start_time = asyncio.get_event_loop().time()
        
        try:
            # 1. Generate embedding
            query_vector = self.encode(query)
            if not query_vector:
                return []

            # 2. Vector Search in Postgres
            # Using Cosine Similarity (<=> operator in pgvector is distance, so we order by it ASC)
            # 1 - distance = similarity
            sql = """
                SELECT 
                    cpe_string, 
                    (1 - (embedding <=> $1::vector)) as similarity 
                FROM cpe_embeddings
                ORDER BY embedding <=> $1::vector
                LIMIT $2
            """
            
            # Convert list to string representation for Postgres vector
            vector_str = str(query_vector)
            
            results = await db_manager.fetch_all(sql, vector_str, limit)
            
            matches = []
            for row in results:
                matches.append({
                    "cpe": row["cpe_string"],
                    "score": float(row["similarity"]),
                    "method": "semantic"
                })

            duration = (asyncio.get_event_loop().time() - start_time) * 1000
            logger.info("Semantic search completed", matches=len(matches), duration_ms=duration)
            
            return matches

        except Exception as e:
            logger.error("Semantic search failed", error=str(e))
            return []

    async def match_software(self, vendor: str, product: str, version: str) -> List[Dict]:
        """
        High-level matching logic:
        1. Try CPE Guesser first (L1 - exact match)
        2. Fallback to semantic search (L2 - if pgvector enabled)
        3. Return empty if both fail
        """
        # Try CPE Guesser first (L1 - exact match via direct Valkey access)
        # Use direct Valkey access since data is already seeded in DB 8
        if settings.cpe_guesser_enabled:
            try:
                # Direct Valkey access (faster than HTTP service)
                import valkey
                from ..core.config import settings as app_settings
                
                # Connect to Valkey DB 8 where CPE data is stored
                rdb = valkey.Valkey(
                    host=app_settings.redis_host,
                    port=app_settings.redis_port,
                    db=8,  # CPE Guesser uses DB 8
                    password=app_settings.redis_password,
                    decode_responses=True,
                    socket_timeout=2.0
                )
                
                # Construct search words from vendor and product ONLY
                # Version is NOT used for CPE lookup - it's matched against CVE version ranges later
                words = []
                if vendor:
                    # Split vendor into words (e.g., "Microsoft Corporation" -> ["microsoft", "corporation"])
                    words.extend([w.lower() for w in vendor.split() if len(w) > 2])
                if product:
                    # Split product into words (e.g., "nginx server" -> ["nginx", "server"])
                    words.extend([w.lower() for w in product.split() if len(w) > 2])
                
                # Remove duplicates while preserving order
                words = list(dict.fromkeys(words))
                
                if words:
                    # CPE Guesser algorithm: find intersection of word sets
                    word_keys = [f"w:{word}" for word in words]
                    
                    if len(word_keys) == 1:
                        # Single word - get sorted set directly (s: prefix)
                        results = rdb.zrevrange(f"s:{words[0]}", 0, 9, withscores=True)
                        if results:
                            matches = []
                            for cpe, score in results[:5]:
                                matches.append({
                                    "cpe": cpe,
                                    "score": float(score) / 1000.0 if score else 0.9,
                                    "method": "cpe_guesser_direct"
                                })
                            if matches:
                                logger.info("CPE Guesser (direct) found matches", count=len(matches))
                                return matches
                    else:
                        # Multiple words - use sinter on w: keys (regular sets)
                        matched_cpes = rdb.sinter(*word_keys)
                        
                        if matched_cpes:
                            # Get ranks from rank:cpe sorted set
                            ranked_results = []
                            for cpe in matched_cpes:
                                try:
                                    rank = rdb.zrank("rank:cpe", cpe)
                                    if rank is not None:
                                        ranked_results.append((rank, cpe))
                                    else:
                                        ranked_results.append((0, cpe))
                                except:
                                    ranked_results.append((0, cpe))
                            
                            # Sort by rank (higher is better)
                            ranked_results.sort(key=lambda x: x[0], reverse=True)
                            
                            matches = []
                            for rank, cpe in ranked_results[:5]:
                                matches.append({
                                    "cpe": cpe,
                                    "score": 0.95 - (0.05 * matches.__len__()),  # High confidence for exact matches
                                    "method": "cpe_guesser_direct"
                                })
                            if matches:
                                logger.info("CPE Guesser (direct) found matches", count=len(matches))
                                return matches
            except Exception as e:
                logger.debug("CPE Guesser (direct) not available, trying HTTP service", error=str(e))
                # Fallback to HTTP service if direct access fails
                try:
                    from ..cpe_guesser_client import get_cpe_guesser_client
                    words = []
                    if vendor:
                        words.append(vendor.lower())
                    if product:
                        words.append(product.lower())
                    if version:
                        words.append(version.lower())
                    if words:
                        cpe_client = await get_cpe_guesser_client()
                        results = await cpe_client.guess_cpe(words)
                        if results:
                            matches = []
                            for rank, cpe in results[:5]:
                                matches.append({
                                    "cpe": cpe,
                                    "score": 1.0 - (rank / 1000.0) if rank else 0.9,
                                    "method": "cpe_guesser_http"
                                })
                            if matches:
                                logger.info("CPE Guesser (HTTP) found matches", count=len(matches))
                                return matches
                except Exception as e2:
                    logger.debug("CPE Guesser HTTP also failed", error=str(e2))
        
        # Fallback to semantic search (L2) only if pgvector is enabled
        if settings.pgvector_enabled:
            query = f"{vendor} {product} {version}".strip()
            return await self.search(query)
        
        # No matches found
        logger.debug("No CPE matches found", vendor=vendor, product=product, version=version)
        return []

# Global instance
semantic_matcher = SemanticCPEMatcher()


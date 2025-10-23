"""
CPE Matching Engine using AI/ML techniques
Matches software packages to CPE identifiers for better CVE correlation
Enhanced with version range matching and confidence scoring
"""

import json
import logging
import re
from typing import Dict, List, Optional, Tuple
from pathlib import Path
import numpy as np
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import pickle
from ..version_matcher import cpe_version_matcher, version_range

logger = logging.getLogger(__name__)

class CPEMatcher:
    def __init__(self, data_dir: str = None):
        self.data_dir = Path(data_dir or Path(__file__).parent.parent)
        self.cpe_index_file = self.data_dir / 'cpe_index.pkl'
        self.vectorizer = TfidfVectorizer(
            max_features=10000,
            stop_words='english',
            ngram_range=(1, 3)
        )
        self.cpe_vectors = None
        self.cpe_data = []
        self._load_or_build_index()
    
    def _load_or_build_index(self):
        """Load existing CPE index or build new one"""
        if self.cpe_index_file.exists():
            try:
                with open(self.cpe_index_file, 'rb') as f:
                    index_data = pickle.load(f)
                    self.cpe_vectors = index_data['vectors']
                    self.cpe_data = index_data['data']
                    self.vectorizer = index_data['vectorizer']
                logger.info(f"Loaded CPE index with {len(self.cpe_data)} entries")
            except Exception as e:
                logger.error(f"Error loading CPE index: {e}")
                self._build_index()
        else:
            self._build_index()
    
    def _build_index(self):
        """Build CPE index from CVE data"""
        logger.info("Building CPE index...")
        
        # Load CVE data
        cve_file = self.data_dir / 'cve_data.json'
        if not cve_file.exists():
            logger.warning("No CVE data found, creating empty index")
            self.cpe_data = []
            self.cpe_vectors = np.array([])
            return
        
        with open(cve_file, 'r') as f:
            cve_data = json.load(f)
        
        # Extract CPE information from CVE data
        cpe_entries = []
        for cve in cve_data:
            cpe_info = self._extract_cpe_from_cve(cve)
            if cpe_info:
                cpe_entries.extend(cpe_info)
        
        if not cpe_entries:
            logger.warning("No CPE data found in CVE database")
            self.cpe_data = []
            self.cpe_vectors = np.array([])
            return
        
        # Create text representations for vectorization
        texts = []
        for entry in cpe_entries:
            text = f"{entry['vendor']} {entry['product']} {entry['version']} {entry['description']}"
            texts.append(text)
            entry['text'] = text
        
        # Fit vectorizer and create vectors
        self.cpe_vectors = self.vectorizer.fit_transform(texts)
        self.cpe_data = cpe_entries
        
        # Save index
        self._save_index()
        
        logger.info(f"Built CPE index with {len(cpe_entries)} entries")
    
    def _extract_cpe_from_cve(self, cve: Dict) -> List[Dict]:
        """Extract CPE information from CVE data"""
        cpe_entries = []
        
        # Extract from CVE configurations
        configurations = cve.get('configurations', [])
        for config in configurations:
            nodes = config.get('nodes', [])
            for node in nodes:
                cpe_matches = node.get('cpeMatch', [])
                for match in cpe_matches:
                    cpe_string = match.get('criteria', '')
                    if cpe_string:
                        cpe_info = self._parse_cpe_string(cpe_string)
                        if cpe_info:
                            cpe_info['cve_id'] = cve.get('id', '')
                            cpe_info['description'] = cve.get('descriptions', [{}])[0].get('value', '')
                            cpe_entries.append(cpe_info)
        
        return cpe_entries
    
    def _parse_cpe_string(self, cpe_string: str) -> Optional[Dict]:
        """Parse CPE string into components"""
        # CPE format: cpe:2.3:a:vendor:product:version:update:edition:language:sw_edition:target_sw:target_hw:other
        parts = cpe_string.split(':')
        if len(parts) < 6:
            return None
        
        return {
            'cpe_string': cpe_string,
            'vendor': parts[3] if len(parts) > 3 else '',
            'product': parts[4] if len(parts) > 4 else '',
            'version': parts[5] if len(parts) > 5 else '',
            'update': parts[6] if len(parts) > 6 else '',
            'edition': parts[7] if len(parts) > 7 else '',
        }
    
    def _save_index(self):
        """Save CPE index to file"""
        try:
            index_data = {
                'vectors': self.cpe_vectors,
                'data': self.cpe_data,
                'vectorizer': self.vectorizer
            }
            with open(self.cpe_index_file, 'wb') as f:
                pickle.dump(index_data, f)
            logger.info("Saved CPE index")
        except Exception as e:
            logger.error(f"Error saving CPE index: {e}")
    
    def match_software_to_cpe(self, software_name: str, version: str = None, vendor: str = None) -> List[Dict]:
        """
        Match software to CPE identifiers using AI/ML similarity with version range matching
        """
        if not self.cpe_data or self.cpe_vectors is None:
            return []
        
        # Create query text
        query_text = f"{vendor or ''} {software_name} {version or ''}"
        
        # Vectorize query
        query_vector = self.vectorizer.transform([query_text])
        
        # Calculate similarities
        similarities = cosine_similarity(query_vector, self.cpe_vectors).flatten()
        
        # Get top matches
        top_indices = np.argsort(similarities)[::-1][:20]  # Top 20 matches for version filtering
        
        matches = []
        for idx in top_indices:
            if similarities[idx] > 0.1:  # Minimum similarity threshold
                match = self.cpe_data[idx].copy()
                match['similarity_score'] = float(similarities[idx])
                
                # Enhanced version matching
                version_match, version_confidence = self._match_version_range(
                    version, match.get('version', ''), match.get('cpe_string', '')
                )
                
                # Combine similarity and version confidence
                combined_confidence = (similarities[idx] * 0.6) + (version_confidence * 0.4)
                match['version_match'] = version_match
                match['version_confidence'] = version_confidence
                match['combined_confidence'] = combined_confidence
                match['confidence'] = self._calculate_confidence(combined_confidence)
                
                matches.append(match)
        
        # Sort by combined confidence and filter
        matches.sort(key=lambda x: x['combined_confidence'], reverse=True)
        return matches[:10]  # Return top 10 matches
    
    def _match_version_range(self, software_version: str, cpe_version: str, cpe_string: str) -> Tuple[bool, float]:
        """Match software version against CPE version range"""
        if not software_version or not cpe_version:
            return True, 0.5  # No version info, medium confidence
        
        # Use CPE version matcher for range matching
        return cpe_version_matcher.match_cpe_version(software_version, cpe_string)
    
    def _calculate_confidence(self, score: float) -> str:
        """Calculate confidence level based on combined score"""
        if score > 0.8:
            return 'HIGH'
        elif score > 0.6:
            return 'MEDIUM'
        elif score > 0.4:
            return 'LOW'
        else:
            return 'VERY_LOW'
    
    def get_cpe_for_software(self, software_name: str, version: str = None, vendor: str = None) -> Optional[str]:
        """
        Get the best CPE identifier for software with enhanced version matching
        """
        matches = self.match_software_to_cpe(software_name, version, vendor)
        
        if not matches:
            return None
        
        # Return the highest confidence match with version compatibility
        for match in matches:
            if (match['confidence'] in ['HIGH', 'MEDIUM'] and 
                match.get('version_match', True) and 
                match.get('version_confidence', 0) > 0.5):
                return match['cpe_string']
        
        # Fallback to best match even if version doesn't match perfectly
        best_match = matches[0]
        if best_match['confidence'] in ['HIGH', 'MEDIUM']:
            return best_match['cpe_string']
        
        return None
    
    def update_index(self):
        """Update CPE index with new data"""
        logger.info("Updating CPE index...")
        self._build_index()

# Global instance
cpe_matcher = CPEMatcher()


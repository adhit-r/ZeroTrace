"""
Semantic Version Comparison and Matching
Handles various versioning schemes and range operators
"""

import re
from typing import List, Tuple, Optional, Union
from enum import Enum

class VersionScheme(Enum):
    SEMVER = "semver"  # 1.2.3
    CALVER = "calver"  # 2023.12.01
    PEP440 = "pep440"  # Python versioning
    DEBIAN = "debian"  # 1:2.3.4-5
    RPM = "rpm"       # 1.2.3-4.el8
    CUSTOM = "custom"  # Other schemes

class VersionComparison:
    """Handles version comparison for different schemes"""
    
    def __init__(self):
        self.semver_pattern = re.compile(r'^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$')
        self.calver_pattern = re.compile(r'^(\d{4})\.(\d{1,2})\.(\d{1,2})(?:\.(\d+))?(?:-([0-9A-Za-z-]+))?$')
        self.pep440_pattern = re.compile(r'^(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:([ab]|rc)(\d+))?(?:\.(post|dev)(\d+))?$')
    
    def detect_scheme(self, version: str) -> VersionScheme:
        """Detect the versioning scheme used"""
        if self.semver_pattern.match(version):
            return VersionScheme.SEMVER
        elif self.calver_pattern.match(version):
            return VersionScheme.CALVER
        elif self.pep440_pattern.match(version):
            return VersionScheme.PEP440
        elif ':' in version and '-' in version:
            return VersionScheme.DEBIAN
        elif version.count('-') >= 2:
            return VersionScheme.RPM
        else:
            return VersionScheme.CUSTOM
    
    def parse_semver(self, version: str) -> Tuple[int, int, int, str, str]:
        """Parse semantic version (1.2.3-alpha+build)"""
        match = self.semver_pattern.match(version)
        if not match:
            return (0, 0, 0, '', '')
        
        major, minor, patch, prerelease, build = match.groups()
        return (
            int(major or 0),
            int(minor or 0), 
            int(patch or 0),
            prerelease or '',
            build or ''
        )
    
    def parse_calver(self, version: str) -> Tuple[int, int, int, int, str]:
        """Parse calendar version (2023.12.01.1-alpha)"""
        match = self.calver_pattern.match(version)
        if not match:
            return (0, 0, 0, 0, '')
        
        year, month, day, micro, prerelease = match.groups()
        return (
            int(year or 0),
            int(month or 0),
            int(day or 0),
            int(micro or 0),
            prerelease or ''
        )
    
    def compare_semver(self, v1: str, v2: str) -> int:
        """Compare two semantic versions. Returns -1, 0, or 1"""
        p1 = self.parse_semver(v1)
        p2 = self.parse_semver(v2)
        
        # Compare major, minor, patch
        for i in range(3):
            if p1[i] < p2[i]:
                return -1
            elif p1[i] > p2[i]:
                return 1
        
        # Compare prerelease
        if p1[3] and not p2[3]:
            return -1  # prerelease is less than release
        elif not p1[3] and p2[3]:
            return 1
        elif p1[3] and p2[3]:
            return -1 if p1[3] < p2[3] else 1 if p1[3] > p2[3] else 0
        
        return 0
    
    def compare_calver(self, v1: str, v2: str) -> int:
        """Compare two calendar versions"""
        p1 = self.parse_calver(v1)
        p2 = self.parse_calver(v2)
        
        # Compare year, month, day, micro
        for i in range(4):
            if p1[i] < p2[i]:
                return -1
            elif p1[i] > p2[i]:
                return 1
        
        # Compare prerelease
        if p1[4] and not p2[4]:
            return -1
        elif not p1[4] and p2[4]:
            return 1
        elif p1[4] and p2[4]:
            return -1 if p1[4] < p2[4] else 1 if p1[4] > p2[4] else 0
        
        return 0
    
    def compare_versions(self, v1: str, v2: str) -> int:
        """Compare two versions using appropriate scheme"""
        scheme1 = self.detect_scheme(v1)
        scheme2 = self.detect_scheme(v2)
        
        # Use the more specific scheme if they differ
        if scheme1 == scheme2:
            if scheme1 == VersionScheme.SEMVER:
                return self.compare_semver(v1, v2)
            elif scheme1 == VersionScheme.CALVER:
                return self.compare_calver(v1, v2)
        
        # Fallback to string comparison
        return -1 if v1 < v2 else 1 if v1 > v2 else 0

class VersionRange:
    """Handles version range matching with various operators"""
    
    def __init__(self):
        self.comparison = VersionComparison()
        self.range_pattern = re.compile(r'([<>=!]+)\s*([^\s,]+)')
    
    def parse_range(self, range_str: str) -> List[Tuple[str, str]]:
        """Parse version range string into operator-version pairs"""
        if not range_str or range_str.strip() == '*':
            return []
        
        # Handle comma-separated ranges
        ranges = []
        for part in range_str.split(','):
            part = part.strip()
            if not part:
                continue
            
            # Handle "up to excluding" and "up to including"
            if part.startswith('up to '):
                version = part[6:].strip()
                if version.startswith('excluding '):
                    ranges.append(('<', version[10:].strip()))
                else:
                    ranges.append(('<=', version))
            else:
                # Parse operator and version
                match = self.range_pattern.match(part)
                if match:
                    operator, version = match.groups()
                    ranges.append((operator, version))
        
        return ranges
    
    def matches_version(self, version: str, range_str: str) -> bool:
        """Check if a version matches the given range"""
        if not range_str or range_str.strip() == '*':
            return True
        
        ranges = self.parse_range(range_str)
        if not ranges:
            return True
        
        for operator, range_version in ranges:
            if not self._check_constraint(version, operator, range_version):
                return False
        
        return True
    
    def _check_constraint(self, version: str, operator: str, range_version: str) -> bool:
        """Check if version satisfies a single constraint"""
        try:
            comparison = self.comparison.compare_versions(version, range_version)
            
            if operator == '<':
                return comparison < 0
            elif operator == '<=':
                return comparison <= 0
            elif operator == '>':
                return comparison > 0
            elif operator == '>=':
                return comparison >= 0
            elif operator == '=' or operator == '==':
                return comparison == 0
            elif operator == '!=':
                return comparison != 0
            elif operator.startswith('<='):
                return comparison <= 0
            elif operator.startswith('>='):
                return comparison >= 0
            elif operator.startswith('!='):
                return comparison != 0
            else:
                return False
                
        except Exception:
            # Fallback to string comparison for unknown formats
            if operator == '<':
                return version < range_version
            elif operator == '<=':
                return version <= range_version
            elif operator == '>':
                return version > range_version
            elif operator == '>=':
                return version >= range_version
            elif operator == '=' or operator == '==':
                return version == range_version
            elif operator == '!=':
                return version != range_version
            else:
                return False

class CPEVersionMatcher:
    """Matches software versions against CPE version ranges"""
    
    def __init__(self):
        self.version_range = VersionRange()
        self.comparison = VersionComparison()
    
    def extract_cpe_versions(self, cpe_string: str) -> List[str]:
        """Extract version components from CPE string"""
        # CPE format: cpe:2.3:a:vendor:product:version:update:edition:language:sw_edition:target_sw:target_hw:other
        parts = cpe_string.split(':')
        if len(parts) < 5:
            return []
        
        versions = []
        # Version is typically at index 4
        if len(parts) > 4 and parts[4] != '*':
            versions.append(parts[4])
        
        return versions
    
    def match_cpe_version(self, software_version: str, cpe_string: str) -> Tuple[bool, float]:
        """
        Match software version against CPE version.
        Returns (matches, confidence_score)
        """
        if not cpe_string or cpe_string == '*':
            return True, 0.5  # Wildcard match with medium confidence
        
        # Extract version from CPE
        cpe_versions = self.extract_cpe_versions(cpe_string)
        if not cpe_versions:
            return True, 0.5  # No version specified in CPE
        
        best_match = False
        best_confidence = 0.0
        
        for cpe_version in cpe_versions:
            # Direct version match
            if software_version == cpe_version:
                return True, 1.0  # Exact match
            
            # Check if software version matches CPE version range
            if self.version_range.matches_version(software_version, cpe_version):
                return True, 0.9  # Range match
            
            # Calculate similarity for partial matches
            similarity = self._calculate_version_similarity(software_version, cpe_version)
            if similarity > best_confidence:
                best_confidence = similarity
                best_match = similarity > 0.7  # Threshold for partial match
        
        return best_match, best_confidence
    
    def _calculate_version_similarity(self, v1: str, v2: str) -> float:
        """Calculate similarity between two version strings"""
        try:
            # Try semantic comparison first
            comparison = self.comparison.compare_versions(v1, v2)
            if comparison == 0:
                return 1.0
            
            # Calculate numeric similarity
            v1_parts = re.findall(r'\d+', v1)
            v2_parts = re.findall(r'\d+', v2)
            
            if not v1_parts or not v2_parts:
                return 0.0
            
            # Compare numeric parts
            max_parts = max(len(v1_parts), len(v2_parts))
            matches = 0
            
            for i in range(max_parts):
                v1_num = int(v1_parts[i]) if i < len(v1_parts) else 0
                v2_num = int(v2_parts[i]) if i < len(v2_parts) else 0
                
                if v1_num == v2_num:
                    matches += 1
                else:
                    # Partial credit for close numbers
                    diff = abs(v1_num - v2_num)
                    if diff <= 1:
                        matches += 0.5
            
            return matches / max_parts
            
        except Exception:
            # Fallback to string similarity
            return 1.0 if v1 == v2 else 0.0

# Global instances for easy import
version_comparison = VersionComparison()
version_range = VersionRange()
cpe_version_matcher = CPEVersionMatcher()

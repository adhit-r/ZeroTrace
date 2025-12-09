#!/usr/bin/env python3
"""
Download and store CISA KEV (Known Exploited Vulnerabilities) catalog
"""

import asyncio
import json
import httpx
from pathlib import Path
from typing import Dict, List

KEV_URL = "https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json"
KEV_FILE = Path(__file__).parent.parent / "kev_catalog.json"

async def download_kev():
    """Download CISA KEV catalog"""
    print(f"Downloading CISA KEV catalog from {KEV_URL}...")
    
    async with httpx.AsyncClient() as client:
        response = await client.get(KEV_URL)
        if response.status_code == 200:
            data = response.json()
            vulnerabilities = data.get('vulnerabilities', [])
            
            print(f"Downloaded {len(vulnerabilities)} KEV entries")
            
            # Save to file
            with open(KEV_FILE, 'w') as f:
                json.dump(data, f, indent=2)
            
            print(f"Saved to {KEV_FILE}")
            
            # Create CVE -> KEV mapping summary
            cve_ids = [v.get('cveID') for v in vulnerabilities if v.get('cveID')]
            print(f"Total unique CVEs with known exploits: {len(cve_ids)}")
            
            # Show sample
            if vulnerabilities:
                sample = vulnerabilities[0]
                print(f"\nSample KEV entry:")
                print(f"   CVE: {sample.get('cveID')}")
                print(f"   Vendor: {sample.get('vendorProject')}")
                print(f"   Product: {sample.get('product')}")
                print(f"   Date Added: {sample.get('dateAdded')}")
                print(f"   Description: {sample.get('shortDescription', '')[:80]}...")
            
            return data
        else:
            print(f"ERROR: Failed to download: HTTP {response.status_code}")
            return None

if __name__ == '__main__':
    asyncio.run(download_kev())


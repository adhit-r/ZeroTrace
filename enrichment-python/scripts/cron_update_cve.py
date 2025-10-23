#!/usr/bin/env python3
"""
Cron-compatible CVE data updater
Designed to be run weekly via cron job
"""

import os
import sys
import logging
from pathlib import Path

# Add the parent directory to the path so we can import our modules
sys.path.insert(0, str(Path(__file__).parent.parent))

from scripts.update_cve_data import CVEDataUpdater
import asyncio

# Configure logging for cron
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - CVE-UPDATE - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('/tmp/cve_update.log'),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)

async def main():
    """Main function for cron execution"""
    try:
        # Get configuration from environment
        data_dir = os.getenv('CVE_DATA_DIR', '/opt/zerotrace/cve_data')
        nvd_api_key = os.getenv('NVD_API_KEY', '')
        
        logger.info("Starting weekly CVE data update...")
        
        # Create updater
        updater = CVEDataUpdater(data_dir=data_dir, nvd_api_key=nvd_api_key)
        
        # Run update
        success = await updater.update_cve_data(force_update=False)
        
        if success:
            logger.info("Weekly CVE data update completed successfully")
            return 0
        else:
            logger.error("Weekly CVE data update failed")
            return 1
            
    except Exception as e:
        logger.error(f"Unexpected error in CVE update: {e}")
        return 1

if __name__ == '__main__':
    exit_code = asyncio.run(main())
    sys.exit(exit_code)


#!/bin/bash
# Download all CVEs from NVD (390K+ CVEs)
# This script handles the full download with proper rate limiting

set -e

echo " Starting full NVD CVE download (390K+ CVEs)"
echo ""
echo "️  This will take a long time:"
echo "   - Without API key: ~65 hours (5 requests/30s)"
echo "   - With API key: ~6.5 hours (50 requests/30s)"
echo ""
echo " Tip: Get a free NVD API key at: https://nvd.nist.gov/developers/request-an-api-key"
echo ""

# Check if API key is set
if [ -z "$NVD_API_KEY" ]; then
    echo "️  No NVD_API_KEY found. Using rate-limited access (5 req/30s)"
    echo "   Set NVD_API_KEY environment variable for faster downloads"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
else
    echo " NVD_API_KEY found. Using enhanced rate limits (50 req/30s)"
    echo ""
fi

# Export API key if set
export NVD_API_KEY

# Run the update script
cd enrichment-python/scripts

echo " Starting download..."
echo "   Progress will be logged to /tmp/cve_download.log"
echo ""

# Run in background and show progress
python3 update_cve_data.py --force-full > /tmp/cve_download.log 2>&1 &
PID=$!

echo "Download started (PID: $PID)"
echo ""
echo "Monitor progress with:"
echo "  tail -f /tmp/cve_download.log"
echo ""
echo "Or check current status:"
echo "  tail -20 /tmp/cve_download.log | grep -E 'Fetched|Total|ERROR'"
echo ""

# Wait for process
wait $PID

echo ""
echo " Download complete!"
echo ""
echo "Next step: Migrate to PostgreSQL"
echo "  cd enrichment-python/scripts"
echo "  python3 migrate_to_postgres.py"


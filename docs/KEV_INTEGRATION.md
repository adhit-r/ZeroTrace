# CISA KEV (Known Exploited Vulnerabilities) Integration

## Overview

CISA KEV catalog maps CVEs to known exploits - vulnerabilities that are **actively exploited in the wild**. This is critical threat intelligence for prioritizing remediation.

## What is KEV?

- **Source**: CISA (Cybersecurity and Infrastructure Security Agency)
- **Catalog**: Known Exploited Vulnerabilities
- **Size**: ~1,468 CVEs (as of Dec 2025)
- **Update Frequency**: Daily
- **URL**: https://www.cisa.gov/known-exploited-vulnerabilities-catalog

## Integration Status

**Fully Integrated**

- KEV service added to `threat_intel.py`
- Integrated into enrichment service
- Automatically marks KEV CVEs as `known_exploited: true`
- Elevates KEV CVEs to "critical" severity
- Loads on service startup

## KEV Data Structure

Each KEV entry contains:
- `cveID`: CVE identifier
- `vendorProject`: Vendor name
- `product`: Product name
- `vulnerabilityName`: Vulnerability name
- `dateAdded`: When added to KEV catalog
- `shortDescription`: Description
- `requiredAction`: Required remediation action
- `dueDate`: Remediation deadline
- `knownRansomwareCampaignUse`: If used in ransomware

## Usage

### Automatic Integration

KEV is **enabled by default**. When enriching vulnerabilities:

```python
# Enrichment automatically includes KEV data
vulnerabilities = await enrichment_service.enrich([{
    "name": "nginx",
    "version": "1.18",
    "vendor": "nginx"
}])

# If a CVE is in KEV, it will have:
vuln["known_exploited"] = True
vuln["threat_intel"]["cisa_kev"] = {
    "cveID": "CVE-2021-23017",
    "vendorProject": "F5",
    "product": "nginx",
    "dateAdded": "2021-06-01",
    ...
}
```

### Manual KEV Lookup

```python
from app.services.threat_intel import cisa_kev_service

# Check if CVE is known exploited
is_exploited = await cisa_kev_service.is_known_exploited("CVE-2021-23017")

# Get full KEV info
kev_info = await cisa_kev_service.get_kev_info("CVE-2021-23017")
```

### Download KEV Catalog

```bash
cd enrichment-python/scripts
python3 download_kev.py
```

This downloads the full KEV catalog to `enrichment-python/kev_catalog.json`.

## API Response Example

When a CVE is in KEV, enrichment response includes:

```json
{
  "name": "nginx",
  "version": "1.18",
  "vulnerabilities": [
    {
      "id": "CVE-2021-23017",
      "description": "...",
      "severity": "critical",
      "known_exploited": true,
      "threat_intel": {
        "cisa_kev": {
          "cveID": "CVE-2021-23017",
          "vendorProject": "F5",
          "product": "nginx",
          "vulnerabilityName": "nginx Remote Code Execution",
          "dateAdded": "2021-06-01",
          "shortDescription": "...",
          "requiredAction": "Apply updates per vendor instructions",
          "dueDate": "2021-06-15",
          "knownRansomwareCampaignUse": "Unknown"
        }
      }
    }
  ]
}
```

## Configuration

In `enrichment-python/.env`:

```bash
# KEV is enabled by default
CISA_KEV_ENABLED=true
```

## Benefits

1. **Prioritization**: KEV CVEs are automatically marked as critical
2. **Actionable Intelligence**: Know which CVEs are actively exploited
3. **Compliance**: CISA requires federal agencies to remediate KEV CVEs
4. **Risk Assessment**: Higher risk score for known exploited vulnerabilities

## Statistics

- **Total KEV Entries**: ~1,468 CVEs
- **Update Frequency**: Daily
- **Coverage**: All actively exploited CVEs tracked by CISA

## References

- CISA KEV Catalog: https://www.cisa.gov/known-exploited-vulnerabilities-catalog
- KEV JSON Feed: https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json
- CISA Binding Operational Directive: https://www.cisa.gov/binding-operational-directive-22-01


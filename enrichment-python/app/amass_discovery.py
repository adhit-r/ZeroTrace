import asyncio
import subprocess
import json
import time
from typing import List, Dict, Optional
from dataclasses import dataclass
from datetime import datetime
import aiohttp
import aiofiles
import os

@dataclass
class AmassResult:
    """Represents a discovered asset from OWASP Amass"""
    domain: str
    subdomain: str
    ip_address: str
    asn: Optional[int] = None
    as_name: Optional[str] = None
    country: Optional[str] = None
    city: Optional[str] = None
    service: Optional[str] = None
    port: Optional[int] = None
    protocol: Optional[str] = None
    source: str = "amass"
    discovered_at: datetime = None
    
    def __post_init__(self):
        if self.discovered_at is None:
            self.discovered_at = datetime.now()

class AmassDiscovery:
    """Handles external asset discovery using OWASP Amass"""
    
    def __init__(self, amass_path: str = "amass", output_dir: str = "/tmp/amass"):
        self.amass_path = amass_path
        self.output_dir = output_dir
        os.makedirs(output_dir, exist_ok=True)
    
    async def discover_domain(self, domain: str, company_id: str) -> List[AmassResult]:
        """
        Discover assets for a given domain using OWASP Amass
        
        Args:
            domain: The domain to scan
            company_id: Company identifier for tracking
            
        Returns:
            List of discovered assets
        """
        try:
            # Create output file path
            timestamp = int(time.time())
            output_file = f"{self.output_dir}/amass_{company_id}_{domain}_{timestamp}.json"
            
            # Run Amass enumeration
            cmd = [
                self.amass_path,
                "enum",
                "-d", domain,
                "-json", output_file,
                "-timeout", "5",
                "-max-dns-queries", "1000",
                "-active"
            ]
            
            # Execute Amass
            process = await asyncio.create_subprocess_exec(
                *cmd,
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            
            stdout, stderr = await process.communicate()
            
            if process.returncode != 0:
                print(f"Amass error for {domain}: {stderr.decode()}")
                return []
            
            # Parse results
            results = await self._parse_amass_output(output_file, domain)
            
            # Clean up output file
            try:
                os.remove(output_file)
            except:
                pass
            
            return results
            
        except Exception as e:
            print(f"Error running Amass for {domain}: {e}")
            return []
    
    async def _parse_amass_output(self, output_file: str, domain: str) -> List[AmassResult]:
        """Parse Amass JSON output file"""
        results = []
        
        try:
            async with aiofiles.open(output_file, 'r') as f:
                async for line in f:
                    line = line.strip()
                    if not line:
                        continue
                    
                    try:
                        data = json.loads(line)
                        result = self._parse_amass_entry(data, domain)
                        if result:
                            results.append(result)
                    except json.JSONDecodeError:
                        continue
                        
        except Exception as e:
            print(f"Error parsing Amass output: {e}")
        
        return results
    
    def _parse_amass_entry(self, data: Dict, domain: str) -> Optional[AmassResult]:
        """Parse a single Amass entry"""
        try:
            # Extract basic information
            name = data.get('name', '')
            addresses = data.get('addresses', [])
            
            if not addresses:
                return None
            
            # Get first IP address
            ip_address = addresses[0].get('ip', '')
            if not ip_address:
                return None
            
            # Extract additional information
            asn = None
            as_name = None
            country = None
            city = None
            
            if 'asn' in addresses[0]:
                asn = addresses[0]['asn']
            
            if 'as_name' in addresses[0]:
                as_name = addresses[0]['as_name']
            
            if 'country' in addresses[0]:
                country = addresses[0]['country']
            
            if 'city' in addresses[0]:
                city = addresses[0]['city']
            
            # Determine if it's a subdomain
            subdomain = name if name != domain else domain
            
            # Extract port and service information
            port = None
            protocol = None
            service = None
            
            if 'ports' in data:
                ports = data['ports']
                if ports:
                    port_info = ports[0]
                    port = port_info.get('port')
                    protocol = port_info.get('protocol', 'tcp')
                    
                    # Try to determine service
                    if port:
                        service = self._get_service_name(port)
            
            return AmassResult(
                domain=domain,
                subdomain=subdomain,
                ip_address=ip_address,
                asn=asn,
                as_name=as_name,
                country=country,
                city=city,
                service=service,
                port=port,
                protocol=protocol,
                source="amass"
            )
            
        except Exception as e:
            print(f"Error parsing Amass entry: {e}")
            return None
    
    def _get_service_name(self, port: int) -> str:
        """Get service name for a port"""
        services = {
            21: "ftp", 22: "ssh", 23: "telnet", 25: "smtp", 53: "dns",
            80: "http", 110: "pop3", 143: "imap", 443: "https", 993: "imaps",
            995: "pop3s", 3389: "rdp", 135: "rpc", 139: "netbios", 445: "smb",
            1433: "mssql", 1521: "oracle", 3306: "mysql", 5432: "postgresql",
            6379: "redis", 8080: "http-proxy", 8443: "https-alt",
        }
        
        return services.get(port, "unknown")
    
    async def discover_multiple_domains(self, domains: List[str], company_id: str) -> List[AmassResult]:
        """Discover assets for multiple domains concurrently"""
        tasks = []
        for domain in domains:
            task = self.discover_domain(domain, company_id)
            tasks.append(task)
        
        results = await asyncio.gather(*tasks, return_exceptions=True)
        
        # Flatten results and filter out exceptions
        all_results = []
        for result in results:
            if isinstance(result, list):
                all_results.extend(result)
            elif isinstance(result, Exception):
                print(f"Error in domain discovery: {result}")
        
        return all_results
    
    async def check_amass_installed(self) -> bool:
        """Check if Amass is installed and accessible"""
        try:
            process = await asyncio.create_subprocess_exec(
                self.amass_path, "version",
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            
            stdout, stderr = await process.communicate()
            return process.returncode == 0
            
        except Exception:
            return False
    
    async def get_amass_version(self) -> Optional[str]:
        """Get Amass version"""
        try:
            process = await asyncio.create_subprocess_exec(
                self.amass_path, "version",
                stdout=asyncio.subprocess.PIPE,
                stderr=asyncio.subprocess.PIPE
            )
            
            stdout, stderr = await process.communicate()
            if process.returncode == 0:
                return stdout.decode().strip()
            
        except Exception:
            pass
        
        return None

# Example usage
async def main():
    """Example usage of AmassDiscovery"""
    amass = AmassDiscovery()
    
    # Check if Amass is installed
    if not await amass.check_amass_installed():
        print("Amass is not installed or not accessible")
        return
    
    version = await amass.get_amass_version()
    print(f"Amass version: {version}")
    
    # Discover assets for a domain
    domain = "example.com"
    results = await amass.discover_domain(domain, "company-123")
    
    print(f"Discovered {len(results)} assets for {domain}")
    for result in results:
        print(f"  {result.subdomain} -> {result.ip_address} (ASN: {result.asn})")

if __name__ == "__main__":
    asyncio.run(main())

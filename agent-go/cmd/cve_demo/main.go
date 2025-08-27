package main

import (
	"fmt"
	"log"

	"zerotrace/agent/internal/scanner"
)

func main() {
	fmt.Println("🔍 ZeroTrace Real CVE Data Demo")
	fmt.Println("================================")

	// Initialize CVE sources
	nvdSource := scanner.NewNVDSource("") // No API key for demo (rate limited)
	githubSource := scanner.NewGitHubAdvisorySource()

	fmt.Println("\n📡 Fetching real CVE data from NIST NVD...")

	// Get a specific CVE (Log4Shell as example)
	cveID := "CVE-2021-44228"
	vuln, err := nvdSource.GetCVE(cveID)
	if err != nil {
		log.Printf("Error fetching CVE %s: %v", cveID, err)
		fmt.Println("⚠️  Note: NVD API has rate limits. In production, use an API key.")
	} else {
		fmt.Printf("✅ Found CVE: %s\n", vuln.CVEID)
		fmt.Printf("   Title: %s\n", vuln.Title)
		fmt.Printf("   Severity: %s\n", vuln.Severity)
		if vuln.CVSSScore != nil {
			fmt.Printf("   CVSS Score: %.1f\n", *vuln.CVSSScore)
		}
		fmt.Printf("   Description: %s\n", truncateString(vuln.Description, 100))
	}

	fmt.Println("\n📡 Fetching recent CVEs from GitHub Security Advisories...")

	// Get recent CVEs from GitHub
	recentCVEs, err := githubSource.GetRecentCVEs(5)
	if err != nil {
		log.Printf("Error fetching recent CVEs: %v", err)
	} else {
		fmt.Printf("✅ Found %d recent CVEs:\n", len(recentCVEs))
		for i, vuln := range recentCVEs {
			fmt.Printf("   %d. %s (%s) - %s\n",
				i+1, vuln.CVEID, vuln.Severity, truncateString(vuln.Title, 50))
		}
	}

	fmt.Println("\n🔧 How to integrate real CVE data:")
	fmt.Println("1. Get NVD API key from: https://nvd.nist.gov/developers/request-an-api-key")
	fmt.Println("2. Set environment variable: NVD_API_KEY=your_key_here")
	fmt.Println("3. Agent will automatically fetch real CVE data")
	fmt.Println("4. Risk scores calculated from actual CVSS scores")
	fmt.Println("5. Network topology shows real vulnerability paths")

	fmt.Println("\n📊 Real CVE Sources Available:")
	fmt.Println("• NIST National Vulnerability Database (NVD)")
	fmt.Println("• GitHub Security Advisories")
	fmt.Println("• MITRE CVE Database")
	fmt.Println("• Vendor Security Advisories")
	fmt.Println("• Custom vulnerability feeds")

	fmt.Println("\n🚀 Next Steps:")
	fmt.Println("1. Configure API keys in environment")
	fmt.Println("2. Run agent with real CVE scanning")
	fmt.Println("3. View real vulnerability data in UI")
	fmt.Println("4. Monitor network topology with actual risk scores")
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

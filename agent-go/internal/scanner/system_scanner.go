package scanner

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"zerotrace/agent/internal/config"
)

// SystemScanner collects detailed system information
type SystemScanner struct {
	config *config.Config
}

// NewSystemScanner creates a new system scanner
func NewSystemScanner(cfg *config.Config) *SystemScanner {
	return &SystemScanner{
		config: cfg,
	}
}

// SystemInfo represents comprehensive system information
type SystemInfo struct {
	// Core Asset Information
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Hostname       string    `json:"hostname"`
	IPAddress      string    `json:"ip_address"`
	MACAddress     string    `json:"mac_address"`
	Tags           []string  `json:"tags"`
	RiskScore      float64   `json:"risk_score"`
	LastSeen       time.Time `json:"last_seen"`
	CreatedAt      time.Time `json:"created_at"`

	// Operating System Details
	OSName        string `json:"os_name"`
	OSVersion     string `json:"os_version"`
	OSBuild       string `json:"os_build"`
	KernelVersion string `json:"kernel_version"`

	// Hardware Specifications
	CPUModel       string  `json:"cpu_model"`
	CPUCores       int     `json:"cpu_cores"`
	MemoryTotalGB  float64 `json:"memory_total_gb"`
	StorageTotalGB float64 `json:"storage_total_gb"`
	GPUModel       string  `json:"gpu_model"`
	SerialNumber   string  `json:"serial_number"`
	Platform       string  `json:"platform"`

	// Location Information
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
}

// Scan collects comprehensive system information
func (ss *SystemScanner) Scan() (*SystemInfo, error) {
	log.Printf("[SystemScanner] Starting comprehensive system scan...")

	systemInfo := &SystemInfo{
		ID:             ss.config.AgentID,
		OrganizationID: ss.config.OrganizationID,
		Hostname:       ss.config.Hostname,
		LastSeen:       time.Now(),
		CreatedAt:      time.Now(),
		Tags:           []string{"production", "managed"},
		RiskScore:      0.0,
	}

	// Collect OS information
	if err := ss.collectOSInfo(systemInfo); err != nil {
		log.Printf("[SystemScanner] Error collecting OS info: %v", err)
	}

	// Collect hardware information
	if err := ss.collectHardwareInfo(systemInfo); err != nil {
		log.Printf("[SystemScanner] Error collecting hardware info: %v", err)
	}

	// Collect network information
	if err := ss.collectNetworkInfo(systemInfo); err != nil {
		log.Printf("[SystemScanner] Error collecting network info: %v", err)
	}

	// Collect location information
	if err := ss.collectLocationInfo(systemInfo); err != nil {
		log.Printf("[SystemScanner] Error collecting location info: %v", err)
	}

	log.Printf("[SystemScanner] System scan completed successfully")
	return systemInfo, nil
}

// collectOSInfo gathers operating system details
func (ss *SystemScanner) collectOSInfo(info *SystemInfo) error {
	switch runtime.GOOS {
	case "darwin":
		return ss.collectMacOSInfo(info)
	case "linux":
		return ss.collectLinuxInfo(info)
	case "windows":
		return ss.collectWindowsInfo(info)
	default:
		info.OSName = runtime.GOOS
		info.OSVersion = "unknown"
		return nil
	}
}

// collectMacOSInfo gathers macOS-specific information
func (ss *SystemScanner) collectMacOSInfo(info *SystemInfo) error {
	info.OSName = "macOS"
	info.Platform = "arm64" // or "x86_64"

	// Get macOS version
	if output, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
		info.OSVersion = strings.TrimSpace(string(output))
	}

	// Get build number
	if output, err := exec.Command("sw_vers", "-buildVersion").Output(); err == nil {
		info.OSBuild = strings.TrimSpace(string(output))
	}

	// Get kernel version
	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		info.KernelVersion = strings.TrimSpace(string(output))
	}

	// Get serial number
	if output, err := exec.Command("system_profiler", "SPHardwareDataType", "-json").Output(); err == nil {
		var hardwareData map[string]interface{}
		if err := json.Unmarshal(output, &hardwareData); err == nil {
			if hardware := hardwareData["SPHardwareDataType"].([]interface{}); len(hardware) > 0 {
				if hw := hardware[0].(map[string]interface{}); hw["serial_number"] != nil {
					info.SerialNumber = hw["serial_number"].(string)
				}
			}
		}
	}

	return nil
}

// collectLinuxInfo gathers Linux-specific information
func (ss *SystemScanner) collectLinuxInfo(info *SystemInfo) error {
	info.OSName = "Linux"
	info.Platform = runtime.GOARCH

	// Get OS version from /etc/os-release
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				info.OSVersion = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
		}
	}

	// Get kernel version
	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		info.KernelVersion = strings.TrimSpace(string(output))
	}

	// Get serial number from DMI
	if output, err := exec.Command("dmidecode", "-s", "system-serial-number").Output(); err == nil {
		info.SerialNumber = strings.TrimSpace(string(output))
	}

	return nil
}

// collectWindowsInfo gathers Windows-specific information
func (ss *SystemScanner) collectWindowsInfo(info *SystemInfo) error {
	info.OSName = "Windows"
	info.Platform = runtime.GOARCH

	// Get Windows version
	if output, err := exec.Command("cmd", "/c", "ver").Output(); err == nil {
		info.OSVersion = strings.TrimSpace(string(output))
	}

	// Get build number
	if output, err := exec.Command("cmd", "/c", "wmic os get buildnumber /value").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "BuildNumber=") {
				info.OSBuild = strings.TrimSpace(strings.TrimPrefix(line, "BuildNumber="))
				break
			}
		}
	}

	return nil
}

// collectHardwareInfo gathers hardware specifications
func (ss *SystemScanner) collectHardwareInfo(info *SystemInfo) error {
	switch runtime.GOOS {
	case "darwin":
		return ss.collectMacOSHardware(info)
	case "linux":
		return ss.collectLinuxHardware(info)
	case "windows":
		return ss.collectWindowsHardware(info)
	default:
		return nil
	}
}

// collectMacOSHardware gathers macOS hardware information
func (ss *SystemScanner) collectMacOSHardware(info *SystemInfo) error {
	// Get CPU information
	if output, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
		info.CPUModel = strings.TrimSpace(string(output))
	}

	// Get CPU cores
	if output, err := exec.Command("sysctl", "-n", "hw.ncpu").Output(); err == nil {
		if cores, err := strconv.Atoi(strings.TrimSpace(string(output))); err == nil {
			info.CPUCores = cores
		}
	}

	// Get memory information
	if output, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
		if bytes, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64); err == nil {
			info.MemoryTotalGB = float64(bytes) / (1024 * 1024 * 1024)
		}
	}

	// Get storage information
	if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) > 1 {
				// Parse storage size (e.g., "500Gi" -> 500)
				sizeStr := fields[1]
				if strings.HasSuffix(sizeStr, "Gi") {
					if size, err := strconv.ParseFloat(strings.TrimSuffix(sizeStr, "Gi"), 64); err == nil {
						info.StorageTotalGB = size
					}
				}
			}
		}
	}

	// Get GPU information
	if output, err := exec.Command("system_profiler", "SPDisplaysDataType", "-json").Output(); err == nil {
		var displayData map[string]interface{}
		if err := json.Unmarshal(output, &displayData); err == nil {
			if displays := displayData["SPDisplaysDataType"].([]interface{}); len(displays) > 0 {
				if display := displays[0].(map[string]interface{}); display["_name"] != nil {
					info.GPUModel = display["_name"].(string)
				}
			}
		}
	}

	return nil
}

// collectLinuxHardware gathers Linux hardware information
func (ss *SystemScanner) collectLinuxHardware(info *SystemInfo) error {
	// Get CPU information
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				info.CPUModel = strings.TrimSpace(strings.Split(line, ":")[1])
				break
			}
		}
	}

	// Get CPU cores
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		coreCount := 0
		for _, line := range lines {
			if strings.HasPrefix(line, "processor") {
				coreCount++
			}
		}
		info.CPUCores = coreCount
	}

	// Get memory information
	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) > 1 {
					if kb, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
						info.MemoryTotalGB = float64(kb) / (1024 * 1024)
					}
				}
				break
			}
		}
	}

	// Get storage information
	if output, err := exec.Command("df", "-h", "/").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) > 1 {
				sizeStr := fields[1]
				if strings.HasSuffix(sizeStr, "G") {
					if size, err := strconv.ParseFloat(strings.TrimSuffix(sizeStr, "G"), 64); err == nil {
						info.StorageTotalGB = size
					}
				}
			}
		}
	}

	return nil
}

// collectWindowsHardware gathers Windows hardware information
func (ss *SystemScanner) collectWindowsHardware(info *SystemInfo) error {
	// Get CPU information
	if output, err := exec.Command("wmic", "cpu", "get", "name", "/value").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Name=") {
				info.CPUModel = strings.TrimSpace(strings.TrimPrefix(line, "Name="))
				break
			}
		}
	}

	// Get CPU cores
	if output, err := exec.Command("wmic", "cpu", "get", "NumberOfCores", "/value").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "NumberOfCores=") {
				if cores, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "NumberOfCores="))); err == nil {
					info.CPUCores = cores
				}
				break
			}
		}
	}

	// Get memory information
	if output, err := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/value").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "TotalPhysicalMemory=") {
				if bytes, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, "TotalPhysicalMemory=")), 10, 64); err == nil {
					info.MemoryTotalGB = float64(bytes) / (1024 * 1024 * 1024)
				}
				break
			}
		}
	}

	return nil
}

// collectNetworkInfo gathers network information
func (ss *SystemScanner) collectNetworkInfo(info *SystemInfo) error {
	// Get IP address
	if output, err := exec.Command("hostname", "-I").Output(); err == nil {
		info.IPAddress = strings.TrimSpace(string(output))
	} else if output, err := exec.Command("ip", "route", "get", "1").Output(); err == nil {
		// Alternative method for Linux
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "src") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "src" && i+1 < len(fields) {
						info.IPAddress = fields[i+1]
						break
					}
				}
				break
			}
		}
	}

	// Get MAC address
	if output, err := exec.Command("ifconfig").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "ether") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "ether" && i+1 < len(fields) {
						info.MACAddress = fields[i+1]
						break
					}
				}
				break
			}
		}
	}

	return nil
}

// collectLocationInfo gathers location information (best effort)
func (ss *SystemScanner) collectLocationInfo(info *SystemInfo) error {
	// Get timezone
	if output, err := exec.Command("date", "+%Z").Output(); err == nil {
		info.Timezone = strings.TrimSpace(string(output))
	}

	// Try to get location from system settings (macOS)
	if runtime.GOOS == "darwin" {
		// Try to get location from system preferences
		if output, err := exec.Command("defaults", "read", "/Library/Preferences/com.apple.timezone.auto.plist", "Active").Output(); err == nil {
			timezone := strings.TrimSpace(string(output))
			if timezone != "" {
				// Parse timezone to get location info
				info.City = "Unknown"
				info.Region = "Unknown"
				info.Country = "Unknown"

				// Try to get more specific location info
				if output, err := exec.Command("system_profiler", "SPLocationDataType").Output(); err == nil {
					// Parse location data if available
					lines := strings.Split(string(output), "\n")
					for _, line := range lines {
						if strings.Contains(line, "City") {
							parts := strings.Split(line, ":")
							if len(parts) > 1 {
								info.City = strings.TrimSpace(parts[1])
							}
						}
						if strings.Contains(line, "Region") {
							parts := strings.Split(line, ":")
							if len(parts) > 1 {
								info.Region = strings.TrimSpace(parts[1])
							}
						}
						if strings.Contains(line, "Country") {
							parts := strings.Split(line, ":")
							if len(parts) > 1 {
								info.Country = strings.TrimSpace(parts[1])
							}
						}
					}
				}
			}
		}
	}

	// Fallback to timezone-based location detection
	if info.City == "" || info.City == "Unknown" {
		// Use timezone to infer location
		timezone := info.Timezone
		switch timezone {
		case "IST":
			info.City = "Mumbai"
			info.Region = "Maharashtra"
			info.Country = "India"
		case "PST", "PDT":
			info.City = "San Francisco"
			info.Region = "California"
			info.Country = "United States"
		case "EST", "EDT":
			info.City = "New York"
			info.Region = "New York"
			info.Country = "United States"
		case "GMT", "UTC":
			info.City = "London"
			info.Region = "England"
			info.Country = "United Kingdom"
		default:
			info.City = "Unknown"
			info.Region = "Unknown"
			info.Country = "Unknown"
		}
	}

	return nil
}

package scanner

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// IoTOTScanner handles IoT and OT security scanning
type IoTOTScanner struct {
	config *config.Config
}

// IoTOTFinding represents an IoT/OT security finding
type IoTOTFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // device, protocol, firmware, network, access
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	DeviceID      string                 `json:"device_id,omitempty"`
	DeviceType    string                 `json:"device_type,omitempty"`
	Protocol      string                 `json:"protocol,omitempty"`
	IPAddress     string                 `json:"ip_address,omitempty"`
	MACAddress    string                 `json:"mac_address,omitempty"`
	Port          int                    `json:"port,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// DeviceInfo represents IoT/OT device information
type DeviceInfo struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	Type              string                 `json:"type"` // sensor, actuator, gateway, controller
	Manufacturer      string                 `json:"manufacturer"`
	Model             string                 `json:"model"`
	FirmwareVersion   string                 `json:"firmware_version"`
	HardwareVersion   string                 `json:"hardware_version"`
	IPAddress         string                 `json:"ip_address"`
	MACAddress        string                 `json:"mac_address"`
	Ports             []int                  `json:"ports"`
	Protocols         []string               `json:"protocols"`
	Services          []string               `json:"services"`
	IsEncrypted       bool                   `json:"is_encrypted"`
	HasAuthentication bool                   `json:"has_authentication"`
	IsShadowDevice    bool                   `json:"is_shadow_device"`
	RiskScore         float64                `json:"risk_score"`
	LastSeen          time.Time              `json:"last_seen"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// ProtocolInfo represents IoT/OT protocol information
type ProtocolInfo struct {
	Name              string   `json:"name"`
	Type              string   `json:"type"` // wireless, wired, industrial
	Port              int      `json:"port"`
	IsEncrypted       bool     `json:"is_encrypted"`
	HasAuthentication bool     `json:"has_authentication"`
	Vulnerabilities   []string `json:"vulnerabilities"`
	SecurityLevel     string   `json:"security_level"` // high, medium, low
	Description       string   `json:"description"`
}

// FirmwareInfo represents firmware information
type FirmwareInfo struct {
	DeviceID        string                 `json:"device_id"`
	Version         string                 `json:"version"`
	BuildDate       time.Time              `json:"build_date"`
	Vendor          string                 `json:"vendor"`
	Size            int64                  `json:"size"`
	Hash            string                 `json:"hash"`
	Vulnerabilities []string               `json:"vulnerabilities"`
	IsLatest        bool                   `json:"is_latest"`
	UpdateAvailable bool                   `json:"update_available"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NewIoTOTScanner creates a new IoT/OT security scanner
func NewIoTOTScanner(cfg *config.Config) *IoTOTScanner {
	return &IoTOTScanner{
		config: cfg,
	}
}

// Scan performs comprehensive IoT/OT security scanning
func (ios *IoTOTScanner) Scan() ([]IoTOTFinding, []DeviceInfo, []ProtocolInfo, []FirmwareInfo, error) {
	var findings []IoTOTFinding
	var devices []DeviceInfo
	var protocols []ProtocolInfo
	var firmwares []FirmwareInfo

	// Discover IoT/OT devices
	discoveredDevices := ios.discoverDevices()
	devices = append(devices, discoveredDevices...)

	// Scan each device
	for _, device := range discoveredDevices {
		deviceFindings := ios.scanDevice(device)
		findings = append(findings, deviceFindings...)
	}

	// Discover protocols
	discoveredProtocols := ios.discoverProtocols()
	protocols = append(protocols, discoveredProtocols...)

	// Scan protocols
	for _, protocol := range discoveredProtocols {
		protocolFindings := ios.scanProtocol(protocol)
		findings = append(findings, protocolFindings...)
	}

	// Discover firmware
	discoveredFirmwares := ios.discoverFirmware(devices)
	firmwares = append(firmwares, discoveredFirmwares...)

	// Scan firmware
	for _, firmware := range discoveredFirmwares {
		firmwareFindings := ios.scanFirmware(firmware)
		findings = append(findings, firmwareFindings...)
	}

	return findings, devices, protocols, firmwares, nil
}

// discoverDevices discovers IoT/OT devices
func (ios *IoTOTScanner) discoverDevices() []DeviceInfo {
	var devices []DeviceInfo

	// Network device discovery
	networkDevices := ios.discoverNetworkDevices()
	devices = append(devices, networkDevices...)

	// Bluetooth device discovery
	bluetoothDevices := ios.discoverBluetoothDevices()
	devices = append(devices, bluetoothDevices...)

	// USB device discovery
	usbDevices := ios.discoverUSBDevices()
	devices = append(devices, usbDevices...)

	// Serial device discovery
	serialDevices := ios.discoverSerialDevices()
	devices = append(devices, serialDevices...)

	return devices
}

// discoverNetworkDevices discovers network-connected IoT/OT devices
func (ios *IoTOTScanner) discoverNetworkDevices() []DeviceInfo {
	var devices []DeviceInfo

	// Scan common IoT ports
	iotPorts := []int{80, 443, 1883, 8883, 5683, 5684, 8080, 8443, 9090, 9091}

	// Get local network range
	networkRange := ios.getLocalNetworkRange()
	if networkRange == "" {
		return devices
	}

	// Scan each IP in the range
	for ip := range ios.generateIPRange(networkRange) {
		for _, port := range iotPorts {
			if ios.isPortOpen(ip, port) {
				device := ios.identifyDevice(ip, port)
				if device != nil {
					devices = append(devices, *device)
				}
			}
		}
	}

	return devices
}

// discoverBluetoothDevices discovers Bluetooth IoT devices
func (ios *IoTOTScanner) discoverBluetoothDevices() []DeviceInfo {
	var devices []DeviceInfo

	// Check if Bluetooth is available
	if !ios.isCommandAvailable("bluetoothctl") {
		return devices
	}

	// Scan for Bluetooth devices
	cmd := exec.Command("bluetoothctl", "scan", "on")
	output, err := cmd.Output()
	if err != nil {
		return devices
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Device") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				device := DeviceInfo{
					ID:                parts[1],
					Name:              "Unknown Bluetooth Device",
					Type:              "bluetooth",
					Manufacturer:      "Unknown",
					Protocols:         []string{"Bluetooth"},
					IsEncrypted:       false,
					HasAuthentication: false,
					RiskScore:         0.5,
					LastSeen:          time.Now(),
				}
				devices = append(devices, device)
			}
		}
	}

	return devices
}

// discoverUSBDevices discovers USB IoT devices
func (ios *IoTOTScanner) discoverUSBDevices() []DeviceInfo {
	var devices []DeviceInfo

	// Check if lsusb is available
	if !ios.isCommandAvailable("lsusb") {
		return devices
	}

	// List USB devices
	cmd := exec.Command("lsusb")
	output, err := cmd.Output()
	if err != nil {
		return devices
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" {
			device := DeviceInfo{
				ID:                ios.extractUSBDeviceID(line),
				Name:              ios.extractUSBDeviceName(line),
				Type:              "usb",
				Manufacturer:      ios.extractUSBManufacturer(line),
				Protocols:         []string{"USB"},
				IsEncrypted:       false,
				HasAuthentication: false,
				RiskScore:         0.3,
				LastSeen:          time.Now(),
			}
			devices = append(devices, device)
		}
	}

	return devices
}

// discoverSerialDevices discovers serial IoT devices
func (ios *IoTOTScanner) discoverSerialDevices() []DeviceInfo {
	var devices []DeviceInfo

	// Check for serial devices
	cmd := exec.Command("ls", "/dev/tty*")
	output, err := cmd.Output()
	if err != nil {
		return devices
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "tty") && !strings.Contains(line, "ttyS") {
			device := DeviceInfo{
				ID:                line,
				Name:              "Serial Device",
				Type:              "serial",
				Manufacturer:      "Unknown",
				Protocols:         []string{"Serial"},
				IsEncrypted:       false,
				HasAuthentication: false,
				RiskScore:         0.4,
				LastSeen:          time.Now(),
			}
			devices = append(devices, device)
		}
	}

	return devices
}

// getLocalNetworkRange gets the local network range
func (ios *IoTOTScanner) getLocalNetworkRange() string {
	// This would typically involve getting the local network interface
	// For now, return a placeholder
	return "192.168.1.0/24"
}

// generateIPRange generates IP addresses in a range
func (ios *IoTOTScanner) generateIPRange(networkRange string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		// Simplified implementation
		// In reality, this would parse the CIDR and generate all IPs
		ch <- "192.168.1.1"
		ch <- "192.168.1.2"
		ch <- "192.168.1.3"
	}()
	return ch
}

// isPortOpen checks if a port is open
func (ios *IoTOTScanner) isPortOpen(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// identifyDevice identifies a device by IP and port
func (ios *IoTOTScanner) identifyDevice(ip string, port int) *DeviceInfo {
	// This would involve banner grabbing and service identification
	// For now, return a placeholder device
	return &DeviceInfo{
		ID:                fmt.Sprintf("%s:%d", ip, port),
		Name:              "Unknown IoT Device",
		Type:              "sensor",
		Manufacturer:      "Unknown",
		IPAddress:         ip,
		Ports:             []int{port},
		Protocols:         []string{"HTTP"},
		IsEncrypted:       false,
		HasAuthentication: false,
		RiskScore:         0.5,
		LastSeen:          time.Now(),
	}
}

// extractUSBDeviceID extracts USB device ID from lsusb output
func (ios *IoTOTScanner) extractUSBDeviceID(line string) string {
	parts := strings.Fields(line)
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

// extractUSBDeviceName extracts USB device name from lsusb output
func (ios *IoTOTScanner) extractUSBDeviceName(line string) string {
	// This would parse the lsusb output to extract device name
	return "USB Device"
}

// extractUSBManufacturer extracts USB manufacturer from lsusb output
func (ios *IoTOTScanner) extractUSBManufacturer(line string) string {
	// This would parse the lsusb output to extract manufacturer
	return "Unknown"
}

// discoverProtocols discovers IoT/OT protocols
func (ios *IoTOTScanner) discoverProtocols() []ProtocolInfo {
	var protocols []ProtocolInfo

	// Common IoT/OT protocols
	commonProtocols := []ProtocolInfo{
		{
			Name:              "MQTT",
			Type:              "wireless",
			Port:              1883,
			IsEncrypted:       false,
			HasAuthentication: false,
			Vulnerabilities:   []string{"unencrypted", "no_auth"},
			SecurityLevel:     "low",
			Description:       "Message Queuing Telemetry Transport",
		},
		{
			Name:              "MQTTS",
			Type:              "wireless",
			Port:              8883,
			IsEncrypted:       true,
			HasAuthentication: true,
			Vulnerabilities:   []string{},
			SecurityLevel:     "high",
			Description:       "MQTT over SSL/TLS",
		},
		{
			Name:              "CoAP",
			Type:              "wireless",
			Port:              5683,
			IsEncrypted:       false,
			HasAuthentication: false,
			Vulnerabilities:   []string{"unencrypted", "no_auth"},
			SecurityLevel:     "low",
			Description:       "Constrained Application Protocol",
		},
		{
			Name:              "CoAPS",
			Type:              "wireless",
			Port:              5684,
			IsEncrypted:       true,
			HasAuthentication: true,
			Vulnerabilities:   []string{},
			SecurityLevel:     "high",
			Description:       "CoAP over SSL/TLS",
		},
		{
			Name:              "Modbus",
			Type:              "industrial",
			Port:              502,
			IsEncrypted:       false,
			HasAuthentication: false,
			Vulnerabilities:   []string{"unencrypted", "no_auth"},
			SecurityLevel:     "low",
			Description:       "Modbus Protocol",
		},
		{
			Name:              "DNP3",
			Type:              "industrial",
			Port:              20000,
			IsEncrypted:       false,
			HasAuthentication: false,
			Vulnerabilities:   []string{"unencrypted", "no_auth"},
			SecurityLevel:     "low",
			Description:       "Distributed Network Protocol",
		},
	}

	protocols = append(protocols, commonProtocols...)

	return protocols
}

// discoverFirmware discovers firmware on devices
func (ios *IoTOTScanner) discoverFirmware(devices []DeviceInfo) []FirmwareInfo {
	var firmwares []FirmwareInfo

	for _, device := range devices {
		firmware := FirmwareInfo{
			DeviceID:        device.ID,
			Version:         device.FirmwareVersion,
			Vendor:          device.Manufacturer,
			Size:            0,  // Would need to extract from device
			Hash:            "", // Would need to calculate
			Vulnerabilities: []string{},
			IsLatest:        false,
			UpdateAvailable: false,
			Metadata:        make(map[string]interface{}),
		}
		firmwares = append(firmwares, firmware)
	}

	return firmwares
}

// scanDevice scans a specific device for security issues
func (ios *IoTOTScanner) scanDevice(device DeviceInfo) []IoTOTFinding {
	var findings []IoTOTFinding

	// Check for unencrypted communication
	if !device.IsEncrypted {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "device",
			Severity:      "high",
			Title:         "Unencrypted IoT Device Communication",
			Description:   fmt.Sprintf("Device %s is not using encryption", device.Name),
			DeviceID:      device.ID,
			DeviceType:    device.Type,
			IPAddress:     device.IPAddress,
			CurrentValue:  "unencrypted",
			RequiredValue: "encrypted",
			Remediation:   "Enable encryption for device communication",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"device_id":   device.ID,
				"device_type": device.Type,
				"encrypted":   false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for missing authentication
	if !device.HasAuthentication {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "device",
			Severity:      "high",
			Title:         "IoT Device Without Authentication",
			Description:   fmt.Sprintf("Device %s has no authentication", device.Name),
			DeviceID:      device.ID,
			DeviceType:    device.Type,
			IPAddress:     device.IPAddress,
			CurrentValue:  "no_auth",
			RequiredValue: "authenticated",
			Remediation:   "Implement authentication for device access",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"device_id":     device.ID,
				"device_type":   device.Type,
				"authenticated": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for shadow devices
	if device.IsShadowDevice {
		finding := IoTOTFinding{
			ID:           uuid.New().String(),
			Type:         "device",
			Severity:     "medium",
			Title:        "Shadow IoT Device Detected",
			Description:  fmt.Sprintf("Device %s is a shadow IoT device", device.Name),
			DeviceID:     device.ID,
			DeviceType:   device.Type,
			IPAddress:    device.IPAddress,
			Remediation:  "Review and secure shadow IoT devices",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"device_id":     device.ID,
				"device_type":   device.Type,
				"shadow_device": true,
			},
		}
		findings = append(findings, finding)
	}

	// Check for high risk score
	if device.RiskScore > 0.7 {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "device",
			Severity:      "medium",
			Title:         "High Risk IoT Device",
			Description:   fmt.Sprintf("Device %s has high risk score: %.2f", device.Name, device.RiskScore),
			DeviceID:      device.ID,
			DeviceType:    device.Type,
			IPAddress:     device.IPAddress,
			CurrentValue:  fmt.Sprintf("%.2f", device.RiskScore),
			RequiredValue: "0.7-",
			Remediation:   "Review device security configuration",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"device_id":   device.ID,
				"device_type": device.Type,
				"risk_score":  device.RiskScore,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanProtocol scans a specific protocol for security issues
func (ios *IoTOTScanner) scanProtocol(protocol ProtocolInfo) []IoTOTFinding {
	var findings []IoTOTFinding

	// Check for unencrypted protocols
	if !protocol.IsEncrypted {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "protocol",
			Severity:      "high",
			Title:         "Unencrypted IoT Protocol",
			Description:   fmt.Sprintf("Protocol %s is not encrypted", protocol.Name),
			Protocol:      protocol.Name,
			Port:          protocol.Port,
			CurrentValue:  "unencrypted",
			RequiredValue: "encrypted",
			Remediation:   "Use encrypted version of the protocol",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"protocol":  protocol.Name,
				"port":      protocol.Port,
				"encrypted": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for missing authentication
	if !protocol.HasAuthentication {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "protocol",
			Severity:      "high",
			Title:         "IoT Protocol Without Authentication",
			Description:   fmt.Sprintf("Protocol %s has no authentication", protocol.Name),
			Protocol:      protocol.Name,
			Port:          protocol.Port,
			CurrentValue:  "no_auth",
			RequiredValue: "authenticated",
			Remediation:   "Implement authentication for the protocol",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"protocol":      protocol.Name,
				"port":          protocol.Port,
				"authenticated": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for known vulnerabilities
	if len(protocol.Vulnerabilities) > 0 {
		finding := IoTOTFinding{
			ID:           uuid.New().String(),
			Type:         "protocol",
			Severity:     "medium",
			Title:        "IoT Protocol Vulnerabilities",
			Description:  fmt.Sprintf("Protocol %s has known vulnerabilities: %s", protocol.Name, strings.Join(protocol.Vulnerabilities, ", ")),
			Protocol:     protocol.Name,
			Port:         protocol.Port,
			Remediation:  "Update protocol implementation and patch vulnerabilities",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"protocol":        protocol.Name,
				"port":            protocol.Port,
				"vulnerabilities": protocol.Vulnerabilities,
			},
		}
		findings = append(findings, finding)
	}

	// Check for low security level
	if protocol.SecurityLevel == "low" {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "protocol",
			Severity:      "medium",
			Title:         "Low Security IoT Protocol",
			Description:   fmt.Sprintf("Protocol %s has low security level", protocol.Name),
			Protocol:      protocol.Name,
			Port:          protocol.Port,
			CurrentValue:  protocol.SecurityLevel,
			RequiredValue: "high",
			Remediation:   "Upgrade to higher security protocol version",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"protocol":       protocol.Name,
				"port":           protocol.Port,
				"security_level": protocol.SecurityLevel,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanFirmware scans firmware for security issues
func (ios *IoTOTScanner) scanFirmware(firmware FirmwareInfo) []IoTOTFinding {
	var findings []IoTOTFinding

	// Check for vulnerabilities
	if len(firmware.Vulnerabilities) > 0 {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "firmware",
			Severity:      "high",
			Title:         "Firmware Vulnerabilities",
			Description:   fmt.Sprintf("Firmware %s has %d vulnerabilities", firmware.Version, len(firmware.Vulnerabilities)),
			DeviceID:      firmware.DeviceID,
			CurrentValue:  fmt.Sprintf("%d vulnerabilities", len(firmware.Vulnerabilities)),
			RequiredValue: "0 vulnerabilities",
			Remediation:   "Update firmware to latest version",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"device_id":       firmware.DeviceID,
				"version":         firmware.Version,
				"vulnerabilities": firmware.Vulnerabilities,
			},
		}
		findings = append(findings, finding)
	}

	// Check for outdated firmware
	if !firmware.IsLatest {
		finding := IoTOTFinding{
			ID:            uuid.New().String(),
			Type:          "firmware",
			Severity:      "medium",
			Title:         "Outdated Firmware",
			Description:   fmt.Sprintf("Firmware %s is not the latest version", firmware.Version),
			DeviceID:      firmware.DeviceID,
			CurrentValue:  firmware.Version,
			RequiredValue: "latest",
			Remediation:   "Update firmware to latest version",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"device_id": firmware.DeviceID,
				"version":   firmware.Version,
				"is_latest": firmware.IsLatest,
			},
		}
		findings = append(findings, finding)
	}

	// Check for update availability
	if firmware.UpdateAvailable {
		finding := IoTOTFinding{
			ID:           uuid.New().String(),
			Type:         "firmware",
			Severity:     "low",
			Title:        "Firmware Update Available",
			Description:  fmt.Sprintf("Firmware update is available for device %s", firmware.DeviceID),
			DeviceID:     firmware.DeviceID,
			Remediation:  "Install available firmware update",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"device_id":        firmware.DeviceID,
				"version":          firmware.Version,
				"update_available": firmware.UpdateAvailable,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// isCommandAvailable checks if a command is available
func (ios *IoTOTScanner) isCommandAvailable(command string) bool {
	cmd := exec.Command("which", command)
	err := cmd.Run()
	return err == nil
}

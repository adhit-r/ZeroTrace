package scanner

import (
	"fmt"
	"strings"
)

// DeviceClassifier identifies device types based on scan results
type DeviceClassifier struct{}

// NewDeviceClassifier creates a new device classifier
func NewDeviceClassifier() *DeviceClassifier {
	return &DeviceClassifier{}
}

// ClassifyDevice identifies the device type based on open ports, services, OS, and banners
func (dc *DeviceClassifier) ClassifyDevice(host string, ports []int, services map[int]string, osInfo string, banners map[int]string) string {
	// Check for switch indicators
	if dc.isSwitch(ports, services, banners) {
		return "switch"
	}

	// Check for router indicators
	if dc.isRouter(ports, services, banners) {
		return "router"
	}

	// Check for IoT device indicators
	if dc.isIoTDevice(ports, services, banners) {
		return "iot"
	}

	// Check for phone/mobile device indicators
	if dc.isPhone(ports, services, banners, osInfo) {
		return "phone"
	}

	// Check for server indicators
	if dc.isServer(ports, services, osInfo) {
		return "server"
	}

	return "unknown"
}

// isSwitch checks if the device is likely a network switch
func (dc *DeviceClassifier) isSwitch(ports []int, services map[int]string, banners map[int]string) bool {
	// SNMP port (161) is common for switch management
	hasSNMP := false
	for _, port := range ports {
		if port == 161 {
			hasSNMP = true
			break
		}
	}

	// Check for switch-specific services
	hasSwitchService := false
	for port, service := range services {
		serviceLower := strings.ToLower(service)
		if port == 161 || // SNMP
			strings.Contains(serviceLower, "snmp") ||
			strings.Contains(serviceLower, "switch") ||
			strings.Contains(serviceLower, "cisco") ||
			strings.Contains(serviceLower, "hp-procurve") ||
			strings.Contains(serviceLower, "juniper") {
			hasSwitchService = true
			break
		}
	}

	// Check banners for switch indicators
	hasSwitchBanner := false
	for _, banner := range banners {
		bannerLower := strings.ToLower(banner)
		if strings.Contains(bannerLower, "cisco ios") ||
			strings.Contains(bannerLower, "switch") ||
			strings.Contains(bannerLower, "procurve") ||
			strings.Contains(bannerLower, "juniper") {
			hasSwitchBanner = true
			break
		}
	}

	return hasSNMP && (hasSwitchService || hasSwitchBanner)
}

// isRouter checks if the device is likely a router
func (dc *DeviceClassifier) isRouter(ports []int, services map[int]string, banners map[int]string) bool {
	// Routing protocol ports
	routingPorts := []int{179, 520, 88} // BGP, RIP, OSPF
	hasRoutingPort := false
	for _, port := range ports {
		for _, rp := range routingPorts {
			if port == rp {
				hasRoutingPort = true
				break
			}
		}
		if hasRoutingPort {
			break
		}
	}

	// Check for router-specific services
	hasRouterService := false
	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "bgp") ||
			strings.Contains(serviceLower, "ospf") ||
			strings.Contains(serviceLower, "rip") ||
			strings.Contains(serviceLower, "router") ||
			strings.Contains(serviceLower, "cisco") ||
			strings.Contains(serviceLower, "juniper") {
			hasRouterService = true
			break
		}
	}

	// Check banners for router indicators
	hasRouterBanner := false
	for _, banner := range banners {
		bannerLower := strings.ToLower(banner)
		if strings.Contains(bannerLower, "cisco ios") ||
			strings.Contains(bannerLower, "router") ||
			strings.Contains(bannerLower, "juniper") ||
			strings.Contains(bannerLower, "mikrotik") {
			hasRouterBanner = true
			break
		}
	}

	return hasRoutingPort || hasRouterService || hasRouterBanner
}

// isIoTDevice checks if the device is likely an IoT device
func (dc *DeviceClassifier) isIoTDevice(ports []int, services map[int]string, banners map[int]string) bool {
	// IoT-specific ports
	iotPorts := []int{1883, 8883, 5683, 5684, 1900, 49152} // MQTT, CoAP, UPnP
	hasIoTPort := false
	for _, port := range ports {
		for _, iotPort := range iotPorts {
			if port == iotPort {
				hasIoTPort = true
				break
			}
		}
		if hasIoTPort {
			break
		}
	}

	// Check for IoT-specific services
	hasIoTService := false
	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "mqtt") ||
			strings.Contains(serviceLower, "coap") ||
			strings.Contains(serviceLower, "upnp") ||
			strings.Contains(serviceLower, "iot") ||
			strings.Contains(serviceLower, "smart") {
			hasIoTService = true
			break
		}
	}

	// Check banners for IoT indicators
	hasIoTBanner := false
	for _, banner := range banners {
		bannerLower := strings.ToLower(banner)
		if strings.Contains(bannerLower, "mqtt") ||
			strings.Contains(bannerLower, "coap") ||
			strings.Contains(bannerLower, "iot") ||
			strings.Contains(bannerLower, "smart") ||
			strings.Contains(bannerLower, "philips hue") ||
			strings.Contains(bannerLower, "nest") ||
			strings.Contains(bannerLower, "echo") {
			hasIoTBanner = true
			break
		}
	}

	return hasIoTPort || hasIoTService || hasIoTBanner
}

// isPhone checks if the device is likely a phone or mobile device
func (dc *DeviceClassifier) isPhone(ports []int, services map[int]string, banners map[int]string, osInfo string) bool {
	// Check OS info for mobile OS
	osLower := strings.ToLower(osInfo)
	if strings.Contains(osLower, "android") ||
		strings.Contains(osLower, "ios") ||
		strings.Contains(osLower, "iphone") ||
		strings.Contains(osLower, "mobile") {
		return true
	}

	// Check for phone-specific services
	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "airplay") ||
			strings.Contains(serviceLower, "airprint") ||
			strings.Contains(serviceLower, "dlna") {
			return true
		}
	}

	// Check banners for phone indicators
	for _, banner := range banners {
		bannerLower := strings.ToLower(banner)
		if strings.Contains(bannerLower, "iphone") ||
			strings.Contains(bannerLower, "android") ||
			strings.Contains(bannerLower, "airplay") {
			return true
		}
	}

	return false
}

// isServer checks if the device is likely a server
func (dc *DeviceClassifier) isServer(ports []int, services map[int]string, osInfo string) bool {
	// Common server ports
	serverPorts := []int{22, 80, 443, 3306, 5432, 1433, 3389, 5985, 5986}
	hasServerPort := false
	for _, port := range ports {
		for _, sp := range serverPorts {
			if port == sp {
				hasServerPort = true
				break
			}
		}
		if hasServerPort {
			break
		}
	}

	// Check for server-specific services
	hasServerService := false
	for _, service := range services {
		serviceLower := strings.ToLower(service)
		if strings.Contains(serviceLower, "http") ||
			strings.Contains(serviceLower, "https") ||
			strings.Contains(serviceLower, "ssh") ||
			strings.Contains(serviceLower, "mysql") ||
			strings.Contains(serviceLower, "postgres") ||
			strings.Contains(serviceLower, "mssql") ||
			strings.Contains(serviceLower, "rdp") {
			hasServerService = true
			break
		}
	}

	// Check OS info for server OS
	osLower := strings.ToLower(osInfo)
	hasServerOS := strings.Contains(osLower, "linux") ||
		strings.Contains(osLower, "windows") ||
		strings.Contains(osLower, "unix") ||
		strings.Contains(osLower, "server")

	return (hasServerPort || hasServerService) && hasServerOS
}

// GetDeviceConfidence returns a confidence score for device classification
func (dc *DeviceClassifier) GetDeviceConfidence(deviceType string, ports []int, services map[int]string, osInfo string) float64 {
	// Simple confidence scoring based on number of indicators
	indicators := 0.0
	totalChecks := 0.0

	switch deviceType {
	case "switch":
		for _, port := range ports {
			if port == 161 {
				indicators++
				break
			}
		}
		totalChecks++
		for _, service := range services {
			if strings.Contains(strings.ToLower(service), "snmp") ||
				strings.Contains(strings.ToLower(service), "switch") {
				indicators++
				break
			}
		}
		totalChecks++

	case "router":
		routingPorts := []int{179, 520, 88}
		for _, port := range ports {
			for _, rp := range routingPorts {
				if port == rp {
					indicators++
					break
				}
			}
		}
		totalChecks++

	case "iot":
		iotPorts := []int{1883, 8883, 5683, 5684}
		for _, port := range ports {
			for _, iotPort := range iotPorts {
				if port == iotPort {
					indicators++
					break
				}
			}
		}
		totalChecks++

	case "phone":
		osLower := strings.ToLower(osInfo)
		if strings.Contains(osLower, "android") || strings.Contains(osLower, "ios") {
			indicators++
		}
		totalChecks++

	case "server":
		serverPorts := []int{22, 80, 443, 3306, 5432}
		for _, port := range ports {
			for _, sp := range serverPorts {
				if port == sp {
					indicators++
					break
				}
			}
		}
		totalChecks++
	}

	if totalChecks == 0 {
		return 0.0
	}

	confidence := indicators / totalChecks
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// FormatDeviceInfo formats device information for display
func (dc *DeviceClassifier) FormatDeviceInfo(deviceType string, host string, osInfo string) string {
	if osInfo != "" {
		return fmt.Sprintf("%s (%s) - %s", deviceType, host, osInfo)
	}
	return fmt.Sprintf("%s (%s)", deviceType, host)
}


package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"zerotrace/api/internal/constants"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"
)

// ConfigParserService handles configuration file parsing
type ConfigParserService struct {
	configFileRepo *repository.ConfigFileRepository
}

// NewConfigParserService creates a new config parser service
func NewConfigParserService(configFileRepo *repository.ConfigFileRepository) *ConfigParserService {
	return &ConfigParserService{
		configFileRepo: configFileRepo,
	}
}

// ParseConfigFile parses a configuration file based on manufacturer and device type
func (s *ConfigParserService) ParseConfigFile(configFile *models.ConfigFile) error {
	// Update status to parsing
	err := s.configFileRepo.UpdateParsingStatus(configFile.ID, constants.StatusParsing, nil, "")
	if err != nil {
		return err
	}

	var parsedData map[string]interface{}
	var parseErr error

	// Parse based on manufacturer
	switch strings.ToLower(configFile.Manufacturer) {
	case "cisco":
		if strings.ToLower(configFile.DeviceType) == "firewall" {
			parsedData, parseErr = s.ParseCiscoASA(configFile.FileContent)
		} else {
			parsedData, parseErr = s.ParseCiscoIOS(configFile.FileContent)
		}
	case "palo alto", "paloalto", "palo alto networks":
		parsedData, parseErr = s.ParsePaloAlto(configFile.FileContent)
	case "fortinet", "fortigate":
		parsedData, parseErr = s.ParseFortinet(configFile.FileContent)
	case "juniper":
		parsedData, parseErr = s.ParseJuniper(configFile.FileContent)
	default:
		parseErr = fmt.Errorf("unsupported manufacturer: %s", configFile.Manufacturer)
	}

	if parseErr != nil {
		err = s.configFileRepo.UpdateParsingStatus(configFile.ID, constants.StatusFailed, nil, parseErr.Error())
		return parseErr
	}

	// Convert to JSONB format
	parsedJSON, err := json.Marshal(parsedData)
	if err != nil {
		err = s.configFileRepo.UpdateParsingStatus(configFile.ID, constants.StatusFailed, nil, err.Error())
		return err
	}

	// Update status to parsed
	err = s.configFileRepo.UpdateParsingStatus(configFile.ID, constants.StatusParsed, parsedJSON, "")
	if err != nil {
		return err
	}

	return nil
}

// ParseCiscoASA parses Cisco ASA configuration
func (s *ConfigParserService) ParseCiscoASA(config []byte) (map[string]interface{}, error) {
	content := string(config)
	lines := strings.Split(content, "\n")

	parsed := map[string]interface{}{
		"interfaces":     []map[string]interface{}{},
		"access_lists":   []map[string]interface{}{},
		"nat_rules":      []map[string]interface{}{},
		"vpn_configs":   []map[string]interface{}{},
		"user_accounts":  []map[string]interface{}{},
		"passwords":      []map[string]interface{}{},
		"logging":        map[string]interface{}{},
		"snmp":           map[string]interface{}{},
		"ssh":            map[string]interface{}{},
		"telnet":         map[string]interface{}{},
		"http_server":    map[string]interface{}{},
		"crypto":         map[string]interface{}{},
		"version":        "",
		"hostname":       "",
		"domain_name":    "",
	}

	var currentInterface map[string]interface{}
	var currentACL map[string]interface{}

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "!") {
			continue
		}

		// Parse version
		if strings.HasPrefix(line, "ASA Version") {
			parsed["version"] = strings.TrimSpace(strings.TrimPrefix(line, "ASA Version"))
		}

		// Parse hostname
		if strings.HasPrefix(line, "hostname ") {
			parsed["hostname"] = strings.TrimSpace(strings.TrimPrefix(line, "hostname "))
		}

		// Parse domain name
		if strings.HasPrefix(line, "domain-name ") {
			parsed["domain_name"] = strings.TrimSpace(strings.TrimPrefix(line, "domain-name "))
		}

		// Parse interfaces
		if strings.HasPrefix(line, "interface ") {
			if currentInterface != nil {
				interfaces := parsed["interfaces"].([]map[string]interface{})
				parsed["interfaces"] = append(interfaces, currentInterface)
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentInterface = map[string]interface{}{
					"name":        parts[1],
					"line_number": i + 1,
					"config":      []string{line},
				}
			}
		} else if currentInterface != nil {
			configLines := currentInterface["config"].([]string)
			currentInterface["config"] = append(configLines, line)
			if strings.HasPrefix(line, "!") || !strings.HasPrefix(line, " ") {
				interfaces := parsed["interfaces"].([]map[string]interface{})
				parsed["interfaces"] = append(interfaces, currentInterface)
				currentInterface = nil
			}
		}

		// Parse access-lists
		if strings.HasPrefix(line, "access-list ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentACL = map[string]interface{}{
					"name":        parts[1],
					"line_number": i + 1,
					"rule":        line,
				}
				acls := parsed["access_lists"].([]map[string]interface{})
				parsed["access_lists"] = append(acls, currentACL)
			}
		}

		// Parse NAT rules
		if strings.HasPrefix(line, "nat ") || strings.HasPrefix(line, "static ") {
			natRule := map[string]interface{}{
				"line_number": i + 1,
				"rule":        line,
			}
			nats := parsed["nat_rules"].([]map[string]interface{})
			parsed["nat_rules"] = append(nats, natRule)
		}

		// Parse user accounts
		if strings.HasPrefix(line, "username ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				user := map[string]interface{}{
					"username":    parts[1],
					"line_number": i + 1,
					"config":      line,
				}
				users := parsed["user_accounts"].([]map[string]interface{})
				parsed["user_accounts"] = append(users, user)
			}
		}

		// Parse passwords (encrypted)
		if strings.HasPrefix(line, "enable password ") || strings.HasPrefix(line, "password ") {
			password := map[string]interface{}{
				"line_number": i + 1,
				"config":      line,
				"encrypted":   true,
			}
			passwords := parsed["passwords"].([]map[string]interface{})
			parsed["passwords"] = append(passwords, password)
		}

		// Parse logging
		if strings.HasPrefix(line, "logging ") {
			logging := parsed["logging"].(map[string]interface{})
			logging["enabled"] = true
			logging["config"] = line
		}

		// Parse SNMP
		if strings.HasPrefix(line, "snmp-server ") {
			snmp := parsed["snmp"].(map[string]interface{})
			snmp["enabled"] = true
			if snmp["config"] == nil {
				snmp["config"] = []string{}
			}
			configs := snmp["config"].([]string)
			snmp["config"] = append(configs, line)
		}

		// Parse SSH
		if strings.HasPrefix(line, "ssh ") {
			ssh := parsed["ssh"].(map[string]interface{})
			ssh["enabled"] = true
			ssh["config"] = line
		}

		// Parse Telnet
		if strings.HasPrefix(line, "telnet ") {
			telnet := parsed["telnet"].(map[string]interface{})
			telnet["enabled"] = true
			telnet["config"] = line
		}

		// Parse HTTP server
		if strings.HasPrefix(line, "http server enable") {
			httpServer := parsed["http_server"].(map[string]interface{})
			httpServer["enabled"] = true
		}

		// Parse crypto settings
		if strings.HasPrefix(line, "crypto ") {
			crypto := parsed["crypto"].(map[string]interface{})
			if crypto["config"] == nil {
				crypto["config"] = []string{}
			}
			configs := crypto["config"].([]string)
			crypto["config"] = append(configs, line)
		}
	}

	// Add final interface if exists
	if currentInterface != nil {
		interfaces := parsed["interfaces"].([]map[string]interface{})
		parsed["interfaces"] = append(interfaces, currentInterface)
	}

	return parsed, nil
}

// ParseCiscoIOS parses Cisco IOS configuration
func (s *ConfigParserService) ParseCiscoIOS(config []byte) (map[string]interface{}, error) {
	content := string(config)
	lines := strings.Split(content, "\n")

	parsed := map[string]interface{}{
		"interfaces":    []map[string]interface{}{},
		"access_lists":  []map[string]interface{}{},
		"routing":       []map[string]interface{}{},
		"snmp":          map[string]interface{}{},
		"logging":       map[string]interface{}{},
		"authentication": map[string]interface{}{},
		"version":       "",
		"hostname":      "",
	}

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "!") {
			continue
		}

		// Parse version
		if strings.Contains(line, "version ") {
			parts := strings.Fields(line)
			for j, part := range parts {
				if part == "version" && j+1 < len(parts) {
					parsed["version"] = parts[j+1]
					break
				}
			}
		}

		// Parse hostname
		if strings.HasPrefix(line, "hostname ") {
			parsed["hostname"] = strings.TrimSpace(strings.TrimPrefix(line, "hostname "))
		}

		// Parse interfaces
		if strings.HasPrefix(line, "interface ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				iface := map[string]interface{}{
					"name":        strings.Join(parts[1:], " "),
					"line_number": i + 1,
				}
				interfaces := parsed["interfaces"].([]map[string]interface{})
				parsed["interfaces"] = append(interfaces, iface)
			}
		}

		// Parse access-lists
		if strings.HasPrefix(line, "access-list ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				acl := map[string]interface{}{
					"name":        parts[1],
					"line_number": i + 1,
					"rule":        line,
				}
				acls := parsed["access_lists"].([]map[string]interface{})
				parsed["access_lists"] = append(acls, acl)
			}
		}

		// Parse SNMP
		if strings.HasPrefix(line, "snmp-server ") {
			snmp := parsed["snmp"].(map[string]interface{})
			snmp["enabled"] = true
		}

		// Parse logging
		if strings.HasPrefix(line, "logging ") {
			logging := parsed["logging"].(map[string]interface{})
			logging["enabled"] = true
		}
	}

	return parsed, nil
}

// ParsePaloAlto parses Palo Alto Networks configuration
func (s *ConfigParserService) ParsePaloAlto(config []byte) (map[string]interface{}, error) {
	content := string(config)

	// Try to parse as XML first
	if strings.HasPrefix(strings.TrimSpace(content), "<") {
		return s.parsePaloAltoXML(content)
	}

	// Try to parse as JSON
	if strings.HasPrefix(strings.TrimSpace(content), "{") {
		var jsonData map[string]interface{}
		err := json.Unmarshal(config, &jsonData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Palo Alto JSON: %w", err)
		}
		return jsonData, nil
	}

	return nil, errors.New("unsupported Palo Alto configuration format")
}

// parsePaloAltoXML parses Palo Alto XML configuration
func (s *ConfigParserService) parsePaloAltoXML(content string) (map[string]interface{}, error) {
	parsed := map[string]interface{}{
		"security_policies": []map[string]interface{}{},
		"zones":              []map[string]interface{}{},
		"interfaces":         []map[string]interface{}{},
		"user_accounts":      []map[string]interface{}{},
		"format":             "xml",
	}

	// Basic XML parsing - extract key sections
	// In production, use proper XML parser

	// Extract security policies
	policyRegex := regexp.MustCompile(`<entry name="([^"]+)">`)
	policyMatches := policyRegex.FindAllStringSubmatch(content, -1)
	for _, match := range policyMatches {
		if len(match) >= 2 {
			policy := map[string]interface{}{
				"name": match[1],
			}
			policies := parsed["security_policies"].([]map[string]interface{})
			parsed["security_policies"] = append(policies, policy)
		}
	}

	// Extract zones
	zoneRegex := regexp.MustCompile(`<zone name="([^"]+)">`)
	zoneMatches := zoneRegex.FindAllStringSubmatch(content, -1)
	for _, match := range zoneMatches {
		if len(match) >= 2 {
			zone := map[string]interface{}{
				"name": match[1],
			}
			zones := parsed["zones"].([]map[string]interface{})
			parsed["zones"] = append(zones, zone)
		}
	}

	return parsed, nil
}

// ParseFortinet parses Fortinet FortiOS configuration
func (s *ConfigParserService) ParseFortinet(config []byte) (map[string]interface{}, error) {
	content := string(config)
	lines := strings.Split(content, "\n")

	parsed := map[string]interface{}{
		"firewall_policies": []map[string]interface{}{},
		"vpn_configs":       []map[string]interface{}{},
		"user_accounts":     []map[string]interface{}{},
		"interfaces":        []map[string]interface{}{},
	}

	var currentSection string
	var currentPolicy map[string]interface{}

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse sections
		if strings.HasPrefix(line, "config ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentSection = parts[1]
			}
		}

		// Parse firewall policies
		if currentSection == "firewall" && strings.HasPrefix(line, "edit ") {
			if currentPolicy != nil {
				policies := parsed["firewall_policies"].([]map[string]interface{})
				parsed["firewall_policies"] = append(policies, currentPolicy)
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentPolicy = map[string]interface{}{
					"id":          parts[1],
					"line_number": i + 1,
				}
			}
		} else if currentPolicy != nil {
			if strings.HasPrefix(line, "set ") {
				parts := strings.Fields(line)
				if len(parts) >= 3 {
					key := parts[1]
					value := strings.Join(parts[2:], " ")
					currentPolicy[key] = value
				}
			} else if strings.HasPrefix(line, "next") || strings.HasPrefix(line, "end") {
				policies := parsed["firewall_policies"].([]map[string]interface{})
				parsed["firewall_policies"] = append(policies, currentPolicy)
				currentPolicy = nil
			}
		}
	}

	if currentPolicy != nil {
		policies := parsed["firewall_policies"].([]map[string]interface{})
		parsed["firewall_policies"] = append(policies, currentPolicy)
	}

	return parsed, nil
}

// ParseJuniper parses Juniper Junos configuration
func (s *ConfigParserService) ParseJuniper(config []byte) (map[string]interface{}, error) {
	content := string(config)
	lines := strings.Split(content, "\n")

	parsed := map[string]interface{}{
		"security_policies": []map[string]interface{}{},
		"zones":             []map[string]interface{}{},
		"interfaces":         []map[string]interface{}{},
	}

	var currentSection []string

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Track indentation depth
		indent := 0
		for _, char := range line {
			if char == ' ' {
				indent++
			} else {
				break
			}
		}
		newDepth := indent / 4

		// Update current section based on depth
		if newDepth < len(currentSection) {
			currentSection = currentSection[:newDepth]
		}

		// Parse security policies
		if strings.Contains(strings.Join(currentSection, " "), "security policies") {
			if strings.HasPrefix(line, "policy ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					policy := map[string]interface{}{
						"name":        parts[1],
						"line_number": i + 1,
					}
					policies := parsed["security_policies"].([]map[string]interface{})
					parsed["security_policies"] = append(policies, policy)
				}
			}
		}

		// Parse zones
		if strings.Contains(strings.Join(currentSection, " "), "security zones") {
			if strings.HasPrefix(line, "security-zone ") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					zone := map[string]interface{}{
						"name":        parts[1],
						"line_number": i + 1,
					}
					zones := parsed["zones"].([]map[string]interface{})
					parsed["zones"] = append(zones, zone)
				}
			}
		}

		// Update section tracking
		if !strings.HasPrefix(line, "}") && !strings.HasPrefix(line, "{") {
			sectionName := strings.Fields(line)[0]
			if newDepth >= len(currentSection) {
				currentSection = append(currentSection, sectionName)
			} else {
				currentSection[newDepth] = sectionName
			}
		}
	}

	return parsed, nil
}


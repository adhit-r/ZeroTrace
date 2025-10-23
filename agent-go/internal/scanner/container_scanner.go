package scanner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"zerotrace/agent/internal/config"

	"github.com/google/uuid"
)

// ContainerScanner handles container and Kubernetes security scanning
type ContainerScanner struct {
	config *config.Config
}

// ContainerFinding represents a container security finding
type ContainerFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // image, runtime, network, storage, config
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	ContainerID   string                 `json:"container_id,omitempty"`
	ImageName     string                 `json:"image_name,omitempty"`
	Namespace     string                 `json:"namespace,omitempty"`
	PodName       string                 `json:"pod_name,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// ContainerInfo represents container information
type ContainerInfo struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Image           string                 `json:"image"`
	ImageID         string                 `json:"image_id"`
	Status          string                 `json:"status"`
	Created         time.Time              `json:"created"`
	Ports           []string               `json:"ports"`
	Mounts          []string               `json:"mounts"`
	Environment     map[string]string      `json:"environment"`
	Labels          map[string]string      `json:"labels"`
	IsRunning       bool                   `json:"is_running"`
	IsPrivileged    bool                   `json:"is_privileged"`
	HasRootUser     bool                   `json:"has_root_user"`
	HasSecrets      bool                   `json:"has_secrets"`
	HasConfigMaps   bool                   `json:"has_config_maps"`
	NetworkMode     string                 `json:"network_mode"`
	SecurityContext map[string]interface{} `json:"security_context"`
}

// KubernetesInfo represents Kubernetes cluster information
type KubernetesInfo struct {
	ClusterName       string   `json:"cluster_name"`
	Version           string   `json:"version"`
	Nodes             int      `json:"nodes"`
	Pods              int      `json:"pods"`
	Namespaces        []string `json:"namespaces"`
	RBACEnabled       bool     `json:"rbac_enabled"`
	NetworkPolicies   []string `json:"network_policies"`
	IngressRules      []string `json:"ingress_rules"`
	ServiceAccounts   []string `json:"service_accounts"`
	Secrets           []string `json:"secrets"`
	ConfigMaps        []string `json:"config_maps"`
	PersistentVolumes []string `json:"persistent_volumes"`
}

// IaCFinding represents Infrastructure as Code security finding
type IaCFinding struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`     // terraform, cloudformation, kubernetes, dockerfile
	Severity      string                 `json:"severity"` // critical, high, medium, low
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	FilePath      string                 `json:"file_path"`
	LineNumber    int                    `json:"line_number,omitempty"`
	ResourceName  string                 `json:"resource_name,omitempty"`
	CurrentValue  string                 `json:"current_value,omitempty"`
	RequiredValue string                 `json:"required_value,omitempty"`
	Remediation   string                 `json:"remediation"`
	DiscoveredAt  time.Time              `json:"discovered_at"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// NewContainerScanner creates a new container security scanner
func NewContainerScanner(cfg *config.Config) *ContainerScanner {
	return &ContainerScanner{
		config: cfg,
	}
}

// Scan performs comprehensive container and Kubernetes security scanning
func (cs *ContainerScanner) Scan() ([]ContainerFinding, []ContainerInfo, KubernetesInfo, []IaCFinding, error) {
	var findings []ContainerFinding
	var containers []ContainerInfo
	var k8sInfo KubernetesInfo
	var iacFindings []IaCFinding

	// Discover containers
	discoveredContainers := cs.discoverContainers()
	containers = append(containers, discoveredContainers...)

	// Scan each container
	for _, container := range discoveredContainers {
		containerFindings := cs.scanContainer(container)
		findings = append(findings, containerFindings...)
	}

	// Scan Kubernetes cluster
	k8sInfo = cs.scanKubernetesCluster()
	k8sFindings := cs.scanKubernetesSecurity(k8sInfo)
	findings = append(findings, k8sFindings...)

	// Scan Infrastructure as Code
	iacFindings = cs.scanIaCFiles()

	return findings, containers, k8sInfo, iacFindings, nil
}

// discoverContainers discovers running containers
func (cs *ContainerScanner) discoverContainers() []ContainerInfo {
	var containers []ContainerInfo

	// Try Docker
	dockerContainers := cs.discoverDockerContainers()
	containers = append(containers, dockerContainers...)

	// Try Podman
	podmanContainers := cs.discoverPodmanContainers()
	containers = append(containers, podmanContainers...)

	// Try containerd
	containerdContainers := cs.discoverContainerdContainers()
	containers = append(containers, containerdContainers...)

	return containers
}

// discoverDockerContainers discovers Docker containers
func (cs *ContainerScanner) discoverDockerContainers() []ContainerInfo {
	var containers []ContainerInfo

	// Check if Docker is available
	if !cs.isCommandAvailable("docker") {
		return containers
	}

	// List running containers
	cmd := exec.Command("docker", "ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}|{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		return containers
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		container := ContainerInfo{
			ID:        parts[0],
			Name:      parts[1],
			Image:     parts[2],
			Status:    parts[3],
			Ports:     strings.Split(parts[4], ","),
			IsRunning: strings.Contains(parts[3], "Up"),
		}

		// Get detailed information
		cs.enrichContainerInfo(&container, "docker")
		containers = append(containers, container)
	}

	return containers
}

// discoverPodmanContainers discovers Podman containers
func (cs *ContainerScanner) discoverPodmanContainers() []ContainerInfo {
	var containers []ContainerInfo

	// Check if Podman is available
	if !cs.isCommandAvailable("podman") {
		return containers
	}

	// List running containers
	cmd := exec.Command("podman", "ps", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}|{{.Ports}}")
	output, err := cmd.Output()
	if err != nil {
		return containers
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		container := ContainerInfo{
			ID:        parts[0],
			Name:      parts[1],
			Image:     parts[2],
			Status:    parts[3],
			Ports:     strings.Split(parts[4], ","),
			IsRunning: strings.Contains(parts[3], "Up"),
		}

		// Get detailed information
		cs.enrichContainerInfo(&container, "podman")
		containers = append(containers, container)
	}

	return containers
}

// discoverContainerdContainers discovers containerd containers
func (cs *ContainerScanner) discoverContainerdContainers() []ContainerInfo {
	var containers []ContainerInfo

	// Check if containerd is available
	if !cs.isCommandAvailable("ctr") {
		return containers
	}

	// List running containers
	cmd := exec.Command("ctr", "containers", "list")
	output, err := cmd.Output()
	if err != nil {
		return containers
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" || strings.Contains(line, "CONTAINER") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		container := ContainerInfo{
			ID:        parts[0],
			Name:      parts[1],
			Image:     parts[2],
			Status:    "running",
			IsRunning: true,
		}

		// Get detailed information
		cs.enrichContainerInfo(&container, "containerd")
		containers = append(containers, container)
	}

	return containers
}

// enrichContainerInfo enriches container information
func (cs *ContainerScanner) enrichContainerInfo(container *ContainerInfo, runtime string) {
	// Get container details
	cmd := exec.Command(runtime, "inspect", container.ID)
	output, err := cmd.Output()
	if err != nil {
		return
	}

	var inspectData []map[string]interface{}
	if err := json.Unmarshal(output, &inspectData); err != nil {
		return
	}

	if len(inspectData) > 0 {
		data := inspectData[0]

		// Extract image ID
		if imageID, ok := data["Image"].(string); ok {
			container.ImageID = imageID
		}

		// Extract creation time
		if created, ok := data["Created"].(string); ok {
			if t, err := time.Parse(time.RFC3339, created); err == nil {
				container.Created = t
			}
		}

		// Extract environment variables
		if env, ok := data["Config"].(map[string]interface{}); ok {
			if envVars, ok := env["Env"].([]interface{}); ok {
				container.Environment = make(map[string]string)
				for _, envVar := range envVars {
					if envStr, ok := envVar.(string); ok {
						parts := strings.SplitN(envStr, "=", 2)
						if len(parts) == 2 {
							container.Environment[parts[0]] = parts[1]
						}
					}
				}
			}
		}

		// Extract labels
		if labels, ok := data["Config"].(map[string]interface{}); ok {
			if labelMap, ok := labels["Labels"].(map[string]interface{}); ok {
				container.Labels = make(map[string]string)
				for k, v := range labelMap {
					if vStr, ok := v.(string); ok {
						container.Labels[k] = vStr
					}
				}
			}
		}

		// Extract mounts
		if mounts, ok := data["Mounts"].([]interface{}); ok {
			for _, mount := range mounts {
				if mountMap, ok := mount.(map[string]interface{}); ok {
					if source, ok := mountMap["Source"].(string); ok {
						container.Mounts = append(container.Mounts, source)
					}
				}
			}
		}

		// Check for privileged mode
		if hostConfig, ok := data["HostConfig"].(map[string]interface{}); ok {
			if privileged, ok := hostConfig["Privileged"].(bool); ok {
				container.IsPrivileged = privileged
			}
		}

		// Check for root user
		container.HasRootUser = cs.checkRootUser(container.ID, runtime)

		// Check for secrets
		container.HasSecrets = cs.checkSecrets(container.ID, runtime)

		// Check for config maps
		container.HasConfigMaps = cs.checkConfigMaps(container.ID, runtime)

		// Extract network mode
		if hostConfig, ok := data["HostConfig"].(map[string]interface{}); ok {
			if networkMode, ok := hostConfig["NetworkMode"].(string); ok {
				container.NetworkMode = networkMode
			}
		}
	}
}

// checkRootUser checks if container runs as root
func (cs *ContainerScanner) checkRootUser(containerID, runtime string) bool {
	cmd := exec.Command(runtime, "exec", containerID, "id", "-u")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	uid := strings.TrimSpace(string(output))
	return uid == "0"
}

// checkSecrets checks if container has secrets
func (cs *ContainerScanner) checkSecrets(containerID, runtime string) bool {
	cmd := exec.Command(runtime, "exec", containerID, "ls", "/run/secrets")
	_, err := cmd.Output()
	return err == nil
}

// checkConfigMaps checks if container has config maps
func (cs *ContainerScanner) checkConfigMaps(containerID, runtime string) bool {
	cmd := exec.Command(runtime, "exec", containerID, "ls", "/etc/config")
	_, err := cmd.Output()
	return err == nil
}

// scanContainer scans a specific container for security issues
func (cs *ContainerScanner) scanContainer(container ContainerInfo) []ContainerFinding {
	var findings []ContainerFinding

	// Check for privileged mode
	if container.IsPrivileged {
		finding := ContainerFinding{
			ID:            uuid.New().String(),
			Type:          "config",
			Severity:      "critical",
			Title:         "Privileged Container",
			Description:   fmt.Sprintf("Container %s is running in privileged mode", container.Name),
			ContainerID:   container.ID,
			ImageName:     container.Image,
			CurrentValue:  "privileged",
			RequiredValue: "non-privileged",
			Remediation:   "Remove privileged mode and use specific capabilities instead",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"container_id": container.ID,
				"image":        container.Image,
				"privileged":   true,
			},
		}
		findings = append(findings, finding)
	}

	// Check for root user
	if container.HasRootUser {
		finding := ContainerFinding{
			ID:            uuid.New().String(),
			Type:          "config",
			Severity:      "high",
			Title:         "Container Running as Root",
			Description:   fmt.Sprintf("Container %s is running as root user", container.Name),
			ContainerID:   container.ID,
			ImageName:     container.Image,
			CurrentValue:  "root",
			RequiredValue: "non-root",
			Remediation:   "Run container as non-root user",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"container_id": container.ID,
				"image":        container.Image,
				"user":         "root",
			},
		}
		findings = append(findings, finding)
	}

	// Check for exposed ports
	if len(container.Ports) > 0 {
		finding := ContainerFinding{
			ID:           uuid.New().String(),
			Type:         "network",
			Severity:     "medium",
			Title:        "Exposed Container Ports",
			Description:  fmt.Sprintf("Container %s has exposed ports: %s", container.Name, strings.Join(container.Ports, ", ")),
			ContainerID:  container.ID,
			ImageName:    container.Image,
			Remediation:  "Review and restrict exposed ports",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"container_id": container.ID,
				"image":        container.Image,
				"ports":        container.Ports,
			},
		}
		findings = append(findings, finding)
	}

	// Check for host network mode
	if container.NetworkMode == "host" {
		finding := ContainerFinding{
			ID:            uuid.New().String(),
			Type:          "network",
			Severity:      "high",
			Title:         "Container Using Host Network",
			Description:   fmt.Sprintf("Container %s is using host network mode", container.Name),
			ContainerID:   container.ID,
			ImageName:     container.Image,
			CurrentValue:  "host",
			RequiredValue: "bridge",
			Remediation:   "Use bridge network mode instead of host",
			DiscoveredAt:  time.Now(),
			Metadata: map[string]interface{}{
				"container_id": container.ID,
				"image":        container.Image,
				"network_mode": "host",
			},
		}
		findings = append(findings, finding)
	}

	// Check for secrets
	if container.HasSecrets {
		finding := ContainerFinding{
			ID:           uuid.New().String(),
			Type:         "config",
			Severity:     "medium",
			Title:        "Container Has Secrets",
			Description:  fmt.Sprintf("Container %s has mounted secrets", container.Name),
			ContainerID:  container.ID,
			ImageName:    container.Image,
			Remediation:  "Review secret usage and ensure proper access controls",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"container_id": container.ID,
				"image":        container.Image,
				"has_secrets":  true,
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanKubernetesCluster scans Kubernetes cluster
func (cs *ContainerScanner) scanKubernetesCluster() KubernetesInfo {
	info := KubernetesInfo{}

	// Check if kubectl is available
	if !cs.isCommandAvailable("kubectl") {
		return info
	}

	// Get cluster info
	cmd := exec.Command("kubectl", "cluster-info")
	output, err := cmd.Output()
	if err != nil {
		return info
	}

	info.ClusterName = "local-cluster"
	info.Version = "unknown"

	// Get nodes
	cmd = exec.Command("kubectl", "get", "nodes", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		info.Nodes = len(lines) - 1 // Subtract 1 for empty line
	}

	// Get pods
	cmd = exec.Command("kubectl", "get", "pods", "--all-namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		info.Pods = len(lines) - 1
	}

	// Get namespaces
	cmd = exec.Command("kubectl", "get", "namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.Namespaces = append(info.Namespaces, parts[0])
				}
			}
		}
	}

	// Check RBAC
	cmd = exec.Command("kubectl", "get", "clusterroles")
	_, err = cmd.Output()
	info.RBACEnabled = err == nil

	// Get network policies
	cmd = exec.Command("kubectl", "get", "networkpolicies", "--all-namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.NetworkPolicies = append(info.NetworkPolicies, parts[0])
				}
			}
		}
	}

	// Get ingress rules
	cmd = exec.Command("kubectl", "get", "ingress", "--all-namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.IngressRules = append(info.IngressRules, parts[0])
				}
			}
		}
	}

	// Get service accounts
	cmd = exec.Command("kubectl", "get", "serviceaccounts", "--all-namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.ServiceAccounts = append(info.ServiceAccounts, parts[0])
				}
			}
		}
	}

	// Get secrets
	cmd = exec.Command("kubectl", "get", "secrets", "--all-namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.Secrets = append(info.Secrets, parts[0])
				}
			}
		}
	}

	// Get config maps
	cmd = exec.Command("kubectl", "get", "configmaps", "--all-namespaces", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.ConfigMaps = append(info.ConfigMaps, parts[0])
				}
			}
		}
	}

	// Get persistent volumes
	cmd = exec.Command("kubectl", "get", "persistentvolumes", "--no-headers")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line != "" {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					info.PersistentVolumes = append(info.PersistentVolumes, parts[0])
				}
			}
		}
	}

	return info
}

// scanKubernetesSecurity scans Kubernetes security
func (cs *ContainerScanner) scanKubernetesSecurity(info KubernetesInfo) []ContainerFinding {
	var findings []ContainerFinding

	// Check for RBAC
	if !info.RBACEnabled {
		finding := ContainerFinding{
			ID:           uuid.New().String(),
			Type:         "config",
			Severity:     "high",
			Title:        "RBAC Not Enabled",
			Description:  "Kubernetes cluster does not have RBAC enabled",
			Remediation:  "Enable RBAC for proper access control",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"rbac_enabled": false,
			},
		}
		findings = append(findings, finding)
	}

	// Check for network policies
	if len(info.NetworkPolicies) == 0 {
		finding := ContainerFinding{
			ID:           uuid.New().String(),
			Type:         "network",
			Severity:     "medium",
			Title:        "No Network Policies",
			Description:  "Kubernetes cluster has no network policies configured",
			Remediation:  "Implement network policies for network segmentation",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"network_policies": 0,
			},
		}
		findings = append(findings, finding)
	}

	// Check for excessive secrets
	if len(info.Secrets) > 50 {
		finding := ContainerFinding{
			ID:           uuid.New().String(),
			Type:         "config",
			Severity:     "low",
			Title:        "Excessive Secrets",
			Description:  fmt.Sprintf("Kubernetes cluster has %d secrets", len(info.Secrets)),
			Remediation:  "Review and clean up unnecessary secrets",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"secrets_count": len(info.Secrets),
			},
		}
		findings = append(findings, finding)
	}

	// Check for excessive config maps
	if len(info.ConfigMaps) > 100 {
		finding := ContainerFinding{
			ID:           uuid.New().String(),
			Type:         "config",
			Severity:     "low",
			Title:        "Excessive Config Maps",
			Description:  fmt.Sprintf("Kubernetes cluster has %d config maps", len(info.ConfigMaps)),
			Remediation:  "Review and consolidate config maps",
			DiscoveredAt: time.Now(),
			Metadata: map[string]interface{}{
				"config_maps_count": len(info.ConfigMaps),
			},
		}
		findings = append(findings, finding)
	}

	return findings
}

// scanIaCFiles scans Infrastructure as Code files
func (cs *ContainerScanner) scanIaCFiles() []IaCFinding {
	var findings []IaCFinding

	// Scan for Terraform files
	terraformFindings := cs.scanTerraformFiles()
	findings = append(findings, terraformFindings...)

	// Scan for CloudFormation files
	cloudformationFindings := cs.scanCloudFormationFiles()
	findings = append(findings, cloudformationFindings...)

	// Scan for Kubernetes manifests
	k8sFindings := cs.scanKubernetesManifests()
	findings = append(findings, k8sFindings...)

	// Scan for Dockerfiles
	dockerfileFindings := cs.scanDockerfiles()
	findings = append(findings, dockerfileFindings...)

	return findings
}

// scanTerraformFiles scans Terraform files
func (cs *ContainerScanner) scanTerraformFiles() []IaCFinding {
	var findings []IaCFinding

	// This would require parsing Terraform files
	// For now, return a placeholder finding
	finding := IaCFinding{
		ID:           uuid.New().String(),
		Type:         "terraform",
		Severity:     "medium",
		Title:        "Terraform Security Scan",
		Description:  "Terraform files should be scanned for security issues",
		Remediation:  "Use tools like tfsec or Checkov to scan Terraform files",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"tool": "terraform",
		},
	}
	findings = append(findings, finding)

	return findings
}

// scanCloudFormationFiles scans CloudFormation files
func (cs *ContainerScanner) scanCloudFormationFiles() []IaCFinding {
	var findings []IaCFinding

	// This would require parsing CloudFormation files
	// For now, return a placeholder finding
	finding := IaCFinding{
		ID:           uuid.New().String(),
		Type:         "cloudformation",
		Severity:     "medium",
		Title:        "CloudFormation Security Scan",
		Description:  "CloudFormation templates should be scanned for security issues",
		Remediation:  "Use tools like cfn-lint or Checkov to scan CloudFormation templates",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"tool": "cloudformation",
		},
	}
	findings = append(findings, finding)

	return findings
}

// scanKubernetesManifests scans Kubernetes manifests
func (cs *ContainerScanner) scanKubernetesManifests() []IaCFinding {
	var findings []IaCFinding

	// This would require parsing Kubernetes YAML files
	// For now, return a placeholder finding
	finding := IaCFinding{
		ID:           uuid.New().String(),
		Type:         "kubernetes",
		Severity:     "medium",
		Title:        "Kubernetes Manifest Security Scan",
		Description:  "Kubernetes manifests should be scanned for security issues",
		Remediation:  "Use tools like kube-score or Polaris to scan Kubernetes manifests",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"tool": "kubernetes",
		},
	}
	findings = append(findings, finding)

	return findings
}

// scanDockerfiles scans Dockerfiles
func (cs *ContainerScanner) scanDockerfiles() []IaCFinding {
	var findings []IaCFinding

	// This would require parsing Dockerfile files
	// For now, return a placeholder finding
	finding := IaCFinding{
		ID:           uuid.New().String(),
		Type:         "dockerfile",
		Severity:     "medium",
		Title:        "Dockerfile Security Scan",
		Description:  "Dockerfiles should be scanned for security issues",
		Remediation:  "Use tools like hadolint or dockerfilelint to scan Dockerfiles",
		DiscoveredAt: time.Now(),
		Metadata: map[string]interface{}{
			"tool": "dockerfile",
		},
	}
	findings = append(findings, finding)

	return findings
}

// isCommandAvailable checks if a command is available
func (cs *ContainerScanner) isCommandAvailable(command string) bool {
	cmd := exec.Command("which", command)
	err := cmd.Run()
	return err == nil
}

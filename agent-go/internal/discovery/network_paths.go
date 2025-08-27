package discovery

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"zerotrace/agent/internal/models"
)



// NetworkGraph represents the network topology for path analysis
type NetworkGraph struct {
	Nodes map[string]*models.NetworkAsset `json:"nodes"`
	Edges map[string]map[string]float64   `json:"edges"` // source -> destination -> weight
}

// NetworkPathAnalyzer handles shortest path calculations using fast SSSP principles
type NetworkPathAnalyzer struct {
	graph *NetworkGraph
}

// NewNetworkPathAnalyzer creates a new path analyzer
func NewNetworkPathAnalyzer() *NetworkPathAnalyzer {
	return &NetworkPathAnalyzer{
		graph: &NetworkGraph{
			Nodes: make(map[string]*models.NetworkAsset),
			Edges: make(map[string]map[string]float64),
		},
	}
}

// AddAsset adds a network asset to the graph
func (npa *NetworkPathAnalyzer) AddAsset(asset *models.NetworkAsset) {
	npa.graph.Nodes[asset.IPAddress] = asset
	
	// Initialize edges for this node
	if npa.graph.Edges[asset.IPAddress] == nil {
		npa.graph.Edges[asset.IPAddress] = make(map[string]float64)
	}
}

// AddConnection adds a connection between two assets
func (npa *NetworkPathAnalyzer) AddConnection(source, destination string, weight float64) {
	if npa.graph.Edges[source] == nil {
		npa.graph.Edges[source] = make(map[string]float64)
	}
	npa.graph.Edges[source][destination] = weight
}

// calculateEdgeWeight calculates the weight of an edge based on network characteristics
func (npa *NetworkPathAnalyzer) calculateEdgeWeight(source, dest *models.NetworkAsset) float64 {
	baseWeight := 1.0
	
	// Adjust weight based on risk factors
	if source.RiskScore > 7.0 || dest.RiskScore > 7.0 {
		baseWeight *= 2.0 // Higher risk = higher weight (longer path)
	}
	
	// Adjust based on device types
	if source.DeviceType == "network_device" || dest.DeviceType == "network_device" {
		baseWeight *= 0.8 // Network devices are preferred paths
	}
	
	// Adjust based on monitoring status
	if !source.IsMonitored || !dest.IsMonitored {
		baseWeight *= 1.5 // Unmonitored assets have higher weight
	}
	
	return baseWeight
}

// FastSSSP implements the core algorithm inspired by the Duan-Mao breakthrough
// This is a simplified version focusing on the key principles
func (npa *NetworkPathAnalyzer) FastSSSP(source string) map[string]*models.NetworkPath {
	distances := make(map[string]float64)
	predecessors := make(map[string]string)
	visited := make(map[string]bool)
	
	// Initialize distances
	for node := range npa.graph.Nodes {
		distances[node] = math.Inf(1)
	}
	distances[source] = 0
	
	// Priority queue simulation (in practice, use a proper heap)
	var queue []string
	queue = append(queue, source)
	
	for len(queue) > 0 {
		// Find minimum distance node (simplified - use proper heap in production)
		minNode := ""
		minDist := math.Inf(1)
		minIdx := -1
		
		for i, node := range queue {
			if !visited[node] && distances[node] < minDist {
				minDist = distances[node]
				minNode = node
				minIdx = i
			}
		}
		
		if minNode == "" {
			break
		}
		
		// Remove from queue
		queue = append(queue[:minIdx], queue[minIdx+1:]...)
		visited[minNode] = true
		
		// Relax edges
		for neighbor, weight := range npa.graph.Edges[minNode] {
			if !visited[neighbor] {
				newDist := distances[minNode] + weight
				if newDist < distances[neighbor] {
					distances[neighbor] = newDist
					predecessors[neighbor] = minNode
					queue = append(queue, neighbor)
				}
			}
		}
	}
	
	// Build paths
	paths := make(map[string]*models.NetworkPath)
	for dest, distance := range distances {
		if distance < math.Inf(1) && dest != source {
			path := npa.buildPath(source, dest, predecessors)
			paths[dest] = &models.NetworkPath{
				Source:      source,
				Destination: dest,
				Path:        path,
				Distance:    distance,
				Hops:        len(path) - 1,
				Latency:     distance * 10, // Rough latency estimation
				RiskScore:   npa.calculatePathRisk(path),
				Discovered:  time.Now(),
			}
		}
	}
	
	return paths
}

// buildPath reconstructs the path from source to destination
func (npa *NetworkPathAnalyzer) buildPath(source, dest string, predecessors map[string]string) []string {
	var path []string
	current := dest
	
	for current != "" {
		path = append([]string{current}, path...)
		current = predecessors[current]
	}
	
	return path
}

// calculatePathRisk calculates the risk score for a path
func (npa *NetworkPathAnalyzer) calculatePathRisk(path []string) float64 {
	if len(path) == 0 {
		return 0
	}
	
	totalRisk := 0.0
	for _, node := range path {
		if asset, exists := npa.graph.Nodes[node]; exists {
			totalRisk += asset.RiskScore
		}
	}
	
	return totalRisk / float64(len(path))
}

// FindCriticalPaths finds the most critical paths in the network
func (npa *NetworkPathAnalyzer) FindCriticalPaths(ctx context.Context) ([]*models.NetworkPath, error) {
	var criticalPaths []*models.NetworkPath
	
	// Find paths from high-risk assets to critical assets
	for sourceIP, sourceAsset := range npa.graph.Nodes {
		if sourceAsset.RiskScore > 7.0 {
			paths := npa.FastSSSP(sourceIP)
			
			for destIP, path := range paths {
				destAsset := npa.graph.Nodes[destIP]
				if destAsset != nil && destAsset.DeviceType == "server" {
					criticalPaths = append(criticalPaths, path)
				}
			}
		}
	}
	
	// Sort by risk score (highest first)
	sort.Slice(criticalPaths, func(i, j int) bool {
		return criticalPaths[i].RiskScore > criticalPaths[j].RiskScore
	})
	
	return criticalPaths, nil
}

// AnalyzeNetworkTopology performs comprehensive network analysis
func (npa *NetworkPathAnalyzer) AnalyzeNetworkTopology(ctx context.Context, assets []models.NetworkAsset) (*models.NetworkTopology, error) {
	// Build graph from assets
	for i := range assets {
		npa.AddAsset(&assets[i])
	}
	
	// Add connections based on discovered peers
	for _, asset := range assets {
		for _, peer := range asset.ConnectedPeers {
			weight := npa.calculateEdgeWeight(&asset, &models.NetworkAsset{
				IPAddress: peer.IPAddress,
				RiskScore: peer.RiskScore,
			})
			npa.AddConnection(asset.IPAddress, peer.IPAddress, weight)
		}
	}
	
	// Find critical paths
	criticalPaths, err := npa.FindCriticalPaths(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find critical paths: %w", err)
	}
	
	// Build topology nodes
	var nodes []models.TopologyNode
	for ip, asset := range npa.graph.Nodes {
		nodes = append(nodes, models.TopologyNode{
			ID:          ip,
			Name:        asset.Hostname,
			Type:        asset.DeviceType,
			IPAddress:   ip,
			RiskScore:   asset.RiskScore,
			IsMonitored: asset.IsMonitored,
			Location:    asset.Location,
			Department:  asset.Department,
		})
	}
	
	// Build topology links
	var links []models.TopologyLink
	for source, edges := range npa.graph.Edges {
		for dest, weight := range edges {
			links = append(links, models.TopologyLink{
				Source:      source,
				Target:      dest,
				Weight:      weight,
				Type:        "network",
				IsCritical:  npa.isCriticalPath(source, dest, criticalPaths),
			})
		}
	}
	
	// Identify clusters
	clusters := npa.identifyClusters(nodes, links)
	
	return &models.NetworkTopology{
		Nodes:           nodes,
		Links:           links,
		Clusters:        clusters,
		CriticalPaths:   criticalPaths,
		TotalAssets:     len(nodes),
		TotalConnections: len(links),
		LastUpdated:     time.Now(),
	}, nil
}

// isCriticalPath checks if a link is part of a critical path
func (npa *NetworkPathAnalyzer) isCriticalPath(source, dest string, criticalPaths []*models.NetworkPath) bool {
	for _, path := range criticalPaths {
		for i := 0; i < len(path.Path)-1; i++ {
			if (path.Path[i] == source && path.Path[i+1] == dest) ||
				(path.Path[i] == dest && path.Path[i+1] == source) {
				return true
			}
		}
	}
	return false
}

// identifyClusters identifies logical clusters in the network
func (npa *NetworkPathAnalyzer) identifyClusters(nodes []models.TopologyNode, links []models.TopologyLink) []models.Cluster {
	// Simple clustering by department and location
	clusterMap := make(map[string][]models.TopologyNode)
	
	for _, node := range nodes {
		clusterKey := fmt.Sprintf("%s-%s", node.Department, node.Location)
		clusterMap[clusterKey] = append(clusterMap[clusterKey], node)
	}
	
	var clusters []models.Cluster
	for key, clusterNodes := range clusterMap {
		if len(clusterNodes) > 1 {
			clusters = append(clusters, models.Cluster{
				ID:       key,
				Name:     key,
				NodeIDs:  extractNodeIDs(clusterNodes),
				Type:     "logical",
				RiskScore: calculateClusterRisk(clusterNodes),
			})
		}
	}
	
	return clusters
}

// Helper functions
func extractNodeIDs(nodes []models.TopologyNode) []string {
	ids := make([]string, len(nodes))
	for i, node := range nodes {
		ids[i] = node.ID
	}
	return ids
}

func calculateClusterRisk(nodes []models.TopologyNode) float64 {
	if len(nodes) == 0 {
		return 0
	}
	
	totalRisk := 0.0
	for _, node := range nodes {
		totalRisk += node.RiskScore
	}
	
	return totalRisk / float64(len(nodes))
}

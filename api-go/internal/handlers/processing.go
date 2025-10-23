package handlers

import (
	"fmt"
	"net/http"
	"time"

	"zerotrace/api/internal/models"
	"zerotrace/api/internal/services"

	"github.com/gin-gonic/gin"
)

// ProcessingStatus represents the current processing status
type ProcessingStatus struct {
	AgentID           string                 `json:"agent_id"`
	Status            string                 `json:"status"` // idle, processing, completed, error
	CurrentApp        string                 `json:"current_app"`
	AppsProcessed     int                    `json:"apps_processed"`
	TotalApps         int                    `json:"total_apps"`
	ProgressPercent   float64                `json:"progress_percent"`
	CPEMatches        int                    `json:"cpe_matches"`
	Vulnerabilities   int                    `json:"vulnerabilities_found"`
	ProcessingLogs    []ProcessingLog        `json:"processing_logs"`
	LastUpdated       time.Time              `json:"last_updated"`
	EstimatedTimeLeft string                 `json:"estimated_time_left"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// ProcessingLog represents a processing log entry
type ProcessingLog struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // info, warning, error
	Message   string    `json:"message"`
	AppName   string    `json:"app_name,omitempty"`
	CPE       string    `json:"cpe,omitempty"`
	Vulns     int       `json:"vulnerabilities,omitempty"`
}

// GetProcessingStatus returns the current processing status for all agents
func GetProcessingStatus(agentService *services.AgentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		agents := agentService.GetAllAgents()

		var processingStatuses []ProcessingStatus
		for _, agent := range agents {
			status := ProcessingStatus{
				AgentID:           agent.ID.String(),
				Status:            getProcessingStatus(agent),
				CurrentApp:        getCurrentApp(agent),
				AppsProcessed:     getAppsProcessed(agent),
				TotalApps:         getTotalApps(agent),
				ProgressPercent:   getProgressPercent(agent),
				CPEMatches:        getCPEMatches(agent),
				Vulnerabilities:   getVulnerabilities(agent),
				ProcessingLogs:    getProcessingLogs(agent),
				LastUpdated:       agent.UpdatedAt,
				EstimatedTimeLeft: getEstimatedTimeLeft(agent),
				Metadata:          agent.Metadata,
			}
			processingStatuses = append(processingStatuses, status)
		}

		c.JSON(http.StatusOK, models.APIResponse{
			Success:   true,
			Data:      processingStatuses,
			Message:   "Processing status retrieved successfully",
			Timestamp: time.Now(),
		})
	}
}

// Helper functions to extract processing information from agent metadata
func getProcessingStatus(agent *models.Agent) string {
	if status, ok := agent.Metadata["processing_status"].(string); ok {
		return status
	}
	return "idle"
}

func getCurrentApp(agent *models.Agent) string {
	if app, ok := agent.Metadata["current_app"].(string); ok {
		return app
	}
	return ""
}

func getAppsProcessed(agent *models.Agent) int {
	if count, ok := agent.Metadata["apps_processed"].(int); ok {
		return count
	}
	if deps, ok := agent.Metadata["dependencies"].([]interface{}); ok {
		return len(deps)
	}
	return 0
}

func getTotalApps(agent *models.Agent) int {
	if total, ok := agent.Metadata["total_apps"].(int); ok {
		return total
	}
	return 100 // Default estimate
}

func getProgressPercent(agent *models.Agent) float64 {
	processed := getAppsProcessed(agent)
	total := getTotalApps(agent)
	if total == 0 {
		return 0
	}
	return float64(processed) / float64(total) * 100
}

func getCPEMatches(agent *models.Agent) int {
	if matches, ok := agent.Metadata["cpe_matches"].(int); ok {
		return matches
	}
	return 0
}

func getVulnerabilities(agent *models.Agent) int {
	if vulns, ok := agent.Metadata["total_vulnerabilities"].(int); ok {
		return vulns
	}
	return 0
}

func getProcessingLogs(agent *models.Agent) []ProcessingLog {
	if logs, ok := agent.Metadata["processing_logs"].([]interface{}); ok {
		var processingLogs []ProcessingLog
		for _, log := range logs {
			if logMap, ok := log.(map[string]interface{}); ok {
				processingLog := ProcessingLog{
					Timestamp: time.Now(), // Default to now if not specified
					Level:     "info",
					Message:   "",
				}

				if timestamp, ok := logMap["timestamp"].(string); ok {
					if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
						processingLog.Timestamp = t
					}
				}
				if level, ok := logMap["level"].(string); ok {
					processingLog.Level = level
				}
				if message, ok := logMap["message"].(string); ok {
					processingLog.Message = message
				}
				if appName, ok := logMap["app_name"].(string); ok {
					processingLog.AppName = appName
				}
				if cpe, ok := logMap["cpe"].(string); ok {
					processingLog.CPE = cpe
				}
				if vulns, ok := logMap["vulnerabilities"].(int); ok {
					processingLog.Vulns = vulns
				}

				processingLogs = append(processingLogs, processingLog)
			}
		}
		return processingLogs
	}
	return []ProcessingLog{}
}

func getEstimatedTimeLeft(agent *models.Agent) string {
	if timeLeft, ok := agent.Metadata["estimated_time_left"].(string); ok {
		return timeLeft
	}

	// Calculate based on processing rate
	processed := getAppsProcessed(agent)
	total := getTotalApps(agent)
	if processed == 0 || total == 0 {
		return "Unknown"
	}

	// Simple estimation: assume 1 app per second
	remaining := total - processed
	if remaining <= 0 {
		return "Completed"
	}

	minutes := remaining / 60
	seconds := remaining % 60

	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

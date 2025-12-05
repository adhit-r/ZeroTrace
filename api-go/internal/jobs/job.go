package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// JobType represents different job types
type JobType string

const (
	JobTypeEnrichment     JobType = "enrichment"
	JobTypeVulnerability  JobType = "vulnerability_scan"
	JobTypeReport         JobType = "report_generation"
	JobTypeNotification   JobType = "notification"
	JobTypeCleanup        JobType = "cleanup"
)

// JobPayload represents job payload
type JobPayload struct {
	JobID      string                 `json:"job_id"`
	CompanyID  string                 `json:"company_id"`
	AgentID    string                 `json:"agent_id,omitempty"`
	Data       map[string]interface{} `json:"data"`
	RetryCount int                    `json:"retry_count"`
	CreatedAt  time.Time              `json:"created_at"`
}

// JobManager manages background jobs using asynq
type JobManager struct {
	client *asynq.Client
	server *asynq.Server
	mux    *asynq.ServeMux
}

// NewJobManager creates a new job manager
// Uses Valkey (Redis-compatible) for job queue storage
func NewJobManager(redisOpt *asynq.RedisClientOpt) *JobManager {
	client := asynq.NewClient(redisOpt)
	
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			RetryDelayFunc: func(n int, e error, task *asynq.Task) time.Duration {
				return time.Duration(n*n) * time.Second
			},
		},
	)
	
	mux := asynq.NewServeMux()
	
	return &JobManager{
		client: client,
		server: server,
		mux:    mux,
	}
}

// EnqueueJob enqueues a new job
func (jm *JobManager) EnqueueJob(
	ctx context.Context,
	jobType JobType,
	payload JobPayload,
	opts ...asynq.Option,
) (*asynq.TaskInfo, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	task := asynq.NewTask(string(jobType), payloadJSON, opts...)
	
	info, err := jm.client.Enqueue(task)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue job: %w", err)
	}
	
	return info, nil
}

// RegisterHandler registers a job handler
func (jm *JobManager) RegisterHandler(jobType JobType, handler asynq.Handler) {
	jm.mux.HandleFunc(string(jobType), handler)
}

// Start starts the job server
func (jm *JobManager) Start() error {
	return jm.server.Run(jm.mux)
}

// Shutdown gracefully shuts down the job server
func (jm *JobManager) Shutdown() {
	jm.server.Shutdown()
	jm.client.Close()
}

// GetRedisClientOpt creates Valkey (Redis-compatible) client options from config
// Note: Uses Redis client library as Valkey is fully Redis-compatible
func GetRedisClientOpt(host string, port int, password string, db int) *asynq.RedisClientOpt {
	return &asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	}
}


package services

import (
	"log"
	"sync"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/constants"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
)

// ConfigJobService handles asynchronous config analysis jobs
type ConfigJobService struct {
	configFileRepo    *repository.ConfigFileRepository
	parserService     *ConfigParserService
	analyzerService   *ConfigAnalyzerService
	jobQueue          chan uuid.UUID
	workerCount       int
	wg                sync.WaitGroup
	stopChan          chan struct{}
}

// NewConfigJobService creates a new config job service
func NewConfigJobService(
	configFileRepo *repository.ConfigFileRepository,
	parserService *ConfigParserService,
	analyzerService *ConfigAnalyzerService,
	cfg *config.Config,
) *ConfigJobService {
	workerCount := cfg.ConfigAuditorWorkerCount
	if workerCount <= 0 {
		workerCount = constants.DefaultWorkerCount // Fallback to constant
	}

	queueBufferSize := cfg.ConfigAuditorQueueBufferSize
	if queueBufferSize <= 0 {
		queueBufferSize = constants.DefaultQueueBufferSize // Fallback to constant
	}

	service := &ConfigJobService{
		configFileRepo:  configFileRepo,
		parserService:   parserService,
		analyzerService: analyzerService,
		jobQueue:        make(chan uuid.UUID, queueBufferSize),
		workerCount:     workerCount,
		stopChan:        make(chan struct{}),
	}

	// Start workers
	service.startWorkers()

	return service
}

// QueueConfigAnalysis queues a config file for analysis
func (s *ConfigJobService) QueueConfigAnalysis(configFileID uuid.UUID) error {
	select {
	case s.jobQueue <- configFileID:
		log.Printf("Queued config analysis for file: %s", configFileID)
		return nil
	default:
		log.Printf("Job queue full, dropping config file: %s", configFileID)
		return nil // Don't error, just log
	}
}

// ProcessConfigAnalysis processes a config file analysis
func (s *ConfigJobService) ProcessConfigAnalysis(configFileID uuid.UUID) error {
	// Get config file
	configFile, err := s.configFileRepo.GetByID(configFileID)
	if err != nil {
		return err
	}

	// Step 1: Parse the config file
	if configFile.ParsingStatus != "parsed" {
		err = s.parserService.ParseConfigFile(configFile)
		if err != nil {
			log.Printf("Failed to parse config file %s: %v", configFileID, err)
			return err
		}

		// Reload config file to get parsed data
		configFile, err = s.configFileRepo.GetByID(configFileID)
		if err != nil {
			return err
		}
	}

	// Step 2: Analyze the config file
	if configFile.ParsingStatus == "parsed" {
		err = s.analyzerService.AnalyzeConfigFile(configFileID)
		if err != nil {
			log.Printf("Failed to analyze config file %s: %v", configFileID, err)
			return err
		}
	}

	return nil
}

// startWorkers starts background workers to process jobs
func (s *ConfigJobService) startWorkers() {
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

// worker processes jobs from the queue
func (s *ConfigJobService) worker(id int) {
	defer s.wg.Done()

	log.Printf("Config analysis worker %d started", id)

	for {
		select {
		case configFileID := <-s.jobQueue:
			log.Printf("Worker %d processing config file: %s", id, configFileID)
			err := s.ProcessConfigAnalysis(configFileID)
			if err != nil {
				log.Printf("Worker %d error processing %s: %v", id, configFileID, err)
			} else {
				log.Printf("Worker %d completed processing: %s", id, configFileID)
			}

		case <-s.stopChan:
			log.Printf("Config analysis worker %d stopping", id)
			return
		}
	}
}

// Stop stops all workers
func (s *ConfigJobService) Stop() {
	close(s.stopChan)
	s.wg.Wait()
	log.Println("All config analysis workers stopped")
}

// GetAnalysisStatus gets the analysis status for a config file
func (s *ConfigJobService) GetAnalysisStatus(configFileID uuid.UUID) (string, error) {
	configFile, err := s.configFileRepo.GetByID(configFileID)
	if err != nil {
		return "", err
	}
	return configFile.AnalysisStatus, nil
}


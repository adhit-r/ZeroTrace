package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EnrollmentService handles agent enrollment and token management
type EnrollmentService struct {
	cfg *config.Config
	db  *repository.Database
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(cfg *config.Config, db *repository.Database) *EnrollmentService {
	return &EnrollmentService{
		cfg: cfg,
		db:  db,
	}
}

// GenerateEnrollmentToken creates a new enrollment token for an organization
func (s *EnrollmentService) GenerateEnrollmentToken(req *models.GenerateEnrollmentTokenRequest, issuedBy uuid.UUID) (*models.EnrollmentToken, error) {
	// Generate a secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash the token for storage
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashStr := hex.EncodeToString(tokenHash[:])

	// Set expiration time
	expiresIn := req.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = 60 // default 60 minutes
	}
	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Minute)

	enrollmentToken := &models.EnrollmentToken{
		ID:             uuid.New(),
		OrganizationID: req.OrganizationID,
		Token:          token,
		TokenHash:      tokenHashStr,
		IssuedBy:       issuedBy,
		IssuedAt:       time.Now(),
		ExpiresAt:      expiresAt,
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Save to database
	if err := s.db.DB.Create(enrollmentToken).Error; err != nil {
		return nil, fmt.Errorf("failed to save enrollment token: %w", err)
	}

	return enrollmentToken, nil
}

// ValidateEnrollmentToken validates an enrollment token
func (s *EnrollmentService) ValidateEnrollmentToken(token string) (*models.EnrollmentToken, error) {
	// Hash the provided token (for database lookup)
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashStr := hex.EncodeToString(tokenHash[:])

	// Look up token in database
	var enrollmentToken models.EnrollmentToken
	if err := s.db.DB.Where("token_hash = ?", tokenHashStr).First(&enrollmentToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid enrollment token")
		}
		return nil, fmt.Errorf("failed to find enrollment token: %w", err)
	}

	// Check expiration
	if time.Now().After(enrollmentToken.ExpiresAt) {
		return nil, fmt.Errorf("enrollment token expired")
	}

	// Check status
	if enrollmentToken.Status != "active" {
		return nil, fmt.Errorf("enrollment token is not active (status: %s)", enrollmentToken.Status)
	}

	return &enrollmentToken, nil
}

// EnrollAgent enrolls an agent using an enrollment token
func (s *EnrollmentService) EnrollAgent(req *models.AgentEnrollmentRequest) (*models.AgentEnrollmentResponse, error) {
	// Validate the enrollment token
	enrollmentToken, err := s.ValidateEnrollmentToken(req.EnrollmentToken)
	if err != nil {
		return nil, fmt.Errorf("invalid enrollment token: %w", err)
	}

	// Create a new agent
	agentID := uuid.New()
	// Save agent to database
	agent := &models.Agent{
		ID:             agentID,
		OrganizationID: enrollmentToken.OrganizationID,
		CompanyID:      uuid.Nil, // This would be looked up from the organization
		Name:           req.AgentInfo.Hostname,
		Status:         "active",
		Version:        req.AgentInfo.Version,
		LastSeen:       time.Now(),
		Hostname:       req.AgentInfo.Hostname,
		OS:             req.AgentInfo.OS,
		Metadata:       req.AgentInfo.Metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.DB.Create(agent).Error; err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Generate a long-lived credential for the agent
	credential, err := s.generateAgentCredential(agentID, enrollmentToken.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate agent credential: %w", err)
	}

	// Mark the enrollment token as used
	enrollmentToken.Status = "used"
	enrollmentToken.UsedAt = &time.Time{}
	*enrollmentToken.UsedAt = time.Now()
	enrollmentToken.UsedBy = &agentID
	enrollmentToken.UpdatedAt = time.Now()

	// Update enrollment token in database
	if err := s.db.DB.Save(enrollmentToken).Error; err != nil {
		return nil, fmt.Errorf("failed to update enrollment token: %w", err)
	}

	// Save agent credential to database
	if err := s.db.DB.Create(credential).Error; err != nil {
		return nil, fmt.Errorf("failed to save agent credential: %w", err)
	}

	return &models.AgentEnrollmentResponse{
		AgentID:        agentID,
		OrganizationID: enrollmentToken.OrganizationID,
		Credential:     credential.CredentialHash,   // In real implementation, this would be the actual credential
		ExpiresAt:      time.Now().AddDate(1, 0, 0), // 1 year
	}, nil
}

// generateAgentCredential generates a long-lived credential for an agent
func (s *EnrollmentService) generateAgentCredential(agentID, organizationID uuid.UUID) (*models.AgentCredential, error) {
	// Generate a secure credential
	credentialBytes := make([]byte, 32)
	if _, err := rand.Read(credentialBytes); err != nil {
		return nil, fmt.Errorf("failed to generate credential: %w", err)
	}
	credential := hex.EncodeToString(credentialBytes)

	// Hash the credential for storage
	credentialHash := sha256.Sum256([]byte(credential))
	credentialHashStr := hex.EncodeToString(credentialHash[:])

	// Set expiration (1 year from now)
	expiresAt := time.Now().AddDate(1, 0, 0)

	agentCredential := &models.AgentCredential{
		ID:             uuid.New(),
		AgentID:        agentID,
		OrganizationID: organizationID,
		CredentialHash: credentialHashStr,
		IssuedAt:       time.Now(),
		ExpiresAt:      &expiresAt,
		LastUsedAt:     time.Now(),
		Status:         "active",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return agentCredential, nil
}

// ValidateAgentCredential validates an agent's credential
func (s *EnrollmentService) ValidateAgentCredential(credential string, agentID uuid.UUID) (*models.AgentCredential, error) {
	// Hash the credential
	credentialHash := sha256.Sum256([]byte(credential))
	credentialHashStr := hex.EncodeToString(credentialHash[:])

	// Look up credential in database
	var agentCredential models.AgentCredential
	if err := s.db.DB.Where("credential_hash = ? AND agent_id = ?", credentialHashStr, agentID).First(&agentCredential).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid agent credential")
		}
		return nil, fmt.Errorf("failed to find agent credential: %w", err)
	}

	// Check status
	if agentCredential.Status != "active" {
		return nil, fmt.Errorf("agent credential is not active")
	}

	// Check expiration
	if agentCredential.ExpiresAt != nil && time.Now().After(*agentCredential.ExpiresAt) {
		return nil, fmt.Errorf("agent credential expired")
	}

	// Update LastUsedAt
	agentCredential.LastUsedAt = time.Now()
	// Update in background or simple save
	go func() {
		s.db.DB.Model(&agentCredential).Update("last_used_at", time.Now())
	}()

	return &agentCredential, nil
}

// RevokeEnrollmentToken revokes an enrollment token
func (s *EnrollmentService) RevokeEnrollmentToken(tokenID uuid.UUID) error {
	// Look up token in database
	var enrollmentToken models.EnrollmentToken
	if err := s.db.DB.Where("id = ?", tokenID).First(&enrollmentToken).Error; err != nil {
		return fmt.Errorf("failed to find enrollment token: %w", err)
	}

	enrollmentToken.Status = "revoked"
	enrollmentToken.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&enrollmentToken).Error; err != nil {
		return fmt.Errorf("failed to revoke enrollment token: %w", err)
	}

	return nil
}

// RevokeAgentCredential revokes an agent's credential
func (s *EnrollmentService) RevokeAgentCredential(credentialID uuid.UUID) error {
	// Look up credential in database
	var agentCredential models.AgentCredential
	if err := s.db.DB.Where("id = ?", credentialID).First(&agentCredential).Error; err != nil {
		return fmt.Errorf("failed to find agent credential: %w", err)
	}

	agentCredential.Status = "revoked"
	agentCredential.UpdatedAt = time.Now()

	if err := s.db.DB.Save(&agentCredential).Error; err != nil {
		return fmt.Errorf("failed to revoke agent credential: %w", err)
	}

	return nil
}

package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/models"

	"github.com/google/uuid"

)

// EnrollmentService handles agent enrollment and token management
type EnrollmentService struct {
	cfg *config.Config
}

// NewEnrollmentService creates a new enrollment service
func NewEnrollmentService(cfg *config.Config) *EnrollmentService {
	return &EnrollmentService{
		cfg: cfg,
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

	// TODO: Save to database
	// if err := s.repo.CreateEnrollmentToken(enrollmentToken); err != nil {
	//     return nil, fmt.Errorf("failed to save enrollment token: %w", err)
	// }

	return enrollmentToken, nil
}

// ValidateEnrollmentToken validates an enrollment token
func (s *EnrollmentService) ValidateEnrollmentToken(token string) (*models.EnrollmentToken, error) {
	// Hash the provided token (for future database lookup)
	_ = sha256.Sum256([]byte(token))

	// TODO: Look up token in database
	// enrollmentToken, err := s.repo.GetEnrollmentTokenByHash(tokenHashStr)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to find enrollment token: %w", err)
	// }

	// Database integration required - enrollment token lookup needs to be implemented
	// Returning error until database integration is complete
	return nil, fmt.Errorf("enrollment token lookup requires database integration")
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
	_ = &models.Agent{
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

	// TODO: Save agent to database
	// if err := s.repo.CreateAgent(agent); err != nil {
	//     return nil, fmt.Errorf("failed to create agent: %w", err)
	// }

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

	// TODO: Update enrollment token in database
	// if err := s.repo.UpdateEnrollmentToken(enrollmentToken); err != nil {
	//     return nil, fmt.Errorf("failed to update enrollment token: %w", err)
	// }

	// TODO: Save agent credential to database
	// if err := s.repo.CreateAgentCredential(credential); err != nil {
	//     return nil, fmt.Errorf("failed to save agent credential: %w", err)
	// }

	return &models.AgentEnrollmentResponse{
		AgentID:        agentID,
		OrganizationID: enrollmentToken.OrganizationID,
		Credential:     credential.CredentialHash, // In real implementation, this would be the actual credential
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
	// Hash the credential (for future database lookup)
	_ = sha256.Sum256([]byte(credential))

	// TODO: Look up credential in database
	// agentCredential, err := s.repo.GetAgentCredentialByHash(credentialHashStr)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to find agent credential: %w", err)
	// }

	// Database integration required - agent credential lookup needs to be implemented
	// Returning error until database integration is complete
	return nil, fmt.Errorf("agent credential lookup requires database integration")
}

// RevokeEnrollmentToken revokes an enrollment token
func (s *EnrollmentService) RevokeEnrollmentToken(tokenID uuid.UUID) error {
	// TODO: Look up token in database and mark as revoked
	// enrollmentToken, err := s.repo.GetEnrollmentTokenByID(tokenID)
	// if err != nil {
	//     return fmt.Errorf("failed to find enrollment token: %w", err)
	// }

	// enrollmentToken.Status = "revoked"
	// enrollmentToken.UpdatedAt = time.Now()

	// if err := s.repo.UpdateEnrollmentToken(enrollmentToken); err != nil {
	//     return fmt.Errorf("failed to revoke enrollment token: %w", err)
	// }

	return nil
}

// RevokeAgentCredential revokes an agent's credential
func (s *EnrollmentService) RevokeAgentCredential(credentialID uuid.UUID) error {
	// TODO: Look up credential in database and mark as revoked
	// agentCredential, err := s.repo.GetAgentCredentialByID(credentialID)
	// if err != nil {
	//     return fmt.Errorf("failed to find agent credential: %w", err)
	// }

	// agentCredential.Status = "revoked"
	// agentCredential.UpdatedAt = time.Now()

	// if err := s.repo.UpdateAgentCredential(agentCredential); err != nil {
	//     return fmt.Errorf("failed to revoke agent credential: %w", err)
	// }

	return nil
}

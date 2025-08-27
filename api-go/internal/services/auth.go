package services

import (
	"errors"
	"time"

	"zerotrace/api/internal/config"
	"zerotrace/api/internal/models"
	"zerotrace/api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	config   *config.Config
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(cfg *config.Config, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		config:   cfg,
		userRepo: userRepo,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// TODO: Implement actual user lookup from database
	// For now, return a mock user for testing

	if email == "admin@zerotrace.com" && password == "password" {
		user := &models.User{
			ID:        uuid.New(),
			Email:     email,
			Name:      "Admin User",
			Role:      models.RoleAdmin,
			CompanyID: uuid.New(),
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		token, err := s.generateToken(user)
		if err != nil {
			return nil, "", err
		}

		return user, token, nil
	}

	return nil, "", errors.New("invalid credentials")
}

// Register creates a new user account
func (s *AuthService) Register(req models.RegisterRequest) (*models.User, error) {
	// TODO: Implement actual user creation with database
	// For now, return a mock user for testing

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Name:      req.Name,
		Role:      models.RoleUser,
		CompanyID: req.CompanyID,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"company_id": user.CompanyID.String(),
		"role":       string(user.Role),
		"exp":        time.Now().Add(s.config.JWTExpiry).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

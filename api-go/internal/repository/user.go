package repository

import (
	"time"

	"zerotrace/api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByCompanyID retrieves users by company ID
func (r *UserRepository) GetByCompanyID(companyID uuid.UUID, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Get total count
	err := r.db.Model(&models.User{}).Where("company_id = ?", companyID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err = r.db.Where("company_id = ?", companyID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// UpdateLastLogin updates the last login time
func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("last_login", time.Now()).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

// ExistsByEmail checks if a user exists by email
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

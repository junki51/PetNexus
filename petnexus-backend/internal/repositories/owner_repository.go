package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

var (
	ErrOwnerProfileNotFound      = errors.New("owner profile not found")
	ErrOwnerProfileAlreadyExists = errors.New("owner profile already exists")
)

// OwnerProfileRepository defines database-only owner profile operations.
type OwnerProfileRepository interface {
	Create(profile *models.OwnerProfile) error
	FindByUserID(userID uuid.UUID) (*models.OwnerProfile, error)
	ExistsByUserID(userID uuid.UUID) (bool, error)
	Update(profile *models.OwnerProfile) error
}

type ownerProfileRepository struct {
	db *gorm.DB
}

// NewOwnerProfileRepository creates a GORM-backed owner profile repository.
func NewOwnerProfileRepository(db *gorm.DB) OwnerProfileRepository {
	return &ownerProfileRepository{db: db}
}

func (r *ownerProfileRepository) Create(profile *models.OwnerProfile) error {
	if err := r.db.Create(profile).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrOwnerProfileAlreadyExists
		}
		return fmt.Errorf("create owner profile: %w", err)
	}
	return nil
}

func (r *ownerProfileRepository) FindByUserID(userID uuid.UUID) (*models.OwnerProfile, error) {
	var profile models.OwnerProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOwnerProfileNotFound
		}
		return nil, fmt.Errorf("find owner profile by user ID: %w", err)
	}
	return &profile, nil
}

func (r *ownerProfileRepository) ExistsByUserID(userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.OwnerProfile{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check owner profile existence: %w", err)
	}
	return count > 0, nil
}

func (r *ownerProfileRepository) Update(profile *models.OwnerProfile) error {
	result := r.db.Model(&models.OwnerProfile{}).
		Where("id = ?", profile.ID).
		Select(
			"first_name",
			"last_name",
			"gender",
			"date_of_birth",
			"phone_number",
			"avatar_url",
			"address_line1",
			"address_line2",
			"province",
			"district",
			"subdistrict",
			"postal_code",
			"updated_at",
		).
		Updates(profile)
	if result.Error != nil {
		return fmt.Errorf("update owner profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrOwnerProfileNotFound
	}
	return nil
}

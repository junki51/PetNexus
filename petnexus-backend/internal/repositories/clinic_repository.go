package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

var (
	ErrClinicProfileNotFound      = errors.New("clinic profile not found")
	ErrClinicProfileAlreadyExists = errors.New("clinic profile already exists")
)

// ClinicProfileRepository defines database-only clinic profile operations.
type ClinicProfileRepository interface {
	Create(profile *models.ClinicProfile) error
	FindByID(id uuid.UUID) (*models.ClinicProfile, error)
	FindByUserID(userID uuid.UUID) (*models.ClinicProfile, error)
	ExistsByUserID(userID uuid.UUID) (bool, error)
	Update(profile *models.ClinicProfile) error
}

func (r *clinicProfileRepository) FindByID(id uuid.UUID) (*models.ClinicProfile, error) {
	var profile models.ClinicProfile
	if err := r.db.Where("id = ?", id).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClinicProfileNotFound
		}
		return nil, fmt.Errorf("find clinic profile by ID: %w", err)
	}
	return &profile, nil
}

type clinicProfileRepository struct {
	db *gorm.DB
}

func NewClinicProfileRepository(db *gorm.DB) ClinicProfileRepository {
	return &clinicProfileRepository{db: db}
}

func (r *clinicProfileRepository) Create(profile *models.ClinicProfile) error {
	if err := r.db.Create(profile).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrClinicProfileAlreadyExists
		}
		return fmt.Errorf("create clinic profile: %w", err)
	}
	return nil
}

func (r *clinicProfileRepository) FindByUserID(userID uuid.UUID) (*models.ClinicProfile, error) {
	var profile models.ClinicProfile
	if err := r.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClinicProfileNotFound
		}
		return nil, fmt.Errorf("find clinic profile by user ID: %w", err)
	}
	return &profile, nil
}

func (r *clinicProfileRepository) ExistsByUserID(userID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.ClinicProfile{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return false, fmt.Errorf("check clinic profile existence: %w", err)
	}
	return count > 0, nil
}

func (r *clinicProfileRepository) Update(profile *models.ClinicProfile) error {
	result := r.db.Model(&models.ClinicProfile{}).
		Where("id = ?", profile.ID).
		Select("clinic_name", "phone_number", "email", "address", "updated_at").
		Updates(profile)
	if result.Error != nil {
		return fmt.Errorf("update clinic profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrClinicProfileNotFound
	}
	return nil
}

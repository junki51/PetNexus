package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

var ErrBreedNotFound = errors.New("breed not found")

// BreedRepository defines database-only breed reference operations.
type BreedRepository interface {
	FindAll(species string) ([]models.Breed, error)
	FindByID(id uuid.UUID) (*models.Breed, error)
}

type breedRepository struct {
	db *gorm.DB
}

func NewBreedRepository(db *gorm.DB) BreedRepository {
	return &breedRepository{db: db}
}

func (r *breedRepository) FindAll(species string) ([]models.Breed, error) {
	breeds := make([]models.Breed, 0)
	query := r.db.Order("species ASC, name ASC")
	if species != "" {
		query = query.Where("species = ?", species)
	}
	if err := query.Find(&breeds).Error; err != nil {
		return nil, fmt.Errorf("find breeds: %w", err)
	}
	return breeds, nil
}

func (r *breedRepository) FindByID(id uuid.UUID) (*models.Breed, error) {
	var breed models.Breed
	if err := r.db.First(&breed, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBreedNotFound
		}
		return nil, fmt.Errorf("find breed by ID: %w", err)
	}
	return &breed, nil
}

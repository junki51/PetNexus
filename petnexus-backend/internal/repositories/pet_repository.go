package repositories

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

var ErrPetNotFound = errors.New("pet not found")

// PetRepository defines database-only pet operations.
type PetRepository interface {
	Create(pet *models.Pet) error
	FindByID(id uuid.UUID) (*models.Pet, error)
	FindByIDAndOwnerProfileID(id, ownerProfileID uuid.UUID) (*models.Pet, error)
	FindAllByOwnerProfileID(ownerProfileID uuid.UUID) ([]models.Pet, error)
	Update(pet *models.Pet) error
}

type petRepository struct {
	db *gorm.DB
}

func NewPetRepository(db *gorm.DB) PetRepository {
	return &petRepository{db: db}
}

func (r *petRepository) Create(pet *models.Pet) error {
	if err := r.db.Omit("Breed").Create(pet).Error; err != nil {
		return fmt.Errorf("create pet: %w", err)
	}
	return nil
}

func (r *petRepository) FindByID(id uuid.UUID) (*models.Pet, error) {
	var pet models.Pet
	if err := r.db.Preload("Breed").First(&pet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPetNotFound
		}
		return nil, fmt.Errorf("find pet by ID: %w", err)
	}
	return &pet, nil
}

func (r *petRepository) FindByIDAndOwnerProfileID(id, ownerProfileID uuid.UUID) (*models.Pet, error) {
	var pet models.Pet
	if err := r.db.Preload("Breed").
		Where("id = ? AND owner_profile_id = ?", id, ownerProfileID).
		First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPetNotFound
		}
		return nil, fmt.Errorf("find pet by ID and owner profile: %w", err)
	}
	return &pet, nil
}

func (r *petRepository) FindAllByOwnerProfileID(ownerProfileID uuid.UUID) ([]models.Pet, error) {
	pets := make([]models.Pet, 0)
	if err := r.db.Preload("Breed").
		Where("owner_profile_id = ?", ownerProfileID).
		Order("created_at DESC").
		Find(&pets).Error; err != nil {
		return nil, fmt.Errorf("find pets by owner profile: %w", err)
	}
	return pets, nil
}

func (r *petRepository) Update(pet *models.Pet) error {
	result := r.db.Model(&models.Pet{}).
		Where("id = ? AND owner_profile_id = ?", pet.ID, pet.OwnerProfileID).
		Select(
			"breed_id",
			"species",
			"name",
			"gender",
			"date_of_birth",
			"weight_kg",
			"microchip_id",
			"avatar_url",
			"color",
			"distinctive_marks",
			"is_neutered",
			"updated_at",
		).
		Updates(pet)
	if result.Error != nil {
		return fmt.Errorf("update pet: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrPetNotFound
	}
	return nil
}

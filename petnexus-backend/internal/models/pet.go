package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	PetGenderMale    = "male"
	PetGenderFemale  = "female"
	PetGenderUnknown = "unknown"
)

// Pet is owner-controlled basic pet identity data. Passport and clinic data
// intentionally belong to later features.
type Pet struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	PublicPetID      string     `gorm:"column:public_pet_id;size:50;not null;uniqueIndex:idx_pets_public_pet_id_unique"`
	OwnerProfileID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	BreedID          *uuid.UUID `gorm:"type:uuid;index"`
	Species          string     `gorm:"size:20;not null;index"`
	Name             string     `gorm:"size:100;not null"`
	Gender           *string    `gorm:"size:30"`
	DateOfBirth      *time.Time `gorm:"type:date"`
	WeightKG         *float64   `gorm:"column:weight_kg;type:numeric(6,2)"`
	MicrochipID      *string    `gorm:"size:100"`
	AvatarURL        *string    `gorm:"type:text"`
	Color            *string    `gorm:"size:100"`
	DistinctiveMarks *string    `gorm:"type:text"`
	IsNeutered       *bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Breed            *Breed        `gorm:"foreignKey:BreedID"`
	OwnerProfile     *OwnerProfile `gorm:"foreignKey:OwnerProfileID"`
}

func (Pet) TableName() string {
	return "pets"
}

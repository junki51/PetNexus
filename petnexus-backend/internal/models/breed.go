package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	SpeciesDog = "dog"
	SpeciesCat = "cat"
)

// Breed is a selectable dog or cat breed reference.
type Breed struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Species   string    `gorm:"size:20;not null;uniqueIndex:idx_breeds_species_name_unique"`
	Name      string    `gorm:"size:100;not null;uniqueIndex:idx_breeds_species_name_unique"`
	NameTH    *string   `gorm:"column:name_th;size:100"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Breed) TableName() string {
	return "breeds"
}

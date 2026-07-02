package models

import (
	"time"

	"github.com/google/uuid"
)

// OwnerProfile contains owner identity data separately from authentication.
// It must not be serialized directly as an API response.
type OwnerProfile struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_owner_profiles_user_id_unique"`
	FirstName    string     `gorm:"size:100;not null"`
	LastName     string     `gorm:"size:100;not null"`
	Gender       *string    `gorm:"size:30"`
	DateOfBirth  *time.Time `gorm:"type:date"`
	PhoneNumber  string     `gorm:"size:30;not null"`
	AvatarURL    *string    `gorm:"type:text"`
	AddressLine1 *string    `gorm:"size:255"`
	AddressLine2 *string    `gorm:"size:255"`
	Province     *string    `gorm:"size:100"`
	District     *string    `gorm:"size:100"`
	Subdistrict  *string    `gorm:"size:100"`
	PostalCode   *string    `gorm:"size:20"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName keeps the database name explicit and stable.
func (OwnerProfile) TableName() string {
	return "owner_profiles"
}

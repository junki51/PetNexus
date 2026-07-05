package models

import (
	"time"

	"github.com/google/uuid"
)

// ClinicProfile contains settings/identity data for one clinic staff account.
// It intentionally excludes access, visit, medical, and staff-member data.
type ClinicProfile struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_clinic_profiles_user_id_unique"`
	ClinicName  string    `gorm:"size:200;not null"`
	PhoneNumber *string   `gorm:"size:30"`
	Email       *string   `gorm:"size:255"`
	Address     *string   `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ClinicProfile) TableName() string {
	return "clinic_profiles"
}

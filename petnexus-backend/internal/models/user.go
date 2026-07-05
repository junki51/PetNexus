package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleOwner       = "owner"
	RoleClinic      = "clinic"
	RoleClinicStaff = "clinic_staff"
	RoleAdmin       = "admin"
)

// User represents a PetNexus login account. It intentionally contains only a
// password hash and must never be serialized directly as an API response.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	Phone        string    `gorm:"size:30"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"type:user_role;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

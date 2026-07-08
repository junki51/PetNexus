package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	AppointmentTypeCheckup      = "checkup"
	AppointmentTypeVaccination  = "vaccination"
	AppointmentTypeConsultation = "consultation"
	AppointmentTypeFollowUp     = "follow_up"
	AppointmentTypeGrooming     = "grooming"
	AppointmentTypeEmergency    = "emergency"
	AppointmentTypeOther        = "other"

	AppointmentStatusRequested = "requested"
	AppointmentStatusScheduled = "scheduled"
	AppointmentStatusCheckedIn = "checked_in"
	AppointmentStatusCompleted = "completed"
	AppointmentStatusCancelled = "cancelled"
)

// Appointment is a scheduled relationship between one owner, pet, and clinic.
// It intentionally contains no medical record, payment, or staff schedule data.
type Appointment struct {
	ID              uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OwnerProfileID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	ClinicProfileID uuid.UUID  `gorm:"type:uuid;not null;index"`
	PetID           uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title           *string    `gorm:"size:150"`
	AppointmentType string     `gorm:"size:50;not null;index"`
	ScheduledAt     time.Time  `gorm:"type:timestamptz;not null;index"`
	DurationMinutes int        `gorm:"not null"`
	Status          string     `gorm:"size:50;not null;index"`
	Note            *string    `gorm:"type:text"`
	CreatedByUserID *uuid.UUID `gorm:"type:uuid"`
	CreatedByRole   string     `gorm:"size:20;not null"`
	CancelledAt     *time.Time `gorm:"type:timestamptz"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Pet             *Pet           `gorm:"foreignKey:PetID"`
	OwnerProfile    *OwnerProfile  `gorm:"foreignKey:OwnerProfileID"`
	ClinicProfile   *ClinicProfile `gorm:"foreignKey:ClinicProfileID"`
}

func (Appointment) TableName() string {
	return "appointments"
}

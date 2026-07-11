package models

import (
	"time"

	"github.com/google/uuid"
)

// MedicalRecord contains clinic-owned clinical notes for one patient visit.
// It intentionally does not include files, lab results, vaccination records,
// prescriptions tables, or version history.
type MedicalRecord struct {
	ID                   uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ClinicProfileID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	PetID                uuid.UUID  `gorm:"type:uuid;not null;index"`
	AppointmentID        *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_medical_records_appointment_id_unique"`
	CreatedByUserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	VisitAt              time.Time  `gorm:"type:timestamptz;not null;index"`
	ChiefComplaint       string     `gorm:"type:text;not null"`
	ClinicalFindings     *string    `gorm:"type:text"`
	Diagnosis            *string    `gorm:"type:text"`
	TreatmentPlan        *string    `gorm:"type:text"`
	Medications          *string    `gorm:"type:text"`
	FollowUpInstructions *string    `gorm:"type:text"`
	NextFollowUpAt       *time.Time `gorm:"type:timestamptz"`
	WeightKG             *float64   `gorm:"column:weight_kg;type:numeric(6,2)"`
	TemperatureC         *float64   `gorm:"column:temperature_c;type:numeric(5,2)"`
	Notes                *string    `gorm:"type:text"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Pet                  *Pet           `gorm:"foreignKey:PetID"`
	ClinicProfile        *ClinicProfile `gorm:"foreignKey:ClinicProfileID"`
	Appointment          *Appointment   `gorm:"foreignKey:AppointmentID"`
	CreatedByUser        *User          `gorm:"foreignKey:CreatedByUserID"`
}

func (MedicalRecord) TableName() string {
	return "medical_records"
}

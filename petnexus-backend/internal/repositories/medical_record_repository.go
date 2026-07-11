package repositories

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

var (
	ErrMedicalRecordNotFound                 = errors.New("medical record not found")
	ErrMedicalRecordAppointmentAlreadyLinked = errors.New("appointment already has a medical record")
)

// MedicalRecordFilters contains normalized repository-level filters.
// DateTo is exclusive so an inclusive UI date can be represented safely.
type MedicalRecordFilters struct {
	PetID    *uuid.UUID
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    int
	Offset   int
}

type MedicalRecordRepository interface {
	Create(record *models.MedicalRecord) error
	FindAllByClinicProfileID(clinicProfileID uuid.UUID, filters MedicalRecordFilters) ([]models.MedicalRecord, error)
	CountByClinicProfileID(clinicProfileID uuid.UUID, filters MedicalRecordFilters) (int64, error)
	FindByIDAndClinicProfileID(id, clinicProfileID uuid.UUID) (*models.MedicalRecord, error)
	Update(record *models.MedicalRecord) error
	PetHasNonCancelledAppointmentWithClinic(clinicProfileID, petID uuid.UUID) (bool, error)
	FindUsableAppointmentForMedicalRecord(appointmentID, clinicProfileID, petID uuid.UUID) (*models.Appointment, error)
	MedicalRecordExistsByAppointmentID(appointmentID uuid.UUID) (bool, error)
}

type medicalRecordRepository struct {
	db *gorm.DB
}

func NewMedicalRecordRepository(db *gorm.DB) MedicalRecordRepository {
	return &medicalRecordRepository{db: db}
}

func (r *medicalRecordRepository) Create(record *models.MedicalRecord) error {
	if err := r.db.Omit("Pet", "ClinicProfile", "Appointment", "CreatedByUser").Create(record).Error; err != nil {
		if isMedicalRecordAppointmentDuplicate(err) {
			return ErrMedicalRecordAppointmentAlreadyLinked
		}
		return fmt.Errorf("create medical record: %w", err)
	}
	return nil
}

func (r *medicalRecordRepository) FindAllByClinicProfileID(clinicProfileID uuid.UUID, filters MedicalRecordFilters) ([]models.MedicalRecord, error) {
	records := make([]models.MedicalRecord, 0)
	query := withMedicalRecordRelations(r.db).Where("clinic_profile_id = ?", clinicProfileID)
	query = applyMedicalRecordFilters(query, filters)
	if err := query.
		Order("visit_at DESC, created_at DESC, id DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&records).Error; err != nil {
		return nil, fmt.Errorf("find medical records: %w", err)
	}
	return records, nil
}

func (r *medicalRecordRepository) CountByClinicProfileID(clinicProfileID uuid.UUID, filters MedicalRecordFilters) (int64, error) {
	var count int64
	query := r.db.Model(&models.MedicalRecord{}).Where("clinic_profile_id = ?", clinicProfileID)
	query = applyMedicalRecordFilters(query, filters)
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count medical records: %w", err)
	}
	return count, nil
}

func (r *medicalRecordRepository) FindByIDAndClinicProfileID(id, clinicProfileID uuid.UUID) (*models.MedicalRecord, error) {
	var record models.MedicalRecord
	if err := withMedicalRecordRelations(r.db).
		Where("id = ? AND clinic_profile_id = ?", id, clinicProfileID).
		First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMedicalRecordNotFound
		}
		return nil, fmt.Errorf("find medical record: %w", err)
	}
	return &record, nil
}

func (r *medicalRecordRepository) Update(record *models.MedicalRecord) error {
	result := r.db.Model(&models.MedicalRecord{}).
		Where("id = ? AND clinic_profile_id = ?", record.ID, record.ClinicProfileID).
		Updates(map[string]any{
			"visit_at":               record.VisitAt,
			"chief_complaint":        record.ChiefComplaint,
			"clinical_findings":      record.ClinicalFindings,
			"diagnosis":              record.Diagnosis,
			"treatment_plan":         record.TreatmentPlan,
			"medications":            record.Medications,
			"follow_up_instructions": record.FollowUpInstructions,
			"next_follow_up_at":      record.NextFollowUpAt,
			"weight_kg":              record.WeightKG,
			"temperature_c":          record.TemperatureC,
			"notes":                  record.Notes,
			"updated_at":             record.UpdatedAt,
		})
	if result.Error != nil {
		return fmt.Errorf("update medical record: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrMedicalRecordNotFound
	}
	return nil
}

func (r *medicalRecordRepository) PetHasNonCancelledAppointmentWithClinic(clinicProfileID, petID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Appointment{}).
		Where("clinic_profile_id = ? AND pet_id = ? AND status <> ?", clinicProfileID, petID, models.AppointmentStatusCancelled).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check clinic patient relationship: %w", err)
	}
	return count > 0, nil
}

func (r *medicalRecordRepository) FindUsableAppointmentForMedicalRecord(appointmentID, clinicProfileID, petID uuid.UUID) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := r.db.
		Where(
			"id = ? AND clinic_profile_id = ? AND pet_id = ? AND status <> ?",
			appointmentID,
			clinicProfileID,
			petID,
			models.AppointmentStatusCancelled,
		).
		First(&appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppointmentNotFound
		}
		return nil, fmt.Errorf("find medical record appointment: %w", err)
	}
	return &appointment, nil
}

func (r *medicalRecordRepository) MedicalRecordExistsByAppointmentID(appointmentID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&models.MedicalRecord{}).
		Where("appointment_id = ?", appointmentID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check appointment medical record existence: %w", err)
	}
	return count > 0, nil
}

func withMedicalRecordRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Pet.Breed").
		Preload("Pet.OwnerProfile").
		Preload("Appointment").
		Preload("CreatedByUser")
}

func applyMedicalRecordFilters(query *gorm.DB, filters MedicalRecordFilters) *gorm.DB {
	if filters.PetID != nil {
		query = query.Where("pet_id = ?", *filters.PetID)
	}
	if filters.DateFrom != nil {
		query = query.Where("visit_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("visit_at < ?", *filters.DateTo)
	}
	return query
}

func isMedicalRecordAppointmentDuplicate(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey) ||
		strings.Contains(err.Error(), "idx_medical_records_appointment_id_unique")
}

package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/models"
)

var ErrAppointmentNotFound = errors.New("appointment not found")

// AppointmentFilters contains normalized repository-level calendar filters.
// DateTo is exclusive so an entire final calendar day can be represented.
type AppointmentFilters struct {
	DateFrom        *time.Time
	DateTo          *time.Time
	Status          string
	AppointmentType string
}

type AppointmentRepository interface {
	Create(appointment *models.Appointment) error
	FindByID(id uuid.UUID) (*models.Appointment, error)
	FindByIDAndOwnerProfileID(id, ownerProfileID uuid.UUID) (*models.Appointment, error)
	FindByIDAndClinicProfileID(id, clinicProfileID uuid.UUID) (*models.Appointment, error)
	FindAllByOwnerProfileID(ownerProfileID uuid.UUID, filters AppointmentFilters) ([]models.Appointment, error)
	FindAllByClinicProfileID(clinicProfileID uuid.UUID, filters AppointmentFilters) ([]models.Appointment, error)
	Update(appointment *models.Appointment) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db: db}
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	if err := r.db.Omit("Pet", "OwnerProfile", "ClinicProfile").Create(appointment).Error; err != nil {
		return fmt.Errorf("create appointment: %w", err)
	}
	return nil
}

func (r *appointmentRepository) FindByID(id uuid.UUID) (*models.Appointment, error) {
	return r.findOne("id = ?", id)
}

func (r *appointmentRepository) FindByIDAndOwnerProfileID(id, ownerProfileID uuid.UUID) (*models.Appointment, error) {
	return r.findOne("id = ? AND owner_profile_id = ?", id, ownerProfileID)
}

func (r *appointmentRepository) FindByIDAndClinicProfileID(id, clinicProfileID uuid.UUID) (*models.Appointment, error) {
	return r.findOne("id = ? AND clinic_profile_id = ?", id, clinicProfileID)
}

func (r *appointmentRepository) FindAllByOwnerProfileID(ownerProfileID uuid.UUID, filters AppointmentFilters) ([]models.Appointment, error) {
	return r.findAll("owner_profile_id = ?", ownerProfileID, filters)
}

func (r *appointmentRepository) FindAllByClinicProfileID(clinicProfileID uuid.UUID, filters AppointmentFilters) ([]models.Appointment, error) {
	return r.findAll("clinic_profile_id = ?", clinicProfileID, filters)
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
	result := r.db.Model(&models.Appointment{}).
		Where("id = ?", appointment.ID).
		Select("status", "cancelled_at", "updated_at").
		Updates(appointment)
	if result.Error != nil {
		return fmt.Errorf("update appointment: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrAppointmentNotFound
	}
	return nil
}

func (r *appointmentRepository) findOne(query string, args ...any) (*models.Appointment, error) {
	var appointment models.Appointment
	if err := withAppointmentRelations(r.db).
		Where(query, args...).
		First(&appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAppointmentNotFound
		}
		return nil, fmt.Errorf("find appointment: %w", err)
	}
	return &appointment, nil
}

func (r *appointmentRepository) findAll(scope string, scopeID uuid.UUID, filters AppointmentFilters) ([]models.Appointment, error) {
	appointments := make([]models.Appointment, 0)
	query := withAppointmentRelations(r.db).Where(scope, scopeID)
	if filters.DateFrom != nil {
		query = query.Where("scheduled_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("scheduled_at < ?", *filters.DateTo)
	}
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.AppointmentType != "" {
		query = query.Where("appointment_type = ?", filters.AppointmentType)
	}
	if err := query.Order("scheduled_at ASC").Find(&appointments).Error; err != nil {
		return nil, fmt.Errorf("find appointments: %w", err)
	}
	return appointments, nil
}

func withAppointmentRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Pet.Breed").
		Preload("OwnerProfile").
		Preload("ClinicProfile")
}

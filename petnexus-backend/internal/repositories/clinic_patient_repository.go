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

var ErrClinicPatientNotFound = errors.New("clinic patient not found")

// ClinicPatientFilters contains normalized database filters for the Clinic Web
// Patients page. Limit and Offset are validated in the service layer.
type ClinicPatientFilters struct {
	Query   string
	Species string
	Status  string
	Limit   int
	Offset  int
	Sort    string
}

type ClinicPatientSummary struct {
	TotalAppointments  int64
	FirstAppointmentAt *time.Time
	LastAppointmentAt  *time.Time
	NextAppointmentAt  *time.Time
	LatestStatus       string
}

type ClinicPatientRecord struct {
	Pet     *models.Pet
	Summary ClinicPatientSummary
}

type ClinicPatientDetail struct {
	Pet                *models.Pet
	Summary            ClinicPatientSummary
	RecentAppointments []models.Appointment
}

type ClinicPatientRepository interface {
	FindPatientsByClinicProfileID(clinicProfileID uuid.UUID, filters ClinicPatientFilters) ([]ClinicPatientRecord, error)
	FindPatientDetailByClinicProfileIDAndPetID(clinicProfileID, petID uuid.UUID) (*ClinicPatientDetail, error)
}

type clinicPatientRepository struct {
	db *gorm.DB
}

func NewClinicPatientRepository(db *gorm.DB) ClinicPatientRepository {
	return &clinicPatientRepository{db: db}
}

func (r *clinicPatientRepository) FindPatientsByClinicProfileID(clinicProfileID uuid.UUID, filters ClinicPatientFilters) ([]ClinicPatientRecord, error) {
	rows := make([]clinicPatientSummaryRow, 0)
	query, args := buildClinicPatientListQuery(clinicProfileID, filters)
	if err := r.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find clinic patients: %w", err)
	}
	if len(rows) == 0 {
		return []ClinicPatientRecord{}, nil
	}

	petIDs := make([]uuid.UUID, 0, len(rows))
	for _, row := range rows {
		petIDs = append(petIDs, row.PetID)
	}
	petsByID, err := r.loadPatientPets(petIDs)
	if err != nil {
		return nil, err
	}

	records := make([]ClinicPatientRecord, 0, len(rows))
	for _, row := range rows {
		pet := petsByID[row.PetID]
		if pet == nil {
			continue
		}
		records = append(records, ClinicPatientRecord{
			Pet:     pet,
			Summary: row.toSummary(),
		})
	}
	return records, nil
}

func (r *clinicPatientRepository) FindPatientDetailByClinicProfileIDAndPetID(clinicProfileID, petID uuid.UUID) (*ClinicPatientDetail, error) {
	rows := make([]clinicPatientSummaryRow, 0, 1)
	if err := r.db.Raw(`
		WITH patient_summary AS (
			SELECT
				a.pet_id,
				COUNT(*) AS total_appointments,
				MIN(a.scheduled_at) AS first_appointment_at,
				MAX(a.scheduled_at) AS last_appointment_at,
				MIN(a.scheduled_at) FILTER (WHERE a.scheduled_at >= NOW()) AS next_appointment_at,
				(ARRAY_AGG(a.status ORDER BY a.scheduled_at DESC, a.created_at DESC, a.id DESC))[1] AS latest_status
			FROM appointments a
			WHERE a.clinic_profile_id = ? AND a.pet_id = ? AND a.status <> ?
			GROUP BY a.pet_id
		)
		SELECT
			pet_id,
			total_appointments,
			first_appointment_at,
			last_appointment_at,
			next_appointment_at,
			latest_status
		FROM patient_summary
	`, clinicProfileID, petID, models.AppointmentStatusCancelled).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("find clinic patient detail summary: %w", err)
	}
	if len(rows) == 0 {
		return nil, ErrClinicPatientNotFound
	}

	var pet models.Pet
	if err := r.db.
		Preload("Breed").
		Preload("OwnerProfile").
		First(&pet, "id = ?", petID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClinicPatientNotFound
		}
		return nil, fmt.Errorf("find clinic patient pet: %w", err)
	}

	appointments := make([]models.Appointment, 0)
	if err := r.db.
		Where("clinic_profile_id = ? AND pet_id = ? AND status <> ?", clinicProfileID, petID, models.AppointmentStatusCancelled).
		Order("scheduled_at DESC, created_at DESC, id DESC").
		Limit(5).
		Find(&appointments).Error; err != nil {
		return nil, fmt.Errorf("find clinic patient recent appointments: %w", err)
	}

	return &ClinicPatientDetail{
		Pet:                &pet,
		Summary:            rows[0].toSummary(),
		RecentAppointments: appointments,
	}, nil
}

func (r *clinicPatientRepository) loadPatientPets(petIDs []uuid.UUID) (map[uuid.UUID]*models.Pet, error) {
	pets := make([]models.Pet, 0, len(petIDs))
	if err := r.db.
		Preload("Breed").
		Preload("OwnerProfile").
		Where("id IN ?", petIDs).
		Find(&pets).Error; err != nil {
		return nil, fmt.Errorf("load clinic patient pets: %w", err)
	}
	result := make(map[uuid.UUID]*models.Pet, len(pets))
	for i := range pets {
		result[pets[i].ID] = &pets[i]
	}
	return result, nil
}

type clinicPatientSummaryRow struct {
	PetID              uuid.UUID  `gorm:"column:pet_id"`
	TotalAppointments  int64      `gorm:"column:total_appointments"`
	FirstAppointmentAt *time.Time `gorm:"column:first_appointment_at"`
	LastAppointmentAt  *time.Time `gorm:"column:last_appointment_at"`
	NextAppointmentAt  *time.Time `gorm:"column:next_appointment_at"`
	LatestStatus       string     `gorm:"column:latest_status"`
}

func (r clinicPatientSummaryRow) toSummary() ClinicPatientSummary {
	return ClinicPatientSummary{
		TotalAppointments:  r.TotalAppointments,
		FirstAppointmentAt: r.FirstAppointmentAt,
		LastAppointmentAt:  r.LastAppointmentAt,
		NextAppointmentAt:  r.NextAppointmentAt,
		LatestStatus:       r.LatestStatus,
	}
}

func buildClinicPatientListQuery(clinicProfileID uuid.UUID, filters ClinicPatientFilters) (string, []any) {
	var builder strings.Builder
	builder.WriteString(`
		WITH patient_summary AS (
			SELECT
				a.clinic_profile_id,
				a.pet_id,
				COUNT(*) AS total_appointments,
				MIN(a.scheduled_at) AS first_appointment_at,
				MAX(a.scheduled_at) AS last_appointment_at,
				MIN(a.scheduled_at) FILTER (WHERE a.scheduled_at >= NOW()) AS next_appointment_at,
				(ARRAY_AGG(a.status ORDER BY a.scheduled_at DESC, a.created_at DESC, a.id DESC))[1] AS latest_status
			FROM appointments a
			WHERE a.clinic_profile_id = ? AND a.status <> ?
			GROUP BY a.clinic_profile_id, a.pet_id
		)
		SELECT
			ps.pet_id,
			ps.total_appointments,
			ps.first_appointment_at,
			ps.last_appointment_at,
			ps.next_appointment_at,
			ps.latest_status
		FROM patient_summary ps
		JOIN pets p ON p.id = ps.pet_id
		WHERE 1 = 1
	`)
	args := []any{clinicProfileID, models.AppointmentStatusCancelled}

	if filters.Query != "" {
		pattern := "%" + strings.ToLower(filters.Query) + "%"
		builder.WriteString(" AND (LOWER(p.name) LIKE ? OR LOWER(p.public_pet_id) LIKE ?)")
		args = append(args, pattern, pattern)
	}
	if filters.Species != "" {
		builder.WriteString(" AND p.species = ?")
		args = append(args, filters.Species)
	}
	if filters.Status != "" {
		builder.WriteString(" AND ps.latest_status = ?")
		args = append(args, filters.Status)
	}

	builder.WriteString(" ORDER BY ")
	builder.WriteString(clinicPatientOrderClause(filters.Sort))
	builder.WriteString(" LIMIT ? OFFSET ?")
	args = append(args, filters.Limit, filters.Offset)

	return builder.String(), args
}

func clinicPatientOrderClause(sort string) string {
	switch sort {
	case "latest_appointment_asc":
		return "ps.last_appointment_at ASC, p.name ASC"
	case "name_asc":
		return "p.name ASC, ps.last_appointment_at DESC"
	case "name_desc":
		return "p.name DESC, ps.last_appointment_at DESC"
	case "next_appointment_asc":
		return "ps.next_appointment_at ASC NULLS LAST, ps.last_appointment_at DESC"
	default:
		return "ps.last_appointment_at DESC, p.name ASC"
	}
}

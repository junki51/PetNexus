package services

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

const (
	appointmentDateLayout = "2006-01-02"
	maxAppointmentTitle   = 150
	maxAppointmentNote    = 1000
	minAppointmentMinutes = 5
	maxAppointmentMinutes = 480
)

var allowedAppointmentTypes = map[string]struct{}{
	models.AppointmentTypeCheckup:      {},
	models.AppointmentTypeVaccination:  {},
	models.AppointmentTypeConsultation: {},
	models.AppointmentTypeFollowUp:     {},
	models.AppointmentTypeGrooming:     {},
	models.AppointmentTypeEmergency:    {},
	models.AppointmentTypeOther:        {},
}

var allowedAppointmentStatuses = map[string]struct{}{
	models.AppointmentStatusRequested: {},
	models.AppointmentStatusScheduled: {},
	models.AppointmentStatusCheckedIn: {},
	models.AppointmentStatusCompleted: {},
	models.AppointmentStatusCancelled: {},
}

type normalizedAppointmentInput struct {
	title           *string
	appointmentType string
	scheduledAt     time.Time
	durationMinutes int
	note            *string
}

func normalizeAppointmentInput(title, appointmentType, scheduledAt string, durationMinutes int, note string, now time.Time) (normalizedAppointmentInput, error) {
	normalizedTitle, err := normalizeOptionalAppointmentField("title", title, maxAppointmentTitle)
	if err != nil {
		return normalizedAppointmentInput{}, err
	}
	normalizedType, err := normalizeAppointmentType(appointmentType)
	if err != nil {
		return normalizedAppointmentInput{}, err
	}
	normalizedScheduledAt, err := parseFutureAppointmentTime(scheduledAt, now)
	if err != nil {
		return normalizedAppointmentInput{}, err
	}
	if durationMinutes < minAppointmentMinutes || durationMinutes > maxAppointmentMinutes {
		return normalizedAppointmentInput{}, appointmentValidationError(
			"duration_minutes must be between " + strconv.Itoa(minAppointmentMinutes) +
				" and " + strconv.Itoa(maxAppointmentMinutes),
		)
	}
	normalizedNote, err := normalizeOptionalAppointmentField("note", note, maxAppointmentNote)
	if err != nil {
		return normalizedAppointmentInput{}, err
	}
	return normalizedAppointmentInput{
		title:           normalizedTitle,
		appointmentType: normalizedType,
		scheduledAt:     normalizedScheduledAt,
		durationMinutes: durationMinutes,
		note:            normalizedNote,
	}, nil
}

func normalizeOptionalAppointmentField(field, value string, limit int) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if utf8.RuneCountInString(trimmed) > limit {
		return nil, appointmentValidationError(field + " must not exceed " + strconv.Itoa(limit) + " characters")
	}
	if trimmed == "" {
		return nil, nil
	}
	return &trimmed, nil
}

func normalizeAppointmentType(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", appointmentValidationError("appointment_type is required")
	}
	if _, ok := allowedAppointmentTypes[normalized]; !ok {
		return "", appointmentValidationError("appointment_type is not supported")
	}
	return normalized, nil
}

func normalizeAppointmentStatus(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	if normalized == "" {
		return "", appointmentValidationError("status is required")
	}
	if _, ok := allowedAppointmentStatuses[normalized]; !ok {
		return "", appointmentValidationError("status is not supported")
	}
	return normalized, nil
}

func parseFutureAppointmentTime(value string, now time.Time) (time.Time, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return time.Time{}, appointmentValidationError("scheduled_at is required")
	}
	parsed, err := time.Parse(time.RFC3339, text)
	if err != nil {
		return time.Time{}, appointmentValidationError("scheduled_at must use RFC3339 format")
	}
	if !parsed.After(now) {
		return time.Time{}, appointmentValidationError("scheduled_at must be in the future")
	}
	return parsed.UTC(), nil
}

func normalizeOwnerAppointmentFilters(filters dto.OwnerAppointmentFilters) (repositories.AppointmentFilters, error) {
	return normalizeAppointmentFilters("", filters.DateFrom, filters.DateTo, filters.Status, "")
}

func normalizeClinicAppointmentFilters(filters dto.ClinicAppointmentFilters) (repositories.AppointmentFilters, error) {
	return normalizeAppointmentFilters(filters.Date, filters.DateFrom, filters.DateTo, filters.Status, filters.AppointmentType)
}

func normalizeAppointmentFilters(date, dateFrom, dateTo, status, appointmentType string) (repositories.AppointmentFilters, error) {
	date = strings.TrimSpace(date)
	dateFrom = strings.TrimSpace(dateFrom)
	dateTo = strings.TrimSpace(dateTo)
	if date != "" && (dateFrom != "" || dateTo != "") {
		return repositories.AppointmentFilters{}, appointmentValidationError("date cannot be combined with date_from or date_to")
	}

	var result repositories.AppointmentFilters
	if date != "" {
		day, err := parseAppointmentDate("date", date)
		if err != nil {
			return result, err
		}
		end := day.AddDate(0, 0, 1)
		result.DateFrom = &day
		result.DateTo = &end
	} else {
		if dateFrom != "" {
			from, err := parseAppointmentDate("date_from", dateFrom)
			if err != nil {
				return result, err
			}
			result.DateFrom = &from
		}
		if dateTo != "" {
			to, err := parseAppointmentDate("date_to", dateTo)
			if err != nil {
				return result, err
			}
			end := to.AddDate(0, 0, 1)
			result.DateTo = &end
			if result.DateFrom != nil && result.DateFrom.After(to) {
				return repositories.AppointmentFilters{}, appointmentValidationError("date_from must not be after date_to")
			}
		}
	}

	if strings.TrimSpace(status) != "" {
		normalized, err := normalizeAppointmentStatus(status)
		if err != nil {
			return repositories.AppointmentFilters{}, err
		}
		result.Status = normalized
	}
	if strings.TrimSpace(appointmentType) != "" {
		normalized, err := normalizeAppointmentType(appointmentType)
		if err != nil {
			return repositories.AppointmentFilters{}, err
		}
		result.AppointmentType = normalized
	}
	return result, nil
}

func parseAppointmentDate(field, value string) (time.Time, error) {
	parsed, err := time.Parse(appointmentDateLayout, value)
	if err != nil {
		return time.Time{}, appointmentValidationError(field + " must use YYYY-MM-DD format")
	}
	return parsed.UTC(), nil
}

func parseAppointmentUserID(value string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Unauthorized",
			"Invalid authenticated user",
			err,
		)
	}
	return id, nil
}

func setAppointmentStatus(appointment *models.Appointment, status string, now time.Time) {
	appointment.Status = status
	appointment.UpdatedAt = now.UTC()
	if status == models.AppointmentStatusCancelled {
		cancelledAt := now.UTC()
		appointment.CancelledAt = &cancelledAt
	} else {
		appointment.CancelledAt = nil
	}
}

func toAppointmentResponse(appointment *models.Appointment) dto.AppointmentResponse {
	var cancelledAt *string
	if appointment.CancelledAt != nil {
		value := appointment.CancelledAt.UTC().Format(time.RFC3339)
		cancelledAt = &value
	}

	petSummary := dto.AppointmentPetSummary{}
	if appointment.Pet != nil {
		petSummary = dto.AppointmentPetSummary{
			ID:          appointment.Pet.ID.String(),
			PublicPetID: appointment.Pet.PublicPetID,
			Name:        appointment.Pet.Name,
			Species:     appointment.Pet.Species,
			AvatarURL:   appointment.Pet.AvatarURL,
		}
		if appointment.Pet.Breed != nil {
			breed := toBreedResponse(appointment.Pet.Breed)
			petSummary.Breed = &breed
		}
	}

	ownerSummary := dto.AppointmentOwnerSummary{}
	if appointment.OwnerProfile != nil {
		ownerSummary.DisplayName = strings.TrimSpace(appointment.OwnerProfile.FirstName + " " + appointment.OwnerProfile.LastName)
		ownerSummary.MaskedPhone = maskOwnerPhone(appointment.OwnerProfile.PhoneNumber)
	}

	clinicSummary := dto.AppointmentClinicSummary{}
	if appointment.ClinicProfile != nil {
		clinicSummary = dto.AppointmentClinicSummary{
			ID:          appointment.ClinicProfile.ID.String(),
			ClinicName:  appointment.ClinicProfile.ClinicName,
			PhoneNumber: appointment.ClinicProfile.PhoneNumber,
			Email:       appointment.ClinicProfile.Email,
		}
	}

	return dto.AppointmentResponse{
		ID:              appointment.ID.String(),
		Title:           appointment.Title,
		AppointmentType: appointment.AppointmentType,
		ScheduledAt:     appointment.ScheduledAt.UTC().Format(time.RFC3339),
		DurationMinutes: appointment.DurationMinutes,
		Status:          appointment.Status,
		Note:            appointment.Note,
		CreatedByRole:   appointment.CreatedByRole,
		CancelledAt:     cancelledAt,
		CreatedAt:       appointment.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:       appointment.UpdatedAt.UTC().Format(time.RFC3339),
		Pet:             petSummary,
		Owner:           ownerSummary,
		Clinic:          clinicSummary,
	}
}

func appointmentValidationError(details string) *utils.AppError {
	return utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details, nil)
}

func appointmentNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"APPOINTMENT_NOT_FOUND",
		"Appointment not found",
		"The appointment does not exist or is outside the authenticated account",
		cause,
	)
}

func appointmentPetNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, "PET_NOT_FOUND", "Pet not found", "The selected pet was not found", cause)
}

func appointmentClinicNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, "CLINIC_PROFILE_NOT_FOUND", "Clinic profile not found", "The selected clinic was not found", cause)
}

func appointmentOwnerProfileRequiredError(cause error) *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, "OWNER_PROFILE_REQUIRED", "Owner profile required", "Create an owner profile before managing appointments", cause)
}

func appointmentClinicProfileRequiredError(cause error) *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, "CLINIC_PROFILE_REQUIRED", "Clinic profile required", "Create a clinic profile before managing appointments", cause)
}

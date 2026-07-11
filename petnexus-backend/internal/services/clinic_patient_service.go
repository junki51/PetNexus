package services

import (
	"errors"
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
	defaultClinicPatientLimit = 20
	maxClinicPatientLimit     = 100
	maxClinicPatientSearch    = 100
)

var allowedClinicPatientSorts = map[string]struct{}{
	"latest_appointment_desc": {},
	"latest_appointment_asc":  {},
	"name_asc":                {},
	"name_desc":               {},
	"next_appointment_asc":    {},
}

type ClinicPatientService interface {
	ListClinicPatients(currentUserID string, filters dto.ClinicPatientFilters) ([]dto.ClinicPatientListItemResponse, error)
	GetClinicPatient(currentUserID string, petID uuid.UUID) (*dto.ClinicPatientDetailResponse, error)
}

type clinicPatientService struct {
	patientRepo repositories.ClinicPatientRepository
	clinicRepo  repositories.ClinicProfileRepository
}

func NewClinicPatientService(
	patientRepo repositories.ClinicPatientRepository,
	clinicRepo repositories.ClinicProfileRepository,
) ClinicPatientService {
	return &clinicPatientService{patientRepo: patientRepo, clinicRepo: clinicRepo}
}

func (s *clinicPatientService) ListClinicPatients(currentUserID string, filters dto.ClinicPatientFilters) ([]dto.ClinicPatientListItemResponse, error) {
	clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	normalizedFilters, err := normalizeClinicPatientFilters(filters)
	if err != nil {
		return nil, err
	}
	records, err := s.patientRepo.FindPatientsByClinicProfileID(clinicProfile.ID, normalizedFilters)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := make([]dto.ClinicPatientListItemResponse, 0, len(records))
	for _, record := range records {
		response = append(response, toClinicPatientListItem(record))
	}
	return response, nil
}

func (s *clinicPatientService) GetClinicPatient(currentUserID string, petID uuid.UUID) (*dto.ClinicPatientDetailResponse, error) {
	clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	detail, err := s.patientRepo.FindPatientDetailByClinicProfileIDAndPetID(clinicProfile.ID, petID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicPatientNotFound) {
			return nil, clinicPatientNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toClinicPatientDetailResponse(detail)
	return &response, nil
}

func (s *clinicPatientService) currentClinicProfile(currentUserID string) (*models.ClinicProfile, error) {
	userID, err := parseAppointmentUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	profile, err := s.clinicRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return nil, clinicPatientClinicProfileRequiredError(err)
		}
		return nil, internalServerError(err)
	}
	return profile, nil
}

func normalizeClinicPatientFilters(filters dto.ClinicPatientFilters) (repositories.ClinicPatientFilters, error) {
	query := strings.TrimSpace(filters.Q)
	if utf8.RuneCountInString(query) > maxClinicPatientSearch {
		return repositories.ClinicPatientFilters{}, clinicPatientValidationError(
			"q must not exceed " + strconv.Itoa(maxClinicPatientSearch) + " characters",
		)
	}

	species := strings.ToLower(strings.TrimSpace(filters.Species))
	if species != "" && species != models.SpeciesDog && species != models.SpeciesCat {
		return repositories.ClinicPatientFilters{}, clinicPatientValidationError("species must be dog or cat")
	}

	status := strings.ToLower(strings.TrimSpace(filters.Status))
	if status != "" {
		if _, ok := allowedAppointmentStatuses[status]; !ok {
			return repositories.ClinicPatientFilters{}, clinicPatientValidationError("status is not supported")
		}
	}

	limit, err := normalizeClinicPatientLimit(filters.Limit)
	if err != nil {
		return repositories.ClinicPatientFilters{}, err
	}
	offset, err := normalizeClinicPatientOffset(filters.Offset)
	if err != nil {
		return repositories.ClinicPatientFilters{}, err
	}
	sort := strings.ToLower(strings.TrimSpace(filters.Sort))
	if sort == "" {
		sort = "latest_appointment_desc"
	}
	if _, ok := allowedClinicPatientSorts[sort]; !ok {
		return repositories.ClinicPatientFilters{}, clinicPatientValidationError("sort is not supported")
	}

	return repositories.ClinicPatientFilters{
		Query:   query,
		Species: species,
		Status:  status,
		Limit:   limit,
		Offset:  offset,
		Sort:    sort,
	}, nil
}

func normalizeClinicPatientLimit(value string) (int, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return defaultClinicPatientLimit, nil
	}
	limit, err := strconv.Atoi(text)
	if err != nil || limit <= 0 {
		return 0, clinicPatientValidationError("limit must be a positive integer")
	}
	if limit > maxClinicPatientLimit {
		return 0, clinicPatientValidationError("limit must be at most " + strconv.Itoa(maxClinicPatientLimit))
	}
	return limit, nil
}

func normalizeClinicPatientOffset(value string) (int, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return 0, nil
	}
	offset, err := strconv.Atoi(text)
	if err != nil || offset < 0 {
		return 0, clinicPatientValidationError("offset must be a non-negative integer")
	}
	return offset, nil
}

func toClinicPatientListItem(record repositories.ClinicPatientRecord) dto.ClinicPatientListItemResponse {
	summary := toClinicPatientAppointmentSummary(record.Summary)
	return dto.ClinicPatientListItemResponse{
		Pet:                toClinicPatientPetSummary(record.Pet),
		Owner:              toClinicPatientOwnerSummary(record.Pet),
		AppointmentSummary: summary,
		FirstSeenAt:        formatClinicPatientTime(record.Summary.FirstAppointmentAt),
	}
}

func toClinicPatientDetailResponse(detail *repositories.ClinicPatientDetail) dto.ClinicPatientDetailResponse {
	recentAppointments := make([]dto.ClinicPatientRecentAppointmentResponse, 0, len(detail.RecentAppointments))
	for _, appointment := range detail.RecentAppointments {
		recentAppointments = append(recentAppointments, dto.ClinicPatientRecentAppointmentResponse{
			ID:              appointment.ID.String(),
			ScheduledAt:     appointment.ScheduledAt.UTC().Format(time.RFC3339),
			AppointmentType: appointment.AppointmentType,
			Status:          appointment.Status,
			Title:           appointment.Title,
		})
	}

	return dto.ClinicPatientDetailResponse{
		Pet:   toClinicPatientPetDetail(detail.Pet),
		Owner: toClinicPatientOwnerSummary(detail.Pet),
		ClinicRelationship: dto.ClinicPatientRelationshipSummary{
			FirstAppointmentAt: formatClinicPatientTime(detail.Summary.FirstAppointmentAt),
			LastAppointmentAt:  formatClinicPatientTime(detail.Summary.LastAppointmentAt),
			NextAppointmentAt:  formatClinicPatientTime(detail.Summary.NextAppointmentAt),
			TotalAppointments:  detail.Summary.TotalAppointments,
		},
		RecentAppointments: recentAppointments,
	}
}

func toClinicPatientPetSummary(pet *models.Pet) dto.ClinicPatientPetSummary {
	if pet == nil {
		return dto.ClinicPatientPetSummary{}
	}
	response := dto.ClinicPatientPetSummary{
		ID:          pet.ID.String(),
		PublicPetID: pet.PublicPetID,
		Name:        pet.Name,
		Species:     pet.Species,
		AvatarURL:   pet.AvatarURL,
	}
	if pet.Breed != nil {
		breed := toBreedResponse(pet.Breed)
		response.Breed = &breed
	}
	return response
}

func toClinicPatientPetDetail(pet *models.Pet) dto.ClinicPatientPetDetail {
	if pet == nil {
		return dto.ClinicPatientPetDetail{}
	}
	response := dto.ClinicPatientPetDetail{
		ID:               pet.ID.String(),
		PublicPetID:      pet.PublicPetID,
		Name:             pet.Name,
		Species:          pet.Species,
		Gender:           pet.Gender,
		DateOfBirth:      formatClinicPatientDate(pet.DateOfBirth),
		WeightKG:         pet.WeightKG,
		MicrochipID:      pet.MicrochipID,
		AvatarURL:        pet.AvatarURL,
		Color:            pet.Color,
		DistinctiveMarks: pet.DistinctiveMarks,
		IsNeutered:       pet.IsNeutered,
	}
	if pet.Breed != nil {
		breed := toBreedResponse(pet.Breed)
		response.Breed = &breed
	}
	return response
}

func toClinicPatientOwnerSummary(pet *models.Pet) dto.ClinicPatientOwnerSummary {
	if pet == nil || pet.OwnerProfile == nil {
		return dto.ClinicPatientOwnerSummary{}
	}
	return dto.ClinicPatientOwnerSummary{
		DisplayName: strings.TrimSpace(pet.OwnerProfile.FirstName + " " + pet.OwnerProfile.LastName),
		MaskedPhone: maskOwnerPhone(pet.OwnerProfile.PhoneNumber),
	}
}

func toClinicPatientAppointmentSummary(summary repositories.ClinicPatientSummary) dto.ClinicPatientAppointmentSummary {
	return dto.ClinicPatientAppointmentSummary{
		TotalAppointments: summary.TotalAppointments,
		LastAppointmentAt: formatClinicPatientTime(summary.LastAppointmentAt),
		NextAppointmentAt: formatClinicPatientTime(summary.NextAppointmentAt),
		LatestStatus:      summary.LatestStatus,
	}
}

func formatClinicPatientTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func formatClinicPatientDate(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.Format(dateOnlyLayout)
	return &formatted
}

func clinicPatientValidationError(details string) *utils.AppError {
	return utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details, nil)
}

func clinicPatientNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"CLINIC_PATIENT_NOT_FOUND",
		"Patient not found",
		"The patient does not exist or is outside the authenticated clinic",
		cause,
	)
}

func clinicPatientClinicProfileRequiredError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"CLINIC_PROFILE_REQUIRED",
		"Clinic profile required",
		"Create a clinic profile before viewing patients",
		cause,
	)
}

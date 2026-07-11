package services

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

const (
	defaultMedicalRecordPage  = 1
	defaultMedicalRecordLimit = 20
	maxMedicalRecordLimit     = 100
)

type MedicalRecordService interface {
	CreateMedicalRecord(currentUserID string, petID uuid.UUID, req dto.CreateMedicalRecordRequest) (*dto.MedicalRecordDetailResponse, error)
	ListMedicalRecords(currentUserID string, filters dto.MedicalRecordFilters) (*dto.MedicalRecordListResponse, error)
	GetMedicalRecord(currentUserID string, recordID uuid.UUID) (*dto.MedicalRecordDetailResponse, error)
	UpdateMedicalRecord(currentUserID string, recordID uuid.UUID, req dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordDetailResponse, error)
}

type medicalRecordService struct {
	recordRepo repositories.MedicalRecordRepository
	clinicRepo repositories.ClinicProfileRepository
	now        func() time.Time
}

func NewMedicalRecordService(
	recordRepo repositories.MedicalRecordRepository,
	clinicRepo repositories.ClinicProfileRepository,
) MedicalRecordService {
	return &medicalRecordService{
		recordRepo: recordRepo,
		clinicRepo: clinicRepo,
		now:        time.Now,
	}
}

func (s *medicalRecordService) CreateMedicalRecord(currentUserID string, petID uuid.UUID, req dto.CreateMedicalRecordRequest) (*dto.MedicalRecordDetailResponse, error) {
	userID, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	input, err := normalizeCreateMedicalRecordInput(req)
	if err != nil {
		return nil, err
	}
	if err := s.ensureClinicPatient(clinicProfile.ID, petID); err != nil {
		return nil, err
	}
	if input.appointmentID != nil {
		if err := s.ensureUsableAppointment(*input.appointmentID, clinicProfile.ID, petID); err != nil {
			return nil, err
		}
	}

	now := s.now().UTC()
	record := &models.MedicalRecord{
		ClinicProfileID:      clinicProfile.ID,
		PetID:                petID,
		AppointmentID:        input.appointmentID,
		CreatedByUserID:      userID,
		VisitAt:              input.visitAt,
		ChiefComplaint:       input.chiefComplaint,
		ClinicalFindings:     input.clinicalFindings,
		Diagnosis:            input.diagnosis,
		TreatmentPlan:        input.treatmentPlan,
		Medications:          input.medications,
		FollowUpInstructions: input.followUpInstructions,
		NextFollowUpAt:       input.nextFollowUpAt,
		WeightKG:             input.weightKG,
		TemperatureC:         input.temperatureC,
		Notes:                input.notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := s.recordRepo.Create(record); err != nil {
		if errors.Is(err, repositories.ErrMedicalRecordAppointmentAlreadyLinked) {
			return nil, medicalRecordAppointmentConflictError(err)
		}
		return nil, internalServerError(err)
	}
	created, err := s.recordRepo.FindByIDAndClinicProfileID(record.ID, clinicProfile.ID)
	if err != nil {
		return nil, internalServerError(err)
	}
	response := toMedicalRecordDetailResponse(created)
	return &response, nil
}

func (s *medicalRecordService) ListMedicalRecords(currentUserID string, filters dto.MedicalRecordFilters) (*dto.MedicalRecordListResponse, error) {
	_, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	normalizedFilters, pagination, err := normalizeMedicalRecordFilters(filters)
	if err != nil {
		return nil, err
	}
	total, err := s.recordRepo.CountByClinicProfileID(clinicProfile.ID, normalizedFilters)
	if err != nil {
		return nil, internalServerError(err)
	}
	records, err := s.recordRepo.FindAllByClinicProfileID(clinicProfile.ID, normalizedFilters)
	if err != nil {
		return nil, internalServerError(err)
	}

	items := make([]dto.MedicalRecordListItemResponse, 0, len(records))
	for i := range records {
		items = append(items, toMedicalRecordListItemResponse(&records[i]))
	}
	pagination.Total = total
	if total > 0 {
		pagination.TotalPages = int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	}
	return &dto.MedicalRecordListResponse{Items: items, Pagination: pagination}, nil
}

func (s *medicalRecordService) GetMedicalRecord(currentUserID string, recordID uuid.UUID) (*dto.MedicalRecordDetailResponse, error) {
	_, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	record, err := s.recordRepo.FindByIDAndClinicProfileID(recordID, clinicProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrMedicalRecordNotFound) {
			return nil, medicalRecordNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	response := toMedicalRecordDetailResponse(record)
	return &response, nil
}

func (s *medicalRecordService) UpdateMedicalRecord(currentUserID string, recordID uuid.UUID, req dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordDetailResponse, error) {
	if !hasMedicalRecordUpdate(req) {
		return nil, medicalRecordValidationError("Request body must contain at least one medical record field")
	}
	_, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	record, err := s.recordRepo.FindByIDAndClinicProfileID(recordID, clinicProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrMedicalRecordNotFound) {
			return nil, medicalRecordNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	if err := applyMedicalRecordUpdate(record, req); err != nil {
		return nil, err
	}
	record.UpdatedAt = s.now().UTC()
	if err := s.recordRepo.Update(record); err != nil {
		if errors.Is(err, repositories.ErrMedicalRecordNotFound) {
			return nil, medicalRecordNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	updated, err := s.recordRepo.FindByIDAndClinicProfileID(recordID, clinicProfile.ID)
	if err != nil {
		return nil, internalServerError(err)
	}
	response := toMedicalRecordDetailResponse(updated)
	return &response, nil
}

func (s *medicalRecordService) currentClinicProfile(currentUserID string) (uuid.UUID, *models.ClinicProfile, error) {
	userID, err := parseAppointmentUserID(currentUserID)
	if err != nil {
		return uuid.Nil, nil, err
	}
	profile, err := s.clinicRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return uuid.Nil, nil, medicalRecordClinicProfileRequiredError(err)
		}
		return uuid.Nil, nil, internalServerError(err)
	}
	return userID, profile, nil
}

func (s *medicalRecordService) ensureClinicPatient(clinicProfileID, petID uuid.UUID) error {
	exists, err := s.recordRepo.PetHasNonCancelledAppointmentWithClinic(clinicProfileID, petID)
	if err != nil {
		return internalServerError(err)
	}
	if !exists {
		return medicalRecordPatientNotFoundError(repositories.ErrClinicPatientNotFound)
	}
	return nil
}

func (s *medicalRecordService) ensureUsableAppointment(appointmentID, clinicProfileID, petID uuid.UUID) error {
	if _, err := s.recordRepo.FindUsableAppointmentForMedicalRecord(appointmentID, clinicProfileID, petID); err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return medicalRecordAppointmentNotFoundError(err)
		}
		return internalServerError(err)
	}
	exists, err := s.recordRepo.MedicalRecordExistsByAppointmentID(appointmentID)
	if err != nil {
		return internalServerError(err)
	}
	if exists {
		return medicalRecordAppointmentConflictError(repositories.ErrMedicalRecordAppointmentAlreadyLinked)
	}
	return nil
}

type normalizedMedicalRecordInput struct {
	appointmentID        *uuid.UUID
	visitAt              time.Time
	chiefComplaint       string
	clinicalFindings     *string
	diagnosis            *string
	treatmentPlan        *string
	medications          *string
	followUpInstructions *string
	nextFollowUpAt       *time.Time
	weightKG             *float64
	temperatureC         *float64
	notes                *string
}

func normalizeCreateMedicalRecordInput(req dto.CreateMedicalRecordRequest) (normalizedMedicalRecordInput, error) {
	appointmentID, err := normalizeOptionalMedicalRecordUUID("appointmentId", req.AppointmentID)
	if err != nil {
		return normalizedMedicalRecordInput{}, err
	}
	visitAt, err := parseRequiredMedicalRecordTime("visitAt", req.VisitAt)
	if err != nil {
		return normalizedMedicalRecordInput{}, err
	}
	chiefComplaint, err := normalizeRequiredMedicalRecordText("chiefComplaint", req.ChiefComplaint)
	if err != nil {
		return normalizedMedicalRecordInput{}, err
	}
	nextFollowUpAt, err := parseOptionalMedicalRecordTime("nextFollowUpAt", req.NextFollowUpAt)
	if err != nil {
		return normalizedMedicalRecordInput{}, err
	}
	if err := validateMedicalRecordFollowUp(visitAt, nextFollowUpAt); err != nil {
		return normalizedMedicalRecordInput{}, err
	}
	if err := validateMedicalRecordVitals(req.WeightKG, req.TemperatureC); err != nil {
		return normalizedMedicalRecordInput{}, err
	}
	return normalizedMedicalRecordInput{
		appointmentID:        appointmentID,
		visitAt:              visitAt,
		chiefComplaint:       chiefComplaint,
		clinicalFindings:     normalizeOptionalMedicalRecordText(req.ClinicalFindings),
		diagnosis:            normalizeOptionalMedicalRecordText(req.Diagnosis),
		treatmentPlan:        normalizeOptionalMedicalRecordText(req.TreatmentPlan),
		medications:          normalizeOptionalMedicalRecordText(req.Medications),
		followUpInstructions: normalizeOptionalMedicalRecordText(req.FollowUpInstructions),
		nextFollowUpAt:       nextFollowUpAt,
		weightKG:             req.WeightKG,
		temperatureC:         req.TemperatureC,
		notes:                normalizeOptionalMedicalRecordText(req.Notes),
	}, nil
}

func normalizeMedicalRecordFilters(filters dto.MedicalRecordFilters) (repositories.MedicalRecordFilters, dto.PaginationMeta, error) {
	page, err := normalizePositiveInt(filters.Page, defaultMedicalRecordPage, "page")
	if err != nil {
		return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, err
	}
	limit, err := normalizePositiveInt(filters.Limit, defaultMedicalRecordLimit, "limit")
	if err != nil {
		return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, err
	}
	if limit > maxMedicalRecordLimit {
		return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, medicalRecordValidationError("limit must be at most " + strconv.Itoa(maxMedicalRecordLimit))
	}

	var petID *uuid.UUID
	if strings.TrimSpace(filters.PetID) != "" {
		id, err := uuid.Parse(strings.TrimSpace(filters.PetID))
		if err != nil {
			return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, medicalRecordValidationError("pet_id must be a valid UUID")
		}
		petID = &id
	}

	var dateFrom *time.Time
	var dateTo *time.Time
	if strings.TrimSpace(filters.From) != "" {
		from, err := parseMedicalRecordDate("from", filters.From)
		if err != nil {
			return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, err
		}
		dateFrom = &from
	}
	if strings.TrimSpace(filters.To) != "" {
		to, err := parseMedicalRecordDate("to", filters.To)
		if err != nil {
			return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, err
		}
		if dateFrom != nil && dateFrom.After(to) {
			return repositories.MedicalRecordFilters{}, dto.PaginationMeta{}, medicalRecordValidationError("from must not be after to")
		}
		toExclusive := to.AddDate(0, 0, 1)
		dateTo = &toExclusive
	}

	return repositories.MedicalRecordFilters{
			PetID:    petID,
			DateFrom: dateFrom,
			DateTo:   dateTo,
			Limit:    limit,
			Offset:   (page - 1) * limit,
		},
		dto.PaginationMeta{Page: page, Limit: limit},
		nil
}

func normalizePositiveInt(value string, defaultValue int, field string) (int, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(text)
	if err != nil || parsed <= 0 {
		return 0, medicalRecordValidationError(field + " must be a positive integer")
	}
	return parsed, nil
}

func parseMedicalRecordDate(field, value string) (time.Time, error) {
	parsed, err := time.Parse(dateOnlyLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, medicalRecordValidationError(field + " must use YYYY-MM-DD format")
	}
	return parsed.UTC(), nil
}

func normalizeOptionalMedicalRecordUUID(field, value string) (*uuid.UUID, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(text)
	if err != nil {
		return nil, medicalRecordValidationError(field + " must be a valid UUID")
	}
	return &parsed, nil
}

func parseRequiredMedicalRecordTime(field, value string) (time.Time, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return time.Time{}, medicalRecordValidationError(field + " is required")
	}
	return parseMedicalRecordTime(field, text)
}

func parseOptionalMedicalRecordTime(field, value string) (*time.Time, error) {
	text := strings.TrimSpace(value)
	if text == "" {
		return nil, nil
	}
	parsed, err := parseMedicalRecordTime(field, text)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseMedicalRecordTime(field, value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, medicalRecordValidationError(field + " must use RFC3339 format")
	}
	return parsed.UTC(), nil
}

func normalizeRequiredMedicalRecordText(field, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", medicalRecordValidationError(field + " is required")
	}
	return trimmed, nil
}

func normalizeOptionalMedicalRecordText(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func validateMedicalRecordFollowUp(visitAt time.Time, nextFollowUpAt *time.Time) error {
	if nextFollowUpAt != nil && nextFollowUpAt.Before(visitAt) {
		return medicalRecordValidationError("nextFollowUpAt cannot be earlier than visitAt")
	}
	return nil
}

func validateMedicalRecordVitals(weightKG, temperatureC *float64) error {
	if weightKG != nil && (math.IsNaN(*weightKG) || math.IsInf(*weightKG, 0) || *weightKG <= 0) {
		return medicalRecordValidationError("weightKg must be greater than zero")
	}
	if temperatureC != nil && (math.IsNaN(*temperatureC) || math.IsInf(*temperatureC, 0) || *temperatureC <= 0) {
		return medicalRecordValidationError("temperatureC must be greater than zero")
	}
	return nil
}

func hasMedicalRecordUpdate(req dto.UpdateMedicalRecordRequest) bool {
	return req.VisitAt != nil ||
		req.ChiefComplaint != nil ||
		req.ClinicalFindings != nil ||
		req.Diagnosis != nil ||
		req.TreatmentPlan != nil ||
		req.Medications != nil ||
		req.FollowUpInstructions != nil ||
		req.NextFollowUpAt != nil ||
		req.WeightKG != nil ||
		req.TemperatureC != nil ||
		req.Notes != nil
}

func applyMedicalRecordUpdate(record *models.MedicalRecord, req dto.UpdateMedicalRecordRequest) error {
	if req.VisitAt != nil {
		visitAt, err := parseRequiredMedicalRecordTime("visitAt", *req.VisitAt)
		if err != nil {
			return err
		}
		record.VisitAt = visitAt
	}
	if req.ChiefComplaint != nil {
		chiefComplaint, err := normalizeRequiredMedicalRecordText("chiefComplaint", *req.ChiefComplaint)
		if err != nil {
			return err
		}
		record.ChiefComplaint = chiefComplaint
	}
	if req.ClinicalFindings != nil {
		record.ClinicalFindings = normalizeOptionalMedicalRecordText(*req.ClinicalFindings)
	}
	if req.Diagnosis != nil {
		record.Diagnosis = normalizeOptionalMedicalRecordText(*req.Diagnosis)
	}
	if req.TreatmentPlan != nil {
		record.TreatmentPlan = normalizeOptionalMedicalRecordText(*req.TreatmentPlan)
	}
	if req.Medications != nil {
		record.Medications = normalizeOptionalMedicalRecordText(*req.Medications)
	}
	if req.FollowUpInstructions != nil {
		record.FollowUpInstructions = normalizeOptionalMedicalRecordText(*req.FollowUpInstructions)
	}
	if req.NextFollowUpAt != nil {
		nextFollowUpAt, err := parseOptionalMedicalRecordTime("nextFollowUpAt", *req.NextFollowUpAt)
		if err != nil {
			return err
		}
		record.NextFollowUpAt = nextFollowUpAt
	}
	if err := validateMedicalRecordVitals(req.WeightKG, req.TemperatureC); err != nil {
		return err
	}
	if req.WeightKG != nil {
		record.WeightKG = req.WeightKG
	}
	if req.TemperatureC != nil {
		record.TemperatureC = req.TemperatureC
	}
	if req.Notes != nil {
		record.Notes = normalizeOptionalMedicalRecordText(*req.Notes)
	}
	return validateMedicalRecordFollowUp(record.VisitAt, record.NextFollowUpAt)
}

func toMedicalRecordListItemResponse(record *models.MedicalRecord) dto.MedicalRecordListItemResponse {
	return dto.MedicalRecordListItemResponse{
		ID:             record.ID.String(),
		VisitAt:        record.VisitAt.UTC().Format(time.RFC3339),
		ChiefComplaint: record.ChiefComplaint,
		Diagnosis:      record.Diagnosis,
		CreatedAt:      record.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      record.UpdatedAt.UTC().Format(time.RFC3339),
		Pet:            toMedicalRecordPetSummary(record.Pet),
		Owner:          toMedicalRecordOwnerSummary(record.Pet),
		Appointment:    toMedicalRecordAppointmentSummary(record.Appointment),
	}
}

func toMedicalRecordDetailResponse(record *models.MedicalRecord) dto.MedicalRecordDetailResponse {
	return dto.MedicalRecordDetailResponse{
		ID:                   record.ID.String(),
		VisitAt:              record.VisitAt.UTC().Format(time.RFC3339),
		ChiefComplaint:       record.ChiefComplaint,
		ClinicalFindings:     record.ClinicalFindings,
		Diagnosis:            record.Diagnosis,
		TreatmentPlan:        record.TreatmentPlan,
		Medications:          record.Medications,
		FollowUpInstructions: record.FollowUpInstructions,
		NextFollowUpAt:       formatMedicalRecordTime(record.NextFollowUpAt),
		WeightKG:             record.WeightKG,
		TemperatureC:         record.TemperatureC,
		Notes:                record.Notes,
		CreatedAt:            record.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:            record.UpdatedAt.UTC().Format(time.RFC3339),
		Pet:                  toMedicalRecordPetSummary(record.Pet),
		Owner:                toMedicalRecordOwnerSummary(record.Pet),
		Appointment:          toMedicalRecordAppointmentSummary(record.Appointment),
		CreatedBy:            toMedicalRecordCreatedBySummary(record.CreatedByUser),
	}
}

func toMedicalRecordPetSummary(pet *models.Pet) dto.MedicalRecordPetSummary {
	if pet == nil {
		return dto.MedicalRecordPetSummary{}
	}
	response := dto.MedicalRecordPetSummary{
		ID:          pet.ID.String(),
		PublicPetID: pet.PublicPetID,
		Name:        pet.Name,
		Species:     pet.Species,
	}
	if pet.Breed != nil {
		breed := toBreedResponse(pet.Breed)
		response.Breed = &breed
	}
	return response
}

func toMedicalRecordOwnerSummary(pet *models.Pet) dto.MedicalRecordOwnerSummary {
	if pet == nil || pet.OwnerProfile == nil {
		return dto.MedicalRecordOwnerSummary{}
	}
	return dto.MedicalRecordOwnerSummary{
		ID:          pet.OwnerProfile.ID.String(),
		FullName:    strings.TrimSpace(pet.OwnerProfile.FirstName + " " + pet.OwnerProfile.LastName),
		PhoneNumber: pet.OwnerProfile.PhoneNumber,
	}
}

func toMedicalRecordAppointmentSummary(appointment *models.Appointment) *dto.MedicalRecordAppointmentSummary {
	if appointment == nil || appointment.ID == uuid.Nil {
		return nil
	}
	return &dto.MedicalRecordAppointmentSummary{
		ID:          appointment.ID.String(),
		ScheduledAt: appointment.ScheduledAt.UTC().Format(time.RFC3339),
		Status:      appointment.Status,
	}
}

func toMedicalRecordCreatedBySummary(user *models.User) *dto.MedicalRecordCreatedBySummary {
	if user == nil || user.ID == uuid.Nil {
		return nil
	}
	return &dto.MedicalRecordCreatedBySummary{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}
}

func formatMedicalRecordTime(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func medicalRecordValidationError(details string) *utils.AppError {
	return utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details, nil)
}

func medicalRecordClinicProfileRequiredError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"CLINIC_PROFILE_REQUIRED",
		"Clinic profile required",
		"Create a clinic profile before managing medical records",
		cause,
	)
}

func medicalRecordPatientNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"CLINIC_PATIENT_NOT_FOUND",
		"Patient not found",
		"The pet is not a patient of the authenticated clinic",
		cause,
	)
}

func medicalRecordAppointmentNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"APPOINTMENT_NOT_FOUND",
		"Appointment not found",
		"The appointment does not exist, is cancelled, or is outside the authenticated clinic and pet",
		cause,
	)
}

func medicalRecordNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"MEDICAL_RECORD_NOT_FOUND",
		"Medical record not found",
		"The medical record does not exist or is outside the authenticated clinic",
		cause,
	)
}

func medicalRecordAppointmentConflictError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusConflict,
		"APPOINTMENT_MEDICAL_RECORD_EXISTS",
		"Appointment already has a medical record",
		"The selected appointment already has a medical record",
		cause,
	)
}

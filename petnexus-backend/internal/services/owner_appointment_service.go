package services

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
)

type OwnerAppointmentService interface {
	CreateOwnerAppointment(currentUserID string, req dto.CreateOwnerAppointmentRequest) (*dto.AppointmentResponse, error)
	ListOwnerAppointments(currentUserID string, filters dto.OwnerAppointmentFilters) ([]dto.AppointmentResponse, error)
	GetOwnerAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error)
	CancelOwnerAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error)
}

type ownerAppointmentService struct {
	appointmentRepo repositories.AppointmentRepository
	ownerRepo       repositories.OwnerProfileRepository
	clinicRepo      repositories.ClinicProfileRepository
	petRepo         repositories.PetRepository
	now             func() time.Time
}

func NewOwnerAppointmentService(
	appointmentRepo repositories.AppointmentRepository,
	ownerRepo repositories.OwnerProfileRepository,
	clinicRepo repositories.ClinicProfileRepository,
	petRepo repositories.PetRepository,
) OwnerAppointmentService {
	return &ownerAppointmentService{
		appointmentRepo: appointmentRepo,
		ownerRepo:       ownerRepo,
		clinicRepo:      clinicRepo,
		petRepo:         petRepo,
		now:             time.Now,
	}
}

func (s *ownerAppointmentService) CreateOwnerAppointment(currentUserID string, req dto.CreateOwnerAppointmentRequest) (*dto.AppointmentResponse, error) {
	userID, ownerProfile, err := s.currentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	clinicID, err := uuid.Parse(req.ClinicProfileID)
	if err != nil {
		return nil, appointmentValidationError("clinic_profile_id must be a valid UUID")
	}
	petID, err := uuid.Parse(req.PetID)
	if err != nil {
		return nil, appointmentValidationError("pet_id must be a valid UUID")
	}
	clinicProfile, err := s.clinicRepo.FindByID(clinicID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return nil, appointmentClinicNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	pet, err := s.petRepo.FindByIDAndOwnerProfileID(petID, ownerProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrPetNotFound) {
			return nil, appointmentPetNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	input, err := normalizeAppointmentInput(
		req.Title,
		req.AppointmentType,
		req.ScheduledAt,
		req.DurationMinutes,
		req.Note,
		s.now(),
	)
	if err != nil {
		return nil, err
	}

	appointment := &models.Appointment{
		OwnerProfileID:  ownerProfile.ID,
		ClinicProfileID: clinicProfile.ID,
		PetID:           pet.ID,
		Title:           input.title,
		AppointmentType: input.appointmentType,
		ScheduledAt:     input.scheduledAt,
		DurationMinutes: input.durationMinutes,
		Status:          models.AppointmentStatusRequested,
		Note:            input.note,
		CreatedByUserID: &userID,
		CreatedByRole:   models.RoleOwner,
		Pet:             pet,
		OwnerProfile:    ownerProfile,
		ClinicProfile:   clinicProfile,
	}
	if err := s.appointmentRepo.Create(appointment); err != nil {
		return nil, internalServerError(err)
	}
	response := toAppointmentResponse(appointment)
	return &response, nil
}

func (s *ownerAppointmentService) ListOwnerAppointments(currentUserID string, filters dto.OwnerAppointmentFilters) ([]dto.AppointmentResponse, error) {
	_, ownerProfile, err := s.currentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	normalizedFilters, err := normalizeOwnerAppointmentFilters(filters)
	if err != nil {
		return nil, err
	}
	appointments, err := s.appointmentRepo.FindAllByOwnerProfileID(ownerProfile.ID, normalizedFilters)
	if err != nil {
		return nil, internalServerError(err)
	}
	return mapAppointmentResponses(appointments), nil
}

func (s *ownerAppointmentService) GetOwnerAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	_, ownerProfile, err := s.currentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	appointment, err := s.appointmentRepo.FindByIDAndOwnerProfileID(appointmentID, ownerProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return nil, appointmentNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	response := toAppointmentResponse(appointment)
	return &response, nil
}

func (s *ownerAppointmentService) CancelOwnerAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	_, ownerProfile, err := s.currentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	appointment, err := s.appointmentRepo.FindByIDAndOwnerProfileID(appointmentID, ownerProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return nil, appointmentNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	if appointment.Status == models.AppointmentStatusCancelled && appointment.CancelledAt != nil {
		response := toAppointmentResponse(appointment)
		return &response, nil
	}
	setAppointmentStatus(appointment, models.AppointmentStatusCancelled, s.now())
	if err := s.appointmentRepo.Update(appointment); err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return nil, appointmentNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	response := toAppointmentResponse(appointment)
	return &response, nil
}

func (s *ownerAppointmentService) currentOwnerProfile(currentUserID string) (uuid.UUID, *models.OwnerProfile, error) {
	userID, err := parseAppointmentUserID(currentUserID)
	if err != nil {
		return uuid.Nil, nil, err
	}
	profile, err := s.ownerRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrOwnerProfileNotFound) {
			return uuid.Nil, nil, appointmentOwnerProfileRequiredError(err)
		}
		return uuid.Nil, nil, internalServerError(err)
	}
	return userID, profile, nil
}

func mapAppointmentResponses(appointments []models.Appointment) []dto.AppointmentResponse {
	responses := make([]dto.AppointmentResponse, 0, len(appointments))
	for i := range appointments {
		responses = append(responses, toAppointmentResponse(&appointments[i]))
	}
	return responses
}

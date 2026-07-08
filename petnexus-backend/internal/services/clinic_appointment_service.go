package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
)

type ClinicAppointmentService interface {
	CreateClinicAppointment(currentUserID string, req dto.CreateClinicAppointmentRequest) (*dto.AppointmentResponse, error)
	ListClinicAppointments(currentUserID string, filters dto.ClinicAppointmentFilters) ([]dto.AppointmentResponse, error)
	GetClinicAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error)
	UpdateClinicAppointmentStatus(currentUserID string, appointmentID uuid.UUID, req dto.UpdateAppointmentStatusRequest) (*dto.AppointmentResponse, error)
	CancelClinicAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error)
}

type clinicAppointmentService struct {
	appointmentRepo repositories.AppointmentRepository
	clinicRepo      repositories.ClinicProfileRepository
	petRepo         repositories.PetRepository
	now             func() time.Time
}

func NewClinicAppointmentService(
	appointmentRepo repositories.AppointmentRepository,
	clinicRepo repositories.ClinicProfileRepository,
	petRepo repositories.PetRepository,
) ClinicAppointmentService {
	return &clinicAppointmentService{
		appointmentRepo: appointmentRepo,
		clinicRepo:      clinicRepo,
		petRepo:         petRepo,
		now:             time.Now,
	}
}

func (s *clinicAppointmentService) CreateClinicAppointment(currentUserID string, req dto.CreateClinicAppointmentRequest) (*dto.AppointmentResponse, error) {
	userID, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	pet, err := s.resolveAppointmentPet(req.PetID, req.PublicPetID)
	if err != nil {
		return nil, err
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
		OwnerProfileID:  pet.OwnerProfileID,
		ClinicProfileID: clinicProfile.ID,
		PetID:           pet.ID,
		Title:           input.title,
		AppointmentType: input.appointmentType,
		ScheduledAt:     input.scheduledAt,
		DurationMinutes: input.durationMinutes,
		Status:          models.AppointmentStatusScheduled,
		Note:            input.note,
		CreatedByUserID: &userID,
		CreatedByRole:   models.RoleClinic,
		Pet:             pet,
		OwnerProfile:    pet.OwnerProfile,
		ClinicProfile:   clinicProfile,
	}
	if err := s.appointmentRepo.Create(appointment); err != nil {
		return nil, internalServerError(err)
	}
	response := toAppointmentResponse(appointment)
	return &response, nil
}

func (s *clinicAppointmentService) ListClinicAppointments(currentUserID string, filters dto.ClinicAppointmentFilters) ([]dto.AppointmentResponse, error) {
	_, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	normalizedFilters, err := normalizeClinicAppointmentFilters(filters)
	if err != nil {
		return nil, err
	}
	appointments, err := s.appointmentRepo.FindAllByClinicProfileID(clinicProfile.ID, normalizedFilters)
	if err != nil {
		return nil, internalServerError(err)
	}
	return mapAppointmentResponses(appointments), nil
}

func (s *clinicAppointmentService) GetClinicAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	_, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	appointment, err := s.appointmentRepo.FindByIDAndClinicProfileID(appointmentID, clinicProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return nil, appointmentNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	response := toAppointmentResponse(appointment)
	return &response, nil
}

func (s *clinicAppointmentService) UpdateClinicAppointmentStatus(currentUserID string, appointmentID uuid.UUID, req dto.UpdateAppointmentStatusRequest) (*dto.AppointmentResponse, error) {
	status, err := normalizeAppointmentStatus(req.Status)
	if err != nil {
		return nil, err
	}
	return s.updateClinicAppointmentStatus(currentUserID, appointmentID, status)
}

func (s *clinicAppointmentService) CancelClinicAppointment(currentUserID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	return s.updateClinicAppointmentStatus(currentUserID, appointmentID, models.AppointmentStatusCancelled)
}

func (s *clinicAppointmentService) updateClinicAppointmentStatus(currentUserID string, appointmentID uuid.UUID, status string) (*dto.AppointmentResponse, error) {
	_, clinicProfile, err := s.currentClinicProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	appointment, err := s.appointmentRepo.FindByIDAndClinicProfileID(appointmentID, clinicProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return nil, appointmentNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	if status == models.AppointmentStatusCancelled &&
		appointment.Status == models.AppointmentStatusCancelled &&
		appointment.CancelledAt != nil {
		response := toAppointmentResponse(appointment)
		return &response, nil
	}
	setAppointmentStatus(appointment, status, s.now())
	if err := s.appointmentRepo.Update(appointment); err != nil {
		if errors.Is(err, repositories.ErrAppointmentNotFound) {
			return nil, appointmentNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	response := toAppointmentResponse(appointment)
	return &response, nil
}

func (s *clinicAppointmentService) currentClinicProfile(currentUserID string) (uuid.UUID, *models.ClinicProfile, error) {
	userID, err := parseAppointmentUserID(currentUserID)
	if err != nil {
		return uuid.Nil, nil, err
	}
	profile, err := s.clinicRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrClinicProfileNotFound) {
			return uuid.Nil, nil, appointmentClinicProfileRequiredError(err)
		}
		return uuid.Nil, nil, internalServerError(err)
	}
	return userID, profile, nil
}

func (s *clinicAppointmentService) resolveAppointmentPet(petIDValue, publicPetIDValue string) (*models.Pet, error) {
	petIDValue = strings.TrimSpace(petIDValue)
	publicPetIDValue = strings.ToUpper(strings.TrimSpace(publicPetIDValue))
	if petIDValue == "" && publicPetIDValue == "" {
		return nil, appointmentValidationError("pet_id or public_pet_id is required")
	}
	if petIDValue != "" && publicPetIDValue != "" {
		return nil, appointmentValidationError("provide only one of pet_id or public_pet_id")
	}

	var (
		pet *models.Pet
		err error
	)
	if publicPetIDValue != "" {
		pet, err = s.petRepo.FindByPublicPetID(publicPetIDValue)
	} else {
		petID, parseErr := uuid.Parse(petIDValue)
		if parseErr != nil {
			return nil, appointmentValidationError("pet_id must be a valid UUID")
		}
		pet, err = s.petRepo.FindByID(petID)
	}
	if err != nil {
		if errors.Is(err, repositories.ErrPetNotFound) {
			return nil, appointmentPetNotFoundError(err)
		}
		return nil, internalServerError(err)
	}
	return pet, nil
}

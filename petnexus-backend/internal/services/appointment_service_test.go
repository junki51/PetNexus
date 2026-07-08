package services

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type appointmentRepositoryStub struct {
	appointments map[uuid.UUID]*models.Appointment
	lastFilters  repositories.AppointmentFilters
}

func newAppointmentRepositoryStub() *appointmentRepositoryStub {
	return &appointmentRepositoryStub{appointments: make(map[uuid.UUID]*models.Appointment)}
}

func (r *appointmentRepositoryStub) Create(appointment *models.Appointment) error {
	if appointment.ID == uuid.Nil {
		appointment.ID = uuid.New()
	}
	if appointment.CreatedAt.IsZero() {
		appointment.CreatedAt = time.Date(2026, 7, 8, 2, 0, 0, 0, time.UTC)
		appointment.UpdatedAt = appointment.CreatedAt
	}
	r.appointments[appointment.ID] = appointment
	return nil
}

func (r *appointmentRepositoryStub) FindByID(id uuid.UUID) (*models.Appointment, error) {
	appointment, ok := r.appointments[id]
	if !ok {
		return nil, repositories.ErrAppointmentNotFound
	}
	return appointment, nil
}

func (r *appointmentRepositoryStub) FindByIDAndOwnerProfileID(id, ownerProfileID uuid.UUID) (*models.Appointment, error) {
	appointment, err := r.FindByID(id)
	if err != nil || appointment.OwnerProfileID != ownerProfileID {
		return nil, repositories.ErrAppointmentNotFound
	}
	return appointment, nil
}

func (r *appointmentRepositoryStub) FindByIDAndClinicProfileID(id, clinicProfileID uuid.UUID) (*models.Appointment, error) {
	appointment, err := r.FindByID(id)
	if err != nil || appointment.ClinicProfileID != clinicProfileID {
		return nil, repositories.ErrAppointmentNotFound
	}
	return appointment, nil
}

func (r *appointmentRepositoryStub) FindAllByOwnerProfileID(ownerProfileID uuid.UUID, filters repositories.AppointmentFilters) ([]models.Appointment, error) {
	r.lastFilters = filters
	return r.findAll(func(appointment *models.Appointment) bool {
		return appointment.OwnerProfileID == ownerProfileID
	}), nil
}

func (r *appointmentRepositoryStub) FindAllByClinicProfileID(clinicProfileID uuid.UUID, filters repositories.AppointmentFilters) ([]models.Appointment, error) {
	r.lastFilters = filters
	return r.findAll(func(appointment *models.Appointment) bool {
		return appointment.ClinicProfileID == clinicProfileID
	}), nil
}

func (r *appointmentRepositoryStub) Update(appointment *models.Appointment) error {
	if _, ok := r.appointments[appointment.ID]; !ok {
		return repositories.ErrAppointmentNotFound
	}
	r.appointments[appointment.ID] = appointment
	return nil
}

func (r *appointmentRepositoryStub) findAll(include func(*models.Appointment) bool) []models.Appointment {
	result := make([]models.Appointment, 0)
	for _, appointment := range r.appointments {
		if include(appointment) {
			result = append(result, *appointment)
		}
	}
	return result
}

func TestOwnerAppointmentCreateUsesJWTProfileAndOwnPet(t *testing.T) {
	now := time.Date(2026, 7, 8, 8, 0, 0, 0, time.UTC)
	userID := uuid.New()
	owner := &models.OwnerProfile{
		ID: uuid.New(), UserID: userID, FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678",
	}
	clinic := &models.ClinicProfile{ID: uuid.New(), ClinicName: "Happy Paws"}
	pet := &models.Pet{
		ID: uuid.New(), PublicPetID: "PNX-PET-ABC123", OwnerProfileID: owner.ID,
		Name: "Milo", Species: "dog", OwnerProfile: owner,
	}
	appointmentRepo := newAppointmentRepositoryStub()
	service := NewOwnerAppointmentService(
		appointmentRepo,
		&petOwnerProfileRepositoryFake{profiles: map[uuid.UUID]*models.OwnerProfile{userID: owner}},
		&clinicProfileRepositoryStub{profile: clinic, exists: true},
		&petRepositoryFake{pets: map[uuid.UUID]*models.Pet{pet.ID: pet}},
	)
	service.(*ownerAppointmentService).now = func() time.Time { return now }

	response, err := service.CreateOwnerAppointment(userID.String(), dto.CreateOwnerAppointmentRequest{
		ClinicProfileID: clinic.ID.String(),
		PetID:           pet.ID.String(),
		Title:           "  Annual checkup  ",
		AppointmentType: "checkup",
		ScheduledAt:     now.Add(24 * time.Hour).Format(time.RFC3339),
		DurationMinutes: 30,
		Note:            "  Bring vaccine card  ",
	})
	if err != nil {
		t.Fatalf("CreateOwnerAppointment() error = %v", err)
	}
	if response.Status != models.AppointmentStatusRequested || response.CreatedByRole != models.RoleOwner {
		t.Fatalf("unexpected status/creator: %#v", response)
	}
	if response.Owner.MaskedPhone != "081****678" || response.Pet.PublicPetID != pet.PublicPetID || response.Clinic.ID != clinic.ID.String() {
		t.Fatalf("unexpected summaries: %#v", response)
	}
	for _, appointment := range appointmentRepo.appointments {
		if appointment.OwnerProfileID != owner.ID || appointment.PetID != pet.ID ||
			appointment.ClinicProfileID != clinic.ID || appointment.CreatedByUserID == nil ||
			*appointment.CreatedByUserID != userID {
			t.Fatalf("appointment ownership was not resolved from server data: %#v", appointment)
		}
	}
}

func TestOwnerAppointmentRejectsAnotherOwnersPet(t *testing.T) {
	userID := uuid.New()
	owner := &models.OwnerProfile{ID: uuid.New(), UserID: userID}
	otherPet := &models.Pet{ID: uuid.New(), OwnerProfileID: uuid.New()}
	clinic := &models.ClinicProfile{ID: uuid.New()}
	service := NewOwnerAppointmentService(
		newAppointmentRepositoryStub(),
		&petOwnerProfileRepositoryFake{profiles: map[uuid.UUID]*models.OwnerProfile{userID: owner}},
		&clinicProfileRepositoryStub{profile: clinic, exists: true},
		&petRepositoryFake{pets: map[uuid.UUID]*models.Pet{otherPet.ID: otherPet}},
	)

	_, err := service.CreateOwnerAppointment(userID.String(), dto.CreateOwnerAppointmentRequest{
		ClinicProfileID: clinic.ID.String(),
		PetID:           otherPet.ID.String(),
		AppointmentType: "checkup",
		ScheduledAt:     time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		DurationMinutes: 30,
	})
	assertAppointmentServiceError(t, err, http.StatusNotFound, "PET_NOT_FOUND")
}

func TestClinicAppointmentCreateByPublicPetIDAndValidation(t *testing.T) {
	now := time.Date(2026, 7, 8, 8, 0, 0, 0, time.UTC)
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID, ClinicName: "Happy Paws"}
	owner := &models.OwnerProfile{ID: uuid.New(), FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678"}
	pet := &models.Pet{
		ID: uuid.New(), PublicPetID: "PNX-PET-ABC123", OwnerProfileID: owner.ID,
		Name: "Milo", Species: "dog", OwnerProfile: owner,
	}
	appointmentRepo := newAppointmentRepositoryStub()
	service := NewClinicAppointmentService(
		appointmentRepo,
		&clinicProfileRepositoryStub{profile: clinic, exists: true},
		&petRepositoryFake{pets: map[uuid.UUID]*models.Pet{pet.ID: pet}},
	)
	service.(*clinicAppointmentService).now = func() time.Time { return now }

	response, err := service.CreateClinicAppointment(userID.String(), dto.CreateClinicAppointmentRequest{
		PublicPetID:     " pnx-pet-abc123 ",
		AppointmentType: "vaccination",
		ScheduledAt:     now.Add(48 * time.Hour).Format(time.RFC3339),
		DurationMinutes: 45,
	})
	if err != nil {
		t.Fatalf("CreateClinicAppointment() error = %v", err)
	}
	if response.Status != models.AppointmentStatusScheduled || response.CreatedByRole != models.RoleClinic {
		t.Fatalf("unexpected clinic-created response: %#v", response)
	}

	_, err = service.CreateClinicAppointment(userID.String(), dto.CreateClinicAppointmentRequest{
		PetID: pet.ID.String(), PublicPetID: pet.PublicPetID,
	})
	assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestClinicAppointmentScopeStatusAndCalendarFilters(t *testing.T) {
	now := time.Date(2026, 7, 8, 8, 0, 0, 0, time.UTC)
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	appointment := &models.Appointment{
		ID: uuid.New(), ClinicProfileID: clinic.ID, OwnerProfileID: uuid.New(), PetID: uuid.New(),
		Status: models.AppointmentStatusRequested, ScheduledAt: now.Add(24 * time.Hour),
		CreatedAt: now, UpdatedAt: now,
	}
	repo := newAppointmentRepositoryStub()
	repo.appointments[appointment.ID] = appointment
	service := NewClinicAppointmentService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true}, &petRepositoryFake{})
	service.(*clinicAppointmentService).now = func() time.Time { return now }

	updated, err := service.UpdateClinicAppointmentStatus(
		userID.String(),
		appointment.ID,
		dto.UpdateAppointmentStatusRequest{Status: "checked_in"},
	)
	if err != nil || updated.Status != models.AppointmentStatusCheckedIn {
		t.Fatalf("UpdateClinicAppointmentStatus() = %#v, %v", updated, err)
	}

	_, err = service.ListClinicAppointments(userID.String(), dto.ClinicAppointmentFilters{
		DateFrom: "2026-07-08", DateTo: "2026-07-10", Status: "checked_in", AppointmentType: "checkup",
	})
	if err != nil {
		t.Fatalf("ListClinicAppointments() error = %v", err)
	}
	if repo.lastFilters.DateFrom == nil || repo.lastFilters.DateTo == nil ||
		repo.lastFilters.Status != "checked_in" || repo.lastFilters.AppointmentType != "checkup" {
		t.Fatalf("normalized filters = %#v", repo.lastFilters)
	}

	_, err = service.ListClinicAppointments(userID.String(), dto.ClinicAppointmentFilters{
		Date: "2026-07-08", DateFrom: "2026-07-08",
	})
	assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")

	otherAppointmentID := uuid.New()
	_, err = service.GetClinicAppointment(userID.String(), otherAppointmentID)
	assertAppointmentServiceError(t, err, http.StatusNotFound, "APPOINTMENT_NOT_FOUND")
}

func TestAppointmentInputValidation(t *testing.T) {
	now := time.Date(2026, 7, 8, 8, 0, 0, 0, time.UTC)
	tests := []struct {
		name            string
		appointmentType string
		scheduledAt     string
		duration        int
	}{
		{name: "invalid type", appointmentType: "surgery", scheduledAt: now.Add(time.Hour).Format(time.RFC3339), duration: 30},
		{name: "invalid datetime", appointmentType: "checkup", scheduledAt: "tomorrow", duration: 30},
		{name: "past", appointmentType: "checkup", scheduledAt: now.Add(-time.Hour).Format(time.RFC3339), duration: 30},
		{name: "duration too short", appointmentType: "checkup", scheduledAt: now.Add(time.Hour).Format(time.RFC3339), duration: 4},
		{name: "duration too long", appointmentType: "checkup", scheduledAt: now.Add(time.Hour).Format(time.RFC3339), duration: 481},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := normalizeAppointmentInput("", tt.appointmentType, tt.scheduledAt, tt.duration, "", now)
			assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
		})
	}
}

func assertAppointmentServiceError(t *testing.T, err error, status int, code string) {
	t.Helper()
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != status || appErr.Code != code {
		t.Fatalf("status/code = %d/%s, want %d/%s", appErr.HTTPStatus, appErr.Code, status, code)
	}
}

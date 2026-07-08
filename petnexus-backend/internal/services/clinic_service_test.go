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

type clinicProfileRepositoryStub struct {
	profile *models.ClinicProfile
	exists  bool
}

func (r *clinicProfileRepositoryStub) Create(profile *models.ClinicProfile) error {
	if r.exists {
		return repositories.ErrClinicProfileAlreadyExists
	}
	profile.ID = uuid.New()
	profile.CreatedAt = time.Date(2026, 7, 5, 1, 2, 3, 0, time.UTC)
	profile.UpdatedAt = profile.CreatedAt
	r.profile = profile
	r.exists = true
	return nil
}

func (r *clinicProfileRepositoryStub) FindByUserID(userID uuid.UUID) (*models.ClinicProfile, error) {
	if r.profile == nil || r.profile.UserID != userID {
		return nil, repositories.ErrClinicProfileNotFound
	}
	return r.profile, nil
}

func (r *clinicProfileRepositoryStub) FindByID(id uuid.UUID) (*models.ClinicProfile, error) {
	if r.profile == nil || r.profile.ID != id {
		return nil, repositories.ErrClinicProfileNotFound
	}
	return r.profile, nil
}

func (r *clinicProfileRepositoryStub) ExistsByUserID(userID uuid.UUID) (bool, error) {
	return r.exists, nil
}

func (r *clinicProfileRepositoryStub) Update(profile *models.ClinicProfile) error {
	if r.profile == nil {
		return repositories.ErrClinicProfileNotFound
	}
	profile.UpdatedAt = profile.UpdatedAt.Add(time.Minute)
	r.profile = profile
	return nil
}

func TestClinicProfileServiceCreateUsesAuthenticatedUserAndNormalizesFields(t *testing.T) {
	userID := uuid.New()
	repo := &clinicProfileRepositoryStub{}
	service := NewClinicProfileService(repo)

	response, err := service.CreateClinicProfile(userID.String(), dto.CreateClinicProfileRequest{
		ClinicName:  "  Happy Paws Clinic  ",
		PhoneNumber: " 02-123-4567 ",
		Email:       " CLINIC@EXAMPLE.COM ",
		Address:     " 123 Pet Street ",
	})
	if err != nil {
		t.Fatalf("CreateClinicProfile() error = %v", err)
	}
	if repo.profile.UserID != userID {
		t.Fatalf("created UserID = %s, want authenticated user %s", repo.profile.UserID, userID)
	}
	if response.ClinicName != "Happy Paws Clinic" {
		t.Fatalf("ClinicName = %q", response.ClinicName)
	}
	if response.Email == nil || *response.Email != "clinic@example.com" {
		t.Fatalf("Email = %v, want normalized email", response.Email)
	}
}

func TestClinicProfileServiceCreateDuplicateReturnsConflict(t *testing.T) {
	service := NewClinicProfileService(&clinicProfileRepositoryStub{exists: true})
	_, err := service.CreateClinicProfile(uuid.NewString(), dto.CreateClinicProfileRequest{ClinicName: "Happy Paws"})
	assertClinicAppError(t, err, http.StatusConflict, "CLINIC_PROFILE_ALREADY_EXISTS")
}

func TestClinicProfileServiceGetMissingReturnsNotFound(t *testing.T) {
	service := NewClinicProfileService(&clinicProfileRepositoryStub{})
	_, err := service.GetMyClinicProfile(uuid.NewString())
	assertClinicAppError(t, err, http.StatusNotFound, "CLINIC_PROFILE_NOT_FOUND")
}

func TestClinicProfileServicePatchUpdatesOnlyProvidedFields(t *testing.T) {
	userID := uuid.New()
	phone := "02-111-1111"
	email := "old@example.com"
	profile := &models.ClinicProfile{
		ID:          uuid.New(),
		UserID:      userID,
		ClinicName:  "Happy Paws",
		PhoneNumber: &phone,
		Email:       &email,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	repo := &clinicProfileRepositoryStub{profile: profile, exists: true}
	service := NewClinicProfileService(repo)
	newName := " Happy Paws Bangkok "
	clearPhone := ""

	response, err := service.UpdateMyClinicProfile(userID.String(), dto.UpdateClinicProfileRequest{
		ClinicName:  &newName,
		PhoneNumber: &clearPhone,
	})
	if err != nil {
		t.Fatalf("UpdateMyClinicProfile() error = %v", err)
	}
	if response.ClinicName != "Happy Paws Bangkok" || response.PhoneNumber != nil {
		t.Fatalf("updated response = %#v", response)
	}
	if response.Email == nil || *response.Email != "old@example.com" {
		t.Fatalf("unspecified email changed: %#v", response.Email)
	}
}

func TestClinicProfileServiceRejectsEmptyPatchAndInvalidInput(t *testing.T) {
	service := NewClinicProfileService(&clinicProfileRepositoryStub{})
	_, err := service.UpdateMyClinicProfile(uuid.NewString(), dto.UpdateClinicProfileRequest{})
	assertClinicAppError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")

	_, err = service.CreateClinicProfile(uuid.NewString(), dto.CreateClinicProfileRequest{ClinicName: "  "})
	assertClinicAppError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")

	_, err = service.CreateClinicProfile(uuid.NewString(), dto.CreateClinicProfileRequest{
		ClinicName: "Happy Paws",
		Email:      "not-an-email",
	})
	assertClinicAppError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
}

func assertClinicAppError(t *testing.T, err error, status int, code string) {
	t.Helper()
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != status || appErr.Code != code {
		t.Fatalf("error status/code = %d/%s, want %d/%s", appErr.HTTPStatus, appErr.Code, status, code)
	}
}

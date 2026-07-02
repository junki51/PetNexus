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

type ownerProfileRepositoryStub struct {
	profile *models.OwnerProfile
	exists  bool
}

func (r *ownerProfileRepositoryStub) Create(profile *models.OwnerProfile) error {
	if r.exists {
		return repositories.ErrOwnerProfileAlreadyExists
	}
	profile.ID = uuid.New()
	profile.CreatedAt = time.Date(2026, 7, 2, 1, 2, 3, 0, time.UTC)
	profile.UpdatedAt = profile.CreatedAt
	r.profile = profile
	r.exists = true
	return nil
}

func (r *ownerProfileRepositoryStub) FindByUserID(userID uuid.UUID) (*models.OwnerProfile, error) {
	if r.profile == nil || r.profile.UserID != userID {
		return nil, repositories.ErrOwnerProfileNotFound
	}
	return r.profile, nil
}

func (r *ownerProfileRepositoryStub) ExistsByUserID(userID uuid.UUID) (bool, error) {
	return r.exists, nil
}

func (r *ownerProfileRepositoryStub) Update(profile *models.OwnerProfile) error {
	if r.profile == nil {
		return repositories.ErrOwnerProfileNotFound
	}
	profile.UpdatedAt = profile.UpdatedAt.Add(time.Minute)
	r.profile = profile
	return nil
}

func TestOwnerProfileServiceCreateUsesAuthenticatedUserAndNormalizesFields(t *testing.T) {
	userID := uuid.New()
	repo := &ownerProfileRepositoryStub{}
	service := NewOwnerProfileService(repo)

	response, err := service.CreateProfile(userID.String(), dto.CreateOwnerProfileRequest{
		FirstName:   "  Sunny  ",
		LastName:    " Example ",
		Gender:      " MALE ",
		DateOfBirth: "2008-01-01",
		PhoneNumber: " 0812345678 ",
		AvatarURL:   " https://example.com/avatar.png ",
	})
	if err != nil {
		t.Fatalf("CreateProfile() error = %v", err)
	}
	if repo.profile.UserID != userID {
		t.Fatalf("created UserID = %s, want authenticated user %s", repo.profile.UserID, userID)
	}
	if response.FirstName != "Sunny" || response.LastName != "Example" || response.DisplayName != "Sunny Example" {
		t.Fatalf("unexpected normalized response: %#v", response)
	}
	if response.Gender == nil || *response.Gender != "male" {
		t.Fatalf("Gender = %v, want male", response.Gender)
	}
	if response.DateOfBirth == nil || *response.DateOfBirth != "2008-01-01" {
		t.Fatalf("DateOfBirth = %v, want 2008-01-01", response.DateOfBirth)
	}
}

func TestOwnerProfileServiceCreateDuplicateReturnsConflict(t *testing.T) {
	service := NewOwnerProfileService(&ownerProfileRepositoryStub{exists: true})
	_, err := service.CreateProfile(uuid.NewString(), validCreateOwnerProfileRequest())
	assertOwnerAppErrorStatus(t, err, http.StatusConflict)
}

func TestOwnerProfileServiceGetMissingReturnsNotFound(t *testing.T) {
	service := NewOwnerProfileService(&ownerProfileRepositoryStub{})
	_, err := service.GetProfile(uuid.NewString())
	assertOwnerAppErrorStatus(t, err, http.StatusNotFound)
}

func TestOwnerProfileServicePatchUpdatesOnlyProvidedFields(t *testing.T) {
	userID := uuid.New()
	gender := "male"
	profile := &models.OwnerProfile{
		ID:          uuid.New(),
		UserID:      userID,
		FirstName:   "Sunny",
		LastName:    "Example",
		Gender:      &gender,
		PhoneNumber: "0812345678",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	repo := &ownerProfileRepositoryStub{profile: profile, exists: true}
	service := NewOwnerProfileService(repo)
	newFirstName := " Sunny Updated "
	newPhone := " 0899999999 "

	response, err := service.UpdateProfile(userID.String(), dto.UpdateOwnerProfileRequest{
		FirstName:   &newFirstName,
		PhoneNumber: &newPhone,
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if response.FirstName != "Sunny Updated" || response.PhoneNumber != "0899999999" {
		t.Fatalf("updated response = %#v", response)
	}
	if response.LastName != "Example" || response.Gender == nil || *response.Gender != "male" {
		t.Fatalf("unspecified fields were changed: %#v", response)
	}
}

func TestOwnerProfileServiceRejectsEmptyPatch(t *testing.T) {
	service := NewOwnerProfileService(&ownerProfileRepositoryStub{})
	_, err := service.UpdateProfile(uuid.NewString(), dto.UpdateOwnerProfileRequest{})
	assertOwnerAppErrorStatus(t, err, http.StatusBadRequest)
}

func TestBuildOwnerProfileRejectsFutureDateAndInvalidGender(t *testing.T) {
	now := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	req := validCreateOwnerProfileRequest()
	req.DateOfBirth = "2026-07-03"
	_, err := buildOwnerProfile(uuid.New(), req, now)
	assertOwnerAppErrorStatus(t, err, http.StatusBadRequest)

	req.DateOfBirth = "2008-01-01"
	req.Gender = "unknown"
	_, err = buildOwnerProfile(uuid.New(), req, now)
	assertOwnerAppErrorStatus(t, err, http.StatusBadRequest)
}

func validCreateOwnerProfileRequest() dto.CreateOwnerProfileRequest {
	return dto.CreateOwnerProfileRequest{
		FirstName:   "Sunny",
		LastName:    "Example",
		PhoneNumber: "0812345678",
	}
}

func assertOwnerAppErrorStatus(t *testing.T, err error, want int) {
	t.Helper()
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != want {
		t.Fatalf("HTTPStatus = %d, want %d", appErr.HTTPStatus, want)
	}
}

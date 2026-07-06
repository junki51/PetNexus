package services

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

func TestClinicLookupByPublicPetIDReturnsLimitedMaskedData(t *testing.T) {
	petID := uuid.New()
	breedID := uuid.New()
	gender := "male"
	dateOfBirth := time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)
	avatarURL := "https://example.com/milo.png"
	repo := &petRepositoryFake{pets: map[uuid.UUID]*models.Pet{
		petID: {
			ID:          petID,
			PublicPetID: "PNX-PET-8F3K2A",
			Name:        "Milo",
			Species:     "dog",
			Gender:      &gender,
			DateOfBirth: &dateOfBirth,
			AvatarURL:   &avatarURL,
			Breed:       &models.Breed{ID: breedID, Species: "dog", Name: "Golden Retriever"},
			OwnerProfile: &models.OwnerProfile{
				FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678",
			},
		},
	}}
	service := NewClinicPetLookupService(repo)

	result, err := service.LookupPetForClinic(dto.ClinicPetLookupQuery{PetID: " pnx-pet-8f3k2a "})
	if err != nil {
		t.Fatalf("LookupPetForClinic() error = %v", err)
	}
	response, ok := result.(*dto.ClinicPetLookupItemResponse)
	if !ok {
		t.Fatalf("result type = %T", result)
	}
	if response.PublicPetID != "PNX-PET-8F3K2A" || response.Owner.DisplayName != "Sunny Example" || response.Owner.MaskedPhone != "081****678" {
		t.Fatalf("unexpected lookup response: %#v", response)
	}
}

func TestClinicLookupByOwnerPhoneIsExactAndReturnsEmptyList(t *testing.T) {
	petID := uuid.New()
	repo := &petRepositoryFake{pets: map[uuid.UUID]*models.Pet{
		petID: {
			ID: petID, PublicPetID: "PNX-PET-ABC123", Name: "Milo", Species: "dog",
			OwnerProfile: &models.OwnerProfile{FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678"},
		},
	}}
	service := NewClinicPetLookupService(repo)

	result, err := service.LookupPetForClinic(dto.ClinicPetLookupQuery{OwnerPhone: " 0812345678 "})
	if err != nil {
		t.Fatalf("exact phone lookup error = %v", err)
	}
	list := result.(*dto.ClinicPetLookupListResponse)
	if len(list.Items) != 1 || list.Items[0].PublicPetID != "PNX-PET-ABC123" {
		t.Fatalf("exact phone lookup result = %#v", list)
	}

	result, err = service.LookupPetForClinic(dto.ClinicPetLookupQuery{OwnerPhone: "081"})
	if err != nil {
		t.Fatalf("partial phone lookup error = %v", err)
	}
	if len(result.(*dto.ClinicPetLookupListResponse).Items) != 0 {
		t.Fatal("partial phone must not return pets")
	}
}

func TestClinicLookupValidatesQueryAndUnknownPet(t *testing.T) {
	service := NewClinicPetLookupService(&petRepositoryFake{pets: make(map[uuid.UUID]*models.Pet)})

	_, err := service.LookupPetForClinic(dto.ClinicPetLookupQuery{})
	assertLookupAppError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")

	_, err = service.LookupPetForClinic(dto.ClinicPetLookupQuery{PetID: "PNX-PET-ABC123", OwnerPhone: "0812345678"})
	assertLookupAppError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")

	_, err = service.LookupPetForClinic(dto.ClinicPetLookupQuery{PetID: "PNX-PET-UNKNOWN"})
	assertLookupAppError(t, err, http.StatusNotFound, "PET_NOT_FOUND")
}

func TestMaskOwnerPhone(t *testing.T) {
	if got := maskOwnerPhone("0812345678"); got != "081****678" {
		t.Fatalf("maskOwnerPhone() = %q", got)
	}
	if got := maskOwnerPhone("12345"); got != "*****" {
		t.Fatalf("short maskOwnerPhone() = %q", got)
	}
}

func assertLookupAppError(t *testing.T, err error, status int, code string) {
	t.Helper()
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != status || appErr.Code != code {
		t.Fatalf("status/code = %d/%s, want %d/%s", appErr.HTTPStatus, appErr.Code, status, code)
	}
}

package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

// ClinicPetLookupService provides a privacy-limited pet lookup for clinics.
type ClinicPetLookupService interface {
	LookupPetForClinic(query dto.ClinicPetLookupQuery) (any, error)
}

type clinicPetLookupService struct {
	petRepo repositories.PetRepository
}

func NewClinicPetLookupService(petRepo repositories.PetRepository) ClinicPetLookupService {
	return &clinicPetLookupService{petRepo: petRepo}
}

func (s *clinicPetLookupService) LookupPetForClinic(query dto.ClinicPetLookupQuery) (any, error) {
	publicPetID := strings.ToUpper(strings.TrimSpace(query.PetID))
	ownerPhone := strings.TrimSpace(query.OwnerPhone)

	if publicPetID == "" && ownerPhone == "" {
		return nil, petValidationError("pet_id or owner_phone is required")
	}
	if publicPetID != "" && ownerPhone != "" {
		return nil, petValidationError("provide only one of pet_id or owner_phone")
	}

	if publicPetID != "" {
		pet, err := s.petRepo.FindByPublicPetID(publicPetID)
		if err != nil {
			if errors.Is(err, repositories.ErrPetNotFound) {
				return nil, clinicLookupPetNotFoundError(err)
			}
			return nil, internalServerError(err)
		}
		response := toClinicPetLookupItem(pet)
		return &response, nil
	}

	pets, err := s.petRepo.FindByOwnerPhone(ownerPhone)
	if err != nil {
		return nil, internalServerError(err)
	}
	items := make([]dto.ClinicPetLookupItemResponse, 0, len(pets))
	for i := range pets {
		items = append(items, toClinicPetLookupItem(&pets[i]))
	}
	return &dto.ClinicPetLookupListResponse{Items: items}, nil
}

func toClinicPetLookupItem(pet *models.Pet) dto.ClinicPetLookupItemResponse {
	var breed *dto.BreedResponse
	if pet.Breed != nil {
		mapped := toBreedResponse(pet.Breed)
		breed = &mapped
	}

	var dateOfBirth *string
	if pet.DateOfBirth != nil {
		formatted := pet.DateOfBirth.Format(dateOnlyLayout)
		dateOfBirth = &formatted
	}

	owner := dto.ClinicPetLookupOwnerResponse{}
	if pet.OwnerProfile != nil {
		owner.DisplayName = strings.TrimSpace(pet.OwnerProfile.FirstName + " " + pet.OwnerProfile.LastName)
		owner.MaskedPhone = maskOwnerPhone(pet.OwnerProfile.PhoneNumber)
	}

	return dto.ClinicPetLookupItemResponse{
		ID:          pet.ID.String(),
		PublicPetID: pet.PublicPetID,
		Name:        pet.Name,
		Species:     pet.Species,
		Breed:       breed,
		Gender:      pet.Gender,
		DateOfBirth: dateOfBirth,
		AvatarURL:   pet.AvatarURL,
		Owner:       owner,
	}
}

func maskOwnerPhone(phone string) string {
	trimmed := strings.TrimSpace(phone)
	runes := []rune(trimmed)
	if len(runes) <= 6 {
		return strings.Repeat("*", len(runes))
	}
	return string(runes[:3]) + "****" + string(runes[len(runes)-3:])
}

func clinicLookupPetNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"PET_NOT_FOUND",
		"Pet not found",
		"No pet was found for the provided public pet ID",
		cause,
	)
}

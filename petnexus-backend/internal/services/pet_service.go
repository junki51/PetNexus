package services

import (
	"errors"
	"math"
	"net/http"
	"net/url"
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

var allowedPetGenders = map[string]struct{}{
	models.PetGenderMale:    {},
	models.PetGenderFemale:  {},
	models.PetGenderUnknown: {},
}

// PetService owns pet validation, ownership, and breed consistency rules.
type PetService interface {
	CreatePet(currentUserID string, req dto.CreatePetRequest) (*dto.PetResponse, error)
	ListMyPets(currentUserID string) ([]dto.PetResponse, error)
	GetMyPet(currentUserID string, petID uuid.UUID) (*dto.PetResponse, error)
	UpdateMyPet(currentUserID string, petID uuid.UUID, req dto.UpdatePetRequest) (*dto.PetResponse, error)
}

type petService struct {
	petRepo          repositories.PetRepository
	breedRepo        repositories.BreedRepository
	ownerProfileRepo repositories.OwnerProfileRepository
}

func NewPetService(
	petRepo repositories.PetRepository,
	breedRepo repositories.BreedRepository,
	ownerProfileRepo repositories.OwnerProfileRepository,
) PetService {
	return &petService{
		petRepo:          petRepo,
		breedRepo:        breedRepo,
		ownerProfileRepo: ownerProfileRepo,
	}
}

func (s *petService) CreatePet(currentUserID string, req dto.CreatePetRequest) (*dto.PetResponse, error) {
	ownerProfile, err := s.findCurrentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}

	species, err := normalizePetSpecies(req.Species, true)
	if err != nil {
		return nil, err
	}
	name, err := normalizeRequiredPetField("name", req.Name, 100)
	if err != nil {
		return nil, err
	}
	breedID, breed, err := s.resolveBreed(req.BreedID, species)
	if err != nil {
		return nil, err
	}
	gender, err := normalizePetGender(req.Gender)
	if err != nil {
		return nil, err
	}
	dateOfBirth, err := parsePetDateOfBirth(req.DateOfBirth, time.Now())
	if err != nil {
		return nil, err
	}
	if err := validatePetWeight(req.WeightKG); err != nil {
		return nil, err
	}
	microchipID, err := normalizePetOptionalField("microchip_id", req.MicrochipID, 100)
	if err != nil {
		return nil, err
	}
	avatarURL, err := normalizePetAvatarURL(req.AvatarURL)
	if err != nil {
		return nil, err
	}
	color, err := normalizePetOptionalField("color", req.Color, 100)
	if err != nil {
		return nil, err
	}
	distinctiveMarks, err := normalizePetOptionalField("distinctive_marks", req.DistinctiveMarks, 1000)
	if err != nil {
		return nil, err
	}

	pet := &models.Pet{
		OwnerProfileID:   ownerProfile.ID,
		BreedID:          breedID,
		Species:          species,
		Name:             name,
		Gender:           gender,
		DateOfBirth:      dateOfBirth,
		WeightKG:         req.WeightKG,
		MicrochipID:      microchipID,
		AvatarURL:        avatarURL,
		Color:            color,
		DistinctiveMarks: distinctiveMarks,
		IsNeutered:       req.IsNeutered,
		Breed:            breed,
	}
	if err := s.petRepo.Create(pet); err != nil {
		return nil, internalServerError(err)
	}

	response := toPetResponse(pet, time.Now())
	return &response, nil
}

func (s *petService) ListMyPets(currentUserID string) ([]dto.PetResponse, error) {
	ownerProfile, err := s.findCurrentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	pets, err := s.petRepo.FindAllByOwnerProfileID(ownerProfile.ID)
	if err != nil {
		return nil, internalServerError(err)
	}

	now := time.Now()
	response := make([]dto.PetResponse, 0, len(pets))
	for i := range pets {
		response = append(response, toPetResponse(&pets[i], now))
	}
	return response, nil
}

func (s *petService) GetMyPet(currentUserID string, petID uuid.UUID) (*dto.PetResponse, error) {
	ownerProfile, err := s.findCurrentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	pet, err := s.petRepo.FindByIDAndOwnerProfileID(petID, ownerProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrPetNotFound) {
			return nil, petNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toPetResponse(pet, time.Now())
	return &response, nil
}

func (s *petService) UpdateMyPet(currentUserID string, petID uuid.UUID, req dto.UpdatePetRequest) (*dto.PetResponse, error) {
	if !hasPetUpdate(req) {
		return nil, petValidationError("Request body must contain at least one pet field")
	}
	ownerProfile, err := s.findCurrentOwnerProfile(currentUserID)
	if err != nil {
		return nil, err
	}
	pet, err := s.petRepo.FindByIDAndOwnerProfileID(petID, ownerProfile.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrPetNotFound) {
			return nil, petNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	if err := s.applyPetUpdate(pet, req, time.Now()); err != nil {
		return nil, err
	}
	if err := s.petRepo.Update(pet); err != nil {
		if errors.Is(err, repositories.ErrPetNotFound) {
			return nil, petNotFoundError(err)
		}
		return nil, internalServerError(err)
	}

	response := toPetResponse(pet, time.Now())
	return &response, nil
}

func (s *petService) findCurrentOwnerProfile(currentUserID string) (*models.OwnerProfile, error) {
	userID, err := uuid.Parse(currentUserID)
	if err != nil {
		return nil, utils.NewAppError(
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"Unauthorized",
			"Invalid authenticated user",
			err,
		)
	}
	profile, err := s.ownerProfileRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrOwnerProfileNotFound) {
			return nil, ownerProfileRequiredError(err)
		}
		return nil, internalServerError(err)
	}
	return profile, nil
}

func (s *petService) resolveBreed(input *string, species string) (*uuid.UUID, *models.Breed, error) {
	if input == nil || strings.TrimSpace(*input) == "" {
		return nil, nil, nil
	}
	breedID, err := uuid.Parse(strings.TrimSpace(*input))
	if err != nil {
		return nil, nil, petValidationError("breed_id must be a valid UUID")
	}
	breed, err := s.breedRepo.FindByID(breedID)
	if err != nil {
		if errors.Is(err, repositories.ErrBreedNotFound) {
			return nil, nil, breedNotFoundError(err)
		}
		return nil, nil, internalServerError(err)
	}
	if breed.Species != species {
		return nil, nil, breedSpeciesMismatchError()
	}
	return &breedID, breed, nil
}

func (s *petService) applyPetUpdate(pet *models.Pet, req dto.UpdatePetRequest, now time.Time) error {
	effectiveSpecies := pet.Species
	if req.Species != nil {
		species, err := normalizePetSpecies(*req.Species, true)
		if err != nil {
			return err
		}
		effectiveSpecies = species
		pet.Species = species
	}
	if req.Name != nil {
		name, err := normalizeRequiredPetField("name", *req.Name, 100)
		if err != nil {
			return err
		}
		pet.Name = name
	}
	if req.BreedID != nil {
		breedID, breed, err := s.resolveBreed(req.BreedID, effectiveSpecies)
		if err != nil {
			return err
		}
		pet.BreedID = breedID
		pet.Breed = breed
	} else if pet.BreedID != nil {
		breed := pet.Breed
		if breed == nil {
			var err error
			breed, err = s.breedRepo.FindByID(*pet.BreedID)
			if err != nil {
				if errors.Is(err, repositories.ErrBreedNotFound) {
					return breedNotFoundError(err)
				}
				return internalServerError(err)
			}
		}
		if breed.Species != effectiveSpecies {
			return breedSpeciesMismatchError()
		}
		pet.Breed = breed
	}
	if req.Gender != nil {
		gender, err := normalizePetGender(*req.Gender)
		if err != nil {
			return err
		}
		pet.Gender = gender
	}
	if req.DateOfBirth != nil {
		dateOfBirth, err := parsePetDateOfBirth(*req.DateOfBirth, now)
		if err != nil {
			return err
		}
		pet.DateOfBirth = dateOfBirth
	}
	if req.WeightKG != nil {
		if err := validatePetWeight(req.WeightKG); err != nil {
			return err
		}
		pet.WeightKG = req.WeightKG
	}
	if req.MicrochipID != nil {
		value, err := normalizePetOptionalField("microchip_id", *req.MicrochipID, 100)
		if err != nil {
			return err
		}
		pet.MicrochipID = value
	}
	if req.AvatarURL != nil {
		value, err := normalizePetAvatarURL(*req.AvatarURL)
		if err != nil {
			return err
		}
		pet.AvatarURL = value
	}
	if req.Color != nil {
		value, err := normalizePetOptionalField("color", *req.Color, 100)
		if err != nil {
			return err
		}
		pet.Color = value
	}
	if req.DistinctiveMarks != nil {
		value, err := normalizePetOptionalField("distinctive_marks", *req.DistinctiveMarks, 1000)
		if err != nil {
			return err
		}
		pet.DistinctiveMarks = value
	}
	if req.IsNeutered != nil {
		pet.IsNeutered = req.IsNeutered
	}
	return nil
}

func hasPetUpdate(req dto.UpdatePetRequest) bool {
	return req.Species != nil || req.Name != nil || req.BreedID != nil ||
		req.Gender != nil || req.DateOfBirth != nil || req.WeightKG != nil ||
		req.MicrochipID != nil || req.AvatarURL != nil || req.Color != nil ||
		req.DistinctiveMarks != nil || req.IsNeutered != nil
}

func normalizePetSpecies(value string, required bool) (string, error) {
	species := strings.ToLower(strings.TrimSpace(value))
	if species == "" && required {
		return "", petValidationError("species is required")
	}
	if species != models.SpeciesDog && species != models.SpeciesCat {
		return "", petValidationError("species must be dog or cat")
	}
	return species, nil
}

func normalizeRequiredPetField(field, value string, maxLength int) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", petValidationError(field + " is required")
	}
	if utf8.RuneCountInString(trimmed) > maxLength {
		return "", petValidationError(field + " must not exceed " + strconv.Itoa(maxLength) + " characters")
	}
	return trimmed, nil
}

func normalizePetOptionalField(field, value string, maxLength int) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if utf8.RuneCountInString(trimmed) > maxLength {
		return nil, petValidationError(field + " must not exceed " + strconv.Itoa(maxLength) + " characters")
	}
	if trimmed == "" {
		return nil, nil
	}
	return &trimmed, nil
}

func normalizePetGender(value string) (*string, error) {
	gender := strings.ToLower(strings.TrimSpace(value))
	if gender == "" {
		return nil, nil
	}
	if _, allowed := allowedPetGenders[gender]; !allowed {
		return nil, petValidationError("gender must be male, female, or unknown")
	}
	return &gender, nil
}

func parsePetDateOfBirth(value string, now time.Time) (*time.Time, error) {
	dateText := strings.TrimSpace(value)
	if dateText == "" {
		return nil, nil
	}
	date, err := time.Parse(dateOnlyLayout, dateText)
	if err != nil {
		return nil, petValidationError("date_of_birth must use YYYY-MM-DD format")
	}
	today, _ := time.Parse(dateOnlyLayout, now.UTC().Format(dateOnlyLayout))
	if date.After(today) {
		return nil, petValidationError("date_of_birth must not be in the future")
	}
	return &date, nil
}

func validatePetWeight(weight *float64) error {
	if weight == nil {
		return nil
	}
	if math.IsNaN(*weight) || math.IsInf(*weight, 0) || *weight <= 0 || *weight > 200 {
		return petValidationError("weight_kg must be greater than 0 and at most 200")
	}
	return nil
}

func normalizePetAvatarURL(value string) (*string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, petValidationError("avatar_url must be a valid HTTP or HTTPS URL")
	}
	return &trimmed, nil
}

func toPetResponse(pet *models.Pet, now time.Time) dto.PetResponse {
	var dateOfBirth *string
	var ageYears *int
	if pet.DateOfBirth != nil {
		formatted := pet.DateOfBirth.Format(dateOnlyLayout)
		dateOfBirth = &formatted
		age := now.Year() - pet.DateOfBirth.Year()
		birthdayThisYear := time.Date(now.Year(), pet.DateOfBirth.Month(), pet.DateOfBirth.Day(), 0, 0, 0, 0, now.Location())
		if now.Before(birthdayThisYear) {
			age--
		}
		if age >= 0 {
			ageYears = &age
		}
	}

	var breed *dto.BreedResponse
	if pet.Breed != nil {
		mapped := toBreedResponse(pet.Breed)
		breed = &mapped
	}

	return dto.PetResponse{
		ID:               pet.ID.String(),
		Species:          pet.Species,
		Name:             pet.Name,
		Gender:           pet.Gender,
		DateOfBirth:      dateOfBirth,
		AgeYears:         ageYears,
		Breed:            breed,
		WeightKG:         pet.WeightKG,
		MicrochipID:      pet.MicrochipID,
		AvatarURL:        pet.AvatarURL,
		Color:            pet.Color,
		DistinctiveMarks: pet.DistinctiveMarks,
		IsNeutered:       pet.IsNeutered,
		CreatedAt:        pet.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        pet.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func petValidationError(details string) *utils.AppError {
	return utils.NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details, nil)
}

func ownerProfileRequiredError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"OWNER_PROFILE_REQUIRED",
		"Owner profile required",
		"Create an owner profile before managing pets",
		cause,
	)
}

func breedNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(http.StatusNotFound, "BREED_NOT_FOUND", "Breed not found", "The selected breed does not exist", cause)
}

func breedSpeciesMismatchError() *utils.AppError {
	return utils.NewAppError(
		http.StatusBadRequest,
		"BREED_SPECIES_MISMATCH",
		"Breed does not match species",
		"The selected breed species must match the pet species",
		nil,
	)
}

func petNotFoundError(cause error) *utils.AppError {
	return utils.NewAppError(
		http.StatusNotFound,
		"PET_NOT_FOUND",
		"Pet not found",
		"The pet does not exist or does not belong to the authenticated owner",
		cause,
	)
}

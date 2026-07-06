package services

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type breedRepositoryFake struct {
	breeds        map[uuid.UUID]*models.Breed
	lastSpecies   string
	findAllResult []models.Breed
}

func (r *breedRepositoryFake) FindAll(species string) ([]models.Breed, error) {
	r.lastSpecies = species
	return r.findAllResult, nil
}

func (r *breedRepositoryFake) FindByID(id uuid.UUID) (*models.Breed, error) {
	breed, ok := r.breeds[id]
	if !ok {
		return nil, repositories.ErrBreedNotFound
	}
	return breed, nil
}

type petRepositoryFake struct {
	pets map[uuid.UUID]*models.Pet
}

func (r *petRepositoryFake) Create(pet *models.Pet) error {
	pet.ID = uuid.New()
	pet.CreatedAt = time.Date(2026, 7, 4, 1, 2, 3, 0, time.UTC)
	pet.UpdatedAt = pet.CreatedAt
	r.pets[pet.ID] = pet
	return nil
}

func (r *petRepositoryFake) FindByID(id uuid.UUID) (*models.Pet, error) {
	pet, ok := r.pets[id]
	if !ok {
		return nil, repositories.ErrPetNotFound
	}
	return pet, nil
}

func (r *petRepositoryFake) FindByIDAndOwnerProfileID(id, ownerProfileID uuid.UUID) (*models.Pet, error) {
	pet, ok := r.pets[id]
	if !ok || pet.OwnerProfileID != ownerProfileID {
		return nil, repositories.ErrPetNotFound
	}
	return pet, nil
}

func (r *petRepositoryFake) FindAllByOwnerProfileID(ownerProfileID uuid.UUID) ([]models.Pet, error) {
	result := make([]models.Pet, 0)
	for _, pet := range r.pets {
		if pet.OwnerProfileID == ownerProfileID {
			result = append(result, *pet)
		}
	}
	return result, nil
}

func (r *petRepositoryFake) FindByPublicPetID(publicPetID string) (*models.Pet, error) {
	for _, pet := range r.pets {
		if pet.PublicPetID == publicPetID {
			return pet, nil
		}
	}
	return nil, repositories.ErrPetNotFound
}

func (r *petRepositoryFake) FindByOwnerPhone(phone string) ([]models.Pet, error) {
	result := make([]models.Pet, 0)
	for _, pet := range r.pets {
		if pet.OwnerProfile != nil && pet.OwnerProfile.PhoneNumber == phone {
			result = append(result, *pet)
		}
	}
	return result, nil
}

func (r *petRepositoryFake) Update(pet *models.Pet) error {
	stored, ok := r.pets[pet.ID]
	if !ok || stored.OwnerProfileID != pet.OwnerProfileID {
		return repositories.ErrPetNotFound
	}
	pet.UpdatedAt = pet.UpdatedAt.Add(time.Minute)
	r.pets[pet.ID] = pet
	return nil
}

type petOwnerProfileRepositoryFake struct {
	profiles map[uuid.UUID]*models.OwnerProfile
}

func (r *petOwnerProfileRepositoryFake) Create(profile *models.OwnerProfile) error {
	return errors.New("not implemented in pet service test")
}

func (r *petOwnerProfileRepositoryFake) FindByUserID(userID uuid.UUID) (*models.OwnerProfile, error) {
	profile, ok := r.profiles[userID]
	if !ok {
		return nil, repositories.ErrOwnerProfileNotFound
	}
	return profile, nil
}

func (r *petOwnerProfileRepositoryFake) ExistsByUserID(userID uuid.UUID) (bool, error) {
	_, ok := r.profiles[userID]
	return ok, nil
}

func (r *petOwnerProfileRepositoryFake) Update(profile *models.OwnerProfile) error {
	return errors.New("not implemented in pet service test")
}

func TestBreedServiceFiltersAndValidatesSpecies(t *testing.T) {
	repo := &breedRepositoryFake{findAllResult: []models.Breed{{ID: uuid.New(), Species: models.SpeciesDog, Name: "Poodle"}}}
	service := NewBreedService(repo)

	response, err := service.ListBreeds(" DOG ")
	if err != nil {
		t.Fatalf("ListBreeds() error = %v", err)
	}
	if repo.lastSpecies != models.SpeciesDog || len(response) != 1 || response[0].Name != "Poodle" {
		t.Fatalf("unexpected breed result: filter=%q response=%#v", repo.lastSpecies, response)
	}
	_, err = service.ListBreeds("bird")
	assertPetAppErrorStatus(t, err, http.StatusBadRequest)
}

func TestPetServiceCreateUsesCurrentOwnerProfileAndMatchingBreed(t *testing.T) {
	userID := uuid.New()
	profileID := uuid.New()
	breedID := uuid.New()
	breed := &models.Breed{ID: breedID, Species: models.SpeciesDog, Name: "Poodle"}
	petRepo := &petRepositoryFake{pets: make(map[uuid.UUID]*models.Pet)}
	service := NewPetService(
		petRepo,
		&breedRepositoryFake{breeds: map[uuid.UUID]*models.Breed{breedID: breed}},
		&petOwnerProfileRepositoryFake{profiles: map[uuid.UUID]*models.OwnerProfile{userID: {ID: profileID, UserID: userID}}},
	)
	weight := 12.5
	neutered := true
	breedText := breedID.String()

	response, err := service.CreatePet(userID.String(), dto.CreatePetRequest{
		Species:     " DOG ",
		Name:        " Milo ",
		BreedID:     &breedText,
		Gender:      " MALE ",
		DateOfBirth: "2022-05-10",
		WeightKG:    &weight,
		IsNeutered:  &neutered,
	})
	if err != nil {
		t.Fatalf("CreatePet() error = %v", err)
	}
	created := petRepo.pets[uuid.MustParse(response.ID)]
	if created.OwnerProfileID != profileID {
		t.Fatalf("OwnerProfileID = %s, want current owner's profile %s", created.OwnerProfileID, profileID)
	}
	if response.Name != "Milo" || response.Species != models.SpeciesDog || response.Breed == nil || response.Breed.ID != breedID.String() {
		t.Fatalf("unexpected response: %#v", response)
	}
	if !strings.HasPrefix(response.PublicPetID, utils.PublicPetIDPrefix) {
		t.Fatalf("PublicPetID = %q, want generated prefix %q", response.PublicPetID, utils.PublicPetIDPrefix)
	}
}

func TestPetServiceRequiresOwnerProfile(t *testing.T) {
	service := NewPetService(
		&petRepositoryFake{pets: make(map[uuid.UUID]*models.Pet)},
		&breedRepositoryFake{},
		&petOwnerProfileRepositoryFake{profiles: make(map[uuid.UUID]*models.OwnerProfile)},
	)
	_, err := service.CreatePet(uuid.NewString(), dto.CreatePetRequest{Species: "dog", Name: "Milo"})
	assertPetAppErrorCode(t, err, http.StatusNotFound, "OWNER_PROFILE_REQUIRED")
}

func TestPetServiceRejectsBreedSpeciesMismatch(t *testing.T) {
	userID := uuid.New()
	breedID := uuid.New()
	breedText := breedID.String()
	service := NewPetService(
		&petRepositoryFake{pets: make(map[uuid.UUID]*models.Pet)},
		&breedRepositoryFake{breeds: map[uuid.UUID]*models.Breed{
			breedID: {ID: breedID, Species: models.SpeciesCat, Name: "Persian"},
		}},
		&petOwnerProfileRepositoryFake{profiles: map[uuid.UUID]*models.OwnerProfile{
			userID: {ID: uuid.New(), UserID: userID},
		}},
	)
	_, err := service.CreatePet(userID.String(), dto.CreatePetRequest{
		Species: "dog",
		Name:    "Milo",
		BreedID: &breedText,
	})
	assertPetAppErrorCode(t, err, http.StatusBadRequest, "BREED_SPECIES_MISMATCH")
}

func TestPetServiceHidesAnotherOwnersPetAsNotFound(t *testing.T) {
	ownerUserID := uuid.New()
	otherUserID := uuid.New()
	ownerProfileID := uuid.New()
	otherProfileID := uuid.New()
	petID := uuid.New()
	service := NewPetService(
		&petRepositoryFake{pets: map[uuid.UUID]*models.Pet{
			petID: {ID: petID, OwnerProfileID: ownerProfileID, Species: "dog", Name: "Milo"},
		}},
		&breedRepositoryFake{},
		&petOwnerProfileRepositoryFake{profiles: map[uuid.UUID]*models.OwnerProfile{
			ownerUserID: {ID: ownerProfileID, UserID: ownerUserID},
			otherUserID: {ID: otherProfileID, UserID: otherUserID},
		}},
	)
	_, err := service.GetMyPet(otherUserID.String(), petID)
	assertPetAppErrorCode(t, err, http.StatusNotFound, "PET_NOT_FOUND")
}

func TestPetServicePatchIsPartialAndValidatesExistingBreedAgainstNewSpecies(t *testing.T) {
	userID := uuid.New()
	profileID := uuid.New()
	petID := uuid.New()
	breedID := uuid.New()
	gender := "male"
	breed := &models.Breed{ID: breedID, Species: "dog", Name: "Poodle"}
	pet := &models.Pet{
		ID:             petID,
		OwnerProfileID: profileID,
		BreedID:        &breedID,
		Breed:          breed,
		Species:        "dog",
		Name:           "Milo",
		Gender:         &gender,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	petRepo := &petRepositoryFake{pets: map[uuid.UUID]*models.Pet{petID: pet}}
	service := NewPetService(
		petRepo,
		&breedRepositoryFake{breeds: map[uuid.UUID]*models.Breed{breedID: breed}},
		&petOwnerProfileRepositoryFake{profiles: map[uuid.UUID]*models.OwnerProfile{userID: {ID: profileID, UserID: userID}}},
	)
	updatedName := "Milo Updated"
	weight := 13.2
	response, err := service.UpdateMyPet(userID.String(), petID, dto.UpdatePetRequest{Name: &updatedName, WeightKG: &weight})
	if err != nil {
		t.Fatalf("UpdateMyPet() error = %v", err)
	}
	if response.Name != updatedName || response.WeightKG == nil || *response.WeightKG != weight || response.Gender == nil || *response.Gender != gender {
		t.Fatalf("unexpected partial update response: %#v", response)
	}

	cat := "cat"
	_, err = service.UpdateMyPet(userID.String(), petID, dto.UpdatePetRequest{Species: &cat})
	assertPetAppErrorCode(t, err, http.StatusBadRequest, "BREED_SPECIES_MISMATCH")

	emptyBreed := ""
	response, err = service.UpdateMyPet(userID.String(), petID, dto.UpdatePetRequest{Species: &cat, BreedID: &emptyBreed})
	if err != nil {
		t.Fatalf("UpdateMyPet(clear breed) error = %v", err)
	}
	if response.Species != "cat" || response.Breed != nil {
		t.Fatalf("cleared breed response = %#v", response)
	}
}

func TestPetServiceRejectsEmptyPatchAndInvalidWeight(t *testing.T) {
	service := NewPetService(nil, nil, nil)
	_, err := service.UpdateMyPet(uuid.NewString(), uuid.New(), dto.UpdatePetRequest{})
	assertPetAppErrorStatus(t, err, http.StatusBadRequest)

	weight := 201.0
	err = validatePetWeight(&weight)
	assertPetAppErrorStatus(t, err, http.StatusBadRequest)
}

func assertPetAppErrorStatus(t *testing.T, err error, status int) {
	t.Helper()
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != status {
		t.Fatalf("HTTPStatus = %d, want %d", appErr.HTTPStatus, status)
	}
}

func assertPetAppErrorCode(t *testing.T, err error, status int, code string) {
	t.Helper()
	var appErr *utils.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want *utils.AppError", err)
	}
	if appErr.HTTPStatus != status || appErr.Code != code {
		t.Fatalf("error status/code = %d/%s, want %d/%s", appErr.HTTPStatus, appErr.Code, status, code)
	}
}

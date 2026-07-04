package services

import (
	"strings"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
)

// BreedService owns breed-list filtering and response mapping.
type BreedService interface {
	ListBreeds(species string) ([]dto.BreedResponse, error)
}

type breedService struct {
	breedRepo repositories.BreedRepository
}

func NewBreedService(breedRepo repositories.BreedRepository) BreedService {
	return &breedService{breedRepo: breedRepo}
}

func (s *breedService) ListBreeds(species string) ([]dto.BreedResponse, error) {
	filter := strings.ToLower(strings.TrimSpace(species))
	if filter != "" && filter != models.SpeciesDog && filter != models.SpeciesCat {
		return nil, petValidationError("species query must be dog or cat")
	}

	breeds, err := s.breedRepo.FindAll(filter)
	if err != nil {
		return nil, internalServerError(err)
	}

	response := make([]dto.BreedResponse, 0, len(breeds))
	for i := range breeds {
		response = append(response, toBreedResponse(&breeds[i]))
	}
	return response, nil
}

func toBreedResponse(breed *models.Breed) dto.BreedResponse {
	return dto.BreedResponse{
		ID:      breed.ID.String(),
		Species: breed.Species,
		Name:    breed.Name,
		NameTH:  breed.NameTH,
	}
}

package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/handlers"
	"github.com/phonlakitz/petnexus-backend/internal/models"
)

type breedServiceSpy struct {
	species string
}

func (s *breedServiceSpy) ListBreeds(species string) ([]dto.BreedResponse, error) {
	s.species = species
	return []dto.BreedResponse{{ID: uuid.NewString(), Species: "dog", Name: "Poodle"}}, nil
}

type petServiceSpy struct {
	currentUserID string
	createRequest dto.CreatePetRequest
}

func (s *petServiceSpy) CreatePet(userID string, req dto.CreatePetRequest) (*dto.PetResponse, error) {
	s.currentUserID = userID
	s.createRequest = req
	return sprint5PetResponse(), nil
}

func (s *petServiceSpy) ListMyPets(userID string) ([]dto.PetResponse, error) {
	s.currentUserID = userID
	return []dto.PetResponse{*sprint5PetResponse()}, nil
}

func (s *petServiceSpy) GetMyPet(userID string, petID uuid.UUID) (*dto.PetResponse, error) {
	s.currentUserID = userID
	return sprint5PetResponse(), nil
}

func (s *petServiceSpy) UpdateMyPet(userID string, petID uuid.UUID, req dto.UpdatePetRequest) (*dto.PetResponse, error) {
	s.currentUserID = userID
	return sprint5PetResponse(), nil
}

func TestBreedRouteIsPublicAndPassesSpeciesFilter(t *testing.T) {
	router, breedSpy, _ := newSprint5RouteTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/breeds?species=dog", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if breedSpy.species != "dog" {
		t.Fatalf("species filter = %q, want dog", breedSpy.species)
	}
}

func TestPetRoutesRequireAuthenticationAndOwnerRole(t *testing.T) {
	router, _, _ := newSprint5RouteTestRouter()

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/pets", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	clinicToken := ownerRouteTestToken(t, uuid.NewString(), models.RoleClinicStaff)
	clinicRequest := httptest.NewRequest(http.MethodGet, "/api/pets", nil)
	clinicRequest.Header.Set("Authorization", "Bearer "+clinicToken)
	clinicResponse := httptest.NewRecorder()
	router.ServeHTTP(clinicResponse, clinicRequest)
	if clinicResponse.Code != http.StatusForbidden {
		t.Fatalf("clinic status = %d, want 403", clinicResponse.Code)
	}
}

func TestCreatePetUsesJWTIdentityAndDoesNotExposeOwnerIDs(t *testing.T) {
	router, _, petSpy := newSprint5RouteTestRouter()
	userID := uuid.NewString()
	token := ownerRouteTestToken(t, userID, models.RoleOwner)
	body := `{
		"user_id":"` + uuid.NewString() + `",
		"owner_profile_id":"` + uuid.NewString() + `",
		"species":"dog",
		"name":"Milo"
	}`
	request := httptest.NewRequest(http.MethodPost, "/api/pets", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", response.Code, response.Body.String())
	}
	if petSpy.currentUserID != userID {
		t.Fatalf("service user ID = %q, want JWT user ID %q", petSpy.currentUserID, userID)
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := payload["data"].(map[string]any)
	if _, found := data["user_id"]; found {
		t.Fatal("response must not expose user_id")
	}
	if _, found := data["owner_profile_id"]; found {
		t.Fatal("response must not expose owner_profile_id")
	}
}

func TestPetDetailRejectsInvalidUUID(t *testing.T) {
	router, _, _ := newSprint5RouteTestRouter()
	token := ownerRouteTestToken(t, uuid.NewString(), models.RoleOwner)
	request := httptest.NewRequest(http.MethodGet, "/api/pets/not-a-uuid", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", response.Code, response.Body.String())
	}
}

func newSprint5RouteTestRouter() (*gin.Engine, *breedServiceSpy, *petServiceSpy) {
	gin.SetMode(gin.TestMode)
	breedSpy := &breedServiceSpy{}
	petSpy := &petServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config:       config.Config{JWTSecret: "owner-route-test-secret", JWTExpiresIn: "1h"},
		BreedHandler: handlers.NewBreedHandler(breedSpy),
		PetHandler:   handlers.NewPetHandler(petSpy),
	})
	return router, breedSpy, petSpy
}

func sprint5PetResponse() *dto.PetResponse {
	return &dto.PetResponse{
		ID:        uuid.NewString(),
		Species:   "dog",
		Name:      "Milo",
		CreatedAt: "2026-07-04T01:02:03Z",
		UpdatedAt: "2026-07-04T01:02:03Z",
	}
}

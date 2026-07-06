package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/handlers"
	"github.com/phonlakitz/petnexus-backend/internal/models"
)

type clinicPetLookupServiceSpy struct {
	query dto.ClinicPetLookupQuery
}

func (s *clinicPetLookupServiceSpy) LookupPetForClinic(query dto.ClinicPetLookupQuery) (any, error) {
	s.query = query
	return &dto.ClinicPetLookupItemResponse{
		ID: uuid.NewString(), PublicPetID: "PNX-PET-ABC123", Name: "Milo", Species: "dog",
	}, nil
}

func TestClinicPetLookupRouteRequiresAuthenticationAndClinicRole(t *testing.T) {
	router, _ := newClinicPetLookupRouteTestRouter()

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/clinic/pet-lookup?pet_id=PNX-PET-ABC123", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	ownerRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/pet-lookup?pet_id=PNX-PET-ABC123", nil)
	ownerRequest.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleOwner))
	ownerResponse := httptest.NewRecorder()
	router.ServeHTTP(ownerResponse, ownerRequest)
	if ownerResponse.Code != http.StatusForbidden {
		t.Fatalf("owner status = %d, want 403", ownerResponse.Code)
	}
}

func TestClinicPetLookupRoutePassesQueryToService(t *testing.T) {
	router, spy := newClinicPetLookupRouteTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/clinic/pet-lookup?owner_phone=0812345678", nil)
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleClinic))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if spy.query.OwnerPhone != "0812345678" || spy.query.PetID != "" {
		t.Fatalf("query = %#v", spy.query)
	}
}

func newClinicPetLookupRouteTestRouter() (*gin.Engine, *clinicPetLookupServiceSpy) {
	gin.SetMode(gin.TestMode)
	spy := &clinicPetLookupServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config:              config.Config{JWTSecret: "clinic-route-test-secret", JWTExpiresIn: "1h"},
		ClinicLookupHandler: handlers.NewClinicPetLookupHandler(spy),
	})
	return router, spy
}

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

type clinicPatientServiceSpy struct {
	userID  string
	petID   uuid.UUID
	filters dto.ClinicPatientFilters
}

func (s *clinicPatientServiceSpy) ListClinicPatients(userID string, filters dto.ClinicPatientFilters) ([]dto.ClinicPatientListItemResponse, error) {
	s.userID = userID
	s.filters = filters
	return []dto.ClinicPatientListItemResponse{{Pet: dto.ClinicPatientPetSummary{ID: uuid.NewString(), Name: "Milo"}}}, nil
}

func (s *clinicPatientServiceSpy) GetClinicPatient(userID string, petID uuid.UUID) (*dto.ClinicPatientDetailResponse, error) {
	s.userID = userID
	s.petID = petID
	return &dto.ClinicPatientDetailResponse{Pet: dto.ClinicPatientPetDetail{ID: petID.String(), Name: "Milo"}}, nil
}

func TestClinicPatientRoutesRequireAuthAndClinicRole(t *testing.T) {
	router, _ := newClinicPatientRouteTestRouter()

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/clinic/patients", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	ownerRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/patients", nil)
	ownerRequest.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleOwner))
	ownerResponse := httptest.NewRecorder()
	router.ServeHTTP(ownerResponse, ownerRequest)
	if ownerResponse.Code != http.StatusForbidden {
		t.Fatalf("owner status = %d, want 403", ownerResponse.Code)
	}
}

func TestClinicPatientListPassesJWTUserAndFilters(t *testing.T) {
	router, spy := newClinicPatientRouteTestRouter()
	userID := uuid.NewString()
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/clinic/patients?q=milo&species=dog&status=scheduled&limit=20&offset=5&sort=name_asc",
		nil,
	)
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, userID, models.RoleClinic))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if spy.userID != userID ||
		spy.filters.Q != "milo" ||
		spy.filters.Species != "dog" ||
		spy.filters.Status != "scheduled" ||
		spy.filters.Limit != "20" ||
		spy.filters.Offset != "5" ||
		spy.filters.Sort != "name_asc" {
		t.Fatalf("service call user=%q filters=%#v", spy.userID, spy.filters)
	}
}

func TestClinicPatientDetailPassesPetID(t *testing.T) {
	router, spy := newClinicPatientRouteTestRouter()
	userID := uuid.NewString()
	petID := uuid.New()
	request := httptest.NewRequest(http.MethodGet, "/api/clinic/patients/"+petID.String(), nil)
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, userID, models.RoleClinicStaff))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if spy.userID != userID || spy.petID != petID {
		t.Fatalf("service call user=%q pet=%s", spy.userID, spy.petID)
	}
}

func TestClinicPatientRouteRejectsInvalidPetID(t *testing.T) {
	router, _ := newClinicPatientRouteTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/clinic/patients/not-a-uuid", nil)
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleClinic))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}

func newClinicPatientRouteTestRouter() (*gin.Engine, *clinicPatientServiceSpy) {
	gin.SetMode(gin.TestMode)
	spy := &clinicPatientServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config:               config.Config{JWTSecret: "clinic-route-test-secret", JWTExpiresIn: "1h"},
		ClinicPatientHandler: handlers.NewClinicPatientHandler(spy),
	})
	return router, spy
}

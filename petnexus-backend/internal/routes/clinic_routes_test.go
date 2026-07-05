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
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

type clinicProfileServiceSpy struct {
	currentUserID string
	createRequest dto.CreateClinicProfileRequest
}

func (s *clinicProfileServiceSpy) CreateClinicProfile(userID string, req dto.CreateClinicProfileRequest) (*dto.ClinicProfileResponse, error) {
	s.currentUserID = userID
	s.createRequest = req
	return clinicRouteTestResponse(), nil
}

func (s *clinicProfileServiceSpy) GetMyClinicProfile(userID string) (*dto.ClinicProfileResponse, error) {
	s.currentUserID = userID
	return clinicRouteTestResponse(), nil
}

func (s *clinicProfileServiceSpy) UpdateMyClinicProfile(userID string, req dto.UpdateClinicProfileRequest) (*dto.ClinicProfileResponse, error) {
	s.currentUserID = userID
	return clinicRouteTestResponse(), nil
}

func TestClinicProfileRoutesRequireAuthenticationAndClinicStaffRole(t *testing.T) {
	router, _ := newClinicRouteTestRouter()

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/clinic/profile", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	ownerToken := clinicRouteTestToken(t, uuid.NewString(), models.RoleOwner)
	ownerRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/profile", nil)
	ownerRequest.Header.Set("Authorization", "Bearer "+ownerToken)
	ownerResponse := httptest.NewRecorder()
	router.ServeHTTP(ownerResponse, ownerRequest)
	if ownerResponse.Code != http.StatusForbidden {
		t.Fatalf("owner status = %d, want 403; body=%s", ownerResponse.Code, ownerResponse.Body.String())
	}
}

func TestClinicStaffCanGetOwnClinicProfile(t *testing.T) {
	router, spy := newClinicRouteTestRouter()
	userID := uuid.NewString()
	token := clinicRouteTestToken(t, userID, models.RoleClinicStaff)
	request := httptest.NewRequest(http.MethodGet, "/api/clinic/profile", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", response.Code, response.Body.String())
	}
	if spy.currentUserID != userID {
		t.Fatalf("service user ID = %q, want JWT user ID %q", spy.currentUserID, userID)
	}
}

func TestCreateClinicProfileUsesJWTIdentityAndDoesNotExposeUserID(t *testing.T) {
	router, spy := newClinicRouteTestRouter()
	userID := uuid.NewString()
	token := clinicRouteTestToken(t, userID, models.RoleClinicStaff)
	body := `{
		"user_id":"` + uuid.NewString() + `",
		"clinic_name":"Happy Paws Clinic",
		"email":"clinic@example.com"
	}`
	request := httptest.NewRequest(http.MethodPost, "/api/clinic/profile", strings.NewReader(body))
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", response.Code, response.Body.String())
	}
	if spy.currentUserID != userID {
		t.Fatalf("service user ID = %q, want JWT user ID %q", spy.currentUserID, userID)
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data := payload["data"].(map[string]any)
	if _, found := data["user_id"]; found {
		t.Fatal("response must not expose user_id")
	}
}

func newClinicRouteTestRouter() (*gin.Engine, *clinicProfileServiceSpy) {
	gin.SetMode(gin.TestMode)
	spy := &clinicProfileServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config:        config.Config{JWTSecret: "clinic-route-test-secret", JWTExpiresIn: "1h"},
		ClinicHandler: handlers.NewClinicProfileHandler(spy),
	})
	return router, spy
}

func clinicRouteTestToken(t *testing.T, userID, role string) string {
	t.Helper()
	token, err := utils.GenerateAccessToken(userID, role, "clinic-route-test-secret", "1h")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	return token
}

func clinicRouteTestResponse() *dto.ClinicProfileResponse {
	return &dto.ClinicProfileResponse{
		ID:         uuid.NewString(),
		ClinicName: "Happy Paws Clinic",
		CreatedAt:  "2026-07-05T01:02:03Z",
		UpdatedAt:  "2026-07-05T01:02:03Z",
	}
}

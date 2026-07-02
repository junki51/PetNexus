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

type ownerProfileServiceSpy struct {
	currentUserID string
	createRequest dto.CreateOwnerProfileRequest
}

func (s *ownerProfileServiceSpy) CreateProfile(userID string, req dto.CreateOwnerProfileRequest) (*dto.OwnerProfileResponse, error) {
	s.currentUserID = userID
	s.createRequest = req
	return testOwnerProfileResponse(), nil
}

func (s *ownerProfileServiceSpy) GetProfile(userID string) (*dto.OwnerProfileResponse, error) {
	s.currentUserID = userID
	return testOwnerProfileResponse(), nil
}

func (s *ownerProfileServiceSpy) UpdateProfile(userID string, req dto.UpdateOwnerProfileRequest) (*dto.OwnerProfileResponse, error) {
	s.currentUserID = userID
	return testOwnerProfileResponse(), nil
}

func TestOwnerProfileRouteRequiresAuthentication(t *testing.T) {
	router, _ := newOwnerRouteTestRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/api/owner/profile", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnauthorized, response.Body.String())
	}
}

func TestOwnerProfileRouteRejectsNonOwner(t *testing.T) {
	router, _ := newOwnerRouteTestRouter(t)
	token := ownerRouteTestToken(t, uuid.NewString(), models.RoleClinicStaff)
	request := httptest.NewRequest(http.MethodGet, "/api/owner/profile", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusForbidden, response.Body.String())
	}
}

func TestCreateOwnerProfileUsesJWTIdentityAndDoesNotExposeUserID(t *testing.T) {
	router, spy := newOwnerRouteTestRouter(t)
	authenticatedUserID := uuid.NewString()
	token := ownerRouteTestToken(t, authenticatedUserID, models.RoleOwner)
	body := `{
		"user_id":"` + uuid.NewString() + `",
		"first_name":"Sunny",
		"last_name":"Example",
		"phone_number":"0812345678"
	}`
	request := httptest.NewRequest(http.MethodPost, "/api/owner/profile", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusCreated, response.Body.String())
	}
	if spy.currentUserID != authenticatedUserID {
		t.Fatalf("service user ID = %q, want JWT user ID %q", spy.currentUserID, authenticatedUserID)
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("response data = %#v", payload["data"])
	}
	if _, exposed := data["user_id"]; exposed {
		t.Fatal("response must not expose user_id")
	}
}

func newOwnerRouteTestRouter(t *testing.T) (*gin.Engine, *ownerProfileServiceSpy) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	spy := &ownerProfileServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config: config.Config{
			JWTSecret:    "owner-route-test-secret",
			JWTExpiresIn: "1h",
		},
		OwnerHandler: handlers.NewOwnerProfileHandler(spy),
	})
	return router, spy
}

func ownerRouteTestToken(t *testing.T, userID, role string) string {
	t.Helper()
	token, err := utils.GenerateAccessToken(userID, role, "owner-route-test-secret", "1h")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	return token
}

func testOwnerProfileResponse() *dto.OwnerProfileResponse {
	return &dto.OwnerProfileResponse{
		ID:          uuid.NewString(),
		FirstName:   "Sunny",
		LastName:    "Example",
		DisplayName: "Sunny Example",
		PhoneNumber: "0812345678",
		CreatedAt:   "2026-07-02T01:02:03Z",
		UpdatedAt:   "2026-07-02T01:02:03Z",
	}
}

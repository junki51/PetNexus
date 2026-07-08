package routes

import (
	"bytes"
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

type ownerAppointmentServiceSpy struct {
	userID string
	create dto.CreateOwnerAppointmentRequest
}

func (s *ownerAppointmentServiceSpy) CreateOwnerAppointment(userID string, req dto.CreateOwnerAppointmentRequest) (*dto.AppointmentResponse, error) {
	s.userID = userID
	s.create = req
	return appointmentRouteResponse(), nil
}

func (s *ownerAppointmentServiceSpy) ListOwnerAppointments(userID string, filters dto.OwnerAppointmentFilters) ([]dto.AppointmentResponse, error) {
	s.userID = userID
	return []dto.AppointmentResponse{*appointmentRouteResponse()}, nil
}

func (s *ownerAppointmentServiceSpy) GetOwnerAppointment(userID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	s.userID = userID
	return appointmentRouteResponse(), nil
}

func (s *ownerAppointmentServiceSpy) CancelOwnerAppointment(userID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	s.userID = userID
	return appointmentRouteResponse(), nil
}

type clinicAppointmentServiceSpy struct {
	userID  string
	filters dto.ClinicAppointmentFilters
	status  string
}

func (s *clinicAppointmentServiceSpy) CreateClinicAppointment(userID string, req dto.CreateClinicAppointmentRequest) (*dto.AppointmentResponse, error) {
	s.userID = userID
	return appointmentRouteResponse(), nil
}

func (s *clinicAppointmentServiceSpy) ListClinicAppointments(userID string, filters dto.ClinicAppointmentFilters) ([]dto.AppointmentResponse, error) {
	s.userID = userID
	s.filters = filters
	return []dto.AppointmentResponse{*appointmentRouteResponse()}, nil
}

func (s *clinicAppointmentServiceSpy) GetClinicAppointment(userID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	s.userID = userID
	return appointmentRouteResponse(), nil
}

func (s *clinicAppointmentServiceSpy) UpdateClinicAppointmentStatus(userID string, appointmentID uuid.UUID, req dto.UpdateAppointmentStatusRequest) (*dto.AppointmentResponse, error) {
	s.userID = userID
	s.status = req.Status
	return appointmentRouteResponse(), nil
}

func (s *clinicAppointmentServiceSpy) CancelClinicAppointment(userID string, appointmentID uuid.UUID) (*dto.AppointmentResponse, error) {
	s.userID = userID
	return appointmentRouteResponse(), nil
}

func TestAppointmentRoutesRequireCorrectRole(t *testing.T) {
	router, _, _ := newAppointmentRouteTestRouter()

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/owner/appointments", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	clinicRequest := httptest.NewRequest(http.MethodGet, "/api/owner/appointments", nil)
	clinicRequest.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleClinic))
	clinicResponse := httptest.NewRecorder()
	router.ServeHTTP(clinicResponse, clinicRequest)
	if clinicResponse.Code != http.StatusForbidden {
		t.Fatalf("clinic on owner route status = %d, want 403", clinicResponse.Code)
	}

	ownerRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/appointments", nil)
	ownerRequest.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleOwner))
	ownerResponse := httptest.NewRecorder()
	router.ServeHTTP(ownerResponse, ownerRequest)
	if ownerResponse.Code != http.StatusForbidden {
		t.Fatalf("owner on clinic route status = %d, want 403", ownerResponse.Code)
	}
}

func TestOwnerAppointmentCreateUsesJWTUser(t *testing.T) {
	router, ownerSpy, _ := newAppointmentRouteTestRouter()
	userID := uuid.NewString()
	body := []byte(`{"clinic_profile_id":"` + uuid.NewString() + `","pet_id":"` + uuid.NewString() + `","appointment_type":"checkup","scheduled_at":"2026-07-10T10:00:00Z","duration_minutes":30}`)
	request := httptest.NewRequest(http.MethodPost, "/api/owner/appointments", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, userID, models.RoleOwner))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", response.Code, response.Body.String())
	}
	if ownerSpy.userID != userID || ownerSpy.create.AppointmentType != "checkup" {
		t.Fatalf("service call = user %q request %#v", ownerSpy.userID, ownerSpy.create)
	}
}

func TestClinicAppointmentRoutesPassFiltersAndStatus(t *testing.T) {
	router, _, clinicSpy := newAppointmentRouteTestRouter()
	userID := uuid.NewString()
	token := clinicRouteTestToken(t, userID, models.RoleClinic)

	listRequest := httptest.NewRequest(
		http.MethodGet,
		"/api/clinic/appointments?date=2026-07-10&status=scheduled&appointment_type=vaccination",
		nil,
	)
	listRequest.Header.Set("Authorization", "Bearer "+token)
	listResponse := httptest.NewRecorder()
	router.ServeHTTP(listResponse, listRequest)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, want 200", listResponse.Code)
	}
	if clinicSpy.filters.Date != "2026-07-10" || clinicSpy.filters.Status != "scheduled" ||
		clinicSpy.filters.AppointmentType != "vaccination" {
		t.Fatalf("filters = %#v", clinicSpy.filters)
	}

	appointmentID := uuid.NewString()
	statusRequest := httptest.NewRequest(
		http.MethodPatch,
		"/api/clinic/appointments/"+appointmentID+"/status",
		bytes.NewBufferString(`{"status":"checked_in"}`),
	)
	statusRequest.Header.Set("Content-Type", "application/json")
	statusRequest.Header.Set("Authorization", "Bearer "+token)
	statusResponse := httptest.NewRecorder()
	router.ServeHTTP(statusResponse, statusRequest)
	if statusResponse.Code != http.StatusOK || clinicSpy.status != "checked_in" {
		t.Fatalf("status update response=%d service status=%q", statusResponse.Code, clinicSpy.status)
	}
}

func TestAppointmentRouteRejectsInvalidPathID(t *testing.T) {
	router, _, _ := newAppointmentRouteTestRouter()
	request := httptest.NewRequest(http.MethodGet, "/api/owner/appointments/not-a-uuid", nil)
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleOwner))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.Code)
	}
}

func newAppointmentRouteTestRouter() (*gin.Engine, *ownerAppointmentServiceSpy, *clinicAppointmentServiceSpy) {
	gin.SetMode(gin.TestMode)
	ownerSpy := &ownerAppointmentServiceSpy{}
	clinicSpy := &clinicAppointmentServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config:                   config.Config{JWTSecret: "clinic-route-test-secret", JWTExpiresIn: "1h"},
		OwnerAppointmentHandler:  handlers.NewOwnerAppointmentHandler(ownerSpy),
		ClinicAppointmentHandler: handlers.NewClinicAppointmentHandler(clinicSpy),
	})
	return router, ownerSpy, clinicSpy
}

func appointmentRouteResponse() *dto.AppointmentResponse {
	return &dto.AppointmentResponse{ID: uuid.NewString(), Status: models.AppointmentStatusScheduled}
}

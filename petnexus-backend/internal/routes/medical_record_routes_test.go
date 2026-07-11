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

type medicalRecordServiceSpy struct {
	userID     string
	petID      uuid.UUID
	recordID   uuid.UUID
	create     dto.CreateMedicalRecordRequest
	update     dto.UpdateMedicalRecordRequest
	filters    dto.MedicalRecordFilters
	createCall bool
}

func (s *medicalRecordServiceSpy) CreateMedicalRecord(userID string, petID uuid.UUID, req dto.CreateMedicalRecordRequest) (*dto.MedicalRecordDetailResponse, error) {
	s.userID = userID
	s.petID = petID
	s.create = req
	s.createCall = true
	return medicalRecordRouteDetailResponse(), nil
}

func (s *medicalRecordServiceSpy) ListMedicalRecords(userID string, filters dto.MedicalRecordFilters) (*dto.MedicalRecordListResponse, error) {
	s.userID = userID
	s.filters = filters
	return &dto.MedicalRecordListResponse{
		Items:      []dto.MedicalRecordListItemResponse{{ID: uuid.NewString(), ChiefComplaint: "Coughing"}},
		Pagination: dto.PaginationMeta{Page: 1, Limit: 20, Total: 1, TotalPages: 1},
	}, nil
}

func (s *medicalRecordServiceSpy) GetMedicalRecord(userID string, recordID uuid.UUID) (*dto.MedicalRecordDetailResponse, error) {
	s.userID = userID
	s.recordID = recordID
	return medicalRecordRouteDetailResponse(), nil
}

func (s *medicalRecordServiceSpy) UpdateMedicalRecord(userID string, recordID uuid.UUID, req dto.UpdateMedicalRecordRequest) (*dto.MedicalRecordDetailResponse, error) {
	s.userID = userID
	s.recordID = recordID
	s.update = req
	return medicalRecordRouteDetailResponse(), nil
}

func TestMedicalRecordRoutesRequireAuthAndClinicRole(t *testing.T) {
	router, _ := newMedicalRecordRouteTestRouter()

	withoutToken := httptest.NewRecorder()
	router.ServeHTTP(withoutToken, httptest.NewRequest(http.MethodGet, "/api/clinic/medical-records", nil))
	if withoutToken.Code != http.StatusUnauthorized {
		t.Fatalf("without token status = %d, want 401", withoutToken.Code)
	}

	ownerRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/medical-records", nil)
	ownerRequest.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, uuid.NewString(), models.RoleOwner))
	ownerResponse := httptest.NewRecorder()
	router.ServeHTTP(ownerResponse, ownerRequest)
	if ownerResponse.Code != http.StatusForbidden {
		t.Fatalf("owner status = %d, want 403", ownerResponse.Code)
	}
}

func TestMedicalRecordCreateRoutePassesJWTUserPetIDAndBody(t *testing.T) {
	router, spy := newMedicalRecordRouteTestRouter()
	userID := uuid.NewString()
	petID := uuid.New()
	appointmentID := uuid.NewString()
	body := []byte(`{"appointmentId":"` + appointmentID + `","visitAt":"2026-07-11T08:00:00Z","chiefComplaint":"Coughing"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/clinic/patients/"+petID.String()+"/medical-records", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, userID, models.RoleClinic))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body=%s", response.Code, response.Body.String())
	}
	if spy.userID != userID || spy.petID != petID || spy.create.AppointmentID != appointmentID || spy.create.ChiefComplaint != "Coughing" {
		t.Fatalf("service call = user=%q pet=%s req=%#v", spy.userID, spy.petID, spy.create)
	}
}

func TestMedicalRecordListRoutePassesFilters(t *testing.T) {
	router, spy := newMedicalRecordRouteTestRouter()
	userID := uuid.NewString()
	petID := uuid.NewString()
	request := httptest.NewRequest(http.MethodGet, "/api/clinic/medical-records?pet_id="+petID+"&from=2026-07-01&to=2026-07-31&page=2&limit=10", nil)
	request.Header.Set("Authorization", "Bearer "+clinicRouteTestToken(t, userID, models.RoleClinicStaff))
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}
	if spy.userID != userID || spy.filters.PetID != petID || spy.filters.From != "2026-07-01" ||
		spy.filters.To != "2026-07-31" || spy.filters.Page != "2" || spy.filters.Limit != "10" {
		t.Fatalf("filters = %#v", spy.filters)
	}
}

func TestMedicalRecordDetailAndPatchRoutesPassRecordID(t *testing.T) {
	router, spy := newMedicalRecordRouteTestRouter()
	userID := uuid.NewString()
	recordID := uuid.New()
	token := clinicRouteTestToken(t, userID, models.RoleClinic)

	getRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/medical-records/"+recordID.String(), nil)
	getRequest.Header.Set("Authorization", "Bearer "+token)
	getResponse := httptest.NewRecorder()
	router.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK || spy.recordID != recordID {
		t.Fatalf("get response=%d record=%s", getResponse.Code, spy.recordID)
	}

	patchRequest := httptest.NewRequest(
		http.MethodPatch,
		"/api/clinic/medical-records/"+recordID.String(),
		bytes.NewBufferString(`{"chiefComplaint":"Updated"}`),
	)
	patchRequest.Header.Set("Content-Type", "application/json")
	patchRequest.Header.Set("Authorization", "Bearer "+token)
	patchResponse := httptest.NewRecorder()
	router.ServeHTTP(patchResponse, patchRequest)
	if patchResponse.Code != http.StatusOK || spy.update.ChiefComplaint == nil || *spy.update.ChiefComplaint != "Updated" {
		t.Fatalf("patch response=%d update=%#v", patchResponse.Code, spy.update)
	}
}

func TestMedicalRecordRoutesRejectInvalidIDs(t *testing.T) {
	router, _ := newMedicalRecordRouteTestRouter()
	token := clinicRouteTestToken(t, uuid.NewString(), models.RoleClinic)

	createRequest := httptest.NewRequest(http.MethodPost, "/api/clinic/patients/not-a-uuid/medical-records", bytes.NewBufferString(`{}`))
	createRequest.Header.Set("Authorization", "Bearer "+token)
	createResponse := httptest.NewRecorder()
	router.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusBadRequest {
		t.Fatalf("create invalid pet status = %d, want 400", createResponse.Code)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/api/clinic/medical-records/not-a-uuid", nil)
	getRequest.Header.Set("Authorization", "Bearer "+token)
	getResponse := httptest.NewRecorder()
	router.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusBadRequest {
		t.Fatalf("get invalid record status = %d, want 400", getResponse.Code)
	}
}

func newMedicalRecordRouteTestRouter() (*gin.Engine, *medicalRecordServiceSpy) {
	gin.SetMode(gin.TestMode)
	spy := &medicalRecordServiceSpy{}
	router := gin.New()
	Register(router, Dependencies{
		Config:               config.Config{JWTSecret: "clinic-route-test-secret", JWTExpiresIn: "1h"},
		MedicalRecordHandler: handlers.NewMedicalRecordHandler(spy),
	})
	return router, spy
}

func medicalRecordRouteDetailResponse() *dto.MedicalRecordDetailResponse {
	return &dto.MedicalRecordDetailResponse{
		ID:             uuid.NewString(),
		VisitAt:        "2026-07-11T08:00:00Z",
		ChiefComplaint: "Coughing",
	}
}

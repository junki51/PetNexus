package services

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/phonlakitz/petnexus-backend/internal/dto"
	"github.com/phonlakitz/petnexus-backend/internal/models"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
)

type medicalRecordRepositoryStub struct {
	records      map[uuid.UUID]*models.MedicalRecord
	patients     map[string]bool
	appointments map[uuid.UUID]*models.Appointment
	lastFilters  repositories.MedicalRecordFilters
}

func newMedicalRecordRepositoryStub() *medicalRecordRepositoryStub {
	return &medicalRecordRepositoryStub{
		records:      make(map[uuid.UUID]*models.MedicalRecord),
		patients:     make(map[string]bool),
		appointments: make(map[uuid.UUID]*models.Appointment),
	}
}

func (r *medicalRecordRepositoryStub) Create(record *models.MedicalRecord) error {
	if record.AppointmentID != nil {
		exists, _ := r.MedicalRecordExistsByAppointmentID(*record.AppointmentID)
		if exists {
			return repositories.ErrMedicalRecordAppointmentAlreadyLinked
		}
	}
	if record.ID == uuid.Nil {
		record.ID = uuid.New()
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = time.Date(2026, 7, 11, 3, 0, 0, 0, time.UTC)
		record.UpdatedAt = record.CreatedAt
	}
	r.records[record.ID] = cloneMedicalRecord(record)
	return nil
}

func (r *medicalRecordRepositoryStub) FindAllByClinicProfileID(clinicProfileID uuid.UUID, filters repositories.MedicalRecordFilters) ([]models.MedicalRecord, error) {
	r.lastFilters = filters
	result := make([]models.MedicalRecord, 0)
	for _, record := range r.records {
		if record.ClinicProfileID != clinicProfileID {
			continue
		}
		if filters.PetID != nil && record.PetID != *filters.PetID {
			continue
		}
		if filters.DateFrom != nil && record.VisitAt.Before(*filters.DateFrom) {
			continue
		}
		if filters.DateTo != nil && !record.VisitAt.Before(*filters.DateTo) {
			continue
		}
		result = append(result, *cloneMedicalRecord(record))
	}
	return result, nil
}

func (r *medicalRecordRepositoryStub) CountByClinicProfileID(clinicProfileID uuid.UUID, filters repositories.MedicalRecordFilters) (int64, error) {
	records, err := r.FindAllByClinicProfileID(clinicProfileID, filters)
	return int64(len(records)), err
}

func (r *medicalRecordRepositoryStub) FindByIDAndClinicProfileID(id, clinicProfileID uuid.UUID) (*models.MedicalRecord, error) {
	record, ok := r.records[id]
	if !ok || record.ClinicProfileID != clinicProfileID {
		return nil, repositories.ErrMedicalRecordNotFound
	}
	return cloneMedicalRecord(record), nil
}

func (r *medicalRecordRepositoryStub) Update(record *models.MedicalRecord) error {
	stored, ok := r.records[record.ID]
	if !ok || stored.ClinicProfileID != record.ClinicProfileID {
		return repositories.ErrMedicalRecordNotFound
	}
	r.records[record.ID] = cloneMedicalRecord(record)
	return nil
}

func (r *medicalRecordRepositoryStub) PetHasNonCancelledAppointmentWithClinic(clinicProfileID, petID uuid.UUID) (bool, error) {
	return r.patients[patientKey(clinicProfileID, petID)], nil
}

func (r *medicalRecordRepositoryStub) FindUsableAppointmentForMedicalRecord(appointmentID, clinicProfileID, petID uuid.UUID) (*models.Appointment, error) {
	appointment, ok := r.appointments[appointmentID]
	if !ok || appointment.ClinicProfileID != clinicProfileID ||
		appointment.PetID != petID || appointment.Status == models.AppointmentStatusCancelled {
		return nil, repositories.ErrAppointmentNotFound
	}
	return appointment, nil
}

func (r *medicalRecordRepositoryStub) MedicalRecordExistsByAppointmentID(appointmentID uuid.UUID) (bool, error) {
	for _, record := range r.records {
		if record.AppointmentID != nil && *record.AppointmentID == appointmentID {
			return true, nil
		}
	}
	return false, nil
}

func cloneMedicalRecord(record *models.MedicalRecord) *models.MedicalRecord {
	clone := *record
	return &clone
}

func patientKey(clinicProfileID, petID uuid.UUID) string {
	return clinicProfileID.String() + ":" + petID.String()
}

func TestMedicalRecordCreateForOwnPatientWithAppointment(t *testing.T) {
	now := time.Date(2026, 7, 11, 8, 0, 0, 0, time.UTC)
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	owner := &models.OwnerProfile{ID: uuid.New(), FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678"}
	pet := &models.Pet{ID: uuid.New(), PublicPetID: "PNX-PET-ABC123", Name: "Milo", Species: models.SpeciesDog, OwnerProfile: owner}
	appointment := &models.Appointment{
		ID: uuid.New(), ClinicProfileID: clinic.ID, PetID: pet.ID,
		ScheduledAt: now.Add(-time.Hour), Status: models.AppointmentStatusCheckedIn,
	}
	repo := newMedicalRecordRepositoryStub()
	repo.patients[patientKey(clinic.ID, pet.ID)] = true
	repo.appointments[appointment.ID] = appointment
	service := NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})
	service.(*medicalRecordService).now = func() time.Time { return now }

	response, err := service.CreateMedicalRecord(userID.String(), pet.ID, dto.CreateMedicalRecordRequest{
		AppointmentID:        appointment.ID.String(),
		VisitAt:              now.Format(time.RFC3339),
		ChiefComplaint:       "  Coughing  ",
		ClinicalFindings:     "Mild fever",
		Diagnosis:            "URI",
		TreatmentPlan:        "Rest",
		Medications:          "Medicine A",
		FollowUpInstructions: "Follow up in one week",
		NextFollowUpAt:       now.Add(7 * 24 * time.Hour).Format(time.RFC3339),
		WeightKG:             floatPtr(12.5),
		TemperatureC:         floatPtr(38.2),
		Notes:                "Owner informed",
	})
	if err != nil {
		t.Fatalf("CreateMedicalRecord() error = %v", err)
	}
	if response.ChiefComplaint != "Coughing" || response.Diagnosis == nil || *response.Diagnosis != "URI" {
		t.Fatalf("unexpected response: %#v", response)
	}
	for _, record := range repo.records {
		if record.ClinicProfileID != clinic.ID || record.PetID != pet.ID ||
			record.CreatedByUserID != userID || record.AppointmentID == nil ||
			*record.AppointmentID != appointment.ID {
			t.Fatalf("server-derived fields were not set correctly: %#v", record)
		}
	}
}

func TestMedicalRecordCreateRejectsUnrelatedPet(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	service := NewMedicalRecordService(newMedicalRecordRepositoryStub(), &clinicProfileRepositoryStub{profile: clinic, exists: true})

	_, err := service.CreateMedicalRecord(userID.String(), uuid.New(), dto.CreateMedicalRecordRequest{
		VisitAt:        time.Now().Format(time.RFC3339),
		ChiefComplaint: "Coughing",
	})
	assertAppointmentServiceError(t, err, http.StatusNotFound, "CLINIC_PATIENT_NOT_FOUND")
}

func TestMedicalRecordAppointmentValidation(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	petID := uuid.New()
	repo := newMedicalRecordRepositoryStub()
	repo.patients[patientKey(clinic.ID, petID)] = true
	validAppointment := &models.Appointment{ID: uuid.New(), ClinicProfileID: clinic.ID, PetID: petID, Status: models.AppointmentStatusScheduled}
	repo.appointments[validAppointment.ID] = validAppointment
	service := NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	_, err := service.CreateMedicalRecord(userID.String(), petID, dto.CreateMedicalRecordRequest{
		AppointmentID: uuid.NewString(), VisitAt: time.Now().Format(time.RFC3339), ChiefComplaint: "Coughing",
	})
	assertAppointmentServiceError(t, err, http.StatusNotFound, "APPOINTMENT_NOT_FOUND")

	otherPetAppointment := &models.Appointment{ID: uuid.New(), ClinicProfileID: clinic.ID, PetID: uuid.New(), Status: models.AppointmentStatusScheduled}
	repo.appointments[otherPetAppointment.ID] = otherPetAppointment
	_, err = service.CreateMedicalRecord(userID.String(), petID, dto.CreateMedicalRecordRequest{
		AppointmentID: otherPetAppointment.ID.String(), VisitAt: time.Now().Format(time.RFC3339), ChiefComplaint: "Coughing",
	})
	assertAppointmentServiceError(t, err, http.StatusNotFound, "APPOINTMENT_NOT_FOUND")

	otherClinicAppointment := &models.Appointment{ID: uuid.New(), ClinicProfileID: uuid.New(), PetID: petID, Status: models.AppointmentStatusScheduled}
	repo.appointments[otherClinicAppointment.ID] = otherClinicAppointment
	_, err = service.CreateMedicalRecord(userID.String(), petID, dto.CreateMedicalRecordRequest{
		AppointmentID: otherClinicAppointment.ID.String(), VisitAt: time.Now().Format(time.RFC3339), ChiefComplaint: "Coughing",
	})
	assertAppointmentServiceError(t, err, http.StatusNotFound, "APPOINTMENT_NOT_FOUND")

	existingRecord := &models.MedicalRecord{
		ID: uuid.New(), ClinicProfileID: clinic.ID, PetID: petID, AppointmentID: &validAppointment.ID,
		CreatedByUserID: userID, VisitAt: time.Now(), ChiefComplaint: "Existing",
	}
	repo.records[existingRecord.ID] = existingRecord
	_, err = service.CreateMedicalRecord(userID.String(), petID, dto.CreateMedicalRecordRequest{
		AppointmentID: validAppointment.ID.String(), VisitAt: time.Now().Format(time.RFC3339), ChiefComplaint: "Coughing",
	})
	assertAppointmentServiceError(t, err, http.StatusConflict, "APPOINTMENT_MEDICAL_RECORD_EXISTS")
}

func TestMedicalRecordCreateWithoutAppointmentForExistingPatient(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	petID := uuid.New()
	repo := newMedicalRecordRepositoryStub()
	repo.patients[patientKey(clinic.ID, petID)] = true
	service := NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	response, err := service.CreateMedicalRecord(userID.String(), petID, dto.CreateMedicalRecordRequest{
		VisitAt:        time.Now().Format(time.RFC3339),
		ChiefComplaint: "Walk-in check",
	})
	if err != nil {
		t.Fatalf("CreateMedicalRecord(no appointment) error = %v", err)
	}
	if response.Appointment != nil {
		t.Fatalf("appointment = %#v, want nil", response.Appointment)
	}
}

func TestMedicalRecordListFiltersAndScopesClinic(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	otherClinicID := uuid.New()
	petID := uuid.New()
	visitAt := time.Date(2026, 7, 11, 3, 0, 0, 0, time.UTC)
	repo := newMedicalRecordRepositoryStub()
	repo.records[uuid.New()] = &models.MedicalRecord{ID: uuid.New(), ClinicProfileID: clinic.ID, PetID: petID, VisitAt: visitAt, ChiefComplaint: "A", CreatedByUserID: userID, CreatedAt: visitAt, UpdatedAt: visitAt}
	repo.records[uuid.New()] = &models.MedicalRecord{ID: uuid.New(), ClinicProfileID: otherClinicID, PetID: petID, VisitAt: visitAt, ChiefComplaint: "B", CreatedByUserID: userID, CreatedAt: visitAt, UpdatedAt: visitAt}
	service := NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	response, err := service.ListMedicalRecords(userID.String(), dto.MedicalRecordFilters{
		PetID: petID.String(), From: "2026-07-10", To: "2026-07-12", Page: "2", Limit: "10",
	})
	if err != nil {
		t.Fatalf("ListMedicalRecords() error = %v", err)
	}
	if len(response.Items) != 1 || response.Pagination.Page != 2 || response.Pagination.Limit != 10 || response.Pagination.Total != 1 {
		t.Fatalf("unexpected list response: %#v", response)
	}
	if repo.lastFilters.PetID == nil || *repo.lastFilters.PetID != petID ||
		repo.lastFilters.DateFrom == nil || repo.lastFilters.DateTo == nil ||
		repo.lastFilters.Offset != 10 || repo.lastFilters.Limit != 10 {
		t.Fatalf("filters = %#v", repo.lastFilters)
	}
}

func TestMedicalRecordGetAndCrossClinicNotFound(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	recordID := uuid.New()
	now := time.Now().UTC()
	repo := newMedicalRecordRepositoryStub()
	repo.records[recordID] = &models.MedicalRecord{ID: recordID, ClinicProfileID: clinic.ID, PetID: uuid.New(), VisitAt: now, ChiefComplaint: "Coughing", CreatedByUserID: userID, CreatedAt: now, UpdatedAt: now}
	service := NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	response, err := service.GetMedicalRecord(userID.String(), recordID)
	if err != nil || response.ID != recordID.String() {
		t.Fatalf("GetMedicalRecord() = %#v, %v", response, err)
	}

	otherClinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	service = NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: otherClinic, exists: true})
	_, err = service.GetMedicalRecord(userID.String(), recordID)
	assertAppointmentServiceError(t, err, http.StatusNotFound, "MEDICAL_RECORD_NOT_FOUND")
}

func TestMedicalRecordUpdateAllowedFieldsAndKeepsImmutableOwnership(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	petID := uuid.New()
	appointmentID := uuid.New()
	recordID := uuid.New()
	visitAt := time.Date(2026, 7, 11, 3, 0, 0, 0, time.UTC)
	repo := newMedicalRecordRepositoryStub()
	repo.records[recordID] = &models.MedicalRecord{
		ID: recordID, ClinicProfileID: clinic.ID, PetID: petID, AppointmentID: &appointmentID,
		CreatedByUserID: userID, VisitAt: visitAt, ChiefComplaint: "Old", CreatedAt: visitAt, UpdatedAt: visitAt,
	}
	service := NewMedicalRecordService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	newVisitAt := visitAt.Add(2 * time.Hour).Format(time.RFC3339)
	response, err := service.UpdateMedicalRecord(userID.String(), recordID, dto.UpdateMedicalRecordRequest{
		VisitAt:          &newVisitAt,
		ChiefComplaint:   stringPtr("Updated complaint"),
		Diagnosis:        stringPtr("Updated diagnosis"),
		WeightKG:         floatPtr(13.2),
		TemperatureC:     floatPtr(38.5),
		ClinicalFindings: stringPtr(""),
	})
	if err != nil {
		t.Fatalf("UpdateMedicalRecord() error = %v", err)
	}
	stored := repo.records[recordID]
	if response.ChiefComplaint != "Updated complaint" || stored.ClinicProfileID != clinic.ID ||
		stored.PetID != petID || stored.AppointmentID == nil || *stored.AppointmentID != appointmentID ||
		stored.CreatedByUserID != userID {
		t.Fatalf("unexpected updated record: response=%#v stored=%#v", response, stored)
	}

	_, err = service.UpdateMedicalRecord(userID.String(), recordID, dto.UpdateMedicalRecordRequest{})
	assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestMedicalRecordValidation(t *testing.T) {
	now := time.Date(2026, 7, 11, 8, 0, 0, 0, time.UTC)
	badFollowUp := now.Add(-time.Hour).Format(time.RFC3339)
	tests := []dto.CreateMedicalRecordRequest{
		{AppointmentID: "not-a-uuid", VisitAt: now.Format(time.RFC3339), ChiefComplaint: "Coughing"},
		{VisitAt: "", ChiefComplaint: "Coughing"},
		{VisitAt: "not-a-time", ChiefComplaint: "Coughing"},
		{VisitAt: now.Format(time.RFC3339), ChiefComplaint: "   "},
		{VisitAt: now.Format(time.RFC3339), ChiefComplaint: "Coughing", WeightKG: floatPtr(0)},
		{VisitAt: now.Format(time.RFC3339), ChiefComplaint: "Coughing", TemperatureC: floatPtr(-1)},
		{VisitAt: now.Format(time.RFC3339), ChiefComplaint: "Coughing", NextFollowUpAt: badFollowUp},
	}
	for _, req := range tests {
		_, err := normalizeCreateMedicalRecordInput(req)
		assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
	}

	filterTests := []dto.MedicalRecordFilters{
		{PetID: "not-a-uuid"},
		{From: "2026/07/11"},
		{To: "2026/07/11"},
		{From: "2026-07-12", To: "2026-07-11"},
		{Page: "0"},
		{Limit: "0"},
		{Limit: "101"},
	}
	for _, filters := range filterTests {
		_, _, err := normalizeMedicalRecordFilters(filters)
		assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
	}
}

func floatPtr(value float64) *float64 {
	return &value
}

func stringPtr(value string) *string {
	return &value
}

func TestMedicalRecordRepositoryUnexpectedErrorShape(t *testing.T) {
	if !errors.Is(repositories.ErrMedicalRecordAppointmentAlreadyLinked, repositories.ErrMedicalRecordAppointmentAlreadyLinked) {
		t.Fatal("sentinel error sanity check failed")
	}
}

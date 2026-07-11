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

type clinicPatientRepositoryStub struct {
	records             []repositories.ClinicPatientRecord
	detail              *repositories.ClinicPatientDetail
	err                 error
	lastClinicProfileID uuid.UUID
	lastPetID           uuid.UUID
	lastFilters         repositories.ClinicPatientFilters
}

func (r *clinicPatientRepositoryStub) FindPatientsByClinicProfileID(clinicProfileID uuid.UUID, filters repositories.ClinicPatientFilters) ([]repositories.ClinicPatientRecord, error) {
	r.lastClinicProfileID = clinicProfileID
	r.lastFilters = filters
	if r.err != nil {
		return nil, r.err
	}
	return r.records, nil
}

func (r *clinicPatientRepositoryStub) FindPatientDetailByClinicProfileIDAndPetID(clinicProfileID, petID uuid.UUID) (*repositories.ClinicPatientDetail, error) {
	r.lastClinicProfileID = clinicProfileID
	r.lastPetID = petID
	if r.err != nil {
		return nil, r.err
	}
	if r.detail == nil {
		return nil, repositories.ErrClinicPatientNotFound
	}
	return r.detail, nil
}

func TestClinicPatientListUsesCurrentClinicAndMapsSafeResponse(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	breedNameTH := "โกลเด้น รีทรีฟเวอร์"
	breed := &models.Breed{ID: uuid.New(), Species: models.SpeciesDog, Name: "Golden Retriever", NameTH: &breedNameTH}
	owner := &models.OwnerProfile{ID: uuid.New(), FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678"}
	pet := &models.Pet{
		ID: uuid.New(), PublicPetID: "PNX-PET-ABC123", Name: "Milo", Species: models.SpeciesDog,
		Breed: breed, OwnerProfile: owner,
	}
	firstSeen := time.Date(2026, 7, 1, 3, 0, 0, 0, time.UTC)
	lastSeen := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	nextSeen := time.Date(2026, 7, 12, 3, 0, 0, 0, time.UTC)
	repo := &clinicPatientRepositoryStub{
		records: []repositories.ClinicPatientRecord{{
			Pet: pet,
			Summary: repositories.ClinicPatientSummary{
				TotalAppointments:  2,
				FirstAppointmentAt: &firstSeen,
				LastAppointmentAt:  &lastSeen,
				NextAppointmentAt:  &nextSeen,
				LatestStatus:       models.AppointmentStatusScheduled,
			},
		}},
	}
	service := NewClinicPatientService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	response, err := service.ListClinicPatients(userID.String(), dto.ClinicPatientFilters{
		Q: "  milo  ", Species: " DOG ", Status: " SCHEDULED ", Limit: "20", Offset: "5", Sort: "name_asc",
	})
	if err != nil {
		t.Fatalf("ListClinicPatients() error = %v", err)
	}
	if repo.lastClinicProfileID != clinic.ID {
		t.Fatalf("clinic scope = %s, want %s", repo.lastClinicProfileID, clinic.ID)
	}
	if repo.lastFilters.Query != "milo" || repo.lastFilters.Species != models.SpeciesDog ||
		repo.lastFilters.Status != models.AppointmentStatusScheduled || repo.lastFilters.Limit != 20 ||
		repo.lastFilters.Offset != 5 || repo.lastFilters.Sort != "name_asc" {
		t.Fatalf("filters = %#v", repo.lastFilters)
	}
	if len(response) != 1 {
		t.Fatalf("len(response) = %d, want 1", len(response))
	}
	item := response[0]
	if item.Pet.ID != pet.ID.String() || item.Pet.PublicPetID != pet.PublicPetID || item.Pet.Breed == nil {
		t.Fatalf("unexpected pet summary: %#v", item.Pet)
	}
	if item.Owner.MaskedPhone != "081****678" || item.Owner.DisplayName != "Sunny Example" {
		t.Fatalf("unexpected owner summary: %#v", item.Owner)
	}
	if item.AppointmentSummary.TotalAppointments != 2 ||
		item.AppointmentSummary.LatestStatus != models.AppointmentStatusScheduled ||
		item.FirstSeenAt == nil || *item.FirstSeenAt != "2026-07-01T03:00:00Z" {
		t.Fatalf("unexpected appointment summary: %#v", item)
	}
}

func TestClinicPatientListValidationAndMissingProfile(t *testing.T) {
	service := NewClinicPatientService(&clinicPatientRepositoryStub{}, &clinicProfileRepositoryStub{})
	_, err := service.ListClinicPatients(uuid.NewString(), dto.ClinicPatientFilters{})
	assertAppointmentServiceError(t, err, http.StatusNotFound, "CLINIC_PROFILE_REQUIRED")

	userID := uuid.New()
	service = NewClinicPatientService(
		&clinicPatientRepositoryStub{},
		&clinicProfileRepositoryStub{profile: &models.ClinicProfile{ID: uuid.New(), UserID: userID}, exists: true},
	)
	tests := []dto.ClinicPatientFilters{
		{Species: "bird"},
		{Status: "unknown"},
		{Limit: "0"},
		{Limit: "101"},
		{Offset: "-1"},
		{Sort: "created_at"},
	}
	for _, tt := range tests {
		_, err := service.ListClinicPatients(userID.String(), tt)
		assertAppointmentServiceError(t, err, http.StatusBadRequest, "VALIDATION_ERROR")
	}
}

func TestClinicPatientDetailMapsRecentAppointmentsAndCrossClinicNotFound(t *testing.T) {
	userID := uuid.New()
	clinic := &models.ClinicProfile{ID: uuid.New(), UserID: userID}
	dateOfBirth := time.Date(2022, 5, 10, 0, 0, 0, 0, time.UTC)
	weight := 12.5
	neutered := true
	microchip := "MC-123"
	color := "Brown"
	marks := "White spot"
	owner := &models.OwnerProfile{ID: uuid.New(), FirstName: "Sunny", LastName: "Example", PhoneNumber: "0812345678"}
	pet := &models.Pet{
		ID: uuid.New(), PublicPetID: "PNX-PET-ABC123", Name: "Milo", Species: models.SpeciesDog,
		DateOfBirth: &dateOfBirth, WeightKG: &weight, MicrochipID: &microchip, Color: &color,
		DistinctiveMarks: &marks, IsNeutered: &neutered, OwnerProfile: owner,
	}
	firstSeen := time.Date(2026, 7, 1, 3, 0, 0, 0, time.UTC)
	lastSeen := time.Date(2026, 7, 9, 3, 0, 0, 0, time.UTC)
	title := "Annual checkup"
	appointmentID := uuid.New()
	repo := &clinicPatientRepositoryStub{
		detail: &repositories.ClinicPatientDetail{
			Pet: pet,
			Summary: repositories.ClinicPatientSummary{
				TotalAppointments:  2,
				FirstAppointmentAt: &firstSeen,
				LastAppointmentAt:  &lastSeen,
				LatestStatus:       models.AppointmentStatusCheckedIn,
			},
			RecentAppointments: []models.Appointment{{
				ID: appointmentID, ScheduledAt: lastSeen, AppointmentType: models.AppointmentTypeCheckup,
				Status: models.AppointmentStatusCheckedIn, Title: &title,
			}},
		},
	}
	service := NewClinicPatientService(repo, &clinicProfileRepositoryStub{profile: clinic, exists: true})

	response, err := service.GetClinicPatient(userID.String(), pet.ID)
	if err != nil {
		t.Fatalf("GetClinicPatient() error = %v", err)
	}
	if repo.lastClinicProfileID != clinic.ID || repo.lastPetID != pet.ID {
		t.Fatalf("scope clinic=%s pet=%s", repo.lastClinicProfileID, repo.lastPetID)
	}
	if response.Pet.MicrochipID == nil || *response.Pet.MicrochipID != microchip ||
		response.Pet.DateOfBirth == nil || *response.Pet.DateOfBirth != "2022-05-10" {
		t.Fatalf("unexpected pet detail: %#v", response.Pet)
	}
	if response.Owner.MaskedPhone != "081****678" ||
		response.ClinicRelationship.TotalAppointments != 2 ||
		len(response.RecentAppointments) != 1 ||
		response.RecentAppointments[0].ID != appointmentID.String() {
		t.Fatalf("unexpected patient detail: %#v", response)
	}

	repo.err = repositories.ErrClinicPatientNotFound
	_, err = service.GetClinicPatient(userID.String(), uuid.New())
	assertAppointmentServiceError(t, err, http.StatusNotFound, "CLINIC_PATIENT_NOT_FOUND")
}

func TestClinicPatientUnexpectedRepositoryError(t *testing.T) {
	userID := uuid.New()
	service := NewClinicPatientService(
		&clinicPatientRepositoryStub{err: errors.New("database down")},
		&clinicProfileRepositoryStub{profile: &models.ClinicProfile{ID: uuid.New(), UserID: userID}, exists: true},
	)
	_, err := service.ListClinicPatients(userID.String(), dto.ClinicPatientFilters{})
	assertAppointmentServiceError(t, err, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR")
}

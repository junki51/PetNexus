package dto

// CreateOwnerAppointmentRequest creates an appointment for the current
// owner's own pet at an existing clinic.
type CreateOwnerAppointmentRequest struct {
	ClinicProfileID string `json:"clinic_profile_id"`
	PetID           string `json:"pet_id"`
	Title           string `json:"title"`
	AppointmentType string `json:"appointment_type"`
	ScheduledAt     string `json:"scheduled_at"`
	DurationMinutes int    `json:"duration_minutes"`
	Note            string `json:"note"`
}

// CreateClinicAppointmentRequest creates an appointment under the current
// clinic profile. Exactly one pet lookup value must be supplied.
type CreateClinicAppointmentRequest struct {
	PetID           string `json:"pet_id"`
	PublicPetID     string `json:"public_pet_id"`
	Title           string `json:"title"`
	AppointmentType string `json:"appointment_type"`
	ScheduledAt     string `json:"scheduled_at"`
	DurationMinutes int    `json:"duration_minutes"`
	Note            string `json:"note"`
}

type UpdateAppointmentStatusRequest struct {
	Status string `json:"status"`
}

type OwnerAppointmentFilters struct {
	DateFrom string
	DateTo   string
	Status   string
}

type ClinicAppointmentFilters struct {
	Date            string
	DateFrom        string
	DateTo          string
	Status          string
	AppointmentType string
}

type AppointmentPetSummary struct {
	ID          string         `json:"id"`
	PublicPetID string         `json:"public_pet_id"`
	Name        string         `json:"name"`
	Species     string         `json:"species"`
	AvatarURL   *string        `json:"avatar_url"`
	Breed       *BreedResponse `json:"breed"`
}

type AppointmentOwnerSummary struct {
	DisplayName string `json:"display_name"`
	MaskedPhone string `json:"masked_phone"`
}

type AppointmentClinicSummary struct {
	ID          string  `json:"id"`
	ClinicName  string  `json:"clinic_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
}

// AppointmentResponse deliberately omits internal ownership IDs and user IDs.
type AppointmentResponse struct {
	ID              string                   `json:"id"`
	Title           *string                  `json:"title"`
	AppointmentType string                   `json:"appointment_type"`
	ScheduledAt     string                   `json:"scheduled_at"`
	DurationMinutes int                      `json:"duration_minutes"`
	Status          string                   `json:"status"`
	Note            *string                  `json:"note"`
	CreatedByRole   string                   `json:"created_by_role"`
	CancelledAt     *string                  `json:"cancelled_at"`
	CreatedAt       string                   `json:"created_at"`
	UpdatedAt       string                   `json:"updated_at"`
	Pet             AppointmentPetSummary    `json:"pet"`
	Owner           AppointmentOwnerSummary  `json:"owner"`
	Clinic          AppointmentClinicSummary `json:"clinic"`
}

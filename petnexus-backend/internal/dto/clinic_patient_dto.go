package dto

// ClinicPatientFilters is built from the Clinic Web Patients query string.
type ClinicPatientFilters struct {
	Q       string
	Species string
	Status  string
	Limit   string
	Offset  string
	Sort    string
}

type ClinicPatientPetSummary struct {
	ID          string         `json:"id"`
	PublicPetID string         `json:"public_pet_id"`
	Name        string         `json:"name"`
	Species     string         `json:"species"`
	AvatarURL   *string        `json:"avatar_url"`
	Breed       *BreedResponse `json:"breed"`
}

type ClinicPatientPetDetail struct {
	ID               string         `json:"id"`
	PublicPetID      string         `json:"public_pet_id"`
	Name             string         `json:"name"`
	Species          string         `json:"species"`
	Gender           *string        `json:"gender"`
	DateOfBirth      *string        `json:"date_of_birth"`
	WeightKG         *float64       `json:"weight_kg"`
	MicrochipID      *string        `json:"microchip_id"`
	AvatarURL        *string        `json:"avatar_url"`
	Color            *string        `json:"color"`
	DistinctiveMarks *string        `json:"distinctive_marks"`
	IsNeutered       *bool          `json:"is_neutered"`
	Breed            *BreedResponse `json:"breed"`
}

type ClinicPatientOwnerSummary struct {
	DisplayName string `json:"display_name"`
	MaskedPhone string `json:"masked_phone"`
}

type ClinicPatientAppointmentSummary struct {
	TotalAppointments int64   `json:"total_appointments"`
	LastAppointmentAt *string `json:"last_appointment_at"`
	NextAppointmentAt *string `json:"next_appointment_at"`
	LatestStatus      string  `json:"latest_status"`
}

type ClinicPatientRelationshipSummary struct {
	FirstAppointmentAt *string `json:"first_appointment_at"`
	LastAppointmentAt  *string `json:"last_appointment_at"`
	NextAppointmentAt  *string `json:"next_appointment_at"`
	TotalAppointments  int64   `json:"total_appointments"`
}

type ClinicPatientListItemResponse struct {
	Pet                ClinicPatientPetSummary         `json:"pet"`
	Owner              ClinicPatientOwnerSummary       `json:"owner"`
	AppointmentSummary ClinicPatientAppointmentSummary `json:"appointment_summary"`
	FirstSeenAt        *string                         `json:"first_seen_at"`
}

type ClinicPatientRecentAppointmentResponse struct {
	ID              string  `json:"id"`
	ScheduledAt     string  `json:"scheduled_at"`
	AppointmentType string  `json:"appointment_type"`
	Status          string  `json:"status"`
	Title           *string `json:"title"`
}

type ClinicPatientDetailResponse struct {
	Pet                ClinicPatientPetDetail                   `json:"pet"`
	Owner              ClinicPatientOwnerSummary                `json:"owner"`
	ClinicRelationship ClinicPatientRelationshipSummary         `json:"clinic_relationship"`
	RecentAppointments []ClinicPatientRecentAppointmentResponse `json:"recent_appointments"`
}

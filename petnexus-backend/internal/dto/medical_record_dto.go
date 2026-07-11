package dto

type CreateMedicalRecordRequest struct {
	AppointmentID        string   `json:"appointmentId"`
	VisitAt              string   `json:"visitAt"`
	ChiefComplaint       string   `json:"chiefComplaint"`
	ClinicalFindings     string   `json:"clinicalFindings"`
	Diagnosis            string   `json:"diagnosis"`
	TreatmentPlan        string   `json:"treatmentPlan"`
	Medications          string   `json:"medications"`
	FollowUpInstructions string   `json:"followUpInstructions"`
	NextFollowUpAt       string   `json:"nextFollowUpAt"`
	WeightKG             *float64 `json:"weightKg"`
	TemperatureC         *float64 `json:"temperatureC"`
	Notes                string   `json:"notes"`
}

type UpdateMedicalRecordRequest struct {
	VisitAt              *string  `json:"visitAt"`
	ChiefComplaint       *string  `json:"chiefComplaint"`
	ClinicalFindings     *string  `json:"clinicalFindings"`
	Diagnosis            *string  `json:"diagnosis"`
	TreatmentPlan        *string  `json:"treatmentPlan"`
	Medications          *string  `json:"medications"`
	FollowUpInstructions *string  `json:"followUpInstructions"`
	NextFollowUpAt       *string  `json:"nextFollowUpAt"`
	WeightKG             *float64 `json:"weightKg"`
	TemperatureC         *float64 `json:"temperatureC"`
	Notes                *string  `json:"notes"`
}

type MedicalRecordFilters struct {
	PetID string
	From  string
	To    string
	Page  string
	Limit string
}

type MedicalRecordPetSummary struct {
	ID          string         `json:"id"`
	PublicPetID string         `json:"publicPetId"`
	Name        string         `json:"name"`
	Species     string         `json:"species"`
	Breed       *BreedResponse `json:"breed"`
}

type MedicalRecordOwnerSummary struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
}

type MedicalRecordAppointmentSummary struct {
	ID          string `json:"id"`
	ScheduledAt string `json:"scheduledAt"`
	Status      string `json:"status"`
}

type MedicalRecordCreatedBySummary struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type MedicalRecordListItemResponse struct {
	ID             string                           `json:"id"`
	VisitAt        string                           `json:"visitAt"`
	ChiefComplaint string                           `json:"chiefComplaint"`
	Diagnosis      *string                          `json:"diagnosis"`
	CreatedAt      string                           `json:"createdAt"`
	UpdatedAt      string                           `json:"updatedAt"`
	Pet            MedicalRecordPetSummary          `json:"pet"`
	Owner          MedicalRecordOwnerSummary        `json:"owner"`
	Appointment    *MedicalRecordAppointmentSummary `json:"appointment"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type MedicalRecordListResponse struct {
	Items      []MedicalRecordListItemResponse `json:"items"`
	Pagination PaginationMeta                  `json:"pagination"`
}

type MedicalRecordDetailResponse struct {
	ID                   string                           `json:"id"`
	VisitAt              string                           `json:"visitAt"`
	ChiefComplaint       string                           `json:"chiefComplaint"`
	ClinicalFindings     *string                          `json:"clinicalFindings"`
	Diagnosis            *string                          `json:"diagnosis"`
	TreatmentPlan        *string                          `json:"treatmentPlan"`
	Medications          *string                          `json:"medications"`
	FollowUpInstructions *string                          `json:"followUpInstructions"`
	NextFollowUpAt       *string                          `json:"nextFollowUpAt"`
	WeightKG             *float64                         `json:"weightKg"`
	TemperatureC         *float64                         `json:"temperatureC"`
	Notes                *string                          `json:"notes"`
	CreatedAt            string                           `json:"createdAt"`
	UpdatedAt            string                           `json:"updatedAt"`
	Pet                  MedicalRecordPetSummary          `json:"pet"`
	Owner                MedicalRecordOwnerSummary        `json:"owner"`
	Appointment          *MedicalRecordAppointmentSummary `json:"appointment"`
	CreatedBy            *MedicalRecordCreatedBySummary   `json:"createdBy"`
}

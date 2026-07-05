package dto

// CreateClinicProfileRequest creates settings for the current clinic user.
// UserID is intentionally absent because identity always comes from JWT.
type CreateClinicProfileRequest struct {
	ClinicName  string `json:"clinic_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Address     string `json:"address"`
}

// UpdateClinicProfileRequest uses pointers so PATCH changes only fields that
// were present in the request body.
type UpdateClinicProfileRequest struct {
	ClinicName  *string `json:"clinic_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
}

// ClinicProfileResponse excludes user ID and authentication data.
type ClinicProfileResponse struct {
	ID          string  `json:"id"`
	ClinicName  string  `json:"clinic_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	Address     *string `json:"address"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

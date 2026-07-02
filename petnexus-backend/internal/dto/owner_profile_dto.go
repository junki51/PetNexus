package dto

// CreateOwnerProfileRequest is the first-time owner profile payload.
// UserID is intentionally absent because identity always comes from JWT.
type CreateOwnerProfileRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	PhoneNumber  string `json:"phone_number"`
	AvatarURL    string `json:"avatar_url"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	Province     string `json:"province"`
	District     string `json:"district"`
	Subdistrict  string `json:"subdistrict"`
	PostalCode   string `json:"postal_code"`
}

// UpdateOwnerProfileRequest uses pointers so PATCH updates only fields that
// were present in the request body.
type UpdateOwnerProfileRequest struct {
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Gender       *string `json:"gender"`
	DateOfBirth  *string `json:"date_of_birth"`
	PhoneNumber  *string `json:"phone_number"`
	AvatarURL    *string `json:"avatar_url"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Province     *string `json:"province"`
	District     *string `json:"district"`
	Subdistrict  *string `json:"subdistrict"`
	PostalCode   *string `json:"postal_code"`
}

// OwnerProfileResponse is the safe public owner profile representation.
type OwnerProfileResponse struct {
	ID           string  `json:"id"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	DisplayName  string  `json:"display_name"`
	Gender       *string `json:"gender"`
	DateOfBirth  *string `json:"date_of_birth"`
	PhoneNumber  string  `json:"phone_number"`
	AvatarURL    *string `json:"avatar_url"`
	AddressLine1 *string `json:"address_line1"`
	AddressLine2 *string `json:"address_line2"`
	Province     *string `json:"province"`
	District     *string `json:"district"`
	Subdistrict  *string `json:"subdistrict"`
	PostalCode   *string `json:"postal_code"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

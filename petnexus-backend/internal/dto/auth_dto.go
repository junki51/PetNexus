package dto

// RegisterRequest is the public owner or clinic staff registration payload.
type RegisterRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// LoginRequest is the credential payload used to create an access token.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse is the safe public representation of a user.
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// AuthResponse contains a safe user response and a signed access token.
type AuthResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json:"accessToken"`
}

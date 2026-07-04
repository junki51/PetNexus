package dto

// BreedResponse is the public breed reference shape.
type BreedResponse struct {
	ID      string  `json:"id"`
	Species string  `json:"species"`
	Name    string  `json:"name"`
	NameTH  *string `json:"name_th"`
}

// CreatePetRequest creates basic pet identity data. OwnerProfileID and UserID
// are intentionally absent because ownership always comes from JWT.
type CreatePetRequest struct {
	Species          string   `json:"species"`
	Name             string   `json:"name"`
	BreedID          *string  `json:"breed_id"`
	Gender           string   `json:"gender"`
	DateOfBirth      string   `json:"date_of_birth"`
	WeightKG         *float64 `json:"weight_kg"`
	MicrochipID      string   `json:"microchip_id"`
	AvatarURL        string   `json:"avatar_url"`
	Color            string   `json:"color"`
	DistinctiveMarks string   `json:"distinctive_marks"`
	IsNeutered       *bool    `json:"is_neutered"`
}

// UpdatePetRequest uses pointers so PATCH changes only provided fields.
// Sending an empty breed_id clears the current breed.
type UpdatePetRequest struct {
	Species          *string  `json:"species"`
	Name             *string  `json:"name"`
	BreedID          *string  `json:"breed_id"`
	Gender           *string  `json:"gender"`
	DateOfBirth      *string  `json:"date_of_birth"`
	WeightKG         *float64 `json:"weight_kg"`
	MicrochipID      *string  `json:"microchip_id"`
	AvatarURL        *string  `json:"avatar_url"`
	Color            *string  `json:"color"`
	DistinctiveMarks *string  `json:"distinctive_marks"`
	IsNeutered       *bool    `json:"is_neutered"`
}

// PetResponse excludes owner IDs and all authentication/private owner data.
type PetResponse struct {
	ID               string         `json:"id"`
	Species          string         `json:"species"`
	Name             string         `json:"name"`
	Gender           *string        `json:"gender"`
	DateOfBirth      *string        `json:"date_of_birth"`
	AgeYears         *int           `json:"age_years"`
	Breed            *BreedResponse `json:"breed"`
	WeightKG         *float64       `json:"weight_kg"`
	MicrochipID      *string        `json:"microchip_id"`
	AvatarURL        *string        `json:"avatar_url"`
	Color            *string        `json:"color"`
	DistinctiveMarks *string        `json:"distinctive_marks"`
	IsNeutered       *bool          `json:"is_neutered"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

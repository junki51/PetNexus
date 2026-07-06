package utils

import (
	"crypto/rand"
	"fmt"
)

const (
	PublicPetIDPrefix       = "PNX-PET-"
	publicPetIDSuffixLength = 6
	publicPetIDAlphabet     = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

// GeneratePublicPetID creates a readable backend-owned pet identifier.
func GeneratePublicPetID() (string, error) {
	randomBytes := make([]byte, publicPetIDSuffixLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate public pet ID randomness: %w", err)
	}

	suffix := make([]byte, publicPetIDSuffixLength)
	for i, value := range randomBytes {
		suffix[i] = publicPetIDAlphabet[int(value)%len(publicPetIDAlphabet)]
	}
	return PublicPetIDPrefix + string(suffix), nil
}

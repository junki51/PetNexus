package utils

import (
	"regexp"
	"testing"
)

func TestGeneratePublicPetIDFormatAndPracticalUniqueness(t *testing.T) {
	pattern := regexp.MustCompile(`^PNX-PET-[A-Z0-9]{6}$`)
	seen := make(map[string]struct{}, 200)
	for i := 0; i < 200; i++ {
		publicPetID, err := GeneratePublicPetID()
		if err != nil {
			t.Fatalf("GeneratePublicPetID() error = %v", err)
		}
		if !pattern.MatchString(publicPetID) {
			t.Fatalf("public pet ID %q does not match expected format", publicPetID)
		}
		if _, exists := seen[publicPetID]; exists {
			t.Fatalf("unexpected duplicate public pet ID %q in test sample", publicPetID)
		}
		seen[publicPetID] = struct{}{}
	}
}

package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword returns a bcrypt hash and never stores or logs the plaintext.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CheckPasswordHash reports whether a plaintext password matches a bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

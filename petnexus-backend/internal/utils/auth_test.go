package utils

import (
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "password123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == password {
		t.Fatal("HashPassword() returned the plaintext password")
	}
	if !CheckPasswordHash(password, hash) {
		t.Fatal("CheckPasswordHash() rejected the correct password")
	}
	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("CheckPasswordHash() accepted an incorrect password")
	}
}

func TestGenerateAndParseAccessToken(t *testing.T) {
	token, err := GenerateAccessToken("user-id", "owner", "test-secret", "1h")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := ParseAccessToken(token, "test-secret")
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if claims.UserID != "user-id" || claims.Role != "owner" {
		t.Fatalf("ParseAccessToken() claims = %#v", claims)
	}

	if _, err := ParseAccessToken(token, "wrong-secret"); err == nil {
		t.Fatal("ParseAccessToken() accepted the wrong secret")
	}
}

func TestParseAccessTokenRejectsExpiredToken(t *testing.T) {
	token, err := GenerateAccessToken("user-id", "owner", "test-secret", "1s")
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	time.Sleep(1100 * time.Millisecond)

	if _, err := ParseAccessToken(token, "test-secret"); err == nil {
		t.Fatal("ParseAccessToken() accepted an expired token")
	}
}

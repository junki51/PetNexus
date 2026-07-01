package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultAccessTokenDuration = 24 * time.Hour

// Claims is the authenticated identity stored in a PetNexus access token.
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken signs an HS256 JWT using the configured secret.
func GenerateAccessToken(userID, role, secret, expiresIn string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("JWT secret is required")
	}

	duration := defaultAccessTokenDuration
	if strings.TrimSpace(expiresIn) != "" {
		parsedDuration, err := time.ParseDuration(expiresIn)
		if err != nil || parsedDuration <= 0 {
			return "", fmt.Errorf("invalid JWT expiration duration %q", expiresIn)
		}
		duration = parsedDuration
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "petnexus-backend",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return signedToken, nil
}

// ParseAccessToken validates an HS256 token and returns its typed claims.
func ParseAccessToken(tokenString, secret string) (*Claims, error) {
	if strings.TrimSpace(secret) == "" {
		return nil, errors.New("JWT secret is required")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method %q", token.Method.Alg())
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("parse access token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

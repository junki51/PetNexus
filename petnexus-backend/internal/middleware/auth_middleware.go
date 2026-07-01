// Package middleware contains HTTP request guards.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

const (
	ContextUserIDKey   = "userID"
	ContextUserRoleKey = "userRole"
)

// AuthMiddleware validates a Bearer access token and stores its identity in
// the Gin context for protected handlers.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			utils.Error(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				"UNAUTHORIZED",
				"Missing Authorization header",
			)
			c.Abort()
			return
		}

		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				"UNAUTHORIZED",
				"Invalid Authorization header format",
			)
			c.Abort()
			return
		}

		claims, err := utils.ParseAccessToken(parts[1], jwtSecret)
		if err != nil || claims.UserID == "" || claims.Role == "" {
			utils.Error(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				"UNAUTHORIZED",
				"Invalid or expired token",
			)
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserRoleKey, claims.Role)
		c.Next()
	}
}

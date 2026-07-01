package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

// RequireRole allows only authenticated users whose role is in roles. It must
// be registered after AuthMiddleware.
func RequireRole(roles ...string) gin.HandlerFunc {
	allowedRoles := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowedRoles[role] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get(ContextUserRoleKey)
		role, valid := value.(string)
		if !exists || !valid || role == "" {
			utils.Error(
				c,
				http.StatusUnauthorized,
				"Unauthorized",
				"UNAUTHORIZED",
				"Authenticated user role is missing",
			)
			c.Abort()
			return
		}

		if _, allowed := allowedRoles[role]; !allowed {
			utils.Error(
				c,
				http.StatusForbidden,
				"Forbidden",
				"FORBIDDEN_ROLE",
				"Your role is not allowed to access this resource",
			)
			c.Abort()
			return
		}

		c.Next()
	}
}

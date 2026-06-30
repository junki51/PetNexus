// Package handlers translates HTTP requests into application responses.
package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/database"
	"github.com/phonlakitz/petnexus-backend/internal/utils"
)

// Health reports whether the backend process is running.
func Health(c *gin.Context) {
	utils.Success(c, http.StatusOK, "PetNexus backend is running", gin.H{
		"status":  "ok",
		"service": "petnexus-backend",
	})
}

// DatabaseHealth verifies that the running application can still reach
// PostgreSQL through its configured GORM connection.
func DatabaseHealth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := database.PingPostgres(c.Request.Context(), db); err != nil {
			log.Printf("database health check failed: %v", err)
			utils.Error(
				c,
				http.StatusServiceUnavailable,
				"Database connection is unhealthy",
				"DATABASE_UNAVAILABLE",
				"Unable to reach PostgreSQL",
			)
			return
		}

		utils.Success(c, http.StatusOK, "Database connection is healthy", gin.H{
			"database": "postgresql",
			"status":   "connected",
		})
	}
}

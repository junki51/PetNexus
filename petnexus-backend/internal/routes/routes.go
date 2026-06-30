// Package routes owns HTTP route registration.
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/handlers"
)

// Register attaches all currently available routes to the router.
func Register(router *gin.Engine, db *gorm.DB) {
	router.GET("/health", handlers.Health)
	router.GET("/health/db", handlers.DatabaseHealth(db))

	// TODO: Register future route groups for:
	// /api/auth
	// /api/owner
	// /api/pets
	// /api/breeds
	// /api/clinic
	// /api/authorizations
	// /api/notifications
}

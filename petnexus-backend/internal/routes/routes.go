// Package routes owns HTTP route registration.
package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/handlers"
)

// Register attaches all currently available routes to the router.
func Register(router *gin.Engine) {
	router.GET("/health", handlers.Health)

	// TODO: Register future route groups for:
	// /api/auth
	// /api/owner
	// /api/pets
	// /api/breeds
	// /api/clinic
	// /api/authorizations
	// /api/notifications
}

// Package routes owns HTTP route registration.
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/handlers"
	"github.com/phonlakitz/petnexus-backend/internal/middleware"
)

// Dependencies contains the concrete application dependencies needed by the
// route layer.
type Dependencies struct {
	Config      config.Config
	DB          *gorm.DB
	AuthHandler *handlers.AuthHandler
}

// Register attaches all currently available routes to the router.
func Register(router *gin.Engine, deps Dependencies) {
	router.GET("/health", handlers.Health)
	router.GET("/health/db", handlers.DatabaseHealth(deps.DB))

	api := router.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", deps.AuthHandler.Register)
	auth.POST("/login", deps.AuthHandler.Login)
	api.GET(
		"/me",
		middleware.AuthMiddleware(deps.Config.JWTSecret),
		deps.AuthHandler.Me,
	)

	// TODO: Register future route groups for:
	// /api/owner
	// /api/pets
	// /api/breeds
	// /api/clinic
	// /api/authorizations
	// /api/notifications
}

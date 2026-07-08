// Package routes owns HTTP route registration.
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/handlers"
	"github.com/phonlakitz/petnexus-backend/internal/middleware"
	"github.com/phonlakitz/petnexus-backend/internal/models"
)

// Dependencies contains the concrete application dependencies needed by the
// route layer.
type Dependencies struct {
	Config                   config.Config
	DB                       *gorm.DB
	AuthHandler              *handlers.AuthHandler
	OwnerHandler             *handlers.OwnerProfileHandler
	BreedHandler             *handlers.BreedHandler
	PetHandler               *handlers.PetHandler
	ClinicHandler            *handlers.ClinicProfileHandler
	ClinicLookupHandler      *handlers.ClinicPetLookupHandler
	OwnerAppointmentHandler  *handlers.OwnerAppointmentHandler
	ClinicAppointmentHandler *handlers.ClinicAppointmentHandler
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

	owner := api.Group(
		"/owner",
		middleware.AuthMiddleware(deps.Config.JWTSecret),
		middleware.RequireRole(models.RoleOwner),
	)
	owner.POST("/profile", deps.OwnerHandler.CreateProfile)
	owner.GET("/profile", deps.OwnerHandler.GetProfile)
	owner.PATCH("/profile", deps.OwnerHandler.UpdateProfile)
	owner.POST("/appointments", deps.OwnerAppointmentHandler.CreateAppointment)
	owner.GET("/appointments", deps.OwnerAppointmentHandler.ListAppointments)
	owner.GET("/appointments/:id", deps.OwnerAppointmentHandler.GetAppointment)
	owner.PATCH("/appointments/:id/cancel", deps.OwnerAppointmentHandler.CancelAppointment)

	api.GET("/breeds", deps.BreedHandler.ListBreeds)

	pets := api.Group(
		"/pets",
		middleware.AuthMiddleware(deps.Config.JWTSecret),
		middleware.RequireRole(models.RoleOwner),
	)
	pets.POST("", deps.PetHandler.CreatePet)
	pets.GET("", deps.PetHandler.ListMyPets)
	pets.GET("/:id", deps.PetHandler.GetMyPet)
	pets.PATCH("/:id", deps.PetHandler.UpdateMyPet)

	clinic := api.Group(
		"/clinic",
		middleware.AuthMiddleware(deps.Config.JWTSecret),
		middleware.RequireRole(models.RoleClinic, models.RoleClinicStaff),
	)
	clinic.POST("/profile", deps.ClinicHandler.CreateClinicProfile)
	clinic.GET("/profile", deps.ClinicHandler.GetMyClinicProfile)
	clinic.PATCH("/profile", deps.ClinicHandler.UpdateMyClinicProfile)
	clinic.GET("/pet-lookup", deps.ClinicLookupHandler.LookupPet)
	clinic.POST("/appointments", deps.ClinicAppointmentHandler.CreateAppointment)
	clinic.GET("/appointments", deps.ClinicAppointmentHandler.ListAppointments)
	clinic.GET("/appointments/:id", deps.ClinicAppointmentHandler.GetAppointment)
	clinic.PATCH("/appointments/:id/status", deps.ClinicAppointmentHandler.UpdateAppointmentStatus)
	clinic.PATCH("/appointments/:id/cancel", deps.ClinicAppointmentHandler.CancelAppointment)

	// TODO: Register future route groups for:
	// /api/authorizations
	// /api/notifications
}

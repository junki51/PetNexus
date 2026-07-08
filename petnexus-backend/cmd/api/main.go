package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/database"
	"github.com/phonlakitz/petnexus-backend/internal/handlers"
	"github.com/phonlakitz/petnexus-backend/internal/repositories"
	"github.com/phonlakitz/petnexus-backend/internal/routes"
	"github.com/phonlakitz/petnexus-backend/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	log.Println("database connected successfully")

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}
	log.Println("database migration completed successfully")

	userRepo := repositories.NewUserRepository(db)
	ownerProfileRepo := repositories.NewOwnerProfileRepository(db)
	breedRepo := repositories.NewBreedRepository(db)
	petRepo := repositories.NewPetRepository(db)
	clinicProfileRepo := repositories.NewClinicProfileRepository(db)
	appointmentRepo := repositories.NewAppointmentRepository(db)
	authService := services.NewAuthService(userRepo, cfg)
	ownerProfileService := services.NewOwnerProfileService(ownerProfileRepo)
	breedService := services.NewBreedService(breedRepo)
	petService := services.NewPetService(petRepo, breedRepo, ownerProfileRepo)
	clinicProfileService := services.NewClinicProfileService(clinicProfileRepo)
	clinicPetLookupService := services.NewClinicPetLookupService(petRepo)
	ownerAppointmentService := services.NewOwnerAppointmentService(appointmentRepo, ownerProfileRepo, clinicProfileRepo, petRepo)
	clinicAppointmentService := services.NewClinicAppointmentService(appointmentRepo, clinicProfileRepo, petRepo)
	authHandler := handlers.NewAuthHandler(authService)
	ownerProfileHandler := handlers.NewOwnerProfileHandler(ownerProfileService)
	breedHandler := handlers.NewBreedHandler(breedService)
	petHandler := handlers.NewPetHandler(petService)
	clinicProfileHandler := handlers.NewClinicProfileHandler(clinicProfileService)
	clinicPetLookupHandler := handlers.NewClinicPetLookupHandler(clinicPetLookupService)
	ownerAppointmentHandler := handlers.NewOwnerAppointmentHandler(ownerAppointmentService)
	clinicAppointmentHandler := handlers.NewClinicAppointmentHandler(clinicAppointmentService)

	router := gin.Default()
	routes.Register(router, routes.Dependencies{
		Config:                   cfg,
		DB:                       db,
		AuthHandler:              authHandler,
		OwnerHandler:             ownerProfileHandler,
		BreedHandler:             breedHandler,
		PetHandler:               petHandler,
		ClinicHandler:            clinicProfileHandler,
		ClinicLookupHandler:      clinicPetLookupHandler,
		OwnerAppointmentHandler:  ownerAppointmentHandler,
		ClinicAppointmentHandler: clinicAppointmentHandler,
	})

	log.Printf("PetNexus backend listening on http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

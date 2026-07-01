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
	authService := services.NewAuthService(userRepo, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	router := gin.Default()
	routes.Register(router, routes.Dependencies{
		Config:      cfg,
		DB:          db,
		AuthHandler: authHandler,
	})

	log.Printf("PetNexus backend listening on http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

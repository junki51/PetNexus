package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/phonlakitz/petnexus-backend/internal/config"
	"github.com/phonlakitz/petnexus-backend/internal/routes"
)

func main() {
	cfg := config.Load()

	router := gin.Default()
	routes.Register(router)

	log.Printf("PetNexus backend listening on http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

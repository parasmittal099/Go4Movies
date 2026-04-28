package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/config"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/handlers"
	"github.com/parasmittal099/backend-project/routes"
	"github.com/parasmittal099/backend-project/workers"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	database.Connect(cfg.DBPath)
	database.Migrate()
	database.MigratePlaintextPasswords()
	database.Seed()

	handlers.JWTSecret = cfg.JWTSecret

	workers.StartBookingExpiry(database.DB, 1*time.Minute)

	r := gin.Default()
	routes.RegisterRoutes(r, cfg)

	log.Printf("Server starting on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

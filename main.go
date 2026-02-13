package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ram-ks/meeting-service/config"
	"github.com/ram-ks/meeting-service/controllers"
	"github.com/ram-ks/meeting-service/repository"
	"github.com/ram-ks/meeting-service/service"
)

func runMigrations(db *sql.DB) error {
	log.Println("🔄 Running database migrations...")

	// Check if migration file exists
	migrationPath := "./migrations/001_create_tables_up.sql"
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		log.Printf("⚠️  Migration file not found at %s, skipping...", migrationPath)
		return nil
	}

	// Read migration file
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := db.Exec(string(migrationSQL)); err != nil {
		// Check if error is because tables already exist
		if isTableExistsError(err) {
			log.Println("ℹ️ Tables already exist, skipping migration")
			return nil
		}
		return fmt.Errorf("failed to run migration: %w", err)
	}

	log.Println("✅ Migrations completed successfully")
	return nil
}

func isTableExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "duplicate")
}

func healthCheck(c *gin.Context) {
	db := config.GetDB()

	if db == nil {
		c.JSON(503, gin.H{
			"status":   "unhealthy",
			"database": "not initialized",
		})
		return
	}

	// Ping database
	if err := db.Ping(); err != nil {
		c.JSON(503, gin.H{
			"status":   "unhealthy",
			"database": "connection failed",
			"error":    err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status":   "healthy",
		"database": "connected",
	})
}

func main() {
	config.ConnectDatabase()

	db := config.GetDB()

	fmt.Printf("DB is nil: %v\n", db == nil)
	if db != nil {
		fmt.Println("✅ DB connection is valid!")

		// doesn't work on managed postgres, hence adding a fallback to run migrations from here
		if err := runMigrations(db); err != nil {
			log.Printf("❌ Migration error: %v", err)
			log.Println("⚠️  App will continue, but database operations may fail")
		}
	} else {
		log.Println("❌ Database connection failed!")
		log.Println("⚠️  App will start but all database operations will fail")
	}

	eventRepo := repository.NewEventRepository(db)
	eventCtrl := controllers.NewEventController(eventRepo)

	availabilityRepo := repository.NewAvailabilityRepository(db)
	availabilityService := service.NewAvailabilityService(availabilityRepo, eventRepo)
	availabilityCtrl := controllers.NewAvailabilityController(availabilityService)

	preferredSlotRepo := repository.NewPreferredSlotRepository(db)
	preferredSlotService := service.NewPreferredSlotService(preferredSlotRepo)

	schedulerService := service.NewSchedulerService(eventRepo, availabilityRepo, preferredSlotRepo)
	recommendationCtrl := controllers.NewRecommendationController(schedulerService)
	preferredSlotCtrl := controllers.NewPreferredSlotController(preferredSlotService)

	router := gin.Default()

	router.GET("/health", healthCheck)

	events := router.Group("/events")
	{
		events.POST("", eventCtrl.CreateEvent)
		events.GET("", eventCtrl.ListEvents)
		events.GET("/:id", eventCtrl.GetEvent)
		events.PUT("/:id", eventCtrl.UpdateEvent)
		events.DELETE("/:id", eventCtrl.DeleteEvent)
		events.GET("/:id/recommendations", recommendationCtrl.GetRecommendations)

		availability := events.Group("/:id/availability")
		{
			availability.POST("", availabilityCtrl.SubmitAvailability)
			availability.GET("", availabilityCtrl.GetAvailability)
			availability.GET("/:participant_id", availabilityCtrl.GetParticipantAvailability)
			availability.PUT("/:availability_id", availabilityCtrl.UpdateAvailability)
			availability.DELETE("/:availability_id", availabilityCtrl.DeleteAvailability)
		}
	}

	preferredSlots := router.Group("/preferred-slots")
	{
		preferredSlots.POST("", preferredSlotCtrl.CreatePreferredSlot)
		preferredSlots.GET("/email/:email", preferredSlotCtrl.GetPreferredSlotsByEmail)
		preferredSlots.PUT("/:id", preferredSlotCtrl.UpdatePreferredSlot)
		preferredSlots.DELETE("/:id", preferredSlotCtrl.DeletePreferredSlot)
	}

	log.Println("🚀 Server starting...")
	router.Run()
}

package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ram-ks/meeting-service/config"
	"github.com/ram-ks/meeting-service/controllers"
	"github.com/ram-ks/meeting-service/repository"
	"github.com/ram-ks/meeting-service/service"
)

func main() {
	config.ConnectDatabase()

	db := config.GetDB()

	fmt.Printf("DB is nil: %v\n", db == nil)
	if db != nil {
		fmt.Println("DB connection is valid!")
	}

	eventRepo := repository.NewEventRepository(db)
	eventCtrl := controllers.NewEventController(eventRepo)

	availabilityRepo := repository.NewAvailabilityRepository(db)
	availabilityService := service.NewAvailabilityService(availabilityRepo, eventRepo)
	availabilityCtrl := controllers.NewAvailabilityController(availabilityService)

	schedulerService := service.NewSchedulerService(eventRepo, availabilityRepo)
	recommendationCtrl := controllers.NewRecommendationController(schedulerService)

	router := gin.Default()
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

	router.Run()
}

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/service"
)

type RecommendationController struct {
	schedulerService service.SchedulerService
}

func NewRecommendationController(schedulerService service.SchedulerService) *RecommendationController {
	return &RecommendationController{schedulerService: schedulerService}
}

func (ctrl *RecommendationController) GetRecommendations(context *gin.Context) {
	eventID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	recommendations, err := ctrl.schedulerService.GetRecommendations(context.Request.Context(), eventID)
	if err != nil {
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, recommendations)
}

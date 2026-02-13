package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/service"
)

type AvailabilityController struct {
	availService service.AvailabilityService
}

func NewAvailabilityController(availService service.AvailabilityService) *AvailabilityController {
	return &AvailabilityController{availService: availService}
}

func (ctrl *AvailabilityController) SubmitAvailability(context *gin.Context) {
	eventID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req model.SubmitAvailabilityRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.availService.SubmitAvailability(context.Request.Context(), eventID, req); err != nil {
		log.Printf("❌ [SubmitAvailability] Database error: %v", err)
		log.Printf("❌ [SubmitAvailability] Error type: %T", err)
		log.Printf("❌ [SubmitAvailability] Request data: %+v", req)
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "availability submitted successfully"})
}

func (ctrl *AvailabilityController) GetAvailability(context *gin.Context) {
	eventID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	availabilities, err := ctrl.availService.GetAvailability(context.Request.Context(), eventID)
	if err != nil {
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"availabilities": availabilities})
}

func (ctrl *AvailabilityController) GetParticipantAvailability(context *gin.Context) {
	eventID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	participantID, err := uuid.Parse(context.Param("participant_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant id"})
		return
	}

	availabilities, err := ctrl.availService.GetParticipantAvailability(context.Request.Context(), eventID, participantID)
	if err != nil {
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"availabilities": availabilities})
}

func (ctrl *AvailabilityController) UpdateAvailability(context *gin.Context) {
	availabilityID, err := uuid.Parse(context.Param("availability_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid availability id"})
		return
	}

	var req model.UpdateAvailabilityRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	availability, err := ctrl.availService.UpdateAvailability(context.Request.Context(), availabilityID, req)
	if err != nil {
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, availability)
}

func (ctrl *AvailabilityController) DeleteAvailability(context *gin.Context) {
	availabilityID, err := uuid.Parse(context.Param("availability_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid availability id"})
		return
	}

	if err := ctrl.availService.DeleteAvailability(context.Request.Context(), availabilityID); err != nil {
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}

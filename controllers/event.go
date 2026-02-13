package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	models "github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/repository"
)

var (
	ErrEventNotFound        = errors.New("event not found")
	ErrSlotNotFound         = errors.New("slot not found")
	ErrInvalidStatus        = errors.New("invalid event status for this operation")
	ErrSlotNotInEvent       = errors.New("slot does not belong to this event")
	ErrInvalidTimeFormat    = errors.New("invalid time format")
	ErrParticipantNotFound  = errors.New("participant not found")
	ErrAvailabilityNotFound = errors.New("availability not found")
)

type EventController struct {
	repo repository.EventRepository
}

func NewEventController(repo repository.EventRepository) *EventController {
	return &EventController{repo: repo}
}

func getOrganizerID(c *gin.Context) uuid.UUID {
	if id, exists := c.Get("user_id"); exists {
		if userID, ok := id.(uuid.UUID); ok {
			return userID
		}
	}
	return uuid.MustParse("00000000-0000-0000-0000-000000000001")
}

func parseTime(timeStr, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, timeStr, loc); err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, errors.New("invalid time format")
}

func (ctrl *EventController) CreateEvent(context *gin.Context) {
	var req models.CreateEventRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ [CreateEvent] Failed to parse request body: %v", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organizerID := getOrganizerID(context)
	now := time.Now().UTC()

	event := &models.Event{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		OrganizerID: organizerID,
		Duration:    req.Duration,
		Status:      models.EventStatusOpen,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	for _, slotReq := range req.ProposedSlots {
		startTime, err := parseTime(slotReq.StartTime, slotReq.Timezone)
		if err != nil {
			fmt.Println("invalid time format")
		}
		endTime, err := parseTime(slotReq.EndTime, slotReq.Timezone)
		if err != nil {
			fmt.Println("invalid time format")
		}

		slot := models.TimeSlot{
			ID:        uuid.New(),
			EventID:   event.ID,
			StartTime: startTime,
			EndTime:   endTime,
			Timezone:  slotReq.Timezone,
			CreatedAt: now,
		}
		event.ProposedSlots = append(event.ProposedSlots, slot)
	}

	for _, pReq := range req.Participants {
		participant := models.Participant{
			ID:        uuid.New(),
			EventID:   event.ID,
			Email:     pReq.Email,
			Name:      pReq.Name,
			Status:    models.ParticipantStatusPending,
			CreatedAt: now,
		}
		event.Participants = append(event.Participants, participant)
	}

	if err := ctrl.repo.Create(context.Request.Context(), event); err != nil {
		log.Printf("❌ [CreateEvent] Database error: %v", err)
		log.Printf("❌ [CreateEvent] Error type: %T", err)
		log.Printf("❌ [CreateEvent] Request data: %+v", req)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	log.Printf("✅ [CreateEvent] Successfully created event: %s", event.ID)
	context.JSON(http.StatusCreated, event)
}

func (ctrl *EventController) ListEvents(context *gin.Context) {
	organizerID := getOrganizerID(context)
	events, err := ctrl.repo.List(context.Request.Context(), organizerID)
	if err != nil {
		handleServiceError(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"events": events})
}

func (ctrl *EventController) GetEvent(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := ctrl.repo.GetByID(context.Request.Context(), id)
	if err != nil {
		handleServiceError(context, err)
		return
	}
	context.JSON(http.StatusOK, event)
}

func (ctrl *EventController) DeleteEvent(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := ctrl.repo.Delete(context.Request.Context(), id); err != nil {
		handleServiceError(context, err)
		return
	}
	context.JSON(http.StatusNoContent, nil)
}

func (ctrl *EventController) UpdateEvent(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req model.UpdateEventRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := ctrl.repo.GetByID(context.Request.Context(), id)
	if err != nil {
		handleServiceError(context, err)
	}

	if event.Status == model.EventStatusFinalized || event.Status == model.EventStatusCancelled {
		handleServiceError(context, err)
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.Duration != nil {
		event.Duration = *req.Duration
	}

	if err := ctrl.repo.Update(context.Request.Context(), event); err != nil {
		handleServiceError(context, err)
	}

	context.JSON(http.StatusOK, event)
}

func handleServiceError(context *gin.Context, err error) {
	switch err {
	case ErrEventNotFound:
		context.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
	case ErrSlotNotFound:
		context.JSON(http.StatusNotFound, gin.H{"error": "slot not found"})
	case ErrInvalidStatus:
		context.JSON(http.StatusConflict, gin.H{"error": "invalid event status for this operation"})
	case ErrSlotNotInEvent:
		context.JSON(http.StatusBadRequest, gin.H{"error": "slot does not belong to this event"})
	case ErrInvalidTimeFormat:
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
	case ErrParticipantNotFound:
		context.JSON(http.StatusNotFound, gin.H{"error": "participant not found"})
	case ErrAvailabilityNotFound:
		context.JSON(http.StatusNotFound, gin.H{"error": "availability not found"})
	default:
		context.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/service"
)

type PreferredSlotController struct {
	service service.PreferredSlotService
}

func NewPreferredSlotController(service service.PreferredSlotService) *PreferredSlotController {
	return &PreferredSlotController{service: service}
}

func (ctrl *PreferredSlotController) CreatePreferredSlot(c *gin.Context) {
	var req model.CreatePreferredSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slot, err := ctrl.service.Create(c.Request.Context(), req)
	if err != nil {
		log.Printf("❌ [CreateSlot] Database error: %v", err)
		log.Printf("❌ [CreateSlot] Error type: %T", err)
		log.Printf("❌ [CreateSlot] Request data: %+v", req)
		handlePreferredSlotError(c, err)
		return
	}

	c.JSON(http.StatusCreated, slot)
}

func (ctrl *PreferredSlotController) GetPreferredSlotsByEmail(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	slots, err := ctrl.service.GetByEmail(c.Request.Context(), email)
	if err != nil {
		handlePreferredSlotError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferred_slots": slots})
}

func (ctrl *PreferredSlotController) UpdatePreferredSlot(c *gin.Context) {
	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}

	var req model.UpdatePreferredSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slot, err := ctrl.service.Update(c.Request.Context(), slotID, req)
	if err != nil {
		handlePreferredSlotError(c, err)
		return
	}

	c.JSON(http.StatusOK, slot)
}

func (ctrl *PreferredSlotController) DeletePreferredSlot(c *gin.Context) {
	slotID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid slot id"})
		return
	}

	if err := ctrl.service.Delete(c.Request.Context(), slotID); err != nil {
		handlePreferredSlotError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func handlePreferredSlotError(c *gin.Context, err error) {
	switch err {
	case service.ErrPreferredSlotNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "preferred slot not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

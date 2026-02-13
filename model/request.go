package model

import (
	"github.com/google/uuid"
)

type CreateEventRequest struct {
	Title         string                     `json:"title" binding:"required"`
	Description   string                     `json:"description"`
	Duration      string                     `json:"duration" binding:"required"`
	ProposedSlots []CreateSlotRequest        `json:"proposed_slots" binding:"required,min=1"`
	Participants  []CreateParticipantRequest `json:"participants" binding:"required,min=1"`
}

type CreateSlotRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Timezone  string `json:"timezone" binding:"required"`
}

type CreateParticipantRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name" binding:"required"`
}

type UpdateEventRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Duration    *string `json:"duration"`
}

type AddSlotRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Timezone  string `json:"timezone" binding:"required"`
}

type UpdateSlotRequest struct {
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Timezone  *string `json:"timezone"`
}

type SubmitAvailabilityRequest struct {
	ParticipantID uuid.UUID                 `json:"participant_id" binding:"required"`
	Slots         []SlotAvailabilityRequest `json:"slots" binding:"required,min=1"`
}

type SlotAvailabilityRequest struct {
	SlotID        uuid.UUID          `json:"slot_id" binding:"required"`
	Status        AvailabilityStatus `json:"status" binding:"required"`
	AvailableFrom *string            `json:"available_from,omitempty"`
	AvailableTo   *string            `json:"available_to,omitempty"`
}

type UpdateAvailabilityRequest struct {
	Status        AvailabilityStatus `json:"status" binding:"required"`
	AvailableFrom *string            `json:"available_from,omitempty"`
	AvailableTo   *string            `json:"available_to,omitempty"`
}

type FinalizeEventRequest struct {
	SlotID uuid.UUID `json:"slot_id" binding:"required"`
}

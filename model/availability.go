package model

import (
	"time"

	"github.com/google/uuid"
)

type AvailabilityStatus string

const (
	AvailabilityStatusAvailable   AvailabilityStatus = "available"
	AvailabilityStatusUnavailable AvailabilityStatus = "unavailable"
	AvailabilityStatusPartial     AvailabilityStatus = "partial"
)

type Availability struct {
	ID            uuid.UUID          `json:"id"`
	EventID       uuid.UUID          `json:"event_id"`
	ParticipantID uuid.UUID          `json:"participant_id"`
	SlotID        uuid.UUID          `json:"slot_id"`
	Status        AvailabilityStatus `json:"status"`
	AvailableFrom *time.Time         `json:"available_from,omitempty"`
	AvailableTo   *time.Time         `json:"available_to,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

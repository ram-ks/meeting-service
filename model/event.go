package model

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string
type ParticipantStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusOpen      EventStatus = "open"
	EventStatusFinalized EventStatus = "finalized"
	EventStatusCancelled EventStatus = "cancelled"
)

const (
	ParticipantStatusPending   ParticipantStatus = "pending"
	ParticipantStatusResponded ParticipantStatus = "responded"
	ParticipantStatusDeclined  ParticipantStatus = "declined"
)

type TimeSlot struct {
	ID        uuid.UUID `json:"id"`
	EventID   uuid.UUID `json:"event_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
}

type Participant struct {
	ID        uuid.UUID         `json:"id"`
	EventID   uuid.UUID         `json:"event_id"`
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	Status    ParticipantStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
}

type Event struct {
	ID              uuid.UUID     `json:"id"`
	Title           string        `json:"title"`
	Description     string        `json:"description,omitempty"`
	OrganizerID     uuid.UUID     `json:"organizer_id"`
	Duration        string        `json:"duration"`
	Status          EventStatus   `json:"status"`
	FinalizedSlotID *uuid.UUID    `json:"finalized_slot_id,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	ProposedSlots   []TimeSlot    `json:"proposed_slots,omitempty"`
	Participants    []Participant `json:"participants,omitempty"`
}

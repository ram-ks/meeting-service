package model

import "github.com/google/uuid"

type Recommendation struct {
	SlotID              uuid.UUID `json:"slot_id"`
	Slot                TimeSlot  `json:"slot"`
	AvailableCount      int       `json:"available_count"`
	TotalParticipants   int       `json:"total_participants"`
	AvailabilityPercent float64   `json:"availability_percent"`
	IsPerfectMatch      bool      `json:"is_perfect_match"`
}

type RecommendationResponse struct {
	EventID      uuid.UUID        `json:"event_id"`
	PerfectSlots []Recommendation `json:"perfect_slots"`
	BestMatches  []Recommendation `json:"best_matches"`
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type PreferredSlot struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Timezone  string    `json:"timezone"`
	DayOfWeek *int      `json:"day_of_week,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePreferredSlotRequest struct {
	Email     string `json:"email" binding:"required,email"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Timezone  string `json:"timezone" binding:"required"`
	DayOfWeek *int   `json:"day_of_week,omitempty"`
}

type UpdatePreferredSlotRequest struct {
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Timezone  *string `json:"timezone"`
	DayOfWeek *int    `json:"day_of_week,omitempty"`
}

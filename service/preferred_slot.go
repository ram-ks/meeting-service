package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/repository"
)

var (
	ErrPreferredSlotNotFound = errors.New("preferred slot not found")
)

type PreferredSlotService interface {
	Create(ctx context.Context, req model.CreatePreferredSlotRequest) (*model.PreferredSlot, error)
	GetByEmail(ctx context.Context, email string) ([]model.PreferredSlot, error)
	Update(ctx context.Context, slotID uuid.UUID, req model.UpdatePreferredSlotRequest) (*model.PreferredSlot, error)
	Delete(ctx context.Context, slotID uuid.UUID) error
}

type preferredSlotService struct {
	repo repository.PreferredSlotRepository
}

func NewPreferredSlotService(repo repository.PreferredSlotRepository) PreferredSlotService {
	return &preferredSlotService{repo: repo}
}

func (s *preferredSlotService) Create(ctx context.Context, req model.CreatePreferredSlotRequest) (*model.PreferredSlot, error) {
	startTime, err := parseTime(req.StartTime, req.Timezone)
	if err != nil {
		return nil, err
	}
	endTime, err := parseTime(req.EndTime, req.Timezone)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	slot := &model.PreferredSlot{
		ID:        uuid.New(),
		Email:     req.Email,
		StartTime: startTime,
		EndTime:   endTime,
		Timezone:  req.Timezone,
		DayOfWeek: req.DayOfWeek,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, slot); err != nil {
		return nil, err
	}

	return slot, nil
}

func (s *preferredSlotService) GetByEmail(ctx context.Context, email string) ([]model.PreferredSlot, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *preferredSlotService) Update(ctx context.Context, slotID uuid.UUID, req model.UpdatePreferredSlotRequest) (*model.PreferredSlot, error) {
	slot, err := s.repo.GetByID(ctx, slotID)
	if err != nil {
		return nil, ErrPreferredSlotNotFound
	}

	timezone := slot.Timezone
	if req.Timezone != nil {
		timezone = *req.Timezone
		slot.Timezone = timezone
	}

	if req.StartTime != nil {
		startTime, err := parseTime(*req.StartTime, timezone)
		if err != nil {
			return nil, err
		}
		slot.StartTime = startTime
	}

	if req.EndTime != nil {
		endTime, err := parseTime(*req.EndTime, timezone)
		if err != nil {
			return nil, err
		}
		slot.EndTime = endTime
	}

	if req.DayOfWeek != nil {
		slot.DayOfWeek = req.DayOfWeek
	}

	if err := s.repo.Update(ctx, slot); err != nil {
		return nil, err
	}

	return slot, nil
}

func (s *preferredSlotService) Delete(ctx context.Context, slotID uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, slotID)
	if err != nil {
		return ErrPreferredSlotNotFound
	}
	return s.repo.Delete(ctx, slotID)
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

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
	ErrEventNotFound        = errors.New("event not found")
	ErrSlotNotFound         = errors.New("slot not found")
	ErrInvalidStatus        = errors.New("invalid event status for this operation")
	ErrSlotNotInEvent       = errors.New("slot does not belong to this event")
	ErrInvalidTimeFormat    = errors.New("invalid time format")
	ErrParticipantNotFound  = errors.New("participant not found")
	ErrAvailabilityNotFound = errors.New("availability not found")
)

type AvailabilityService interface {
	SubmitAvailability(ctx context.Context, eventID uuid.UUID, req model.SubmitAvailabilityRequest) error
	GetAvailability(ctx context.Context, eventID uuid.UUID) ([]model.Availability, error)
	GetParticipantAvailability(ctx context.Context, eventID, participantID uuid.UUID) ([]model.Availability, error)
	UpdateAvailability(ctx context.Context, availabilityID uuid.UUID, req model.UpdateAvailabilityRequest) (*model.Availability, error)
	DeleteAvailability(ctx context.Context, availabilityID uuid.UUID) error
}

type availabilityService struct {
	availRepo repository.AvailabilityRepository
	eventRepo repository.EventRepository
}

func NewAvailabilityService(availRepo repository.AvailabilityRepository, eventRepo repository.EventRepository) AvailabilityService {
	return &availabilityService{
		availRepo: availRepo,
		eventRepo: eventRepo,
	}
}

func (s *availabilityService) SubmitAvailability(ctx context.Context, eventID uuid.UUID, req model.SubmitAvailabilityRequest) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return ErrEventNotFound
	}

	participantFound := false
	for _, p := range event.Participants {
		if p.ID == req.ParticipantID {
			participantFound = true
			break
		}
	}
	if !participantFound {
		return ErrParticipantNotFound
	}

	now := time.Now().UTC()

	for _, slotAvail := range req.Slots {
		slotFound := false
		for _, slot := range event.ProposedSlots {
			if slot.ID == slotAvail.SlotID {
				slotFound = true
				break
			}
		}
		if !slotFound {
			return ErrSlotNotInEvent
		}

		availability := &model.Availability{
			ID:            uuid.New(),
			EventID:       eventID,
			ParticipantID: req.ParticipantID,
			SlotID:        slotAvail.SlotID,
			Status:        slotAvail.Status,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if slotAvail.AvailableFrom != nil {
			t, err := time.Parse(time.RFC3339, *slotAvail.AvailableFrom)
			if err != nil {
				return ErrInvalidTimeFormat
			}
			availability.AvailableFrom = &t
		}
		if slotAvail.AvailableTo != nil {
			t, err := time.Parse(time.RFC3339, *slotAvail.AvailableTo)
			if err != nil {
				return ErrInvalidTimeFormat
			}
			availability.AvailableTo = &t
		}

		if err := s.availRepo.Upsert(ctx, availability); err != nil {
			return err
		}
	}

	if err := s.eventRepo.UpdateParticipantStatus(ctx, req.ParticipantID, model.ParticipantStatusResponded); err != nil {
		return err
	}

	return nil
}

func (s *availabilityService) GetAvailability(ctx context.Context, eventID uuid.UUID) ([]model.Availability, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}

	return s.availRepo.GetByEventID(ctx, eventID)
}

func (s *availabilityService) GetParticipantAvailability(ctx context.Context, eventID, participantID uuid.UUID) ([]model.Availability, error) {
	availabilities, err := s.availRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	var result []model.Availability
	for _, a := range availabilities {
		if a.ParticipantID == participantID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (s *availabilityService) UpdateAvailability(ctx context.Context, availabilityID uuid.UUID, req model.UpdateAvailabilityRequest) (*model.Availability, error) {
	availability, err := s.availRepo.GetByID(ctx, availabilityID)
	if err != nil {
		return nil, ErrAvailabilityNotFound
	}

	availability.Status = req.Status

	if req.AvailableFrom != nil {
		t, err := time.Parse(time.RFC3339, *req.AvailableFrom)
		if err != nil {
			return nil, ErrInvalidTimeFormat
		}
		availability.AvailableFrom = &t
	}
	if req.AvailableTo != nil {
		t, err := time.Parse(time.RFC3339, *req.AvailableTo)
		if err != nil {
			return nil, ErrInvalidTimeFormat
		}
		availability.AvailableTo = &t
	}

	if err := s.availRepo.Update(ctx, availability); err != nil {
		return nil, err
	}

	return availability, nil
}

func (s *availabilityService) DeleteAvailability(ctx context.Context, availabilityID uuid.UUID) error {
	_, err := s.availRepo.GetByID(ctx, availabilityID)
	if err != nil {
		return ErrAvailabilityNotFound
	}
	return s.availRepo.Delete(ctx, availabilityID)
}

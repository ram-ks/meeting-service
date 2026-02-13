package service

import (
	"context"
	"sort"

	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/repository"
)

type SchedulerService interface {
	GetRecommendations(ctx context.Context, eventID uuid.UUID) (*model.RecommendationResponse, error)
}

type schedulerService struct {
	eventRepo repository.EventRepository
	availRepo repository.AvailabilityRepository
}

func NewSchedulerService(eventRepo repository.EventRepository, availRepo repository.AvailabilityRepository) SchedulerService {
	return &schedulerService{
		eventRepo: eventRepo,
		availRepo: availRepo,
	}
}

func (s *schedulerService) GetRecommendations(ctx context.Context, eventID uuid.UUID) (*model.RecommendationResponse, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}

	availabilities, err := s.availRepo.GetByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	availBySlot := make(map[uuid.UUID][]model.Availability)
	for _, a := range availabilities {
		availBySlot[a.SlotID] = append(availBySlot[a.SlotID], a)
	}

	totalParticipants := len(event.Participants)
	var recommendations []model.Recommendation

	for _, slot := range event.ProposedSlots {
		slotAvailabilities := availBySlot[slot.ID]
		availableCount := 0

		for _, a := range slotAvailabilities {
			if a.Status == model.AvailabilityStatusAvailable {
				availableCount++
			} else if a.Status == model.AvailabilityStatusPartial {
				availableCount++
			}
		}

		percent := 0.0
		if totalParticipants > 0 {
			percent = float64(availableCount) / float64(totalParticipants) * 100
		}

		rec := model.Recommendation{
			SlotID:              slot.ID,
			Slot:                slot,
			AvailableCount:      availableCount,
			TotalParticipants:   totalParticipants,
			AvailabilityPercent: percent,
			IsPerfectMatch:      availableCount == totalParticipants && totalParticipants > 0,
		}
		recommendations = append(recommendations, rec)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].IsPerfectMatch != recommendations[j].IsPerfectMatch {
			return recommendations[i].IsPerfectMatch
		}
		return recommendations[i].AvailabilityPercent > recommendations[j].AvailabilityPercent
	})

	response := &model.RecommendationResponse{
		EventID:      eventID,
		PerfectSlots: []model.Recommendation{},
		BestMatches:  []model.Recommendation{},
	}

	for _, rec := range recommendations {
		if rec.IsPerfectMatch {
			response.PerfectSlots = append(response.PerfectSlots, rec)
		} else {
			response.BestMatches = append(response.BestMatches, rec)
		}
	}

	return response, nil
}

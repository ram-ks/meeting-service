package service

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/repository"
)

type SchedulerService interface {
	GetRecommendations(ctx context.Context, eventID uuid.UUID) (*model.RecommendationResponse, error)
}

type schedulerService struct {
	eventRepo         repository.EventRepository
	availRepo         repository.AvailabilityRepository
	preferredSlotRepo repository.PreferredSlotRepository
}

func NewSchedulerService(eventRepo repository.EventRepository, availRepo repository.AvailabilityRepository, preferredSlotRepo repository.PreferredSlotRepository) SchedulerService {
	return &schedulerService{
		eventRepo:         eventRepo,
		availRepo:         availRepo,
		preferredSlotRepo: preferredSlotRepo,
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

	emails := make([]string, len(event.Participants))
	for i, p := range event.Participants {
		emails[i] = p.Email
	}

	prefByEmail := make(map[string][]model.PreferredSlot)
	if len(emails) > 0 {
		preferredSlots, err := s.preferredSlotRepo.GetByEmails(ctx, emails)
		if err == nil {
			for _, ps := range preferredSlots {
				key := strings.ToLower(ps.Email)
				prefByEmail[key] = append(prefByEmail[key], ps)
			}
		}
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
		preferredCount := 0

		respondedParticipants := make(map[uuid.UUID]bool)
		for _, a := range slotAvailabilities {
			respondedParticipants[a.ParticipantID] = true
			if a.Status == model.AvailabilityStatusAvailable {
				availableCount++
			} else if a.Status == model.AvailabilityStatusPartial {
				availableCount++
			}
		}

		for _, p := range event.Participants {
			prefs := prefByEmail[strings.ToLower(p.Email)]
			for _, pref := range prefs {
				if slotOverlapsPreference(slot, pref) {
					preferredCount++
					break
				}
			}
		}

		percent := 0.0
		if totalParticipants > 0 {
			percent = float64(availableCount) / float64(totalParticipants) * 100
		}

		preferredPercent := 0.0
		if totalParticipants > 0 {
			preferredPercent = float64(preferredCount) / float64(totalParticipants) * 100
		}

		rec := model.Recommendation{
			SlotID:              slot.ID,
			Slot:                slot,
			AvailableCount:      availableCount,
			TotalParticipants:   totalParticipants,
			AvailabilityPercent: percent,
			PreferredCount:      preferredCount,
			PreferredPercent:    preferredPercent,
			IsPerfectMatch:      availableCount == totalParticipants && totalParticipants > 0,
		}
		recommendations = append(recommendations, rec)
	}

	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].IsPerfectMatch != recommendations[j].IsPerfectMatch {
			return recommendations[i].IsPerfectMatch
		}
		if recommendations[i].AvailabilityPercent != recommendations[j].AvailabilityPercent {
			return recommendations[i].AvailabilityPercent > recommendations[j].AvailabilityPercent
		}
		return recommendations[i].PreferredPercent > recommendations[j].PreferredPercent
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

func slotOverlapsPreference(slot model.TimeSlot, pref model.PreferredSlot) bool {
	if pref.DayOfWeek != nil {
		slotDay := int(slot.StartTime.Weekday())
		if slotDay != *pref.DayOfWeek {
			return false
		}
	}

	slotStart := slot.StartTime.Hour()*60 + slot.StartTime.Minute()
	slotEnd := slot.EndTime.Hour()*60 + slot.EndTime.Minute()
	prefStart := pref.StartTime.Hour()*60 + pref.StartTime.Minute()
	prefEnd := pref.EndTime.Hour()*60 + pref.EndTime.Minute()

	return slotStart >= prefStart && slotEnd <= prefEnd
}

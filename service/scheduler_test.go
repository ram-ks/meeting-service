package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *model.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Event), args.Error(1)
}

func (m *MockEventRepository) List(ctx context.Context, organizerID uuid.UUID) ([]model.Event, error) {
	args := m.Called(ctx, organizerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Event), args.Error(1)
}

func (m *MockEventRepository) Update(ctx context.Context, event *model.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) CreateSlot(ctx context.Context, slot *model.TimeSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockEventRepository) GetSlotsByEventID(ctx context.Context, eventID uuid.UUID) ([]model.TimeSlot, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.TimeSlot), args.Error(1)
}

func (m *MockEventRepository) GetSlotByID(ctx context.Context, id uuid.UUID) (*model.TimeSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TimeSlot), args.Error(1)
}

func (m *MockEventRepository) UpdateSlot(ctx context.Context, slot *model.TimeSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockEventRepository) DeleteSlot(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) CreateParticipant(ctx context.Context, participant *model.Participant) error {
	args := m.Called(ctx, participant)
	return args.Error(0)
}

func (m *MockEventRepository) GetParticipantsByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Participant, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Participant), args.Error(1)
}

func (m *MockEventRepository) GetParticipantByID(ctx context.Context, id uuid.UUID) (*model.Participant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Participant), args.Error(1)
}

func (m *MockEventRepository) UpdateParticipantStatus(ctx context.Context, id uuid.UUID, status model.ParticipantStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockAvailabilityRepository struct {
	mock.Mock
}

func (m *MockAvailabilityRepository) Create(ctx context.Context, availability *model.Availability) error {
	args := m.Called(ctx, availability)
	return args.Error(0)
}

func (m *MockAvailabilityRepository) Upsert(ctx context.Context, availability *model.Availability) error {
	args := m.Called(ctx, availability)
	return args.Error(0)
}

func (m *MockAvailabilityRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Availability, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Availability), args.Error(1)
}

func (m *MockAvailabilityRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Availability, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Availability), args.Error(1)
}

func (m *MockAvailabilityRepository) Update(ctx context.Context, availability *model.Availability) error {
	args := m.Called(ctx, availability)
	return args.Error(0)
}

func (m *MockAvailabilityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockPreferredSlotRepository struct {
	mock.Mock
}

func (m *MockPreferredSlotRepository) Create(ctx context.Context, slot *model.PreferredSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockPreferredSlotRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.PreferredSlot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PreferredSlot), args.Error(1)
}

func (m *MockPreferredSlotRepository) GetByEmail(ctx context.Context, email string) ([]model.PreferredSlot, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PreferredSlot), args.Error(1)
}

func (m *MockPreferredSlotRepository) GetByEmails(ctx context.Context, emails []string) ([]model.PreferredSlot, error) {
	args := m.Called(ctx, emails)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PreferredSlot), args.Error(1)
}

func (m *MockPreferredSlotRepository) Update(ctx context.Context, slot *model.PreferredSlot) error {
	args := m.Called(ctx, slot)
	return args.Error(0)
}

func (m *MockPreferredSlotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSchedulerServiceSuite(t *testing.T) {
	t.Run("GetRecommendations_EventNotFound", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(nil, errors.New("not found"))

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.Nil(t, result)
		assert.Equal(t, ErrEventNotFound, err)
		mockEventRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_NoParticipants_NoSlots", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		event := &model.Event{
			ID:            eventID,
			Participants:  []model.Participant{},
			ProposedSlots: []model.TimeSlot{},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return([]model.Availability{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, eventID, result.EventID)
		assert.Empty(t, result.PerfectSlots)
		assert.Empty(t, result.BestMatches)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_NoAvailabilitySubmitted", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return([]model.Availability{}, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, []string{"alice@example.com", "bob@example.com"}).Return([]model.PreferredSlot{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.PerfectSlots)
		assert.Len(t, result.BestMatches, 1)
		assert.Equal(t, 0, result.BestMatches[0].AvailableCount)
		assert.Equal(t, 2, result.BestMatches[0].TotalParticipants)
		assert.Equal(t, float64(0), result.BestMatches[0].AvailabilityPercent)
		assert.False(t, result.BestMatches[0].IsPerfectMatch)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_PartialAvailability", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()
		participant3 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
				{ID: participant3, Email: "charlie@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant2, SlotID: slotID, Status: model.AvailabilityStatusUnavailable},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return([]model.PreferredSlot{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.PerfectSlots)
		assert.Len(t, result.BestMatches, 1)
		assert.Equal(t, 1, result.BestMatches[0].AvailableCount)
		assert.Equal(t, 3, result.BestMatches[0].TotalParticipants)
		assert.InDelta(t, 33.33, result.BestMatches[0].AvailabilityPercent, 0.01)
		assert.False(t, result.BestMatches[0].IsPerfectMatch)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_PerfectMatch", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant2, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return([]model.PreferredSlot{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Empty(t, result.BestMatches)
		assert.Equal(t, 2, result.PerfectSlots[0].AvailableCount)
		assert.Equal(t, 2, result.PerfectSlots[0].TotalParticipants)
		assert.Equal(t, float64(100), result.PerfectSlots[0].AvailabilityPercent)
		assert.True(t, result.PerfectSlots[0].IsPerfectMatch)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_PartialStatusCountsAsAvailable", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant2, SlotID: slotID, Status: model.AvailabilityStatusPartial},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return([]model.PreferredSlot{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Equal(t, 2, result.PerfectSlots[0].AvailableCount)
		assert.True(t, result.PerfectSlots[0].IsPerfectMatch)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_WithPreferredSlots", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()

		now := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		preferredSlots := []model.PreferredSlot{
			{ID: uuid.New(), Email: "alice@example.com", StartTime: now.Add(-time.Hour), EndTime: now.Add(2 * time.Hour)},
			{ID: uuid.New(), Email: "bob@example.com", StartTime: now.Add(-time.Hour), EndTime: now.Add(2 * time.Hour)},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return(preferredSlots, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.BestMatches, 1)
		assert.Equal(t, 2, result.BestMatches[0].PreferredCount)
		assert.Equal(t, float64(100), result.BestMatches[0].PreferredPercent)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_NoPreferredSlots", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return([]model.PreferredSlot{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Equal(t, 0, result.PerfectSlots[0].PreferredCount)
		assert.Equal(t, float64(0), result.PerfectSlots[0].PreferredPercent)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_PreferredSlotsQueryFails_StillWorks", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Equal(t, 0, result.PerfectSlots[0].PreferredCount)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_MultipleSlots_SortedByAvailability", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID1 := uuid.New()
		slotID2 := uuid.New()
		slotID3 := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()
		participant3 := uuid.New()

		now := time.Now()
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
				{ID: participant3, Email: "charlie@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID1, StartTime: now, EndTime: now.Add(time.Hour)},
				{ID: slotID2, StartTime: now.Add(2 * time.Hour), EndTime: now.Add(3 * time.Hour)},
				{ID: slotID3, StartTime: now.Add(4 * time.Hour), EndTime: now.Add(5 * time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID1, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID2, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant2, SlotID: slotID2, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID3, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant2, SlotID: slotID3, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant3, SlotID: slotID3, Status: model.AvailabilityStatusAvailable},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return([]model.PreferredSlot{}, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Len(t, result.BestMatches, 2)

		assert.Equal(t, slotID3, result.PerfectSlots[0].SlotID)
		assert.True(t, result.PerfectSlots[0].IsPerfectMatch)

		assert.Greater(t, result.BestMatches[0].AvailabilityPercent, result.BestMatches[1].AvailabilityPercent)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_SortByPreferredWhenAvailabilityEqual", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID1 := uuid.New()
		slotID2 := uuid.New()
		participant1 := uuid.New()
		participant2 := uuid.New()

		now := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
				{ID: participant2, Email: "bob@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID1, StartTime: now, EndTime: now.Add(time.Hour)},
				{ID: slotID2, StartTime: now.Add(2 * time.Hour), EndTime: now.Add(3 * time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID1, Status: model.AvailabilityStatusAvailable},
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID2, Status: model.AvailabilityStatusAvailable},
		}

		preferredSlots := []model.PreferredSlot{
			{ID: uuid.New(), Email: "alice@example.com", StartTime: now.Add(-time.Hour), EndTime: now.Add(2 * time.Hour)},
			{ID: uuid.New(), Email: "bob@example.com", StartTime: now.Add(-time.Hour), EndTime: now.Add(2 * time.Hour)},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return(preferredSlots, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.BestMatches, 2)

		assert.Equal(t, result.BestMatches[0].AvailabilityPercent, result.BestMatches[1].AvailabilityPercent)
		assert.Greater(t, result.BestMatches[0].PreferredPercent, result.BestMatches[1].PreferredPercent)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_PreferredSlotWithDayOfWeek", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()

		friday := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
		fridayDayOfWeek := int(friday.Weekday())

		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: friday, EndTime: friday.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		preferredSlots := []model.PreferredSlot{
			{ID: uuid.New(), Email: "alice@example.com", StartTime: friday.Add(-time.Hour), EndTime: friday.Add(2 * time.Hour), DayOfWeek: &fridayDayOfWeek},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return(preferredSlots, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Equal(t, 1, result.PerfectSlots[0].PreferredCount)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_PreferredSlotDayOfWeekMismatch", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()

		friday := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
		monday := 1

		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "alice@example.com"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: friday, EndTime: friday.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		preferredSlots := []model.PreferredSlot{
			{ID: uuid.New(), Email: "alice@example.com", StartTime: friday.Add(-time.Hour), EndTime: friday.Add(2 * time.Hour), DayOfWeek: &monday},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return(preferredSlots, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Equal(t, 0, result.PerfectSlots[0].PreferredCount)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_AvailabilityRepoError", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		event := &model.Event{ID: eventID}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(nil, errors.New("database error"))

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.Nil(t, result)
		assert.Error(t, err)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
	})

	t.Run("GetRecommendations_CaseInsensitiveEmailMatching", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		eventID := uuid.New()
		slotID := uuid.New()
		participant1 := uuid.New()

		now := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
		event := &model.Event{
			ID: eventID,
			Participants: []model.Participant{
				{ID: participant1, Email: "Alice@Example.COM"},
			},
			ProposedSlots: []model.TimeSlot{
				{ID: slotID, StartTime: now, EndTime: now.Add(time.Hour)},
			},
		}

		availabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participant1, SlotID: slotID, Status: model.AvailabilityStatusAvailable},
		}

		preferredSlots := []model.PreferredSlot{
			{ID: uuid.New(), Email: "alice@example.com", StartTime: now.Add(-time.Hour), EndTime: now.Add(2 * time.Hour)},
		}

		mockEventRepo.On("GetByID", mock.Anything, eventID).Return(event, nil)
		mockAvailRepo.On("GetByEventID", mock.Anything, eventID).Return(availabilities, nil)
		mockPrefRepo.On("GetByEmails", mock.Anything, mock.Anything).Return(preferredSlots, nil)

		result, err := svc.GetRecommendations(context.Background(), eventID)

		assert.NoError(t, err)
		assert.Len(t, result.PerfectSlots, 1)
		assert.Equal(t, 1, result.PerfectSlots[0].PreferredCount)

		mockEventRepo.AssertExpectations(t)
		mockAvailRepo.AssertExpectations(t)
		mockPrefRepo.AssertExpectations(t)
	})

	t.Run("NewSchedulerService", func(t *testing.T) {
		mockEventRepo := new(MockEventRepository)
		mockAvailRepo := new(MockAvailabilityRepository)
		mockPrefRepo := new(MockPreferredSlotRepository)

		svc := NewSchedulerService(mockEventRepo, mockAvailRepo, mockPrefRepo)

		assert.NotNil(t, svc)
	})
}

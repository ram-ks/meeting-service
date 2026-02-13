package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/ram-ks/meeting-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSchedulerService struct {
	mock.Mock
}

func (m *MockSchedulerService) GetRecommendations(ctx context.Context, eventID uuid.UUID) (*model.RecommendationResponse, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RecommendationResponse), args.Error(1)
}

func setupRecommendationTestRouter(ctrl *RecommendationController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	events := router.Group("/events")
	{
		events.GET("/:id/recommendations", ctrl.GetRecommendations)
	}

	return router
}

func TestRecommendationControllerSuite(t *testing.T) {
	t.Run("GetRecommendations_Success_EmptyRecommendations", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()

		expectedResponse := &model.RecommendationResponse{
			EventID:      eventID,
			PerfectSlots: []model.Recommendation{},
			BestMatches:  []model.Recommendation{},
		}

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.RecommendationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, eventID, response.EventID)
		assert.Empty(t, response.PerfectSlots)
		assert.Empty(t, response.BestMatches)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_Success_WithPerfectSlots", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()
		slotID := uuid.New()
		now := time.Now()

		expectedResponse := &model.RecommendationResponse{
			EventID: eventID,
			PerfectSlots: []model.Recommendation{
				{
					SlotID: slotID,
					Slot: model.TimeSlot{
						ID:        slotID,
						EventID:   eventID,
						StartTime: now,
						EndTime:   now.Add(time.Hour),
						Timezone:  "UTC",
					},
					AvailableCount:      3,
					TotalParticipants:   3,
					AvailabilityPercent: 100,
					PreferredCount:      2,
					PreferredPercent:    66.67,
					IsPerfectMatch:      true,
				},
			},
			BestMatches: []model.Recommendation{},
		}

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.RecommendationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.PerfectSlots, 1)
		assert.Empty(t, response.BestMatches)
		assert.True(t, response.PerfectSlots[0].IsPerfectMatch)
		assert.Equal(t, float64(100), response.PerfectSlots[0].AvailabilityPercent)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_Success_WithBestMatches", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()
		slotID1 := uuid.New()
		slotID2 := uuid.New()

		expectedResponse := &model.RecommendationResponse{
			EventID:      eventID,
			PerfectSlots: []model.Recommendation{},
			BestMatches: []model.Recommendation{
				{
					SlotID:              slotID1,
					AvailableCount:      2,
					TotalParticipants:   3,
					AvailabilityPercent: 66.67,
					PreferredCount:      2,
					PreferredPercent:    66.67,
					IsPerfectMatch:      false,
				},
				{
					SlotID:              slotID2,
					AvailableCount:      1,
					TotalParticipants:   3,
					AvailabilityPercent: 33.33,
					PreferredCount:      0,
					PreferredPercent:    0,
					IsPerfectMatch:      false,
				},
			},
		}

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.RecommendationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Empty(t, response.PerfectSlots)
		assert.Len(t, response.BestMatches, 2)
		assert.Greater(t, response.BestMatches[0].AvailabilityPercent, response.BestMatches[1].AvailabilityPercent)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_Success_MixedPerfectAndBestMatches", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()
		perfectSlotID := uuid.New()
		partialSlotID := uuid.New()

		expectedResponse := &model.RecommendationResponse{
			EventID: eventID,
			PerfectSlots: []model.Recommendation{
				{
					SlotID:              perfectSlotID,
					AvailableCount:      3,
					TotalParticipants:   3,
					AvailabilityPercent: 100,
					PreferredCount:      1,
					PreferredPercent:    33.33,
					IsPerfectMatch:      true,
				},
			},
			BestMatches: []model.Recommendation{
				{
					SlotID:              partialSlotID,
					AvailableCount:      2,
					TotalParticipants:   3,
					AvailabilityPercent: 66.67,
					PreferredCount:      3,
					PreferredPercent:    100,
					IsPerfectMatch:      false,
				},
			},
		}

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.RecommendationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.PerfectSlots, 1)
		assert.Len(t, response.BestMatches, 1)
		assert.True(t, response.PerfectSlots[0].IsPerfectMatch)
		assert.False(t, response.BestMatches[0].IsPerfectMatch)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_Success_WithPreferredData", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()
		slotID := uuid.New()

		expectedResponse := &model.RecommendationResponse{
			EventID:      eventID,
			PerfectSlots: []model.Recommendation{},
			BestMatches: []model.Recommendation{
				{
					SlotID:              slotID,
					AvailableCount:      1,
					TotalParticipants:   2,
					AvailabilityPercent: 50,
					PreferredCount:      2,
					PreferredPercent:    100,
					IsPerfectMatch:      false,
				},
			},
		}

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.RecommendationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, response.BestMatches[0].PreferredCount)
		assert.Equal(t, float64(100), response.BestMatches[0].PreferredPercent)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_InvalidEventID", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/invalid-uuid/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid event id")
	})

	t.Run("GetRecommendations_EmptyEventID", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events//recommendations", nil)

		router.ServeHTTP(w, httpReq)

		// Empty path segment is treated as invalid UUID, returns 400
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetRecommendations_EventNotFound", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()

		// service.ErrEventNotFound is not handled by handleServiceError, falls through to 500
		mockService.On("GetRecommendations", mock.Anything, eventID).Return(nil, service.ErrEventNotFound)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		// Currently returns 500 as handleServiceError doesn't recognize service.ErrEventNotFound
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_ServiceError", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(nil, errors.New("database error"))

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("GetRecommendations_ResponseContainsSlotDetails", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)
		router := setupRecommendationTestRouter(ctrl)

		eventID := uuid.New()
		slotID := uuid.New()
		now := time.Now().UTC().Truncate(time.Second)

		expectedResponse := &model.RecommendationResponse{
			EventID: eventID,
			PerfectSlots: []model.Recommendation{
				{
					SlotID: slotID,
					Slot: model.TimeSlot{
						ID:        slotID,
						EventID:   eventID,
						StartTime: now,
						EndTime:   now.Add(time.Hour),
						Timezone:  "America/New_York",
					},
					AvailableCount:      2,
					TotalParticipants:   2,
					AvailabilityPercent: 100,
					IsPerfectMatch:      true,
				},
			},
			BestMatches: []model.Recommendation{},
		}

		mockService.On("GetRecommendations", mock.Anything, eventID).Return(expectedResponse, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/recommendations", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response model.RecommendationResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, slotID, response.PerfectSlots[0].Slot.ID)
		assert.Equal(t, eventID, response.PerfectSlots[0].Slot.EventID)
		assert.Equal(t, "America/New_York", response.PerfectSlots[0].Slot.Timezone)

		mockService.AssertExpectations(t)
	})

	t.Run("NewRecommendationController", func(t *testing.T) {
		mockService := new(MockSchedulerService)
		ctrl := NewRecommendationController(mockService)

		assert.NotNil(t, ctrl)
		assert.Equal(t, mockService, ctrl.schedulerService)
	})
}

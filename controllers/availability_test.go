package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ram-ks/meeting-service/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAvailabilityService is a mock implementation of AvailabilityService
type MockAvailabilityService struct {
	mock.Mock
}

func (m *MockAvailabilityService) SubmitAvailability(ctx context.Context, eventID uuid.UUID, req model.SubmitAvailabilityRequest) error {
	args := m.Called(ctx, eventID, req)
	return args.Error(0)
}

func (m *MockAvailabilityService) GetAvailability(ctx context.Context, eventID uuid.UUID) ([]model.Availability, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Availability), args.Error(1)
}

func (m *MockAvailabilityService) GetParticipantAvailability(ctx context.Context, eventID, participantID uuid.UUID) ([]model.Availability, error) {
	args := m.Called(ctx, eventID, participantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Availability), args.Error(1)
}

func (m *MockAvailabilityService) UpdateAvailability(ctx context.Context, availabilityID uuid.UUID, req model.UpdateAvailabilityRequest) (*model.Availability, error) {
	args := m.Called(ctx, availabilityID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Availability), args.Error(1)
}

func (m *MockAvailabilityService) DeleteAvailability(ctx context.Context, availabilityID uuid.UUID) error {
	args := m.Called(ctx, availabilityID)
	return args.Error(0)
}

// setupTestRouter creates a test router with the controller
func setupTestRouter(ctrl *AvailabilityController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	events := router.Group("/events")
	availability := events.Group("/:id/availability")
	{
		availability.POST("", ctrl.SubmitAvailability)
		availability.GET("", ctrl.GetAvailability)
		availability.GET("/:participant_id", ctrl.GetParticipantAvailability)
		availability.PUT("/:availability_id", ctrl.UpdateAvailability)
		availability.DELETE("/:availability_id", ctrl.DeleteAvailability)
	}

	return router
}

func TestAvailabilityControllerSuite(t *testing.T) {
	t.Run("SubmitAvailability_Success", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		from := "2026-02-12T10:00:00Z"
		to := "2026-02-12T11:00:00Z"
		req := model.SubmitAvailabilityRequest{
			ParticipantID: uuid.New(),
			Slots: []model.SlotAvailabilityRequest{
				{SlotID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
				{Status: model.AvailabilityStatusAvailable},
				{AvailableFrom: &from},
				{AvailableTo: &to},
			},
		}

		mockService.On("SubmitAvailability", mock.Anything, eventID, req).Return(nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/events/"+eventID.String()+"/availability", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "availability submitted successfully")
		mockService.AssertExpectations(t)
	})

	t.Run("SubmitAvailability_InvalidEventID", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		req := model.SubmitAvailabilityRequest{
			ParticipantID: uuid.New(),
		}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/events/invalid-uuid/availability", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid")
	})

	t.Run("SubmitAvailability_InvalidJSON", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/events/"+eventID.String()+"/availability", bytes.NewBuffer([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("SubmitAvailability_ServiceError", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		from := "2026-02-12T10:00:00Z"
		to := "2026-02-12T11:00:00Z"
		req := model.SubmitAvailabilityRequest{
			ParticipantID: uuid.New(),
			Slots: []model.SlotAvailabilityRequest{
				{SlotID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
				{Status: model.AvailabilityStatusAvailable},
				{AvailableFrom: &from},
				{AvailableTo: &to},
			},
		}

		mockService.On("SubmitAvailability", mock.Anything, eventID, req).Return(errors.New("service error"))

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("POST", "/events/"+eventID.String()+"/availability", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		// Status code depends on handleServiceError implementation
		assert.NotEqual(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GetAvailability_Success", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		expectedAvailabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID},
			{ID: uuid.New(), EventID: eventID},
		}

		mockService.On("GetAvailability", mock.Anything, eventID).Return(expectedAvailabilities, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/availability", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "availabilities")
		mockService.AssertExpectations(t)
	})

	t.Run("GetAvailability_InvalidEventId", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/invalid-uuid/availability", nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid")
	})

	t.Run("GetAvailability_ServiceError", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		mockService.On("GetAvailability", mock.Anything, eventID).Return(nil, errors.New("service error"))

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/availability", nil)

		router.ServeHTTP(w, httpReq)

		assert.NotEqual(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("GetParticipantAvailability_Success", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		participantID := uuid.New()
		expectedAvailabilities := []model.Availability{
			{ID: uuid.New(), EventID: eventID, ParticipantID: participantID},
		}

		mockService.On("GetParticipantAvailability", mock.Anything, eventID, participantID).Return(expectedAvailabilities, nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/"+eventID.String()+"/availability/"+participantID.String(), nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "availabilities")
		mockService.AssertExpectations(t)
	})

	t.Run("GetParticipantAvailability_InvalidEventID", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		participantID := uuid.New()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest("GET", "/events/invalid-uuid/availability/"+participantID.String(), nil)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetParticipantAvailability_InvalidParticipantID", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"GET",
			"/events/"+eventID.String()+"/availability/invalid-uuid",
			nil,
		)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UpdateAvailability_Success", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		availabilityID := uuid.New()

		availableFrom := "2024-01-01T09:00:00Z"
		availableTo := "2024-01-01T17:00:00Z"

		req := model.UpdateAvailabilityRequest{
			Status:        model.AvailabilityStatus("available"),
			AvailableFrom: &availableFrom,
			AvailableTo:   &availableTo,
		}

		updatedAvailability := &model.Availability{
			ID:      availabilityID,
			EventID: eventID,
		}

		mockService.
			On("UpdateAvailability", mock.Anything, availabilityID, req).
			Return(updatedAvailability, nil)

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"PUT",
			"/events/"+eventID.String()+"/availability/"+availabilityID.String(),
			bytes.NewBuffer(body),
		)
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("UpdateAvailability_InvalidAvailabilityID", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		req := model.UpdateAvailabilityRequest{}

		body, _ := json.Marshal(req)
		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"PUT",
			"/events/"+eventID.String()+"/availability/invalid-uuid",
			bytes.NewBuffer(body),
		)
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UpdateAvailability_InvalidJSON", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		availabilityID := uuid.New()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"PUT",
			"/events/"+eventID.String()+"/availability/"+availabilityID.String(),
			bytes.NewBuffer([]byte("invalid json")),
		)
		httpReq.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteAvailability_Success", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		availabilityID := uuid.New()

		mockService.
			On("DeleteAvailability", mock.Anything, availabilityID).
			Return(nil)

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"DELETE",
			"/events/"+eventID.String()+"/availability/"+availabilityID.String(),
			nil,
		)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("DeleteAvailability_InvalidAvailabilityID", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"DELETE",
			"/events/"+eventID.String()+"/availability/invalid-uuid",
			nil,
		)

		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteAvailability_ServiceError", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)
		router := setupTestRouter(ctrl)

		eventID := uuid.New()
		availabilityID := uuid.New()

		mockService.
			On("DeleteAvailability", mock.Anything, availabilityID).
			Return(errors.New("service error"))

		w := httptest.NewRecorder()
		httpReq, _ := http.NewRequest(
			"DELETE",
			"/events/"+eventID.String()+"/availability/"+availabilityID.String(),
			nil,
		)

		router.ServeHTTP(w, httpReq)

		assert.NotEqual(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("NewAvailabilityController", func(t *testing.T) {
		mockService := new(MockAvailabilityService)
		ctrl := NewAvailabilityController(mockService)

		assert.NotNil(t, ctrl)
		assert.Equal(t, mockService, ctrl.availService)
	})
}

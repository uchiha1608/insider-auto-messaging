package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"insider-auto-messaging/controller"
	"insider-auto-messaging/model"
)

type MockScheduler struct {
	mock.Mock
	running bool
}

func (m *MockScheduler) Start() {
	m.Called()
	m.running = true
}

func (m *MockScheduler) Stop() {
	m.Called()
	m.running = false
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetUnsentMessages(limit int) ([]model.Message, error) {
	args := m.Called(limit)
	return args.Get(0).([]model.Message), args.Error(1)
}

func (m *MockRepo) MarkAsSent(id int64, messageId string) error {
	args := m.Called(id, messageId)
	return args.Error(0)
}

func (m *MockRepo) GetAllSent() ([]model.Message, error) {
	args := m.Called()
	return args.Get(0).([]model.Message), args.Error(1)
}

func TestStartHandler(t *testing.T) {
	mockScheduler := new(MockScheduler)
	mockScheduler.On("Start").Return()

	ctrl := &controller.MessageController{
		Scheduler: mockScheduler,
		Repo:      &MockRepo{},
	}

	req := httptest.NewRequest(http.MethodGet, "/start", nil)
	w := httptest.NewRecorder()

	ctrl.Start(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockScheduler.AssertCalled(t, "Start")
}

func TestStopHandler(t *testing.T) {
	mockScheduler := new(MockScheduler)
	mockScheduler.On("Stop").Return()

	ctrl := &controller.MessageController{
		Scheduler: mockScheduler,
		Repo:      &MockRepo{},
	}

	req := httptest.NewRequest(http.MethodGet, "/stop", nil)
	w := httptest.NewRecorder()

	ctrl.Stop(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockScheduler.AssertCalled(t, "Stop")
}

func TestSentMessagesHandler_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	expectedMessages := []model.Message{
		{
			ID:        1,
			To:        "+123456789",
			Content:   "Hi",
			IsSent:    true,
			MessageID: "abc-123",
			SentAt:    "2025-08-01T12:00:00Z",
		},
	}
	mockRepo.On("GetAllSent").Return(expectedMessages, nil)

	ctrl := &controller.MessageController{
		Scheduler: &MockScheduler{},
		Repo:      mockRepo,
	}

	req := httptest.NewRequest(http.MethodGet, "/sent", nil)
	w := httptest.NewRecorder()

	ctrl.SentMessages(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actual []model.Message
	err := json.NewDecoder(resp.Body).Decode(&actual)
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, actual)

	mockRepo.AssertCalled(t, "GetAllSent")
}

func TestSentMessagesHandler_DBError(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRepo.On("GetAllSent").Return([]model.Message{}, assert.AnError)

	ctrl := &controller.MessageController{
		Scheduler: &MockScheduler{},
		Repo:      mockRepo,
	}

	req := httptest.NewRequest(http.MethodGet, "/sent", nil)
	w := httptest.NewRecorder()

	ctrl.SentMessages(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

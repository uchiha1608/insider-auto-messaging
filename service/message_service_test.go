package service_test

import (
	"bytes"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"testing"

	"insider-auto-messaging/model"
	"insider-auto-messaging/service"
)

// --- Mocks ---

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

// --- Custom HTTP Client (Mock) ---

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// --- Tests ---

func TestSendPendingMessages_Success(t *testing.T) {
	mockRepo := new(MockRepo)
	redisClient, redisMock := redismock.NewClientMock()
	mockMsg := model.Message{
		ID:      1,
		To:      "+123456789",
		Content: "Hello",
	}

	// Mock repository
	mockRepo.On("GetUnsentMessages", 2).Return([]model.Message{mockMsg}, nil)
	mockRepo.On("MarkAsSent", mockMsg.ID, "msg-id-123").Return(nil)

	// Mock Redis SET
	redisMock.ExpectSet("msg-id-123", mock.Anything, 0).SetVal("OK")

	// Mock HTTP
	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			body := `{"message":"Accepted","messageId":"msg-id-123"}`
			return &http.Response{
				StatusCode: http.StatusAccepted,
				Body:       io.NopCloser(bytes.NewBufferString(body)),
			}, nil
		},
	}

	svc := &service.MessageService{
		Repo:       mockRepo,
		Redis:      redisClient,
		HTTPClient: mockHTTP,
	}

	svc.SendPendingMessages()

	mockRepo.AssertExpectations(t)
}

func TestSendPendingMessages_HTTPError(t *testing.T) {
	mockRepo := new(MockRepo)
	redisClient, _ := redismock.NewClientMock()
	mockMsg := model.Message{ID: 1, To: "+1", Content: "fail"}

	mockRepo.On("GetUnsentMessages", 2).Return([]model.Message{mockMsg}, nil)

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	svc := &service.MessageService{
		Repo:       mockRepo,
		Redis:      redisClient,
		HTTPClient: mockHTTP,
	}

	svc.SendPendingMessages()
	mockRepo.AssertExpectations(t)
}

func TestSendPendingMessages_BadStatusCode(t *testing.T) {
	mockRepo := new(MockRepo)
	redisClient, _ := redismock.NewClientMock()
	mockMsg := model.Message{ID: 2, To: "+2", Content: "fail"}

	mockRepo.On("GetUnsentMessages", 2).Return([]model.Message{mockMsg}, nil)

	mockHTTP := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString("")),
			}, nil
		},
	}

	svc := &service.MessageService{
		Repo:       mockRepo,
		Redis:      redisClient,
		HTTPClient: mockHTTP,
	}

	svc.SendPendingMessages()
	mockRepo.AssertExpectations(t)
}

func TestSendPendingMessages_EmptyQueue(t *testing.T) {
	mockRepo := new(MockRepo)
	redisClient, _ := redismock.NewClientMock()

	mockRepo.On("GetUnsentMessages", 2).Return([]model.Message{}, nil)

	svc := &service.MessageService{
		Repo:       mockRepo,
		Redis:      redisClient,
		HTTPClient: &MockHTTPClient{},
	}

	svc.SendPendingMessages()
	mockRepo.AssertExpectations(t)
}

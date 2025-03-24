package usecase

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// MockMatchRepository is a mock implementation of the MatchRepository interface
type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) SaveLike(fromUserID, toUserID int64) error {
	args := m.Called(fromUserID, toUserID)
	return args.Error(0)
}

func (m *MockMatchRepository) CheckMatch(fromUserID, toUserID int64) (bool, error) {
	args := m.Called(fromUserID, toUserID)
	return args.Bool(0), args.Error(1)
}

// MockMatchKafka is a mock implementation of the MatchKafka interface
type MockMatchKafka struct {
	mock.Mock
}

func (m *MockMatchKafka) SendMessage(topic string, message string) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

func (m *MockMatchKafka) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestExample(t *testing.T) {
	// Example test case using the mocks
	repo := new(MockMatchRepository)
	kafka := new(MockMatchKafka)

	// Set up expectations
	repo.On("SaveLike", int64(1), int64(2)).Return(nil)
	repo.On("CheckMatch", int64(1), int64(2)).Return(true, nil)
	kafka.On("SendMessage", "topic", "message").Return(nil)
	kafka.On("Close").Return(nil)

	// Call the methods
	err := repo.SaveLike(1, 2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	match, err := repo.CheckMatch(1, 2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !match {
		t.Errorf("expected match to be true, got %v", match)
	}

	err = kafka.SendMessage("topic", "message")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = kafka.Close()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Assert expectations
	repo.AssertExpectations(t)
	kafka.AssertExpectations(t)
}

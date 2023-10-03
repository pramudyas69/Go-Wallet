package service

import (
	"context"
	"e-wallet/domain"
	mocks "e-wallet/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFindByUser(t *testing.T) {
	mockRepo := new(mocks.MockNotificationRepository)
	service := NewNotification(mockRepo)

	userID := int64(1)
	notifications := []domain.Notification{
		{ID: 1, UserID: userID, Title: "Title 1", Body: "Body 1", Status: 1, IsRead: 0},
		{ID: 2, UserID: userID, Title: "Title 2", Body: "Body 2", Status: 1, IsRead: 0},
	}

	mockRepo.On("FindByUser", mock.Anything, userID).Return(notifications, nil)
	result, err := service.FindByUser(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, result, len(notifications))
	assert.Equal(t, notifications[0].ID, result[0].ID)
	assert.Equal(t, notifications[1].Title, result[1].Title)
	mockRepo.AssertCalled(t, "FindByUser", mock.Anything, userID)
}

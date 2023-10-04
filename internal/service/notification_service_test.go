package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	mocks "e-wallet/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFindByUser(t *testing.T) {
	mockRepo := new(mocks.MockNotificationRepository)
	mockTemplate := new(mocks.MockTemplateRepository)
	mockHub := &dto.Hub{}

	service := NewNotification(mockRepo, mockTemplate, mockHub)

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

func TestInsert(t *testing.T) {
	mockRepo := new(mocks.MockNotificationRepository)
	mockTemplate := new(mocks.MockTemplateRepository)
	mockHub := &dto.Hub{}

	service := NewNotification(mockRepo, mockTemplate, mockHub)

	userID := int64(1)
	code := "test"
	data := map[string]string{
		"test": "test",
	}

	tmpl := domain.Template{
		Code:  code,
		Title: "Title",
		Body:  "Body",
	}

	mockTemplate.On("FindByCode", mock.Anything, code).Return(tmpl, nil)
	mockRepo.On("Insert", mock.Anything, mock.Anything).Return(nil)

	err := service.Insert(context.Background(), userID, code, data)
	assert.NoError(t, err)

	mockTemplate.AssertCalled(t, "FindByCode", mock.Anything, code)
	mockRepo.AssertCalled(t, "Insert", mock.Anything, mock.Anything)
}

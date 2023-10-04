package mocks

import (
	"context"
	"e-wallet/dto"
	"github.com/stretchr/testify/mock"
)

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) FindByUser(ctx context.Context, userID int64) ([]dto.NotificationData, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]dto.NotificationData), args.Error(1)
}

func (m *MockNotificationService) Insert(ctx context.Context, userID int64, code string, data map[string]string) error {
	args := m.Called(ctx, userID, code, data)
	return args.Error(0)
}

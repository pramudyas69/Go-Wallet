package mocks

import (
	"context"
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) FindByUser(ctx context.Context, userID int64) (notifications []domain.Notification, err error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Notification), args.Error(1)
}

func (m *MockNotificationRepository) Insert(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) Update(ctx context.Context, notification *domain.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

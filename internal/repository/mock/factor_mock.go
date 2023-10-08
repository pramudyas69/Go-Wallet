package mocks

import (
	"context"
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

type MockFactorRepository struct {
	mock.Mock
}

func (m *MockFactorRepository) FindByUserID(ctx context.Context, userID int64) (domain.Factor, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(domain.Factor), args.Error(1)
}

func (m *MockFactorRepository) Insert(ctx context.Context, factor *domain.Factor) error {
	args := m.Called(ctx, factor)
	return args.Error(0)
}

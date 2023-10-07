package mocks

import (
	"context"
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

type MockTopupRepository struct {
	mock.Mock
}

func (m *MockTopupRepository) FindByID(ctx context.Context, id string) (topup domain.TopUp, err error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.TopUp), args.Error(1)
}

func (m *MockTopupRepository) Insert(ctx context.Context, topUp *domain.TopUp) error {
	args := m.Called(ctx, topUp)
	return args.Error(0)
}

func (m *MockTopupRepository) Update(ctx context.Context, topUp *domain.TopUp) error {
	args := m.Called(ctx, topUp)
	return args.Error(0)
}

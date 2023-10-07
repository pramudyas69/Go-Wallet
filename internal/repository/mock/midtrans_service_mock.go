package mocks

import (
	"context"
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

type MockMidtransService struct {
	mock.Mock
}

func (m *MockMidtransService) GenerateSnapURL(ctx context.Context, topUp *domain.TopUp) error {
	args := m.Called(ctx, topUp)
	return args.Error(0)
}

func (m *MockMidtransService) VerifyPayment(ctx context.Context, transactionID string) (bool, error) {
	args := m.Called(ctx, transactionID)
	return args.Bool(0), args.Error(1)
}

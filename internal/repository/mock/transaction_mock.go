package mocks

import (
	"context"
	"e-wallet/domain"

	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository adalah struct yang mengimplementasikan TransactionRepository
type MockTransactionRepository struct {
	mock.Mock
}

// Insert adalah implementasi mock untuk Insert
func (m *MockTransactionRepository) Insert(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

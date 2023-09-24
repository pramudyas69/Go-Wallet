package mocks

import (
	"context"
	"e-wallet/domain"

	"github.com/stretchr/testify/mock"
)

// MockAccountRepository adalah struct yang mengimplementasikan AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

// FindByUserID adalah implementasi mock untuk FindByUserID
func (m *MockAccountRepository) FindByUserID(ctx context.Context, id int64) (domain.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Account), args.Error(1)
}

// FindByAccountNumber adalah implementasi mock untuk FindByAccountNumber
func (m *MockAccountRepository) FindByAccountNumber(ctx context.Context, accNumber string) (domain.Account, error) {
	args := m.Called(ctx, accNumber)
	return args.Get(0).(domain.Account), args.Error(1)
}

// Update adalah implementasi mock untuk Update
func (m *MockAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

// Insert adalah implementasi mock untuk Insert
func (m *MockAccountRepository) Insert(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

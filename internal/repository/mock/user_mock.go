package mocks

import (
	"context"
	_ "database/sql"
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository adalah struct yang mengimplementasikan domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

// Implementasikan metode-metode dari domain.UserRepository

// FindByID adalah implementasi mock untuk FindByID
func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.User), args.Error(1)
}

// FindByUsername adalah implementasi mock untuk FindByUsername
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(domain.User), args.Error(1)
}

// Insert adalah implementasi mock untuk Insert
func (m *MockUserRepository) Insert(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// Update adalah implementasi mock untuk Update
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

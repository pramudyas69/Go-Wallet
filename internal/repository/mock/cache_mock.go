package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockCacheRepository adalah struct yang mengimplementasikan CacheRepository
type MockCacheRepository struct {
	mock.Mock
}

// Get adalah implementasi mock untuk Get
func (m *MockCacheRepository) Get(key string) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

// Set adalah implementasi mock untuk Set
func (m *MockCacheRepository) Set(key string, entry []byte) error {
	args := m.Called(key, entry)
	return args.Error(0)
}

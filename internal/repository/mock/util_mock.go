package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockUtilInterface adalah mock object untuk UtilInterface.
type MockUtilInterface struct {
	mock.Mock
}

// GetTokenGenerator adalah mock untuk fungsi GetTokenGenerator.
func (m *MockUtilInterface) GetTokenGenerator(length int) string {
	args := m.Called(length)
	return args.String(0)
}

// GenerateRandomNumber adalah mock untuk fungsi GenerateRandomNumber.
func (m *MockUtilInterface) GenerateRandomNumber(length int) string {
	args := m.Called(length)
	return args.String(0)
}

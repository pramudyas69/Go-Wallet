package mocks

import (
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

// MockJwtInterface adalah mock untuk JwtInterface
type MockJwtInterface struct {
	mock.Mock
}

// GenerateToken adalah implementasi mock untuk GenerateToken
func (m *MockJwtInterface) GenerateToken(userID int64, email string, expirationMinutes int64) (string, error) {
	args := m.Called(userID, email, expirationMinutes)
	return args.String(0), args.Error(1)
}

// VerifyToken adalah implementasi mock untuk VerifyToken
func (m *MockJwtInterface) VerifyToken(tokenString string) (*domain.CustomClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*domain.CustomClaims), args.Error(1)
}

package mocks

import (
	"context"
	"e-wallet/domain"
	"github.com/stretchr/testify/mock"
)

type MockTemplateRepository struct {
	mock.Mock
}

func (m *MockTemplateRepository) FindByCode(ctx context.Context, code string) (domain.Template, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(domain.Template), args.Error(1)
}

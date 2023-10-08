package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	mocks "e-wallet/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestValidatePin(t *testing.T) {
	mockFactorRepo := new(mocks.MockFactorRepository)

	factorService := NewFactor(mockFactorRepo)

	req := dto.ValidatePinReq{
		UserID: 1,
		Pin:    "12345",
	}

	ctx := context.Background()

	mockFactorRepo.On("FindByUserID", ctx, req.UserID).Return(domain.Factor{
		UserID: 1,
		Pin:    "$2a$12$c.IIJywD/G3TzqKMxgDdn.gCJfjZo7v.cckVteZHBCXIggxMVi9HG",
	}, nil)

	err := factorService.ValidatePin(ctx, req)
	assert.NoError(t, err)

	mockFactorRepo.AssertCalled(t, "FindByUserID", ctx, req.UserID)
}

func TestInsertFactory(t *testing.T) {
	mockFactorRepo := new(mocks.MockFactorRepository)

	factorService := NewFactor(mockFactorRepo)

	req := dto.CreatePinReq{
		UserID: 1,
		PIN:    "12345",
	}

	ctx := context.Background()

	mockFactorRepo.On("FindByUserID", ctx, req.UserID).Return(domain.Factor{}, nil)
	mockFactorRepo.On("Insert", ctx, mock.Anything).Return(nil)

	err := factorService.Insert(ctx, req)
	assert.NoError(t, err)

	mockFactorRepo.AssertCalled(t, "FindByUserID", ctx, req.UserID)
	mockFactorRepo.AssertCalled(t, "Insert", ctx, mock.Anything)
}

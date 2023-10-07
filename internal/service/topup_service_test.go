package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	mocks "e-wallet/internal/repository/mock"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInitializeTopUp(t *testing.T) {
	mockTopUpRepo := new(mocks.MockTopupRepository)
	mockMidtransSvc := new(mocks.MockMidtransService)
	mockAccountRepo := new(mocks.MockAccountRepository)
	mockNotificationSvc := new(mocks.MockNotificationService)
	mockTransactionRepo := new(mocks.MockTransactionRepository)

	service := NewTopUpService(mockTopUpRepo, mockMidtransSvc, mockAccountRepo, mockNotificationSvc, mockTransactionRepo)

	reqTopUp := dto.TopUpReq{
		UserID: 1,
		Amount: 10000,
	}

	ctx := context.Background()

	mockMidtransSvc.On("GenerateSnapURL", ctx, mock.Anything).Return(nil)
	mockTopUpRepo.On("Insert", ctx, mock.Anything).Return(nil)

	_, err := service.InitializeTopUp(ctx, reqTopUp)
	assert.NoError(t, err)

	mockMidtransSvc.AssertCalled(t, "GenerateSnapURL", ctx, mock.Anything)
	mockTopUpRepo.AssertCalled(t, "Insert", ctx, mock.Anything)
}

func TestConfirmedTopUp(t *testing.T) {
	mockTopUpRepo := new(mocks.MockTopupRepository)
	mockMidtransSvc := new(mocks.MockMidtransService)
	mockAccountRepo := new(mocks.MockAccountRepository)
	mockNotificationSvc := new(mocks.MockNotificationService)
	mockTransactionRepo := new(mocks.MockTransactionRepository)

	service := NewTopUpService(mockTopUpRepo, mockMidtransSvc, mockAccountRepo, mockNotificationSvc, mockTransactionRepo)

	ctx := context.TODO()

	topUp := domain.TopUp{
		ID:     "1",
		UserID: 1,
		Amount: 10000,
		Status: 0,
	}

	myAccount := domain.Account{
		ID:            1,
		UserId:        topUp.UserID,
		AccountNumber: "1234567890",
	}

	mockTopUpRepo.On("FindByID", ctx, "1").Return(topUp, nil)
	mockAccountRepo.On("FindByUserID", ctx, topUp.UserID).Return(myAccount, nil)
	myAccount.Balance += topUp.Amount
	mockAccountRepo.On("Update", ctx, &myAccount).Return(nil)
	topUp.Status = 1
	mockTopUpRepo.On("Update", ctx, &topUp).Return(nil)

	mockTransactionRepo.On("Insert", ctx, mock.Anything).Return(nil)
	mockNotificationSvc.On("Insert", ctx, topUp.UserID, "TOPUP_SUCCESS", map[string]string{
		"amount": fmt.Sprintf("%.2f", topUp.Amount),
	}).Return(nil)

	err := service.ConfirmedTopUp(ctx, "1")
	assert.NoError(t, err)

	mockTopUpRepo.AssertCalled(t, "FindByID", ctx, "1")
	mockAccountRepo.AssertCalled(t, "FindByUserID", ctx, topUp.UserID)
	mockAccountRepo.AssertCalled(t, "Update", ctx, &myAccount)
	mockTopUpRepo.AssertCalled(t, "Update", ctx, &topUp)
	mockTransactionRepo.AssertCalled(t, "Insert", ctx, mock.Anything)
	mockNotificationSvc.AssertCalled(t, "Insert", ctx, topUp.UserID, "TOPUP_SUCCESS", map[string]string{
		"amount": fmt.Sprintf("%.2f", topUp.Amount),
	})
}

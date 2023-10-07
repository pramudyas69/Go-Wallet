package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type topupService struct {
	topupReposiroty       domain.TopUpRepository
	midtransService       domain.MidtransService
	accountRepository     domain.AccountRepository
	notificationService   domain.NotificationService
	transactionRepository domain.TransactionRepository
}

func NewTopUpService(topupReposiroty domain.TopUpRepository, midtransService domain.MidtransService, accountRepository domain.AccountRepository, notificationService domain.NotificationService, transactionRepository domain.TransactionRepository) domain.TopUpService {
	return &topupService{
		topupReposiroty:       topupReposiroty,
		midtransService:       midtransService,
		accountRepository:     accountRepository,
		notificationService:   notificationService,
		transactionRepository: transactionRepository,
	}
}

func (t topupService) InitializeTopUp(ctx context.Context, req dto.TopUpReq) (dto.TopUpRes, error) {
	topUp := domain.TopUp{
		ID:     uuid.NewString(),
		UserID: req.UserID,
		Amount: req.Amount,
		Status: 0,
	}
	err := t.midtransService.GenerateSnapURL(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	err = t.topupReposiroty.Insert(ctx, &topUp)
	if err != nil {
		return dto.TopUpRes{}, err
	}

	return dto.TopUpRes{
		SnapURL: topUp.SnapURL,
	}, nil
}

func (t topupService) ConfirmedTopUp(ctx context.Context, id string) error {
	topUp, err := t.topupReposiroty.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if topUp == (domain.TopUp{}) {
		return errors.New("topup not found")
	}

	account, err := t.accountRepository.FindByUserID(ctx, topUp.UserID)
	if err != nil {
		return err
	}

	if account == (domain.Account{}) {
		return domain.ErrAccountNotFound
	}

	account.Balance += topUp.Amount
	err = t.accountRepository.Update(ctx, &account)
	if err != nil {
		return err
	}

	topUp.Status = 1
	err = t.topupReposiroty.Update(ctx, &topUp)
	if err != nil {
		return err
	}

	err = t.transactionRepository.Insert(ctx, &domain.Transaction{
		AccountId:           account.ID,
		SofNumber:           "00",
		DofNumber:           account.AccountNumber,
		Amount:              topUp.Amount,
		TransactionType:     "C",
		TransactionDateTime: time.Now(),
	})
	if err != nil {
		return err
	}

	err = t.notificationService.Insert(ctx, topUp.UserID, "TOPUP_SUCCESS", map[string]string{
		"amount": fmt.Sprintf("%.2f", topUp.Amount),
	})
	if err != nil {
		return err
	}

	return nil
}

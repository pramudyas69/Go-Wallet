package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"encoding/json"
	"fmt"
	"time"
)

type transactionService struct {
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
	emailService          domain.EmailService
	userRepository        domain.UserRepository
}

func NewTransaction(accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
	emailService domain.EmailService,
	userRepository domain.UserRepository) domain.TransactionService {
	return &transactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
		emailService:          emailService,
		userRepository:        userRepository,
	}
}

func (t transactionService) TransferInquiry(ctx context.Context, req dto.TransferInquiryReq) (dto.TransferInquiryRes, error) {
	user := ctx.Value("x-users").(dto.UserData)

	myAccount, err := t.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if myAccount.AccountNumber == "" {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, req.AccountNumber)
	if err != nil {
		return dto.TransferInquiryRes{}, err
	}

	if dofAccount.AccountNumber == "" {
		return dto.TransferInquiryRes{}, domain.ErrAccountNotFound
	}

	if myAccount.Balance < req.Amount {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := util.GetTokenGenertaor(32)

	jsonData, _ := json.Marshal(req)
	_ = t.cacheRepository.Set(inquiryKey, jsonData)
	return dto.TransferInquiryRes{
		InquiryKey: inquiryKey,
	}, nil
}

func (t transactionService) TransferExecute(ctx context.Context, req dto.TransferExecuteReq) error {
	val, err := t.cacheRepository.Get(req.InquiryKey)
	if err != nil {
		return domain.ErrInquiryNotFound
	}

	var reqInq dto.TransferInquiryReq
	_ = json.Unmarshal(val, &reqInq)

	if reqInq.AccountNumber == "" {
		return domain.ErrInquiryNotFound
	}

	user := ctx.Value("x-users").(dto.UserData)
	myAccount, err := t.accountRepository.FindByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	dofAccount, err := t.accountRepository.FindByAccountNumber(ctx, reqInq.AccountNumber)
	if err != nil {
		return err
	}

	debitTransaction := domain.Transaction{
		AccountId:           myAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
		TransactionType:     "D",
		Amount:              reqInq.Amount,
		TransactionDateTime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &debitTransaction)
	if err != nil {
		return err
	}

	creditTransaction := domain.Transaction{
		AccountId:           myAccount.ID,
		SofNumber:           myAccount.AccountNumber,
		DofNumber:           dofAccount.AccountNumber,
		TransactionType:     "C",
		Amount:              reqInq.Amount,
		TransactionDateTime: time.Now(),
	}

	err = t.transactionRepository.Insert(ctx, &creditTransaction)
	if err != nil {
		return err
	}

	myAccount.Balance -= reqInq.Amount
	err = t.accountRepository.Update(ctx, &myAccount)
	if err != nil {
		return err
	}

	dofAccount.Balance += reqInq.Amount
	err = t.accountRepository.Update(ctx, &dofAccount)
	if err != nil {
		return err
	}

	myUser, _ := t.userRepository.FindByID(ctx, myAccount.UserId)
	dofUser, _ := t.userRepository.FindByID(ctx, dofAccount.UserId)

	myUserMsg := fmt.Sprintf("Berhasil Transfer Uang Sebesar Rp.%2.f ke Sdr. %s", reqInq.Amount, dofUser.FullName)
	dofUserMsg := fmt.Sprintf("Menerima Uang Sebesar Rp.%2.f Dari Sdr. %s", reqInq.Amount, myUser.FullName)

	_ = t.emailService.Send(myUser.Email, "Berhasil Transfer!", myUserMsg)
	_ = t.emailService.Send(dofUser.Email, "Menerima Dana!", dofUserMsg)

	return nil
}

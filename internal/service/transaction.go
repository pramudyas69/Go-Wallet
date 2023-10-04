package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

type transactionService struct {
	accountRepository     domain.AccountRepository
	transactionRepository domain.TransactionRepository
	cacheRepository       domain.CacheRepository
	emailService          domain.EmailService
	userRepository        domain.UserRepository
	utilInterface         domain.UtilInterface
	notificationService   domain.NotificationService
}

func NewTransaction(accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	cacheRepository domain.CacheRepository,
	emailService domain.EmailService,
	userRepository domain.UserRepository,
	utilInterface domain.UtilInterface,
	notificationService domain.NotificationService) domain.TransactionService {
	return &transactionService{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		cacheRepository:       cacheRepository,
		emailService:          emailService,
		userRepository:        userRepository,
		utilInterface:         utilInterface,
		notificationService:   notificationService,
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

	if myAccount.Balance <= 0 {
		return dto.TransferInquiryRes{}, domain.ErrInsufficientBalance
	}

	inquiryKey := t.utilInterface.GetTokenGenerator(32)

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

	err = t.cacheRepository.Delete(req.InquiryKey)
	if err != nil {
		return err
	}

	go t.notificationAfterTransfer(myAccount, dofAccount, reqInq.Amount)
	go t.sendEmailAfterTransfer(myAccount, dofAccount, myUser, dofUser, reqInq.Amount)

	return nil
}

func (t transactionService) notificationAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, amount float64) {
	err := t.notificationService.Insert(context.Background(), sofAccount.UserId, "TRANSFER", map[string]string{
		"amount": fmt.Sprintf("%2.f", amount),
	})

	if err != nil {
		log.Fatalf("error when send notification: %v", err.Error())
	}

	err = t.notificationService.Insert(context.Background(), dofAccount.UserId, "TRANSFER_DEST", map[string]string{
		"amount": fmt.Sprintf("%2.f", amount),
	})

	if err != nil {
		log.Fatalf("error when send notification: %v", err.Error())
	}
}

func (t transactionService) sendEmailAfterTransfer(sofAccount domain.Account, dofAccount domain.Account, myUser domain.User, dofUser domain.User, amount float64) {
	myUserMsg := fmt.Sprintf("Berhasil Transfer Uang Sebesar Rp.%2.f ke Sdr. %s", amount, myUser.FullName)
	dofUserMsg := fmt.Sprintf("Menerima Uang Sebesar Rp.%2.f Dari Sdr. %s", amount, dofUser.FullName)

	err := t.emailService.Send(myUser.Email, "Berhasil Transfer!", myUserMsg)
	if err != nil {
		return
	}
	err = t.emailService.Send(dofUser.Email, "Menerima Dana!", dofUserMsg)
	if err != nil {
		return
	}
}

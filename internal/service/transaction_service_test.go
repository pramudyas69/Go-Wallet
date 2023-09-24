package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	mocks "e-wallet/internal/repository/mock"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTransactionService_TransferInquiry(t *testing.T) {
	ctx := context.TODO()
	userData := dto.UserData{
		ID:       1,
		FullName: "John Doe",
		Phone:    "1234567890",
		Username: "johndoe",
	}

	t.Run("Valid transfer inquiry", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)

		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		destinationAccount := domain.Account{
			AccountNumber: "987654321",
			Balance:       50.0,
		}
		mockAccountRepository.On("FindByAccountNumber", ctx, mock.AnythingOfType("string")).Return(destinationAccount, nil)
		mockUtil.On("GetTokenGenerator", 32).Return("inquiry123")
		// Mock the cacheRepository to return no error during Set
		mockCacheRepository.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)

		req := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}

		res, err := transactionService.TransferInquiry(ctx, req)

		assert.NoError(t, err)
		assert.NotEmpty(t, res.InquiryKey)
		assert.Equal(t, "inquiry123", res.InquiryKey)

		// Ensure that the cacheRepository Set method was called with the correct arguments
		mockCacheRepository.AssertCalled(t, "Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"))
	})

	t.Run("Invalid source account", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		// Mock the accountRepository to return an error for source account
		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{}, domain.ErrAccountNotFound)

		req := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}

		res, err := transactionService.TransferInquiry(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, res.InquiryKey)
	})

	t.Run("Invalid destination account", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return an error for destination account
		mockAccountRepository.On("FindByAccountNumber", ctx, mock.AnythingOfType("string")).Return(domain.Account{}, domain.ErrAccountNotFound)

		req := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}

		res, err := transactionService.TransferInquiry(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, res.InquiryKey)
	})

	t.Run("Insufficient balance", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		// Mock the accountRepository to return a valid source account with insufficient balance
		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			AccountNumber: "123456789",
			Balance:       20.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		mockAccountRepository.On("FindByAccountNumber", ctx, mock.AnythingOfType("string")).Return(domain.Account{
			AccountNumber: "987654321",
			Balance:       50.0,
		}, nil)

		req := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}

		res, err := transactionService.TransferInquiry(ctx, req)

		assert.Error(t, err)
		assert.Empty(t, res.InquiryKey)
	})
}

func TestTransactionService_TransferExecute(t *testing.T) {
	ctx := context.Background()
	userData := dto.UserData{
		ID:       1,
		FullName: "John Doe",
		Phone:    "1234567890",
		Username: "johndoe",
	}

	t.Run("Valid transfer execution", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		//ctx = context.WithValue(ctx, "x-users", userData)
		ctx = context.WithValue(context.TODO(), "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			ID:            1,
			UserId:        userData.ID,
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		destinationAccount := domain.Account{
			ID:            2,
			UserId:        3, // Different user ID for destination account
			AccountNumber: "987654321",
			Balance:       50.0,
		}
		mockAccountRepository.On("FindByAccountNumber", ctx, inquiryData.AccountNumber).Return(destinationAccount, nil)

		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(nil).Twice() // Twice for debit and credit

		// Mock the accountRepository Update method
		mockAccountRepository.On("Update", ctx, mock.AnythingOfType("*domain.Account")).Return(nil).Twice() // Twice for source and destination accounts

		myUser := domain.User{
			ID:       1,
			FullName: "John Doe",
			Email:    "john@example.com",
		}
		dofUser := domain.User{
			ID:       3,
			FullName: "Jane Doe",
			Email:    "jane@example.com",
		}
		mockUserRepository.On("FindByID", ctx, myUser.ID).Return(myUser, nil)
		mockUserRepository.On("FindByID", ctx, dofUser.ID).Return(dofUser, nil)

		mockEmailService.On("Send", myUser.Email, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
		mockEmailService.On("Send", dofUser.Email, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.NoError(t, err)

		// Ensure that the cacheRepository Get method was called with the correct argument
		mockCacheRepository.AssertCalled(t, "Get", req.InquiryKey)

		// Ensure that the transactionRepository Insert method was called twice
		mockTransactionRepository.AssertNumberOfCalls(t, "Insert", 2)

		// Ensure that the accountRepository Update method was called twice
		mockAccountRepository.AssertNumberOfCalls(t, "Update", 2)

		// Ensure that the userRepository FindByID method was called for both source and destination users
		//mockUserRepository.AssertCalled(t, "FindByID", myUser.ID)
		//mockUserRepository.AssertCalled(t, "FindByID", dofUser.ID)

		// Ensure that the emailService Send method was called for both source and destination users
		mockEmailService.AssertCalled(t, "Send", myUser.Email, mock.AnythingOfType("string"), mock.AnythingOfType("string"))
		mockEmailService.AssertCalled(t, "Send", dofUser.Email, mock.AnythingOfType("string"), mock.AnythingOfType("string"))
	})

	t.Run("Invalid inquiry data", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return([]byte{}, domain.ErrInquiryNotFound)
		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(domain.ErrInquiryNotFound).Twice() // Twice for debit and credit
		mockAccountRepository.On("Update", ctx, mock.AnythingOfType("*domain.Account")).Return(domain.ErrInquiryNotFound).Twice()         // Twice for source and destination accounts
		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInquiryNotFound, err)
	})

	t.Run("Invalid source account", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{}, domain.ErrAccountNotFound)

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrAccountNotFound, err)
	})

	t.Run("Invalid destination account", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			ID:            1,
			UserId:        userData.ID,
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return an error for destination account
		mockAccountRepository.On("FindByAccountNumber", ctx, inquiryData.AccountNumber).Return(domain.Account{}, domain.ErrAccountNotFound)

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrAccountNotFound, err)
	})

	t.Run("Debit transaction insert error", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			ID:            1,
			UserId:        userData.ID,
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		destinationAccount := domain.Account{
			ID:            2,
			UserId:        3, // Different user ID for destination account
			AccountNumber: "987654321",
			Balance:       50.0,
		}
		mockAccountRepository.On("FindByAccountNumber", ctx, inquiryData.AccountNumber).Return(destinationAccount, nil)

		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(errors.New("debit insert error"))

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "debit insert error", err.Error())
	})

	t.Run("Credit transaction insert error", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		// Mock the accountRepository to return a valid source account

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			ID:            1,
			UserId:        userData.ID,
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		destinationAccount := domain.Account{
			ID:            2,
			UserId:        3, // Different user ID for destination account
			AccountNumber: "987654321",
			Balance:       50.0,
		}
		mockAccountRepository.On("FindByAccountNumber", ctx, inquiryData.AccountNumber).Return(destinationAccount, nil)

		// Mock the transactionRepository to return an error during credit transaction insert
		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(nil) // Debit transaction insert success
		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(errors.New("credit insert error"))
		mockAccountRepository.On("Update", ctx, mock.AnythingOfType("*domain.Account")).Return(errors.New("credit insert error"))

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "credit insert error", err.Error())
	})

	t.Run("Source account update error", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		// Mock the accountRepository to return a valid source account

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			ID:            1,
			UserId:        userData.ID,
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		destinationAccount := domain.Account{
			ID:            2,
			UserId:        3, // Different user ID for destination account
			AccountNumber: "987654321",
			Balance:       50.0,
		}
		mockAccountRepository.On("FindByAccountNumber", ctx, inquiryData.AccountNumber).Return(destinationAccount, nil)

		// Mock the transactionRepository to return no error during transaction inserts
		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(nil).Twice() // Debit and credit transactions

		// Mock the accountRepository to return an error during source account update
		mockAccountRepository.On("Update", ctx, mock.AnythingOfType("*domain.Account")).Return(errors.New("source account update error"))

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "source account update error", err.Error())
	})

	t.Run("Destination account update error", func(t *testing.T) {
		mockAccountRepository := new(mocks.MockAccountRepository)
		mockCacheRepository := new(mocks.MockCacheRepository)
		mockTransactionRepository := new(mocks.MockTransactionRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockUserRepository := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)

		transactionService := NewTransaction(mockAccountRepository, mockTransactionRepository, mockCacheRepository, mockEmailService, mockUserRepository, mockUtil)
		// Mock the context to include user data
		ctx = context.WithValue(ctx, "x-users", userData)

		inquiryData := dto.TransferInquiryReq{
			AccountNumber: "987654321",
			Amount:        25.0,
		}
		jsonData, _ := json.Marshal(inquiryData)
		mockCacheRepository.On("Get", mock.AnythingOfType("string")).Return(jsonData, nil)

		// Mock the accountRepository to return a valid source account

		mockAccountRepository.On("FindByUserID", ctx, userData.ID).Return(domain.Account{
			ID:            1,
			UserId:        userData.ID,
			AccountNumber: "123456789",
			Balance:       100.0,
		}, nil)

		// Mock the accountRepository to return a valid destination account
		destinationAccount := domain.Account{
			ID:            2,
			UserId:        3, // Different user ID for destination account
			AccountNumber: "987654321",
			Balance:       50.0,
		}
		mockAccountRepository.On("FindByAccountNumber", ctx, inquiryData.AccountNumber).Return(destinationAccount, nil)

		// Mock the transactionRepository to return no error during transaction inserts
		mockTransactionRepository.On("Insert", ctx, mock.AnythingOfType("*domain.Transaction")).Return(nil).Twice() // Debit and credit transactions

		// Mock the accountRepository to return no error during source account update
		mockAccountRepository.On("Update", ctx, mock.AnythingOfType("*domain.Account")).Return(nil).Once() // Source account update success

		// Mock the accountRepository to return an error during destination account update
		mockAccountRepository.On("Update", ctx, mock.AnythingOfType("*domain.Account")).Return(errors.New("destination account update error")).Once()

		req := dto.TransferExecuteReq{
			InquiryKey: "inquiry123",
		}

		err := transactionService.TransferExecute(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, "destination account update error", err.Error())
	})
}

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

func TestUserService_Authenticate(t *testing.T) {

	authReq := dto.AuthReq{
		Username: "existinguser",
		Password: "password",
	}

	mockUser := domain.User{
		Username: "existinguser",
		Password: "$2a$12$b.GcCEuq64U71a66Xo4/7u6TCmHLXAMPCeX/GBlt.Dgr.ePEW5nqm",
	}

	t.Run("Success", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		mockUserRepo.On("FindByUsername", mock.Anything, authReq.Username).
			Return(mockUser, nil)
		mockCacheRepo.On("Set", mock.Anything, mock.Anything).
			Return(nil)
		mockUtil.On("GetTokenGenerator", mock.Anything).Return("sample_token")

		authRes, err := userSvc.Authenticate(context.Background(), authReq)
		assert.NoError(t, err)
		assert.NotNil(t, authRes.Token)

		mockUserRepo.AssertExpectations(t)
		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		mockUserRepo.On("FindByUsername", mock.Anything, authReq.Username).
			Return(domain.User{}, errors.New("user not found"))
		_, err := userSvc.Authenticate(context.Background(), authReq)
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("PasswordMismatch", func(t *testing.T) {
		mockUserRepo := new(mocks.MockUserRepository)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		mockUserRepo.On("FindByUsername", mock.Anything, authReq.Username).
			Return(mockUser, nil)

		req := dto.AuthReq{
			Username: "existinguser",
			Password: "wrongpassword",
		}
		_, err := userSvc.Authenticate(context.Background(), req)

		assert.Error(t, err)
		assert.EqualError(t, err, "error authentication failed")

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_ValidateToken(t *testing.T) {
	token := "sample_token"

	mockUserData := dto.UserData{
		ID:       1,
		FullName: "John Doe",
		Phone:    "1234567890",
		Username: "johndoe",
	}

	t.Run("Success", func(t *testing.T) {
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(nil, mockCacheRepo, nil, nil, mockUtil)

		req := domain.User{
			ID:       1,
			FullName: "John Doe",
			Phone:    "1234567890",
			Username: "johndoe",
		}

		dataJson, _ := json.Marshal(req)

		mockCacheRepo.On("Get", "user"+token).
			Return(dataJson, nil)

		userData, err := userSvc.ValidateToken(context.Background(), token)

		assert.NoError(t, err)
		assert.Equal(t, mockUserData, userData)

		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("CacheError", func(t *testing.T) {
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(nil, mockCacheRepo, nil, nil, mockUtil)

		mockCacheRepo.On("Get", "user"+token).
			Return([]byte{}, errors.New("cache error"))

		_, err := userSvc.ValidateToken(context.Background(), token)

		assert.Error(t, err)
		assert.EqualError(t, err, domain.ErrAuthFailed.Error())

		mockCacheRepo.AssertExpectations(t)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(nil, mockCacheRepo, nil, nil, mockUtil)
		mockCacheRepo.On("Get", "user"+token).Return([]byte{}, errors.New("token invalid"))

		_, err := userSvc.ValidateToken(context.Background(), token)

		assert.Error(t, err)

		mockCacheRepo.AssertExpectations(t)
	})
}

func TestUserService_Register(t *testing.T) {

	// Data request untuk registrasi
	registerReq := dto.UserRegisterReq{
		FullName: "John Doe",
		Phone:    "1234567890",
		Username: "newuser",
		Password: "password123",
		Email:    "test@example.com",
	}

	t.Run("Success", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockUserRepo := new(mocks.MockUserRepository)
		mockAccountRepo := new(mocks.MockAccountRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockCacheRepo := new(mocks.MockCacheRepository)
		//mockUtil := new(mocks.MockUtilInterface)

		// Inisialisasi service dengan mock repository dan service
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, mockEmailService, mockAccountRepo, mockUtil)

		// Mock behavior untuk metode FindByUsername pada UserRepository
		mockUserRepo.On("FindByUsername", mock.Anything, registerReq.Username).
			Return(domain.User{}, nil)

		// Mock behavior untuk metode Insert pada UserRepository
		mockUserRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(nil)

		// Mock behavior untuk metode Insert pada AccountRepository
		mockAccountRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.Account")).
			Return(nil)

		mockUtil.On("GenerateRandomNumber", mock.Anything).Return("123456")
		mockUtil.On("GetTokenGenerator", mock.Anything).Return("sample_reference_id")

		mockEmailService.On("Send", registerReq.Email, "OTP Code", "otp anda : 123456").
			Return(nil)

		// Mock behavior untuk metode Set pada CacheRepository
		mockCacheRepo.On("Set", "otp:sample_reference_id", []byte("123456")).
			Return(nil)

		mockCacheRepo.On("Set", "user-ref:sample_reference_id", []byte(registerReq.Username)).
			Return(nil)

		// Panggil metode yang ingin diuji
		res, err := userSvc.Register(context.Background(), registerReq)

		// Assert hasil panggilan metode
		assert.NoError(t, err)
		assert.Equal(t, "sample_reference_id", res.ReferenceID)
	})

	t.Run("UsernameExists", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockUserRepo := new(mocks.MockUserRepository)
		mockAccountRepo := new(mocks.MockAccountRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)

		// Inisialisasi service dengan mock repository dan service
		userSvc := NewUser(mockUserRepo, mockCacheRepo, mockEmailService, mockAccountRepo, mockUtil)

		// Mock behavior untuk metode FindByUsername pada UserRepository
		mockUserRepo.On("FindByUsername", mock.Anything, registerReq.Username).
			Return(domain.User{}, domain.ErrUsernameExist)

		// Panggil metode yang ingin diuji
		_, err := userSvc.Register(context.Background(), registerReq)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, domain.ErrUsernameExist.Error())
	})

	t.Run("UserInsertError", func(t *testing.T) {
		/// Inisialisasi mock repository dan mock service
		mockUserRepo := new(mocks.MockUserRepository)
		mockAccountRepo := new(mocks.MockAccountRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)

		// Inisialisasi service dengan mock repository dan service
		userSvc := NewUser(mockUserRepo, mockCacheRepo, mockEmailService, mockAccountRepo, mockUtil)

		// Mock behavior untuk metode FindByUsername pada UserRepository
		mockUserRepo.On("FindByUsername", mock.Anything, registerReq.Username).
			Return(domain.User{}, nil)

		// Mock behavior untuk metode Insert pada UserRepository
		mockUserRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(errors.New("user insert error"))

		// Panggil metode yang ingin diuji
		_, err := userSvc.Register(context.Background(), registerReq)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, "user insert error")
	})

	t.Run("AccountInsertError", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockUserRepo := new(mocks.MockUserRepository)
		mockAccountRepo := new(mocks.MockAccountRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)

		// Inisialisasi service dengan mock repository dan service
		userSvc := NewUser(mockUserRepo, mockCacheRepo, mockEmailService, mockAccountRepo, mockUtil)

		// Mock behavior untuk metode FindByUsername pada UserRepository
		mockUserRepo.On("FindByUsername", mock.Anything, registerReq.Username).
			Return(domain.User{}, nil)

		// Mock behavior untuk metode Insert pada UserRepository
		mockUserRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(nil)

		// Mock behavior untuk metode Insert pada AccountRepository
		mockAccountRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.Account")).
			Return(errors.New("account insert error"))

		// Panggil metode yang ingin diuji
		_, err := userSvc.Register(context.Background(), registerReq)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, "account insert error")
	})

	t.Run("EmailSendError", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockUserRepo := new(mocks.MockUserRepository)
		mockAccountRepo := new(mocks.MockAccountRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)

		// Inisialisasi service dengan mock repository dan service
		userSvc := NewUser(mockUserRepo, mockCacheRepo, mockEmailService, mockAccountRepo, mockUtil)

		// Mock behavior untuk metode FindByUsername pada UserRepository
		mockUserRepo.On("FindByUsername", mock.Anything, registerReq.Username).
			Return(domain.User{}, nil)

		// Mock behavior untuk metode Insert pada UserRepository
		mockUserRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(nil)

		// Mock behavior untuk metode Insert pada AccountRepository
		mockAccountRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.Account")).
			Return(nil)

		// Mock behavior untuk fungsi bcrypt.GenerateFromPassword
		mockUtil.On("GenerateRandomNumber", mock.Anything).Return("123456")
		mockUtil.On("GetTokenGenerator", mock.Anything).Return("sample_reference_id")

		// Mock behavior untuk metode Send pada EmailService
		mockEmailService.On("Send", registerReq.Email, "OTP Code", "otp anda : 123456").
			Return(errors.New("email send error"))

		// Panggil metode yang ingin diuji
		_, err := userSvc.Register(context.Background(), registerReq)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, "email send error")
	})

	t.Run("CacheSetError", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockUserRepo := new(mocks.MockUserRepository)
		mockAccountRepo := new(mocks.MockAccountRepository)
		mockEmailService := new(mocks.MockEmailService)
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUtil := new(mocks.MockUtilInterface)

		// Inisialisasi service dengan mock repository dan service
		userSvc := NewUser(mockUserRepo, mockCacheRepo, mockEmailService, mockAccountRepo, mockUtil)
		// Mock behavior untuk metode FindByUsername pada UserRepository
		mockUserRepo.On("FindByUsername", mock.Anything, registerReq.Username).
			Return(domain.User{}, nil)

		// Mock behavior untuk metode Insert pada UserRepository
		mockUserRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.User")).
			Return(nil)

		// Mock behavior untuk metode Insert pada AccountRepository
		mockAccountRepo.On("Insert", mock.Anything, mock.AnythingOfType("*domain.Account")).
			Return(nil)

		mockUtil.On("GenerateRandomNumber", mock.Anything).Return("123456")
		mockUtil.On("GetTokenGenerator", mock.Anything).Return("sample_reference_id")

		// Mock behavior untuk metode Send pada EmailService
		mockEmailService.On("Send", registerReq.Email, "OTP Code", "otp anda : 123456").
			Return(nil)

		// Mock behavior untuk metode Set pada CacheRepository
		mockCacheRepo.On("Set", "otp:sample_reference_id", []byte("123456")).
			Return(errors.New("cache set error"))

		mockCacheRepo.On("Set", "user-ref:sample_reference_id", []byte(registerReq.Username)).
			Return(errors.New("cache set error"))

		// Panggil metode yang ingin diuji
		_, err := userSvc.Register(context.Background(), registerReq)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, "cache set error")
	})
}

func TestUserService_ValidateOTP(t *testing.T) {
	req := dto.ValidateOtpReq{
		ReferenceID: "sample_reference_id",
		OTP:         "123456",
	}

	res := domain.User{
		ID:       123456,
		FullName: "John Doe",
		Phone:    "1234567890",
		Username: "johndoe",
	}

	t.Run("Success", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "otp:sample_reference_id").
			Return([]byte(req.OTP), nil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "user-ref:sample_reference_id").
			Return([]byte("123456"), nil)

		mockUserRepo.On("FindByUsername", mock.Anything, "123456").Return(res, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)
		// Panggil metode yang ingin diuji
		err := userSvc.ValidateOTP(context.Background(), req)

		// Assert hasil panggilan metode
		assert.NoError(t, err)
	})

	t.Run("CacheError", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "otp:sample_reference_id").
			Return([]byte{}, domain.ErrOtpInvalid)

		// Panggil metode yang ingin diuji
		err := userSvc.ValidateOTP(context.Background(), req)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, domain.ErrOtpInvalid.Error())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "otp:sample_reference_id").
			Return([]byte(req.OTP), nil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "user-ref:sample_reference_id").
			Return([]byte("123456"), nil)

		mockUserRepo.On("FindByUsername", mock.Anything, "123456").Return(domain.User{}, errors.New("User not found"))

		// Panggil metode yang ingin diuji
		err := userSvc.ValidateOTP(context.Background(), req)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, "User not found")
	})

	t.Run("UpdateFailed", func(t *testing.T) {
		// Inisialisasi mock repository dan mock service
		mockCacheRepo := new(mocks.MockCacheRepository)
		mockUserRepo := new(mocks.MockUserRepository)
		mockUtil := new(mocks.MockUtilInterface)
		userSvc := NewUser(mockUserRepo, mockCacheRepo, nil, nil, mockUtil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "otp:sample_reference_id").
			Return([]byte(req.OTP), nil)

		// Mock behavior untuk metode Get pada CacheRepository
		mockCacheRepo.On("Get", "user-ref:sample_reference_id").
			Return([]byte("123456"), nil)

		mockUserRepo.On("FindByUsername", mock.Anything, "123456").Return(res, nil)
		mockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.User")).Return(errors.New("Update failed"))
		// Panggil metode yang ingin diuji
		err := userSvc.ValidateOTP(context.Background(), req)

		// Assert hasil panggilan metode
		assert.Error(t, err)
		assert.EqualError(t, err, "Update failed")
	})
}

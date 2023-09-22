package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"e-wallet/internal/util"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userService struct {
	userRepository    domain.UserRepository
	cacheRepository   domain.CacheRepository
	emailService      domain.EmailService
	accountRepository domain.AccountRepository
}

func NewUser(userRepository domain.UserRepository, cacheRepository domain.CacheRepository, emailService domain.EmailService, accountRepository domain.AccountRepository) domain.UserService {
	return &userService{
		userRepository:    userRepository,
		cacheRepository:   cacheRepository,
		emailService:      emailService,
		accountRepository: accountRepository,
	}
}

func (u userService) Authenticate(ctx context.Context, req dto.AuthReq) (dto.AuthRes, error) {
	user, err := u.userRepository.FindByUsername(ctx, req.Username)
	if err != nil {
		return dto.AuthRes{}, err
	}

	if user == (domain.User{}) {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthRes{}, domain.ErrAuthFailed
	}

	token := util.GetTokenGenertaor(12)
	userJson, err := json.Marshal(user)
	if err != nil {
		return dto.AuthRes{}, err
	}

	err = u.cacheRepository.Set("user"+token, userJson)
	if err != nil {
		return dto.AuthRes{}, err
	}

	return dto.AuthRes{
		Token: token,
	}, nil
}

func (u userService) ValidateToken(ctx context.Context, token string) (dto.UserData, error) {
	data, err := u.cacheRepository.Get("user" + token)
	if err != nil {
		return dto.UserData{}, domain.ErrAuthFailed
	}

	var user domain.User
	_ = json.Unmarshal(data, &user)

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}

func (u userService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
	fmt.Println("Start")
	exist, err := u.userRepository.FindByUsername(ctx, req.Username)

	if err != nil {

		return dto.UserRegisterRes{}, err
	}
	if exist.Username != "" {
		return dto.UserRegisterRes{}, domain.ErrUsernameExist
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	user := domain.User{
		FullName: req.FullName,
		Phone:    req.Phone,
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
	}

	err = u.userRepository.Insert(ctx, &user)

	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	account := domain.Account{
		UserId:        user.ID,
		AccountNumber: util.GenerateRandomNumber(6),
		Balance:       0,
	}

	err = u.accountRepository.Insert(ctx, &account)

	otpCode := util.GenerateRandomNumber(4)
	referenceId := util.GetTokenGenertaor(16)

	err = u.emailService.Send(req.Email, "OTP Code", "otp anda : "+otpCode)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	fmt.Println("your otp code : ", otpCode)
	_ = u.cacheRepository.Set("otp:"+referenceId, []byte(otpCode))
	_ = u.cacheRepository.Set("user-ref:"+referenceId, []byte(user.Username))

	return dto.UserRegisterRes{
		ReferenceID: referenceId,
	}, nil
}

func (u userService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	val, err := u.cacheRepository.Get("otp:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}
	otp := string(val)
	if otp != req.OTP {
		return domain.ErrOtpInvalid
	}

	val, err = u.cacheRepository.Get("user-ref:" + req.ReferenceID)
	if err != nil {
		return domain.ErrOtpInvalid
	}

	user, err := u.userRepository.FindByUsername(ctx, string(val))
	if err != nil {
		return err
	}
	user.EmailVerifiedAt = time.Now()
	_ = u.userRepository.Update(ctx, &user)
	return nil
}

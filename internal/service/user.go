package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userService struct {
	userRepository    domain.UserRepository
	cacheRepository   domain.CacheRepository
	emailService      domain.EmailService
	accountRepository domain.AccountRepository
	utilInterface     domain.UtilInterface
	jwtInterface      domain.JwtInterface
}

func NewUser(userRepository domain.UserRepository,
	cacheRepository domain.CacheRepository,
	emailService domain.EmailService,
	accountRepository domain.AccountRepository,
	utilInterface domain.UtilInterface,
	jwtInterface domain.JwtInterface) domain.UserService {
	return &userService{
		userRepository:    userRepository,
		cacheRepository:   cacheRepository,
		emailService:      emailService,
		accountRepository: accountRepository,
		utilInterface:     utilInterface,
		jwtInterface:      jwtInterface,
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

	token, err := u.jwtInterface.GenerateToken(user.ID, user.Email, 24)
	if err != nil {
		return dto.AuthRes{}, err
	}

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
	if err := json.Unmarshal(data, &user); err != nil {
		return dto.UserData{}, err
	}

	return dto.UserData{
		ID:       user.ID,
		FullName: user.FullName,
		Phone:    user.Phone,
		Username: user.Username,
	}, nil
}

func (u userService) Register(ctx context.Context, req dto.UserRegisterReq) (dto.UserRegisterRes, error) {
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
		AccountNumber: u.utilInterface.GenerateRandomNumber(6),
		Balance:       0,
	}

	err = u.accountRepository.Insert(ctx, &account)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	otpCode := u.utilInterface.GenerateRandomNumber(4)
	referenceId := u.utilInterface.GetTokenGenerator(16)

	err = u.emailService.Send(req.Email, "OTP Code", "otp anda : "+otpCode)
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	err = u.cacheRepository.Set("otp:"+referenceId, []byte(otpCode))
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

	err = u.cacheRepository.Set("user-ref:"+referenceId, []byte(user.Username))
	if err != nil {
		return dto.UserRegisterRes{}, err
	}

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
	err = u.userRepository.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

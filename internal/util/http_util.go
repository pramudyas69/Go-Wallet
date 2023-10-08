package util

import (
	"e-wallet/domain"
	"errors"
)

func GetHttpStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrAuthFailed):
		return 401
	case errors.Is(err, domain.ErrOtpInvalid):
		return 400
	case errors.Is(err, domain.ErrUsernameExist):
		return 400
	case errors.Is(err, domain.ErrInsufficientBalance):
		return 400
	case errors.Is(err, domain.ErrInquiryNotFound):
		return 404
	case errors.Is(err, domain.ErrAccountNotFound):
		return 404
	case errors.Is(err, domain.ErrPinInvalid):
		return 400
	default:
		return 500
	}
}

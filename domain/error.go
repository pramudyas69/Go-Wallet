package domain

import "errors"

var ErrAuthFailed = errors.New("error authentication failed")
var ErrUsernameExist = errors.New("username is already exist")
var ErrOtpInvalid = errors.New("otp invalid")
var ErrAccountNotFound = errors.New("account not found")
var ErrInquiryNotFound = errors.New("inquiry not found")
var ErrInsufficientBalance = errors.New("insufficient balance")

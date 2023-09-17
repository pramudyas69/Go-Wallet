package domain

import "errors"

var ErrAuthFailed = errors.New("error authentication failed")
var ErrUsernameExist = errors.New("username is already exist")
var ErrOtpInvalid = errors.New("otp invalid")

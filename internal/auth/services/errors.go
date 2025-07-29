package services

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountNotVerified = errors.New("account not verified")
	ErrAccountSuspended   = errors.New("account suspended")
	ErrAccountBlocked     = errors.New("account blocked")
	ErrAccountLockout     = errors.New("account has been blocked due to too many failed login attempts")
	ErrRateLimitExceeded  = errors.New("you have requested too many OTPs, please try again later") // <-- Add this line
)

package apperror

import "errors"

var (
	// Repository Expected Error
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUpdateLastLogin    = errors.New("failed to update last login")
	ErrRegisterUser       = errors.New("failed to register user")

	// Service Expercted Error
	ErrInvalidEmailFormat = errors.New("email must be a valid email address")
	ErrInvalidCredential  = errors.New("invalid email or password")

	ErrOTPNotFound = errors.New("OTP is invalid or has expired")
	ErrOTPExpired  = errors.New("OTP has expired")
)

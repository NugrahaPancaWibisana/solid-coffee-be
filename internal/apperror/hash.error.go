package apperror

import "errors"

var (
	ErrEmptyPassword       = errors.New("password cannot be empty")
	ErrEmptyHash           = errors.New("hash cannot be empty")
	ErrInvalidHashFormat   = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible argon2 version")
)

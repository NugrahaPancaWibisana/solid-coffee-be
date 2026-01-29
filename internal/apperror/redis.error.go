package apperror

import "errors"

var (
	ErrSessionExpired = errors.New("session expired, please login again")
	ErrInvalidSession = errors.New("invalid session, please login again")
	ErrInternal       = errors.New("internal server error")
	ErrLogoutFailed   = errors.New("failed to logout")
)

package apperror

import "errors"

var (
	ErrUpdateProfile    = errors.New("failed to update profile")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
	ErrGetPassword      = errors.New("failed to get password")
	ErrUpdatePassword   = errors.New("failed to update password")
	ErrVerifyPassword   = errors.New("old password is incorrect")
)

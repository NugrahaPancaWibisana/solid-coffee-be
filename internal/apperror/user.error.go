package apperror

import "errors"

var (
	ErrUpdateProfile    = errors.New("failed to update profile")
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

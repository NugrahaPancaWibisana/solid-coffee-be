package apperror

import "errors"

var (
	ErrMenuNotFound = errors.New("menu not found")
	ErrGetMenu      = errors.New("failed to get menu")
	ErrUpdateMenu   = errors.New("no fields to update")
	ErrDeleteMenu   = errors.New("failed to delete menu")
)

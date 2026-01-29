package apperror

import "errors"

var (
	ErrSecretNotFound     = errors.New("JWT secret not found in environment")
	ErrIssuerNotFound     = errors.New("JWT issuer not found in environment")
	ErrInvalidIssuer      = errors.New("token issuer does not match expected issuer")
	ErrTokenInvalid       = errors.New("token is invalid")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenClaimsInvalid = errors.New("token claims are invalid")
)
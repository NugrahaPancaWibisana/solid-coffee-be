package apperror

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("User not found")
	ErrEmailAlreadyExists = errors.New("Email already exists")
	ErrUpdateLastLogin    = errors.New("Failed to update last login")
	ErrRegisterUser       = errors.New("Failed to register user")
	ErrInvalidEmailFormat = errors.New("Invalid email format")
	ErrInvalidCredential  = errors.New("Invalid email or password")
	ErrInsertUser         = errors.New("Failed to insert user")
	ErrDeleteUser         = errors.New("Failed to delete user")
	ErrGetUsers           = errors.New("Failed to retrieve users")

	// OTP errors
	ErrOTPNotFound = errors.New("Invalid or expired OTP")
	ErrOTPExpired  = errors.New("OTP has expired")

	// Password/Hash errors
	ErrEmptyPassword       = errors.New("Password cannot be empty")
	ErrEmptyHash           = errors.New("Hash cannot be empty")
	ErrInvalidHashFormat   = errors.New("Invalid hash format")
	ErrIncompatibleVersion = errors.New("Incompatible Argon2 version")

	// JWT errors
	ErrSecretNotFound     = errors.New("JWT secret not found in environment")
	ErrIssuerNotFound     = errors.New("JWT issuer not found in environment")
	ErrInvalidIssuer      = errors.New("Invalid token issuer")
	ErrTokenInvalid       = errors.New("Invalid token")
	ErrTokenExpired       = errors.New("Token has expired")
	ErrTokenClaimsInvalid = errors.New("Invalid token claims")

	// Menu errors
	ErrMenuNotFound = errors.New("Menu not found")
	ErrGetMenu      = errors.New("Failed to retrieve menu")
	ErrUpdateMenu   = errors.New("Failed to update menu")
	ErrDeleteMenu   = errors.New("Failed to delete menu")

	// Session errors
	ErrSessionExpired = errors.New("Session expired, please login again")
	ErrInvalidSession = errors.New("Invalid session, please login again")
	ErrLogoutFailed   = errors.New("Failed to logout")

	// Profile errors
	ErrUpdateProfile    = errors.New("Failed to update profile")
	ErrNoFieldsToUpdate = errors.New("No fields to update")
	ErrGetPassword      = errors.New("Failed to retrieve password")
	ErrUpdatePassword   = errors.New("Failed to update password")
	ErrVerifyPassword   = errors.New("Incorrect old password")
	ErrGetProfile       = errors.New("Failed to retrieve user profile")
	ErrProfileNotFound  = errors.New("User profile not found")

	// Generic errors
	ErrInternal = errors.New("Internal server error")
)
